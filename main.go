package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "deadlinks",
	Short: "Checks for dead links in a website",
	Long:  "Analyzes and detects dead links in a webpage provided by a url. It avoids urls outside the original domain and returns all dead links found in the website.",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		verbose, _ := cmd.Flags().GetBool("verbose")
		maxThreads, _ := cmd.Flags().GetInt("threads")

		scraper, err := NewScraper(url, verbose, maxThreads)
		if err != nil {
			fmt.Println(ErrorStyle.Render(fmt.Sprintf("Invalid url provided: %s", err.Error())))
			return
		}

		start := time.Now()
		r := scraper.Run()
		fmt.Println("Took: ", time.Since(start))

		fmt.Printf("Dead links found: %d/%d\n", len(r.DeadLinks), len(r.LiveLinks)+len(r.DeadLinks))

		if len(r.DeadLinks) <= 0 {
			return
		}

		table := table.New()
		fmt.Println(table)

	},
}

func init() {
	rootCmd.Flags().StringP("url", "u", "", "url for the webpage to be analyzed")
	rootCmd.Flags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.Flags().IntP("threads", "t", 4, "number of concurrent threads to use")

	rootCmd.MarkFlagRequired("url")
}

func main() {
	rootCmd.Execute()
}
