package nag

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dylantientcheu/nbacli/nag/params"
)

type GameLog struct {
	SeasonID            int    `json:"SEASON_ID"`
	PlayerID            int    `json:"Player_ID"`
	GameID              int    `json:"Game_ID"`
	GameDate            string `json:"GAME_DATE"`
	Matchup             string `json:"MATCHUP"`
	WL                  string `json:"WL"`
	Min                 int    `json:"MIN"`
	Points              int    `json:"PTS"`
	Assists             int    `json:"AST"`
	Rebounds            int    `json:"REB"`
	Steals              int    `json:"STL"`
	Blocks              int    `json:"BLK"`
	FieldGoalsMade      int    `json:"FGM"`
	FieldGoalsAttempted int    `json:"FGA"`
}

type PlayerGameLogResponse struct {
	GameLog struct {
		Headers []string        `json:"[]"`
		RowSet  [][]interface{} `json:"[]"`
	} `json:"playergamelog"`
}

// CommonPlayerInfo wraps request to and response from commonplayerinfo endpoint.
type PlayerGameLog struct {
	*Client
	Response   *Response
	PlayerID   string
	LeagueID   string
	Season     string
	SeasonType params.SeasonType
	DateFrom   string
	DateTo     string
}

// NewCommonPlayerInfo creates a default CommonPlayerInfo instance.
func NewPlayerGameLog(id string) *PlayerGameLog {
	return &PlayerGameLog{
		Client:     NewDefaultClient(),
		PlayerID:   id,
		LeagueID:   params.LeagueID.Default(),
		Season:     params.CurrentSeason,
		SeasonType: params.DefaultSeasonType,
	}
}

// Get sends a GET request to commonplayerinfo endpoint.
func (c *PlayerGameLog) Get() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/playergamelog", c.BaseURL.String()), nil)
	if err != nil {
		return err
	}

	req.Header = DefaultStatsHeader

	q := req.URL.Query()
	q.Add("PlayerID", c.PlayerID)
	// q.Add("LeagueID", c.LeagueID)
	q.Add("Season", "2023-24")
	// q.Add("DateFrom", "")
	// q.Add("DateTo", "")
	q.Add("SeasonType", "Regular Season")

	req.URL.RawQuery = q.Encode()
	println(req.URL.String())

	b, err := c.Do(req)
	if err != nil {
		return err
	}

	var res Response
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	c.Response = &res
	return nil
}
