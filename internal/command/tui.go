package command

import (
	"go-app/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive terminal user interface",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(ui.NewModel(TaskClient))
		_, err := p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
