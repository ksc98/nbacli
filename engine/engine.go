package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/ksc98/nbacli/nag"
	"github.com/ksc98/nbacli/nba"
)

var games sync.Map

const (
	DEFAULT_TICKER_INTERVAL = 5
)

type GameEngine struct {
	playbyplay   *nag.PlayByPlayResponse
	boxscore     *nba.BoxScoreSummary
	gameid       string
	tickInterval int
}

type NBAEngine struct {
	games sync.Map
}

func NewGameEngine(id string) *GameEngine {
	return &GameEngine{
		gameid:       id,
		tickInterval: DEFAULT_TICKER_INTERVAL,
		playbyplay:   &nag.PlayByPlayResponse{},
		boxscore:     &nba.BoxScoreSummary{},
	}
}

func (e *GameEngine) InitPlayByPlayTracker() (nag.PlayByPlayResponse, error) {
	pbp := nag.NewPlayByPlayV2(e.gameid)
	pbp.Get()
	if pbp.PlayByPlayResponse == nil {
		return nag.PlayByPlayResponse{}, fmt.Errorf("no play by play data found!")
	}
	e.playbyplay = pbp.PlayByPlayResponse
	go e.startPlayByPlayTicker()
	persist(e)
	return *pbp.PlayByPlayResponse, nil
}

func (e *GameEngine) startPlayByPlayTicker() {
	ticker := time.NewTicker(time.Duration(e.tickInterval) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		pbp := nag.NewPlayByPlayV2(e.gameid)
		pbp.Get()
		// game := pbp.PlayByPlayResponse.Game
		// for i := 0; i < e.counter; i++ {
		// 	arr := game.Actions
		// 	arr = append(arr, nag.PlayByPlayAction{Period: e.counter})
		// 	game.SetActions(arr)
		// }
		// e.counter++
		e.playbyplay = pbp.PlayByPlayResponse
		persist(e)
	}
}

func (e *GameEngine) GetPlayByPlay() *nag.PlayByPlayResponse {
	return e.playbyplay
}

func GetEngine(id string) *GameEngine {
	ge, ok := games.Load(id)
	if ok {
		return ge.(*GameEngine)
	}

	e := NewGameEngine(id)
	if _, err := e.InitPlayByPlayTracker(); err != nil {
		return nil
	}

	games.Store(id, e)
	return e
}

func persist(ge *GameEngine) {
	games.Store(ge.gameid, ge)
}
