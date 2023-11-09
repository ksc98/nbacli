package playbyplay

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/charmbracelet/bubbles/table"
	log "github.com/sirupsen/logrus"
)

// constructor
func New(id string) *PlayByPlay {
	data := PlayByPlayData{}
	data.Game.GameID = id
	return &PlayByPlay{Data: data}
}

// return the list of actions for the game
func (p *PlayByPlay) Actions() []PlayByPlayAction {
	return p.Data.Game.Actions
}

func (p *PlayByPlay) InitPlayByPlayView() {

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
