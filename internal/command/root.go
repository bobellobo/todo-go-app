package command

import (
	"go-app/internal/client"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiURL     string
	TaskClient *client.Client
)

var rootCmd = &cobra.Command{
	Use:   "taskctl",
	Short: "taskctl is a command line tool to manage your tasks API",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// init client before running any command
		TaskClient = client.NewClient(apiURL)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// flag available to all subcommands
	rootCmd.PersistentFlags().StringVarP(&apiURL, "url", "u", "http://localhost:8080", "Base URL of the REST API")
}
