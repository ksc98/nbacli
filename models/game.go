package models

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/ksc98/nbacli/keymaps"
	"github.com/ksc98/nbacli/nba"
	"github.com/ksc98/nbacli/ui/base"
	"github.com/ksc98/nbacli/ui/gameboard/scoretext"
)

type GameModel struct {
	Table                 table.Model
	ActiveGameID          string
	ActiveGame            nba.BoxScoreSummary
	PreviousModel         base.BaseModel
	Help                  help.Model
	Width, Height, Margin int
}

func (m *GameModel) recalculateTable() {
	m.Table = m.Table.WithTargetWidth(m.Width)
}

func (m GameModel) Init() tea.Cmd { return nil }

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			// return to previous page
			return m.PreviousModel, tea.Batch()
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			// TODO: to player view
			return m, tea.Batch()
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.recalculateTable()
	}

	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m GameModel) View() string {
	table := m.Table.View() + "\n"

	helpContainer := lipgloss.NewStyle().
		SetString(m.Help.View(gameKM)).
		Width(m.Width).
		Align(lipgloss.Center).
		PaddingTop(1).
		String()

	return scoretext.RenderScoreText(
		m.ActiveGame.ArenaName,
		m.ActiveGame.GameDate,
		m.ActiveGame.HomeTeamScore,
		m.ActiveGame.VisitorTeamScore,
		m.ActiveGame.HomeTeamName,
		m.ActiveGame.VisitorTeamName) + table + helpContainer
}

var gameKM = keymaps.GameKM{
	Down:     key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "highlight next row")),
	Up:       key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "highlight previous row")),
	Previous: key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("q/esc", "back to games list")),
}
