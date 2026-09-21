package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage task groups",
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active task groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		groups, err := TaskClient.ListGroups()
		if err != nil {
			return err
		}

		if len(groups) == 0 {
			fmt.Println("No active groups found.")
			return nil
		}

		fmt.Println("Group\t\tTask Count")
		fmt.Println("--------------------------")
		for name, count := range groups {
			fmt.Printf("%-12s\t%d\n", name, count)
		}
		return nil
	},
}

func init() {
	groupCmd.AddCommand(groupListCmd)
	rootCmd.AddCommand(groupCmd)
}
