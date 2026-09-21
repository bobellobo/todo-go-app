package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var groupFlag string

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]
		task, err := TaskClient.AddTask(title, groupFlag)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Created task #%d: %s\n", task.ID, task.Title)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&groupFlag, "group", "g", "", "Assign task to a group")
	rootCmd.AddCommand(addCmd)
}
