/*
Copyright © 2024 Kynan Ware

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/jsonpretty"
	"github.com/cli/go-gh/v2/pkg/markdown"
	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/owenrumney/go-sarif/v2/sarif"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view [<analysis-id> | <sarif-file>] [--sarif | --csv | --json]",
	Short: "View GitHub Code Scanning analysis or SARIF results",
	Long: `View results given the GitHub analysis ID or SARIF file.
	
	Use --sarif to get a subset of the analysis SARIF from GitHub.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Setup Repository
		repo, err := GetRepository()
		if err != nil {
			fmt.Println(err)
			return
		}

		terminal := term.FromEnv()
		isTerminal := terminal.IsTerminalOutput()
		var b []byte // SARIF bytes

		// Check if the argument is a file path or an analysis ID
		// If file path, read the file and parse it as SARIF
		// If analysis ID, make a request to the API to get the SARIF
		if f, _ := os.Stat(args[0]); f != nil {
			b, err = os.ReadFile(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			baseURL := fmt.Sprintf("repos/%v/%v/code-scanning/analyses/%v", repo.Owner, repo.Name, args[0])
			u, err := url.Parse(baseURL)
			if err != nil {
				fmt.Println(err)
				return
			}

			var opts api.ClientOptions
			// Always get the SARIF directly unless the JSON meta is requested instead
			if !jsonFlag {
				opts = api.ClientOptions{
					Headers: map[string]string{"Accept": "application/sarif+json"},
				}
			}
			// If GHES, set the host
			if repo.Host != "" {
				opts.Host = repo.Host
			}

			client, err := api.NewRESTClient(opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			response, err := client.Request(http.MethodGet, u.String(), nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			b, err = io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		// Parse the SARIF
		r, err := sarif.FromBytes(b)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Write pretty JSON or SARIF to stdout.
		// JSON is the analysis metadata in JSON, not the actual SARIF.
		// SARIF is the complete SARIF file.
		if jsonFlag || sarifFlag {
			s := string(b)
			if isTerminal {
				if err = jsonpretty.Format(os.Stdout, bytes.NewBufferString(s), "\t", true); err != nil {
					fmt.Println(err)
				}
			} else {
				if err = jsonpretty.Format(os.Stdout, bytes.NewBufferString(s), "\t", false); err != nil {
					fmt.Println(err)
				}
			}
			return
		}

		// Print results to stdout in a table if no other options.
		termWidth, _, _ := terminal.Size()
		t := tableprinter.New(terminal.Out(), terminal.IsTerminalOutput(), termWidth)

		t.AddHeader([]string{"Rule", "Description", "Alert Number", "Severity"})
		if len(r.Runs) <= 0 {
			fmt.Println("No results found.")
			return
		}

		// Print results in a table
		empty := true
		for _, run := range r.Runs {
			// This means this run has no results
			if len(run.Results) <= 0 {
				continue
			}
			if csvFlag {
				fmt.Printf("%v,%v,%v,%v\n", "Rule ID", "Description", "Alert Number", "Severity")
			}
			empty = false
			for _, result := range run.Results {
				if result.Message.Text == nil {
					result.Message.Text = new(string)
					*result.Message.Text = "-"
				}
				if result.RuleID == nil {
					result.RuleID = new(string)
					*result.RuleID = "-"
				}
				if result.Level == nil {
					result.Level = new(string)
					*result.Level = "-"
				}
				if _, ok := result.Properties["github/alertNumber"]; !ok {
					if result.Properties == nil {
						result.Properties = sarif.Properties{}
					}
					result.Properties["github/alertNumber"] = "-"
				}
				// Print to CSV if flag is set
				if csvFlag {
					m := *result.Message.Text
					// Attempts at making the CSV output valid. Might need to revist formally.
					m = strings.ReplaceAll(m, "\n", " ")
					m = strings.ReplaceAll(m, "\r", " ")
					m = strings.ReplaceAll(m, `"`, `""`)
					fmt.Printf("%v,%v,%v,%v\n", *result.RuleID, `"`+m+`"`, result.Properties["github/alertNumber"], *result.Level)
					continue
				}

				m := *result.Message.Text
				if strings.Contains(m, "\n") {
					m = strings.Split(m, "\n")[0]
					m += " ..."
				}
				// Render markdown in the description
				m, err := markdown.Render(m, markdown.WithTheme("dark"))
				if err != nil {
					fmt.Println(err)
					return
				}
				m = strings.ReplaceAll(m, "\n", "")
				m = strings.ReplaceAll(m, "\r", "")
				m = strings.Join(strings.Fields(m), " ")

				t.AddField(*result.RuleID)
				t.AddField(m)
				t.AddField(fmt.Sprintf("%v", result.Properties["github/alertNumber"]))
				t.AddField(*result.Level)
				t.EndRow()
			}
		}
		// If no results in any runs within the analysis...
		if empty {
			fmt.Println("No results found in analysis.")
			return
		}
		// Don't try to render table if output is CSV.
		if csvFlag {
			return
		}
		if err := t.Render(); err != nil {
			fmt.Println(err)
			return
		}

	},
}

var sarifFlag bool
var csvFlag bool

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.Flags().BoolVarP(&sarifFlag, "sarif", "S", false, "Print raw SARIF to stdout")
	viewCmd.Flags().BoolVarP(&csvFlag, "csv", "c", false, "Print results in CSV format")
}
