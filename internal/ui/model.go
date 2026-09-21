package ui

import (
	"fmt"

	"go-app/internal/client"
	"go-app/internal/task"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	modeNormal mode = iota
	modeInput
	modeGroupSelect
)

// Styling definitions
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	statusDoneStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
	statusTodoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87"))
	selectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#75C6FF")).Bold(true)
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	groupBadgeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAF00")).Italic(true)
	activeFilterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
)

// Group item wrapper for bubbles/list
type groupItem struct {
	name  string
	count int
}

func (i groupItem) Title() string {
	if i.name == "" {
		return "All Tasks"
	}
	return i.name
}
func (i groupItem) Description() string {
	if i.name == "" {
		return "Show all tasks across all groups"
	}
	return fmt.Sprintf("%d task(s)", i.count)
}
func (i groupItem) FilterValue() string { return i.name }

type Model struct {
	client        *client.Client
	tasks         []task.Task
	cursor        int
	err           error
	mode          mode
	selectedGroup string
	textInput     textinput.Model
	groupList     list.Model
}

func NewModel(c *client.Client) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter new task title (e.g. 'Deploy app -g devops')..."
	ti.CharLimit = 156
	ti.Width = 50

	// Initialize group selection sub-list
	gl := list.New([]list.Item{}, list.NewDefaultDelegate(), 30, 14)
	gl.Title = "Select Group Filter"
	gl.SetShowHelp(false)

	return Model{
		client:        c,
		tasks:         []task.Task{},
		cursor:        0,
		textInput:     ti,
		groupList:     gl,
		mode:          modeNormal,
		selectedGroup: "",
	}
}

// Custom Messages
type tasksMsg []task.Task
type groupsMsg map[string]int
type errMsg error

func (m Model) fetchTasksCmd() tea.Msg {
	tasks, err := m.client.ListTasks(m.selectedGroup)
	if err != nil {
		return errMsg(err)
	}
	return tasksMsg(tasks)
}

func (m Model) fetchGroupsCmd() tea.Msg {
	groups, err := m.client.ListGroups()
	if err != nil {
		return errMsg(err)
	}
	return groupsMsg(groups)
}

func (m Model) Init() tea.Cmd {
	return m.fetchTasksCmd
}
