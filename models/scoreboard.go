package models

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/ksc98/nbacli/helpers"
	"github.com/ksc98/nbacli/nba"
	play "github.com/ksc98/nbacli/playbyplay"
	"github.com/ksc98/nbacli/ui/constants"
)

var (
	TOGGLE_MAP = map[string]bool{}
)

type mode int

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

		case key.Matches(msg, constants.Keymap.Follow):
			activeGame := m.List.SelectedItem().(nba.BoxScoreSummary)

			if !helpers.Toggle("follow" + activeGame.GameId) {
				activeGame.FollowGame()
			} else {
				activeGame.UnfollowGame()
			}
			m.List.SetItem(m.List.Cursor(), activeGame)
			return m, nil

		case key.Matches(msg, constants.Keymap.PlayByPlay):
			activeGame := m.List.SelectedItem().(nba.BoxScoreSummary)
			pbpView, err := InitPlayByPlayView(activeGame.GameId, activeGame, m)
			pbpView.RecalculateTable()
			// cmds = append(cmds, pbpView.Init())
			// m, _ := pbpView.Update(constants.WindowSize)
			if err == nil {
				return pbpView, pbpView.Init()
				return pbpView.Update(constants.WindowSize)
			}
		case key.Matches(msg, constants.Keymap.Enter):
			m.Gameview = true
			activeGame := m.List.SelectedItem().(nba.BoxScoreSummary)
			gameView := InitGameView(activeGame.GameId, activeGame, m)
			return gameView.Update(constants.WindowSize)

		case msg.String() == "esc":

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
	// return constants.DocStyle.Render(m.List.View() + "\n" + base.StyleSubtle.Render("help"))
	return constants.DocStyle.Render(m.List.View() + "\n")
}

func InitPlayByPlayView(activeGameID string, activeGame nba.BoxScoreSummary, previousModel Model) (*play.PlayByPlayModel, error) {
	model := play.GetPlayByPlayModel(activeGameID)
	model.SetPreviousModel(previousModel)
	return model, nil
}

func InitGameView(activeGameID string, activeGame nba.BoxScoreSummary, previousModel Model) *GameModel {
	columns := []table.Column{
		table.NewFlexColumn("POS", "POS", 2),
		table.NewFlexColumn("NAME", "NAME", 10),
		table.NewFlexColumn("MIN", "MIN", 6),
		table.NewFlexColumn("FG", "FG", 6),
		table.NewFlexColumn("3PT", "3PT", 3),
		table.NewFlexColumn("FT", "FT", 3),
		table.NewFlexColumn("REB", "REB", 3),
		table.NewFlexColumn("AST", "AST", 3),
		table.NewFlexColumn("STL", "STL", 3),
		table.NewFlexColumn("BLK", "BLK", 3),
		table.NewFlexColumn("TO", "TO", 3),
		table.NewFlexColumn("+/-", "+/-", 4),
		table.NewFlexColumn("PTS", "PTS", 3),
	}

	rows := newStatsBoard(nba.Gm, activeGameID)

	t := table.New(columns).WithRows(rows).
		Focused(true).
		Border(constants.CustomTableBorder).WithBaseStyle(baseStyle).WithPageSize(constants.WindowSize.Height / 3)

	m := GameModel{t, activeGameID, activeGame, previousModel, help.New(), constants.WindowSize.Height, constants.WindowSize.Width, 3}
	return &m
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func newStatsBoard(game *nba.BoxScoreRepository, gameID string) []table.Row {
	gameStats := game.GetSingleGameStats(gameID)
	return statsToRows(gameStats)
}

func statsToRows(gameStats []nba.GameStat) []table.Row {
	var rows []table.Row
	areBenchers := false

	rows = append(rows, table.NewRow(renderTeamRow("AWAY TEAM")).
		WithStyle(lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).
			Background(constants.Secondary)))

	for idx, stat := range gameStats {
		// format plus minus
		var plusMinus string
		if stat.PlusMinus > 0 {
			plusMinus = "+" + strconv.FormatInt(stat.PlusMinus, 10)
		} else {
			plusMinus = strconv.FormatInt(stat.PlusMinus, 10)
		}

		if (stat.StartPosition == "") && !areBenchers {
			rows = append(rows, table.NewRow(
				renderBenchRow(),
			).WithStyle(lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Background(lipgloss.AdaptiveColor{Light: "214", Dark: "#181818"})))
			areBenchers = true
		}

		rows = append(rows, table.NewRow(
			table.RowData{
				"POS":  stat.StartPosition,
				"NAME": stat.PlayerName,
				"MIN":  stat.Min,
				"FG":   strconv.FormatInt(stat.Fgm, 10) + "-" + strconv.FormatInt(stat.Fga, 10),
				"3PT":  strconv.FormatInt(stat.Fg3M, 10),
				"FT":   strconv.FormatInt(stat.Ftm, 10),
				"REB":  strconv.FormatInt(stat.Reb, 10),
				"AST":  strconv.FormatInt(stat.AST, 10),
				"STL":  strconv.FormatInt(stat.Stl, 10),
				"BLK":  strconv.FormatInt(stat.Blk, 10),
				"TO":   strconv.FormatInt(stat.To, 10),
				"+/-":  plusMinus,
				"PTS":  strconv.FormatInt(stat.Pts, 10),
			},
		))
		if stat.StartPosition != "" {
			areBenchers = false
		}

		if idx < len(gameStats)-1 && gameStats[idx].TeamID != gameStats[idx+1].TeamID {
			rows = append(rows, table.NewRow(renderTeamRow("HOME TEAM")).WithStyle(lipgloss.NewStyle().
				AlignHorizontal(lipgloss.Center).
				Background(constants.Secondary)))
		}
	}
	return rows
}

func renderBenchRow() table.RowData {
	return table.RowData{
		"POS":  "",
		"NAME": table.NewStyledCell("B E N C H", lipgloss.NewStyle().Foreground(constants.Tertiary).Padding(0)),
		"MIN":  "",
		"FG":   "",
		"3PT":  "",
		"FT":   "",
		"REB":  "",
		"AST":  "",
		"STL":  "",
		"BLK":  "",
		"TO":   "",
		"+/-":  "",
		"PTS":  "",
	}
}

func renderTeamRow(team string) table.RowData {
	return table.RowData{
		"POS":  "",
		"NAME": table.NewStyledCell(team, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))),
		"MIN":  "",
		"FG":   "",
		"3PT":  "",
		"FT":   "",
		"REB":  "",
		"AST":  "",
		"STL":  "",
		"BLK":  "",
		"TO":   "",
		"+/-":  "",
		"PTS":  "",
	}
}

func gamesToItems(games []nba.BoxScoreSummary) []list.Item {
	items := make([]list.Item, len(games))
	for i, proj := range games {
		items[i] = list.Item(proj)
	}
	return items
}

type SelectMsg struct {
	ActiveScorebardID uint
}

func InitScoreboard(date time.Time) tea.Model {
	items := newScoreboardList(nba.Sb, date)
	defaultDelegate := list.NewDefaultDelegate()
	m := Model{
		Mode:        nav,
		CurrentDate: date,
		List:        list.New(items, defaultDelegate, 1, 1),
	}
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
