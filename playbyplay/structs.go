package playbyplay

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
)

type PlayByPlay struct {
	Data  PlayByPlayData
	Rows  []table.Row
	Model PlayByPlayModel
}

type PlayByPlayAction struct {
	ActionNumber            int       `json:"actionNumber"`
	ActionType              string    `json:"actionType"`
	Clock                   string    `json:"clock"`
	Description             string    `json:"description"`
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

type PlayByPlayGame struct {
	GameID  string             `json:"gameId"`
	Actions []PlayByPlayAction `json:"actions"`
}

type PlayByPlayData struct {
	Meta any            `json:"-"`
	Game PlayByPlayGame `json:"game"`
}
