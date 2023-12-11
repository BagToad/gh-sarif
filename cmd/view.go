/*
Copyright © 2023 Kynan Ware
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/jsonpretty"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
)

var sarifFlag bool

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a Code Scanning analysis",
	Long: `View the information about a Code Scanning analysis given the ID.
	
	Use --sarif to view the raw SARIF file associated with the analysis.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Setup Repository
		if repoFlag != "" {
			var err error
			repo, err = repository.Parse(repoFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			var err error
			repo, err = repository.Current()
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		baseURL := fmt.Sprintf("repos/%v/%v/code-scanning/analyses/%v", repo.Owner, repo.Name, args[0])
		u, err := url.Parse(baseURL)
		if err != nil {
			fmt.Println(err)
			return
		}

		params := url.Values{}

		// Add params here

		// if sarifFlag {
		// 	params.Add("sarif", "true")
		// }+

		u.RawQuery = params.Encode()

		var opts api.ClientOptions
		if sarifFlag {
			fmt.Println("SARIF")
			opts = api.ClientOptions{
				Headers: map[string]string{"Accept": "application/sarif+json"},
			}
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

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		bodyString := string(bodyBytes)

		if jsonFlag || sarifFlag {
			writer := os.Stdout
			defer writer.Close()

			if err != nil {
				fmt.Println(err)
				return
			}

			reader := bytes.NewBufferString(bodyString)

			err = jsonpretty.Format(writer, reader, "\t", true)
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.Flags().BoolVarP(&sarifFlag, "sarif", "S", false, "Print raw SARIF to stdout")
}
