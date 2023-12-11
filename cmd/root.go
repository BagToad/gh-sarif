/*
Copyright © 2023 Kynan Ware
*/
package cmd

import (
	"os"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
	// "github.com/cli/go-gh/v2/pkg/api"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-sarif",
	Short: "Work with Code Scanning analyses in GitHub",
	Long:  `Work with Code Scanning analyses in GitHub`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var repoFlag string
var jsonFlag bool
var repo repository.Repository

func init() {
	// ROOT FLAGS
	rootCmd.PersistentFlags().StringVarP(&repoFlag, "repo", "R", "", "GitHub repository (format: owner/repo)")
	rootCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "j", false, "Output JSON instead of text (includes additional fields)")

}
