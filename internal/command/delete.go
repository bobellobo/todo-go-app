package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		if err := TaskClient.DeleteTask(id); err != nil {
			return err
		}

		fmt.Printf("✓ Task #%s deleted successfully\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
