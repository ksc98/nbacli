package playbyplay

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/ksc98/nbacli/engine"
	"github.com/ksc98/nbacli/keymaps"
	"github.com/ksc98/nbacli/nag"
	"github.com/ksc98/nbacli/ui/base"
	"github.com/ksc98/nbacli/ui/constants"
)

type PlayByPlay struct {
	Data  PlayByPlayData
	Model *PlayByPlayModel
}

var (
	styleBase = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a7a#ffefdd")).
			BorderForeground(lipgloss.Color("#a38")).
			Align(lipgloss.Center)
		// Padding(3).

	SearchBoxStyle = lipgloss.NewStyle()

	VERBOSE_VIEW_COLUMNS = []table.Column{
		table.NewFlexColumn("#", "#", 1),
		table.NewFlexColumn("PERIOD", "PERIOD", 1),
		table.NewFlexColumn("CLOCK", "CLOCK", 1),
		table.NewFlexColumn("SCORE", "SCORE", 1),
		table.NewFlexColumn("TYPE", "TYPE", 1).WithFiltered(true),
		table.NewFlexColumn("TEAM", "TEAM", 1).WithFiltered(true),
		table.NewFlexColumn("PLAYER", "PLAYER", 2).WithFiltered(true),
		table.NewFlexColumn("DESCRIPTION", "DESCRIPTION", 5).WithFiltered(true),
	}

	DEFAULT_VIEW_COLUMNS = []table.Column{
		table.NewFlexColumn("PERIOD", "PERIOD", 1),
		table.NewFlexColumn("CLOCK", "CLOCK", 1),
		table.NewFlexColumn("SCORE", "SCORE", 1),
		table.NewFlexColumn("TYPE", "TYPE", 1).WithFiltered(true),
		table.NewFlexColumn("TEAM", "TEAM", 1).WithFiltered(true),
		table.NewFlexColumn("PLAYER", "PLAYER", 2).WithFiltered(true),
		table.NewFlexColumn("DESCRIPTION", "DESCRIPTION", 4).WithFiltered(true),
	}

	CURRENT_VIEW_COLUMNS = DEFAULT_VIEW_COLUMNS

	DEFAULT_UPDATE_INTERVAL_SECS = 5

	inputCounter = 0

	IN_SEARCH bool

	TOGGLE_MAP = map[string]bool{}
)

type PlayByPlayItem struct {
	PlayByPlayAction
}

type PlayByPlayModel struct {
	base.BaseModel
	Table                 table.Model
	Rows                  []table.Row
	PrevModel             base.BaseModel
	Help                  help.Model
	Width, Height, Margin int
	HorizontalMargin      int
	VerticalMargin        int
	filterTextInput       textinput.Model
	gameID                string
	updateIntervalSecs    int
	actions               []nag.PlayByPlayAction
	toggle                int
	marker                int
	secondMarker          int
}

var _ = keymaps.GameKM{
	Down:     key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "highlight next row")),
	Up:       key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "highlight previous row")),
	Previous: key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("q/esc", "back to games list")),
}

var (
	dialogBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(constants.Accent).
		Padding(0, 1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
)

func (m PlayByPlayModel) View() string {
	doc := strings.Builder{}
	prompt := lipgloss.NewStyle().Width(40).Align(lipgloss.Left).Render(m.filterTextInput.View())
	ui := lipgloss.JoinVertical(lipgloss.Left, prompt)
	gameBoard := lipgloss.Place(0, 0,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#212121"}),
	)

	doc.WriteString(gameBoard)

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		doc.String(),
		m.Table.View(),
	) + "\n"

	return view
}

func newPlayByPlayModel() *PlayByPlayModel {
	return &PlayByPlayModel{}
}

var counter = 1

func (m PlayByPlayModel) refreshPlayByPlayRows() tea.Msg {
	actions := GetPlayByPlayActions(m.gameID)
	// if m.marker == len(actions) {
	// 	return "hello"
	// }
	m.actions = actions

	return m.actions
}

func (m PlayByPlayModel) initRefresh() tea.Msg {
	return RefreshActionsEvent{gameID: m.gameID, actions: m.refreshPlayByPlayRows()}
}

func (m PlayByPlayModel) Init() tea.Cmd {
	return m.initRefresh
}

func (m PlayByPlayModel) Reverse() tea.Cmd {
	return nil
}

func (p *PlayByPlayModel) SetTable(t table.Model) {
	p.Table = t
}

func (p *PlayByPlayModel) SetPreviousModel(t base.BaseModel) {
	p.PrevModel = t
}

func (p *PlayByPlayModel) SetRows(rows []table.Row) {
	p.Rows = rows
}

func (p PlayByPlayModel) calculateWidth() int {
	return p.Width - p.HorizontalMargin
}

func (p PlayByPlayModel) calculateHeight() int {
	return p.Height - p.VerticalMargin - 4
}

func (p *PlayByPlayModel) RecalculateTable() {
	p.Table = p.Table.
		WithTargetWidth(p.calculateWidth()).
		WithPageSize(min(p.calculateHeight(), constants.WindowSize.Height-10)).
		WithMinimumHeight(p.calculateHeight())
}

func (m PlayByPlayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd        tea.Cmd
		cmds       []tea.Cmd
		styleGreen = lipgloss.NewStyle().Foreground(lipgloss.Color("#0f0"))
		// styleYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#f00"))
	)

	m.RecalculateTable()
	m.Table, cmd = m.Table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// events when inside text filter box
		if m.filterTextInput.Focused() {
			switch msg.String() {
			case "enter":
				IN_SEARCH = true
				m.filterTextInput.Blur()
				return m, tea.Batch(cmds...)
			case "esc":
				if len(m.filterTextInput.Value()) == 0 {
					m.filterTextInput.Blur()
				} else {
					m.filterTextInput.Reset()
				}
			case "ctrl+c":
				return m, tea.Quit
			default:
				m.filterTextInput, _ = m.filterTextInput.Update(msg)
			}
			// m.Table = m.Table.WithFilterInput(m.filterTextInput)
			return m, tea.Batch(cmds...)
		}
		switch msg.String() {
		case "esc":
			if IN_SEARCH {
				m.filterTextInput.Reset()
				m.filterTextInput.Blur()
				IN_SEARCH = false
			}
		case "q":
			return m.PrevModel, tea.Batch(cmds...)
		case "ctrl+c":
			return m, tea.Quit
		case "A":
			if _, ok := TOGGLE_MAP["A"]; ok {
				m.Table = m.Table.WithColumns(DEFAULT_VIEW_COLUMNS)
				delete(TOGGLE_MAP, "A")
				CURRENT_VIEW_COLUMNS = DEFAULT_VIEW_COLUMNS
			} else {
				m.Table = m.Table.WithColumns(VERBOSE_VIEW_COLUMNS)
				TOGGLE_MAP["A"] = true
				CURRENT_VIEW_COLUMNS = VERBOSE_VIEW_COLUMNS
			}
		case "r":
			// return m, m.refreshPlayByPlayRows
			return m, tea.Batch(cmds...)
		case "/":
			m.filterTextInput.Focus()
			return m, tea.Batch(cmds...)
		case "enter":
		default:
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.RecalculateTable()

	case RefreshActionsEvent:
		if msg.gameID != m.gameID {
			return m, nil
		}
		actions := msg.actions.([]nag.PlayByPlayAction)
		rows := generateRowsFromActions(actions)

		for i := m.marker; i < len(rows); i++ {
			rows[i] = rows[i].WithStyle(styleGreen)
		}

		m.secondMarker = m.marker
		slices.Reverse(rows)
		if len(rows) != m.marker {
			// update table if there are new rows
			m.Table = m.Table.WithRows(rows)
		}
		m.marker = len(rows)
		delay := time.Duration(m.updateIntervalSecs) * time.Second
		cmds = append(cmds, func() tea.Msg {
			time.Sleep(delay)
			return RefreshActionsEvent{
				gameID:  m.gameID,
				actions: m.refreshPlayByPlayRows(),
			}
		})
		return m, tea.Batch(cmds...)
	}
	return m, tea.Batch(cmds...)
}

type RefreshActionsEvent struct {
	gameID  string
	actions tea.Msg
}

func GetPlayByPlayActions(id string) []nag.PlayByPlayAction {
	e := engine.GetEngine(id)
	pbp := e.GetPlayByPlay()
	return pbp.Game.Actions
}

func generateRowsFromActions(actions []nag.PlayByPlayAction) []table.Row {
	rows := []table.Row{}
	last_score := ""
	for _, action := range actions {
		score := fmt.Sprintf("%s - %s", action.ScoreHome, action.ScoreAway)
		if score == last_score {
			score = ""
		} else {
			last_score = score
		}
		row := table.NewRow(table.RowData{
			"#":           action.ActionNumber,
			"PERIOD":      fmt.Sprintf("%d", action.Period),
			"CLOCK":       extractQuarterTime(action.Clock),
			"SCORE":       score,
			"TYPE":        action.SubType,
			"TEAM":        action.TeamTricode,
			"PLAYER":      action.PlayerNameI,
			"DESCRIPTION": action.Desc,
		})
		rows = append(rows, row)
	}
	return rows
}

// func GetPlayByPlayRows(id string) []table.Row {
// 	// populate rows
// 	actions := GetPlayByPlayActions(id)
// 	for i := 0; i < 10; i++ {
// 		actions = append(actions, nag.PlayByPlayAction{Period: 123})
// 	}
// 	rows := generateRowsFromActions(actions)
// 	return rows
// }

// populate data
func GetPlayByPlayModel(id string) *PlayByPlayModel {

	actions := GetPlayByPlayActions(id)
	// populate rows
	rows := generateRowsFromActions(actions)
	slices.Reverse(rows)

	t := createTableFromRows(rows)

	ti := textinput.New()
	ti.Placeholder = `press "/" to start searching`

	model := &PlayByPlayModel{
		Table:              t,
		Rows:               rows,
		Height:             constants.WindowSize.Height,
		Width:              constants.WindowSize.Width,
		HorizontalMargin:   0,
		VerticalMargin:     0,
		Margin:             3,
		filterTextInput:    ti,
		gameID:             id,
		updateIntervalSecs: DEFAULT_UPDATE_INTERVAL_SECS,
		actions:            actions,
		marker:             len(rows),
		secondMarker:       -1,
	}

	model.Table = model.Table.WithPageSize(constants.WindowSize.Height - 10)
	return model
}

func createTableFromRows(rows []table.Row) table.Model {
	return table.New(DEFAULT_VIEW_COLUMNS).WithRows(rows).
		BorderRounded().
		WithBaseStyle(styleBase).
		Filtered(true).
		Focused(true)
}

func extractQuarterTime(s string) string {
	regex := regexp.MustCompile(`PT(\d+)M(\d+)\.\d+S`)
	match := regex.FindStringSubmatch(s)

	if match != nil {
		minutes, _ := strconv.Atoi(match[1])
		seconds, _ := strconv.Atoi(match[2])
		secs := int(seconds)
		if secs < 10 {
			return fmt.Sprintf("%dm  %ds", minutes, secs)
		}
		return fmt.Sprintf("%dm %ds", minutes, secs)

	}
	return ""
}
