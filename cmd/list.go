/*
Copyright © 2023 Kynan Ware
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/jsonpretty"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
)

var repo repository.Repository
var refFlag string
var toolNameFlag string
var pageFlag int
var limitFlag int
var jsonFlag bool

const defaultLimit = 30

type Analysis struct {
	Ref          string `json:"ref"`
	CommitSha    string `json:"commit_sha"`
	AnalysisKey  string `json:"analysis_key"`
	Environment  string `json:"environment"`
	Category     string `json:"category"`
	Error        string `json:"error"`
	CreatedAt    string `json:"created_at"`
	ResultsCount int    `json:"results_count"`
	RulesCount   int    `json:"rules_count"`
	ID           int    `json:"id"`
	URL          string `json:"url"`
	SarifID      string `json:"sarif_id"`
	Tool         struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Guid    string `json:"guid"`
	} `json:"tool"`
	Deletable bool   `json:"deletable"`
	Warning   string `json:"warning"`
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List analyses for a repository",
	Long:  `List analyses for a repository. By default, the most recent 30 analyses are listed.`,
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

		baseURL := fmt.Sprintf("repos/%v/%v/code-scanning/analyses", repo.Owner, repo.Name)
		u, err := url.Parse(baseURL)
		if err != nil {
			fmt.Println(err)
		}

		params := url.Values{}
		params.Add("page", "1")

		if limitFlag != defaultLimit {
			params.Add("per_page", fmt.Sprintf("%v", limitFlag))
		}

		if refFlag != "" {
			params.Add("ref", refFlag)
		}

		if toolNameFlag != "" {
			params.Add("tool_name", toolNameFlag)
		}

		if pageFlag != 1 {
			params.Add("page", fmt.Sprintf("%v", pageFlag))
		}

		u.RawQuery = params.Encode()

		client, err := api.DefaultRESTClient()
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
		bodyString := string(bodyBytes)

		if err != nil {
			fmt.Println(err)
			return
		}

		if jsonFlag {
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

		cyan := func(s string) string {
			return "\u001B[96m" + s + "\u001B[39m"
		}

		red := func(s string) string {
			return "\u001B[91m" + s + "\u001B[39m"
		}

		yellow := func(s string) string {
			return "\u001B[93m" + s + "\u001B[39m"
		}

		var bodyJSON []Analysis
		err = json.Unmarshal(bodyBytes, &bodyJSON)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Table Print
		terminal := term.FromEnv()
		termWidth, _, _ := terminal.Size()
		t := tableprinter.New(terminal.Out(), terminal.IsTerminalOutput(), termWidth)

		if terminal.IsTerminalOutput() {
			fmt.Printf("Showing %d analyses on page %d/?\n\n", len(bodyJSON), pageFlag)
		}

		t.AddField("ID")
		t.AddField("Created At")
		t.AddField("Ref")
		t.AddField("Tool")
		t.AddField("Category")
		t.AddField("Rules Count")
		t.AddField("Results Count")
		t.AddField("Deleteable")
		t.EndRow()
		for _, analysis := range bodyJSON {
			var state func(s string) string
			if analysis.Error != "" {
				state = red
			} else if analysis.Warning != "" {
				state = yellow
			} else {
				state = cyan
			}

			t.AddField(strconv.Itoa(analysis.ID), tableprinter.WithColor(state), tableprinter.WithTruncate(nil))
			t.AddField(analysis.CreatedAt)
			t.AddField(analysis.Ref)
			t.AddField(fmt.Sprintf(`%v@%v`, analysis.Tool.Name, analysis.Tool.Version))
			t.AddField(analysis.Category)
			t.AddField(strconv.Itoa(analysis.RulesCount))
			t.AddField(strconv.Itoa(analysis.ResultsCount))
			t.AddField(strconv.FormatBool(analysis.Deletable))

			t.EndRow()
		}
		if err := t.Render(); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&refFlag, "ref", "r", "", "Git ref (branch, tag, or SHA)")
	listCmd.Flags().StringVarP(&toolNameFlag, "tool", "t", "", "Tool name")
	listCmd.Flags().IntVarP(&pageFlag, "page", "p", 1, "Page number of analyses to return")
	listCmd.Flags().IntVarP(&limitFlag, "limit", "L", defaultLimit, "Number of analyses to return per page (default 30, max 100)")
	listCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output JSON instead of text (includes additional fields)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
