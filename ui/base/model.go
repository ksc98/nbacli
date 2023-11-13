package base

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BaseModel interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

var (
	StyleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))
)
