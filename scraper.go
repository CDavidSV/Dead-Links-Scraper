package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Scraper struct {
	url          *url.URL
	rawUrl       string
	visitedLinks map[string]struct{}
	mu           *sync.RWMutex
	verbose      bool
	httpClient   *http.Client
	semaphore    chan struct{}
}

type Result struct {
	LiveLinks []PageState
	DeadLinks []PageState
}

type PageState struct {
	RawUrl     string
	StatusCode int
	isDead     bool
}

func validateURL(url *url.URL) error {
	if url.Scheme != "http" && url.Scheme != "https" {
		return ErrInvalidSchema
	}

	if url.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

func getHref(n *html.Node) string {
	if n.Type != html.ElementNode || n.Data != "a" {
		return ""
	}

	link := ""
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link = attr.Val
		}
	}

	return link
}

func NewScraper(urlInput string, verbose bool, maxThreads int) (*Scraper, error) {
	// Verify that the link is valid
	url, err := url.Parse(urlInput)
	if err != nil {
		return nil, err
	}

	if err = validateURL(url); err != nil {
		return nil, err
	}

	return &Scraper{
		url:          url,
		rawUrl:       urlInput,
		visitedLinks: make(map[string]struct{}),
		mu:           &sync.RWMutex{},
		verbose:      verbose,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		semaphore: make(chan struct{}, maxThreads),
	}, nil
}

func (s *Scraper) log(msg string) {
	if !s.verbose {
		return
	}

	fmt.Println(msg)
}

func (s *Scraper) findAllLinks(htmlNode *html.Node, linksFound []string) []string {
	for n := htmlNode; n != nil; n = n.NextSibling {
		link := getHref(n)

		if n.FirstChild != nil {
			linksFound = s.findAllLinks(n.FirstChild, linksFound)
		}

		if link != "" {
			linksFound = append(linksFound, link)
		}
	}

	return linksFound
}

func (s *Scraper) processPage(inputURL string, linksChan chan string, resultChan chan PageState) {
	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	var links []string

	// Validate link
	url, err := url.Parse(inputURL)
	if err != nil {
		s.log(ErrorStyle.Render(fmt.Sprintf("Invalid link: %s", err.Error())))
		return
	}

	if url.Host == "" && url.Path != "" {
		inputURL = s.url.Scheme + "://" + s.url.Host + url.Path

		if url.RawQuery != "" {
			inputURL += "?" + url.RawQuery
		}
	}

	if url.Host != "" && url.Host != s.url.Host {
		s.log(WarningStyle.Render(fmt.Sprintf("Skipping external link: %s", inputURL)))
		return
	}

	s.log(InfoStyle.Render("Scanning: " + inputURL))

	response, err := s.httpClient.Get(inputURL)
	if err != nil {
		s.log(ErrorStyle.Render(fmt.Sprintf("Error fetching page: %s", err.Error())))
		resultChan <- PageState{RawUrl: inputURL, isDead: true, StatusCode: -1}
		return
	}
	defer response.Body.Close()

	if response.StatusCode > 400 {
		s.log(ErrorStyle.Render(fmt.Sprintf("%s responded with status code %d", inputURL, response.StatusCode)))
		resultChan <- PageState{RawUrl: inputURL, isDead: true, StatusCode: response.StatusCode}
		return
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		// If the page is not valid HTML, we can't find any links
		s.log(WarningStyle.Render(fmt.Sprintf("Error parsing page: %s", err.Error())))
		resultChan <- PageState{RawUrl: inputURL, isDead: false, StatusCode: response.StatusCode}
		return
	}

	resultChan <- PageState{RawUrl: inputURL, isDead: false, StatusCode: response.StatusCode}
	links = s.findAllLinks(doc, links)
	for _, links := range links {
		linksChan <- links
	}
}

func (s *Scraper) setVisited(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.visitedLinks[url] = struct{}{}
}

func (s *Scraper) visited(url string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.visitedLinks[url]

	return ok
}

func (s *Scraper) Run() Result {
	linksChan := make(chan string)
	processedLinksChan := make(chan PageState)

	var wg sync.WaitGroup

	r := Result{
		LiveLinks: make([]PageState, 0),
		DeadLinks: make([]PageState, 0),
	}

	go func() {
		for link := range linksChan {
			if s.visited(link) {
				continue
			}

			s.setVisited(link)

			wg.Add(1)
			go func(link string) {
				defer wg.Done()
				s.processPage(link, linksChan, processedLinksChan)
			}(link)
		}
	}()

	// Start the process with the initial URL
	linksChan <- s.rawUrl

	go func() {
		// Wait for all gorutines to finish before closing the channel
		wg.Wait()
		close(linksChan)
		close(processedLinksChan)
	}()

	for pageState := range processedLinksChan {
		if pageState.isDead {
			r.DeadLinks = append(r.DeadLinks, pageState)
		} else {
			r.LiveLinks = append(r.LiveLinks, pageState)
		}
	}

	return r
}
