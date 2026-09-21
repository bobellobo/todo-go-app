package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [group]",
	Short: "List all tasks or tasks in a specific group",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		group := ""
		if len(args) == 1 {
			group = args[0]
		}

		tasks, err := TaskClient.ListTasks(group)
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return nil
		}

		fmt.Printf("ID\tStatus\tGroup\t\tTitle\n")
		fmt.Println("--------------------------------------------------")
		for _, t := range tasks {
			status := "[ ]"
			if t.Done {
				status = "[✓]"
			}
			grp := t.Group
			if grp == "" {
				grp = "-"
			}
			fmt.Printf("%d\t%s\t%-12s\t%s\n", t.ID, status, grp, t.Title)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
