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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/jsonpretty"
	"github.com/spf13/cobra"
)

// Represents the response from a successful delete request
type deletedOK struct {
	NextAnalysis  string `json:"next_analysis_url"`
	ConfirmDelete string `json:"confirm_delete_url"`
}

// deleteAnalysis sends a DELETE request to the GitHub API to delete an analysis.
func deleteAnalysis(a string) (*http.Response, error) {
	u, err := url.Parse(a)
	if err != nil {
		return nil, err
	}

	var opts api.ClientOptions
	client, err := api.NewRESTClient(opts)
	if err != nil {
		return nil, err
	}

	response, err := client.Request(http.MethodDelete, u.String(), nil)
	if err != nil {
		// A 400 will indicate that the analysis is not deletable and is probably not the most recent in the set.
		// Response is nil when this happens, so check the error instead to return a more friendly error message.
		if strings.Contains(err.Error(), "HTTP: 400 Analysis is last of its type and deletion may result in the loss of historical alert data.") {
			return nil, fmt.Errorf("HTTP 400: Analysis is last of its type and deletion may result in the loss of historical alert data. Please specify --confirm-delete")
		}
		return nil, err
	}

	return response, nil
}

// deleteAllAnalyses sends DELETE requests to the GitHub API to delete all analyses in a set.
// Respects the --confirm-delete flag if set.
// Returns a slice of the analysis IDs that were deleted.
func deleteAllAnalyses(u string) ([]string, error) {
	var deletedAnalyses []string
	for {
		r, err := deleteAnalysis(u)
		if err != nil {
			// A 400 will indicate that the analysis is not deletable, and there is nothing left to do.
			// This is not an error, and indicates there are no other analyses to delete.
			if strings.Contains(err.Error(), "400") || strings.Contains(err.Error(), "404") {
				break
			}
			// Other errors are unexpected and should be returned.
			return nil, err
		}

		var n deletedOK
		n, err = getDeleteResponse(r)
		if err != nil {
			fmt.Println(err)
			break
		}
		// The last URL segment is the analysis ID.
		// Store this so we can return a list of deleted analyses.
		// We parse the URL to avoid the inclusion of query parameters.
		analysisURL, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		us := strings.Split(analysisURL.Path, "/")
		deletedAnalyses = append(deletedAnalyses, us[len(us)-1])
		// Use confirm delete url if --confirm-delete was used.
		if confirmDeleteFlag {
			// If there are no more analyses to delete, break the loop.
			if n.ConfirmDelete == "" {
				break
			}
			u = n.ConfirmDelete
			continue
		}
		// If there are no more analyses to delete, break the loop.
		if n.NextAnalysis == "" {
			break
		}
		u = n.NextAnalysis
	}
	return deletedAnalyses, nil
}

// getDeleteResponse unmarshals the response from a successful delete request.
func getDeleteResponse(r *http.Response) (deletedOK, error) {
	var d deletedOK
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(b, &d)
	if err != nil {
		return d, err
	}
	return d, nil
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [flags] <analysis-id>...",
	Short: "Delete a GitHub Code Scanning Analysis",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Setup Repository
		repo, err := GetRepository()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Making purge an alias for --delete-all --confirm-delete
		if purgeFlag {
			deleteAllFlag = true
			confirmDeleteFlag = true
		}

		// Cannot use --delete-all or --purge with multiple analysis IDs.
		if deleteAllFlag && len(args) > 1 {
			fmt.Println("Cannot use --delete-all or --purge with multiple analysis IDs.")
			return
		}

		// Delete all analyses provided in args.
		var deletedAnalyses []string
		for _, arg := range args {
			baseURL := fmt.Sprintf("repos/%v/%v/code-scanning/analyses/%v", repo.Owner, repo.Name, arg)

			// Delete all analyses in the set.
			// If --confirm-delete or --purge is not set, the last analysis will not be deleted.
			if deleteAllFlag {
				opts := ""
				if confirmDeleteFlag {
					opts = "?confirm_delete"
				}
				n, err := deleteAllAnalyses(baseURL + opts)
				if err != nil {
					fmt.Println(err)
					return
				}
				deletedAnalyses = append(deletedAnalyses, n...)
				continue
			}

			// Delete a single analysis
			var r *http.Response
			if confirmDeleteFlag {
				r, err = deleteAnalysis(baseURL + `?confirm_delete`)
			} else {
				r, err = deleteAnalysis(baseURL)
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			deletedAnalyses = append(deletedAnalyses, arg)

			d, err := getDeleteResponse(r)
			if err != nil {
				fmt.Println(err)
				return
			}
			if jsonFlag {
				jsonpretty.Format(os.Stdout, r.Body, "\t", true)
				return
			}
			fmt.Printf("Successfully deleted analysis %v", arg)
			if d.NextAnalysis != "" {
				fmt.Printf("Next analysis: %v\n", d.NextAnalysis)
			} else {
				fmt.Println("Next analysis: None (last in set)")
			}
		}
		if jsonFlag {
			j, _ := json.Marshal(deletedAnalyses)
			reader := strings.NewReader(string(j))
			jsonpretty.Format(os.Stdout, reader, "\t", true)
			return
		}
		fmt.Printf("Successfully deleted %v analyses.\n", len(deletedAnalyses))
	},
}

var deleteAllFlag bool
var confirmDeleteFlag bool
var purgeFlag bool

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteAllFlag, "delete-all", false, "Delete all analyses in the set, except the last.")
	deleteCmd.Flags().BoolVar(&confirmDeleteFlag, "confirm-delete", false, "Allow the deletion of the last analysis in the set.")
	deleteCmd.Flags().BoolVar(&purgeFlag, "purge", false, "Alias for --delete-all --confirm-delete .")
}
