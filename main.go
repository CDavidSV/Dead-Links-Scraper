package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
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
		fmt.Println("Total links scanned: ", r.TotalLinks)
		fmt.Printf("Dead links found: %d/%d\n", len(r.DeadLinks), len(r.LiveLinks)+len(r.DeadLinks))

		if len(r.DeadLinks) <= 0 {
			return
		}

		rows := make([][]string, len(r.DeadLinks))

		for i, pageStatus := range r.DeadLinks {
			row := []string{strconv.Itoa(i + 1), fmt.Sprintf("%d %s", pageStatus.StatusCode, http.StatusText(pageStatus.StatusCode)), pageStatus.RawUrl}
			rows[i] = row
		}

		table := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(TableBorderStyle).
			Headers("Num.", "Status", "URL").
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == -1:
					return TableHeaderStyle
				case row%2 == 0 && (col == 0):
					return TableEvenRowStyle.Align(lipgloss.Center)
				case row%2 != 0 && (col == 0):
					return TableOddRowStyle.Align(lipgloss.Center)
				case row%2 == 0:
					return TableEvenRowStyle
				default:
					return TableOddRowStyle
				}
			}).
			Rows(rows...)

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
