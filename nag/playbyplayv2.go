package nag

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const cdnBaseUrl = "https://cdn.nba.com/static/json/liveData/playbyplay/playbyplay"

// PlayByPlayV2 wraps request to and response from playbyplayv2 endpoint.
type PlayByPlayV2 struct {
	*Client
	GameID             string
	PlayByPlayResponse *PlayByPlayResponse
}

type PlayByPlayResponse struct {
	Meta any             `json:"-"`
	Game *PlayByPlayGame `json:"game"`
}

type PlayByPlayGame struct {
	GameID  string             `json:"gameId"`
	Actions []PlayByPlayAction `json:"actions"`
}

func (p *PlayByPlayGame) SetActions(actions []PlayByPlayAction) {
	p.Actions = actions
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

// NewPlayByPlayV2 creates a default PlayByPlayV2 instance.
func NewPlayByPlayV2(id string) *PlayByPlayV2 {
	return &PlayByPlayV2{
		Client: NewDefaultClient(),
		GameID: id,
	}
}

// Get sends a GET request to playbyplayv2 endpoint.
func (c *PlayByPlayV2) Get() error {
	reqUrl := fmt.Sprintf("%s_%s.json", cdnBaseUrl, c.GameID)
	resp, err := http.Get(reqUrl)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var res PlayByPlayResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}

	c.PlayByPlayResponse = &res
	return nil
}
