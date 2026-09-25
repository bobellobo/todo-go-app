package command

import (
	"go-app/internal/client"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiURL     string
	apiToken   string
	TaskClient *client.Client
)

const defaultAPIURL = "http://localhost:8081"

var rootCmd = &cobra.Command{
	Use:   "taskctl",
	Short: "taskctl is a command line tool to manage your tasks API",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// init client before running any command
		TaskClient = client.NewClient(apiURL, apiToken)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// flag available to all subcommands
	if configuredURL := os.Getenv("TASKCTL_API_URL"); configuredURL != "" {
		apiURL = configuredURL
	} else {
		apiURL = defaultAPIURL
	}
	apiToken = os.Getenv("TASKCTL_API_TOKEN")
	rootCmd.PersistentFlags().StringVarP(&apiURL, "url", "u", apiURL, "Base URL of the REST API")
}
