package ui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	var s strings.Builder

	// Header & Active Group Badge
	filterInfo := " [All Groups]"
	if m.selectedGroup != "" {
		filterInfo = fmt.Sprintf(" [Filter: @%s]", m.selectedGroup)
	}
	_, _ = s.WriteString(titleStyle.Render(" TASK MANAGER TUI ") + activeFilterStyle.Render(filterInfo) + "\n\n")

	if m.err != nil {
		s.WriteString(fmt.Sprintf("Error: %v\n\n", m.err))
	}

	// Group Submenu View
	if m.mode == modeGroupSelect {
		_, _ = s.WriteString(m.groupList.View() + "\n")
		s.WriteString(helpStyle.Render("↑/↓: navigate • enter: filter group • esc: back"))
		return s.String()
	}

	// Standard Task List View
	if len(m.tasks) == 0 {
		s.WriteString(" No tasks available in this view. Press 'a' to add one!\n\n")
	} else {
		for i, t := range m.tasks {
			cursor := "  "
			if m.cursor == i {
				cursor = "❯ "
			}

			status := statusTodoStyle.Render("[ ]")
			if t.Done {
				status = statusDoneStyle.Render("[✓]")
			}

			group := ""
			if t.Group != "" {
				group = groupBadgeStyle.Render(fmt.Sprintf("@%s", t.Group)) + " "
			}

			line := fmt.Sprintf("%s%s #%d %s%s", cursor, status, t.ID, group, t.Title)
			if m.cursor == i {
				line = selectedStyle.Render(line)
			}

			_, _ = s.WriteString(line + "\n")
		}
		s.WriteString("\n")
	}

	// Footer Help Text
	if m.mode == modeInput {
		_, _ = s.WriteString("New Task: " + m.textInput.View() + "\n")
		s.WriteString(helpStyle.Render("press enter to save • esc to cancel"))
	} else {
		s.WriteString(helpStyle.Render("j/k/↑/↓: navigate • space: toggle • a: add • g: groups menu • d: delete • q: quit"))
	}

	return s.String()
}
