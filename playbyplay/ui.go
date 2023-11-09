package playbyplay

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksc98/nbacli/ui"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var (
	COLUMNS = []table.Column{
		{Title: "Period", Width: 8},
		{Title: "Clock", Width: 11},
		{Title: "Score", Width: 7},
		{Title: "Type", Width: 15},
		{Title: "Team", Width: 9},
		{Title: "Player", Width: 18},
		{Title: "Description", Width: 70},
	}
)

type PlayByPlayModel struct {
	ui.GameModel
	Table table.Model
}

func (p *PlayByPlay) UI() {
	if _, err := tea.NewProgram(p.Model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (p *PlayByPlay) GetModel() table.Model {
	return p.Model.Table
}

func (m PlayByPlayModel) Init() tea.Cmd { return nil }

func (m PlayByPlayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "g":
			m.Table.GotoTop()
		case "G":
			m.Table.GotoBottom()
		case "enter":
		}
		// case tea.WindowSizeMsg:
		// 	m.width = msg.Width
		// 	m.recalculateTable()
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m PlayByPlayModel) View() string {
	return baseStyle.Render(m.Table.View()) + "\n"
}

func (p *PlayByPlay) LoadModel() PlayByPlayModel {
	t := table.New(
		table.WithColumns(COLUMNS),
		table.WithRows(p.Rows),
		table.WithFocused(true),
		table.WithHeight(35),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := PlayByPlayModel{t}

	p.Model = m
	return m
}

func extractQuarterTime(s string) string {
	regex := regexp.MustCompile(`PT(\d+)M(\d+)\.\d+S`)
	match := regex.FindStringSubmatch(s)

	if match != nil {
		minutes, _ := strconv.Atoi(match[1])
		seconds, _ := strconv.Atoi(match[2])
		// round seconds
		roundedSeconds := float64(float64(seconds)/100.0*10.0) / 10.0
		return fmt.Sprintf("%ds  %.1fs\n", minutes, roundedSeconds)
	}
	return ""
}
