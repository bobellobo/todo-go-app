package ui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tasksMsg:
		m.tasks = msg
		if m.cursor >= len(m.tasks) && len(m.tasks) > 0 {
			m.cursor = len(m.tasks) - 1
		}

	case groupsMsg:
		items := []list.Item{
			groupItem{name: "", count: 0}, // Option for "All Tasks"
		}
		for name, count := range msg {
			items = append(items, groupItem{name: name, count: count})
		}
		m.groupList.SetItems(items)

	case errMsg:
		m.err = msg

	case tea.KeyMsg:
		// MODE 1: Text Input Mode
		if m.mode == modeInput {
			switch msg.String() {
			case "enter":
				val := strings.TrimSpace(m.textInput.Value())
				if val != "" {
					title, group := parseTitleAndGroup(val)
					// If created inside a group filter view, default to current selected group if unspecified
					if group == "" && m.selectedGroup != "" {
						group = m.selectedGroup
					}
					m.client.AddTask(title, group)
					m.textInput.Reset()
					m.mode = modeNormal
					return m, m.fetchTasksCmd
				}
			case "esc":
				m.mode = modeNormal
				m.textInput.Reset()
				return m, nil
			}
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		// MODE 2: Group Submenu Select Mode
		if m.mode == modeGroupSelect {
			switch msg.String() {
			case "enter":
				selected, ok := m.groupList.SelectedItem().(groupItem)
				if ok {
					m.selectedGroup = selected.name
					m.cursor = 0
				}
				m.mode = modeNormal
				return m, m.fetchTasksCmd

			case "esc", "q":
				m.mode = modeNormal
				return m, nil
			}

			m.groupList, cmd = m.groupList.Update(msg)
			return m, cmd
		}

		// MODE 3: Normal Mode Keybindings
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}

		case "space", "enter":
			if len(m.tasks) > 0 {
				targetTask := m.tasks[m.cursor]
				m.client.ToggleTask(strconv.Itoa(targetTask.ID))
				return m, m.fetchTasksCmd
			}

		case "d", "backspace":
			if len(m.tasks) > 0 {
				targetTask := m.tasks[m.cursor]
				m.client.DeleteTask(strconv.Itoa(targetTask.ID))
				return m, m.fetchTasksCmd
			}

		case "a", "n":
			m.mode = modeInput
			m.textInput.Focus()
			return m, textinput.Blink

		case "g":
			m.mode = modeGroupSelect
			return m, m.fetchGroupsCmd

		case "r":
			return m, m.fetchTasksCmd
		}
	}

	return m, nil
}

func parseTitleAndGroup(input string) (string, string) {
	parts := strings.Split(input, " -g ")
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return input, ""
}
