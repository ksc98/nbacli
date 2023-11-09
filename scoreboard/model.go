package scoreboard

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksc98/nbacli/nba"
	"github.com/ksc98/nbacli/ui/constants"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))


type mode int

type SelectMsg struct {
	ActiveScorebardID uint
}

const (
	nav mode = iota
	edit
)

type Model struct {
	Mode        mode
	List        list.Model
	CurrentDate time.Time
	Quitting    bool
	Gameview    bool
}


func InitScoreboard(date time.Time) tea.Model {
	items := newScoreboardList(nba.Sb, date)
	m := Model{Mode: nav, CurrentDate: date, List: list.NewModel(items, list.NewDefaultDelegate(), 8, 8)}
	if constants.WindowSize.Height != 0 {
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.List.SetSize(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	}
	m.List.Title = "NBA Games - " + m.CurrentDate.Format("Monday, 2 Jan 06")
	m.List.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Tomorrow,
			constants.Keymap.Yesterday,
			constants.Keymap.Back,
			constants.Keymap.PlayByPlay,
		}
	}
	return m
}

func newScoreboardList(scbrd *nba.ScoreboardRepository, date time.Time) []list.Item {
	games := scbrd.GetGames(date)
	return gamesToItems(games)
}

func gamesToItems(games []nba.BoxScoreSummary) []list.Item {
	items := make([]list.Item, len(games))
	for i, proj := range games {
		items[i] = list.Item(proj)
	}
	return items
}



// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.List.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Yesterday):
			var previousDay nba.ScoreboardRepository
			m.CurrentDate = m.CurrentDate.AddDate(0, 0, -1)
			games := previousDay.GetGames(m.CurrentDate)
			items := gamesToItems(games)
			m.List.Title = "NBA Games - " + m.CurrentDate.Format("Monday, 2 Jan 06")
			m.List.SetItems(items)
		case key.Matches(msg, constants.Keymap.Tomorrow):
			var nextDay nba.ScoreboardRepository
			m.CurrentDate = m.CurrentDate.AddDate(0, 0, 1)
			games := nextDay.GetGames(m.CurrentDate)
			items := gamesToItems(games)
			m.List.Title = "NBA Games - " + m.CurrentDate.Format("Monday, 2 Jan 06")
			m.List.SetItems(items)
		case key.Matches(msg, constants.Keymap.Quit):
			m.Quitting = true
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.PlayByPlay):
			activeGame := m.List.SelectedItem().(nba.BoxScoreSummary)
			pbpView := InitPlayByPlayView(activeGame.GameId, activeGame, m)
			return pbpView.Update(constants.WindowSize)
		case key.Matches(msg, constants.Keymap.Enter):
			m.Gameview = true
			activeGame := m.List.SelectedItem().(nba.BoxScoreSummary)
			gameView := InitGameView(activeGame.GameId, activeGame, m)
			return gameView.Update(constants.WindowSize)
		default:
			m.List, cmd = m.List.Update(msg)
		}
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}



// View return the text UI to be output to the terminal
func (m Model) View() string {
	if m.Quitting {
		return ""
	}
	return constants.DocStyle.Render(m.List.View() + "\n")
}


