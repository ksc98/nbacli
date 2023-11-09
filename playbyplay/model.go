package playbyplay

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksc98/nbacli/nba"
	"github.com/ksc98/nbacli/scoreboard"
	"github.com/ksc98/nbacli/structs"
	log "github.com/sirupsen/logrus"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type PlayByPlay struct {
	Data  structs.PlayByPlayData
	Rows  []table.Row
	Model PlayByPlayModel
}

type PlayByPlayModel struct {
	Table         table.Model
	PreviousModel scoreboard.Model
}

func (p *PlayByPlay) GetModel() PlayByPlayModel {
	return p.Model
}

func InitPlayByPlayView(activeGameID string, activeGame nba.BoxScoreSummary, previousModel scoreboard.Model) *PlayByPlayModel {
	pbp := New(activeGameID)
	err := pbp.Get()
	if err != nil {
		log.Fatal(err)
	}
	t := pbp.GetModel()

	// m := GameModel{t, activeGameID, activeGame, previousModel, help.New(), constants.WindowSize.Height, constants.WindowSize.Width, 2}
	return &t
}
func (m PlayByPlayModel) View() string {
	return baseStyle.Render(m.Table.View()) + "\n"
}

func (m PlayByPlayModel) Init() tea.Cmd { return nil }

func (m PlayByPlayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m.PreviousModel, tea.Batch()
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

func (p *PlayByPlay) LoadModel() PlayByPlayModel {
	t := table.New(
		table.WithColumns(structs.COLUMNS),
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

	m := PlayByPlayModel{Table: t}

	p.Model = m
	return m
}

func (p *PlayByPlay) UI() {
	if _, err := tea.NewProgram(p.Model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}


// constructor
func New(id string) *PlayByPlay {
	data := structs.PlayByPlayData{}
	data.Game.GameID = id
	return &PlayByPlay{Data: data}
}

// return the list of actions for the game
func (p *PlayByPlay) Actions() []structs.PlayByPlayAction {
	return p.Data.Game.Actions
}

// populate data
func (p *PlayByPlay) Get() error {
	reqUrl := fmt.Sprintf(
		"https://cdn.nba.com/static/json/liveData/playbyplay/playbyplay_%s.json",
		p.Data.Game.GameID)

	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Infof("game with ID %s seems to be invalid, double check it?", p.Data.Game.GameID)
		return fmt.Errorf("error looking up game %s!", p.Data.Game.GameID)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Infof("game with ID %s seems to be invalid, double check it?", p.Data.Game.GameID)
		return fmt.Errorf("error looking up game %s!", p.Data.Game.GameID)
	}

	err = json.Unmarshal(body, &p.Data)
	if err != nil {
		log.Infof("game with ID %s seems to be invalid, double check it?", p.Data.Game.GameID)
		return fmt.Errorf("error looking up game %s!", p.Data.Game.GameID)
	}

	// populate rows
	actions := p.Actions()
	last_score := ""
	for _, action := range actions {
		score := fmt.Sprintf("%s - %s", action.ScoreHome, action.ScoreAway)
		if score == last_score {
			score = ""
		} else {
			last_score = score
		}
		row := table.Row{
			fmt.Sprintf("%d", action.Period),
			extractQuarterTime(action.Clock),
			score,
			action.SubType,
			action.TeamTricode,
			action.PlayerNameI,
			action.Description,
		}
		p.Rows = append(p.Rows, row)
	}

	p.LoadModel()

	return nil
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