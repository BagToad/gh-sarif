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
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/owenrumney/go-sarif/v2/sarif"
	"github.com/spf13/cobra"
)

type sarifUpload struct {
	CommitSha string `json:"commit_sha"`
	Ref       string `json:"ref"`
	Sarif     string `json:"sarif"`
	Validate  bool   `json:"validate"`
}

type uploadedOK struct {
	UploadID  string `json:"id"`
	UploadURL string `json:"url"`
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload [flags] <commit_sha> <ref> <sarif_file>",
	Short: "Upload a SARIF file to a repo",
	Long:  ``,
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		// Setup Repository
		repo, err := GetRepository()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Read the SARIF file
		sarifBytes, err := os.ReadFile(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		// A preliminary check to see if the file is a valid SARIF file.
		if _, err := sarif.FromBytes(sarifBytes); err != nil {
			fmt.Println(err)
			return
		}

		// gzip compress the file
		var gBuff bytes.Buffer
		gWriter := gzip.NewWriter(&gBuff)
		defer gWriter.Close()
		if _, err = gWriter.Write(sarifBytes); err != nil {
			fmt.Println(err)
			return
		}
		gWriter.Close()

		// Base64 encode the compressed file
		var base64Buffer bytes.Buffer
		b64Encoder := base64.NewEncoder(base64.RawStdEncoding, &base64Buffer)
		defer b64Encoder.Close()
		if _, err := b64Encoder.Write(gBuff.Bytes()); err != nil {
			fmt.Println(err)
			return
		}
		b64Encoder.Close()

		// Set the body
		body := sarifUpload{
			CommitSha: args[0],
			Ref:       args[1],
			Sarif:     base64Buffer.String(),
			Validate:  true,
		}
		var requestBody bytes.Buffer
		jsonEncoder := json.NewEncoder(&requestBody)
		jsonEncoder.Encode(body)

		// Upload the SARIF file
		baseURL := fmt.Sprintf("repos/%v/%v/code-scanning/sarifs", repo.Owner, repo.Name)
		u, err := url.Parse(baseURL)
		if err != nil {
			fmt.Println(err)
			return
		}

		var opts api.ClientOptions
		// If GHES, set the host
		if repo.Host != "" {
			opts.Host = repo.Host
		}
		opts.Headers = map[string]string{"Accept": "application/json"}

		client, err := api.NewRESTClient(opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		response, err := client.Request(http.MethodPost, u.String(), &requestBody)
		if err != nil {
			fmt.Println(err)
			return
		}
		if response.StatusCode != http.StatusAccepted {
			fmt.Println("Failed to upload SARIF file.")
			return
		}
		b, err := io.ReadAll(response.Body)
		var uOK uploadedOK
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(b, &uOK)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("SARIF file uploaded successfully.\n\nID: %v\nURL: %v\n", uOK.UploadID, uOK.UploadURL)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
