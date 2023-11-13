package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PlayByPlay struct {
	Data PlayByPlayData
}

// constructor
func New(id string) *PlayByPlay {
	data := PlayByPlayData{}
	data.Game.GameID = id
	return &PlayByPlay{
		Data: data,
	}
}

// return the list of actions for the game
func (p *PlayByPlay) Actions() []PlayByPlayAction {
	return p.Data.Game.Actions
}

func (p *PlayByPlay) Refresh(ms int) {
	ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C
		p.Get()
	}
}

// populate data
func (p *PlayByPlay) Get() error {
	reqUrl := fmt.Sprintf(
		"https://cdn.nba.com/static/json/liveData/playbyplay/playbyplay_%s.json",
		p.Data.Game.GameID)

	resp, err := http.Get(reqUrl)
	if err != nil {
		return fmt.Errorf("error looking up game %s!", p.Data.Game.GameID)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error looking up game %s!", p.Data.Game.GameID)
	}

	err = json.Unmarshal(body, &p.Data)
	if err != nil {
		return fmt.Errorf("error looking up game %s!", p.Data.Game.GameID)
	}

	// populate rows
	// actions := p.Actions()

	return nil
}

type PlayByPlayData struct {
	Meta any            `json:"-"`
	Game PlayByPlayGame `json:"game"`
}

type PlayByPlayGame struct {
	GameID  string             `json:"gameId"`
	Actions []PlayByPlayAction `json:"actions"`
}

type PlayByPlayAction struct {
	ActionNumber            int       `json:"actionNumber"`
	ActionType              string    `json:"actionType"`
	Clock                   string    `json:"clock"`
	Desc                    string    `json:"description"`
	Edited                  time.Time `json:"edited"`
	IsFieldGoal             int       `json:"isFieldGoal"`
	IsTargetScoreLastPeriod bool      `json:"isTargetScoreLastPeriod"`
	OrderNumber             int       `json:"orderNumber"`
	Period                  int       `json:"period"`
	PeriodType              string    `json:"periodType"`
	PersonId                int       `json:"personId"`
	PersonIdsFilter         []int     `json:"personIdsFilter"`
	PlayerName              string    `json:"playerName"`
	PlayerNameI             string    `json:"playerNameI"`
	Possession              int       `json:"possession"`
	Qualifiers              []string  `json:"qualifiers"`
	ReboundDefensiveTotal   int       `json:"reboundDefensiveTotal"`
	ReboundOffensiveTotal   int       `json:"reboundOffensiveTotal"`
	ReboundTotal            int       `json:"reboundTotal"`
	ScoreAway               string    `json:"scoreAway"`
	ScoreHome               string    `json:"scoreHome"`
	ShotActionNumber        int       `json:"shotActionNumber"`
	Side                    string    `json:"side"`
	SubType                 string    `json:"subType"`
	TeamId                  int       `json:"teamId"`
	TeamTricode             string    `json:"teamTricode"`
	TimeActual              string    `json:"timeActual"`
}
