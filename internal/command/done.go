package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [id]",
	Short: "Toggle a task status (done/undone)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		task, err := TaskClient.ToggleTask(id)
		if err != nil {
			return err
		}

		status := "pending"
		if task.Done {
			status = "completed"
		}

		fmt.Printf("✓ Task #%d status updated to %s\n", task.ID, status)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
