package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ksc98/nbacli/nag"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mitchellh/mapstructure"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type GameLogEntry struct {
	AST             int     `json:"AST"`
	BLK             int     `json:"BLK"`
	DREB            int     `json:"DREB"`
	FG3A            int     `json:"FG3A"`
	FG3M            int     `json:"FG3M"`
	FG3_PCT         float64 `json:"FG3_PCT"`
	FGA             int     `json:"FGA"`
	FGM             int     `json:"FGM"`
	FG_PCT          float64 `json:"FG_PCT"`
	FTA             int     `json:"FTA"`
	FTM             int     `json:"FTM"`
	FT_PCT          float64 `json:"FT_PCT"`
	GAME_DATE       string  `json:"GAME_DATE"`
	Game_ID         string  `json:"Game_ID"`
	MATCHUP         string  `json:"MATCHUP"`
	MIN             int     `json:"MIN"`
	OREB            int     `json:"OREB"`
	PF              int     `json:"PF"`
	PLUS_MINUS      int     `json:"PLUS_MINUS"`
	PTS             int     `json:"PTS"`
	Player_ID       int     `json:"Player_ID"`
	REB             int     `json:"REB"`
	SEASON_ID       string  `json:"SEASON_ID"`
	STL             int     `json:"STL"`
	TOV             int     `json:"TOV"`
	VIDEO_AVAILABLE int     `json:"VIDEO_AVAILABLE"`
	WL              string  `json:"WL"`
}

type Player struct {
	DISPLAY_FIRST_LAST        string `json:"DISPLAY_FIRST_LAST"`
	DISPLAY_LAST_COMMA_FIRST  string `json:"DISPLAY_LAST_COMMA_FIRST"`
	FROM_YEAR                 int    `json:"FROM_YEAR"`
	GAMES_PLAYED_FLAG         string `json:"GAMES_PLAYED_FLAG"`
	OTHERLEAGUE_EXPERIENCE_CH string `json:"OTHERLEAGUE_EXPERIENCE_CH"`
	PERSON_ID                 int    `json:"PERSON_ID"`
	PLAYERCODE                string `json:"PLAYERCODE"`
	PLAYER_SLUG               string `json:"PLAYER_SLUG"`
	ROSTERSTATUS              int    `json:"ROSTERSTATUS"`
	TEAM_ABBREVIATION         string `json:"TEAM_ABBREVIATION"`
	TEAM_CITY                 string `json:"TEAM_CITY"`
	TEAM_CODE                 string `json:"TEAM_CODE"`
	TEAM_ID                   int64  `json:"TEAM_ID"`
	TEAM_NAME                 string `json:"TEAM_NAME"`
	TEAM_SLUG                 string `json:"TEAM_SLUG"`
	TO_YEAR                   int    `json:"TO_YEAR"`
}

type CommonAllPlayersResponse struct {
	CommonAllPlayers []Player `json:"CommonAllPlayers"`
}

var playerCmd = &cobra.Command{
	Use:     "players",
	Short:   "Search players for stats",
	Long:    "View or list players",
	Aliases: []string{"player"},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var playerSearchCmd = &cobra.Command{
	Use:     "search <search terms...>",
	Short:   "Search players for stats",
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		searchTerm := strings.Join(args, " ")
		fmt.Println("Searching for:", searchTerm)

		players := nag.NewCommonAllPlayers()
		err := players.Get()
		if err != nil {
			panic(err)
		}

		if players.Response == nil {
			panic("no response")
		}

		n := nag.Map(*players.Response)
		// var result []Player
		var result CommonAllPlayersResponse
		mapstructure.Decode(n, &result)
		var id int
		// playerMap := map[string]any{}
		for _, data := range result.CommonAllPlayers {
			if fuzzy.MatchFold(searchTerm, data.DISPLAY_FIRST_LAST) {
				fmt.Printf("%#v\n", data)
				id = data.PERSON_ID
				break
			}
			// playerMap[fmt.Sprintf("%d", data.PERSON_ID)] = data.DISPLAY_FIRST_LAST
			// playerMap[data.DISPLAY_FIRST_LAST] = data.PERSON_ID
		}

		playerInfo := nag.NewCommonPlayerInfo(fmt.Sprintf("%d", id))
		err = playerInfo.Get()
		if err != nil {
			panic(err)
		}

		if playerInfo.Response == nil {
			panic("no response")
		}

		p := nag.Map(*playerInfo.Response)

		var playerInfoResult map[string]any
		err = mapstructure.Decode(p, &playerInfoResult)
		if err != nil {
			log.Fatal(err)
		}

		playerGameLog := nag.NewPlayerGameLog(strconv.Itoa(id))
		err = playerGameLog.Get()
		if err != nil {
			log.Fatal(err)
		}

		d := nag.Map(*playerGameLog.Response)

		// fmt.Printf("%#v\n", playerGameLog.Response)

		var playerGameLogResult map[string]any
		err = mapstructure.Decode(d, &playerGameLogResult)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("%#v\n", playerGameLogResult["PlayerGameLog"])

		game_logs := []GameLogEntry{}
		for _, d := range playerGameLogResult["PlayerGameLog"].([]map[string]any) {
			entryJSON, _ := json.Marshal(d)
			var e GameLogEntry
			json.Unmarshal(entryJSON, &e)
			game_logs = append(game_logs, e)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"#",
			"MATCHUP",
			"GAME_DATE",
			"WL",
			"MIN",
			"PTS",
			"FGM",
			"FGA",
			"FG%",
			"3PM",
			"3PA",
			"3P%",
			"FTM",
			"FTA",
			"FT%",
			"OREB",
			"DREB",
			"REB",
			"AST",
			"STL",
			"BLK",
			"TOV",
			"PF",
			"+/-",
		})
		rows := [][]string{}
		for i, g := range game_logs {
			// fmt.Printf("%#v\n", g)
			row := []string{
				strconv.Itoa(len(game_logs) - i),
				g.MATCHUP,
				g.GAME_DATE,
				g.WL,
				strconv.Itoa(g.MIN),
				strconv.Itoa(g.PTS),
				strconv.Itoa(g.FGM),
				strconv.Itoa(g.FGA),
				strconv.FormatFloat(g.FG_PCT, 'f', 3, 64),
				strconv.Itoa(g.FG3M),
				strconv.Itoa(g.FG3A),
				strconv.FormatFloat(g.FG3_PCT, 'f', 3, 64),
				strconv.Itoa(g.FTM),
				strconv.Itoa(g.FTA),
				strconv.FormatFloat(g.FT_PCT, 'f', 3, 64),
				strconv.Itoa(g.OREB),
				strconv.Itoa(g.DREB),
				strconv.Itoa(g.REB),
				strconv.Itoa(g.AST),
				strconv.Itoa(g.STL),
				strconv.Itoa(g.BLK),
				strconv.Itoa(g.TOV),
				strconv.Itoa(g.PF),
				strconv.Itoa(g.PLUS_MINUS),
			}
			rows = append(rows, row)
			colors := make([]tablewriter.Colors, len(row))
			switch g.WL {
			case "W":
				for j := range colors {
					colors[j] = tablewriter.Colors{tablewriter.FgGreenColor}
				}
			case "L":
				for j := range colors {
					colors[j] = tablewriter.Colors{tablewriter.FgRedColor}
				}
			}
			table.Rich(row, colors)
		}
		table.Render()
		// fmt.Printf("%#v\n", playerGameLog.Response.ResultSets.)
		// for _, r := range d.ResultSets {
		// 	// fmt.Printf("%#v\n", r.Headers)
		// 	fmt.Printf("%#v\n", r.RowSet)
		// }

		// repr.Println(playerInfoResult)

		// fmt.Printf("%#v\n", playerInfoResult)

		// if results := fuzzy.RankFind(searchTerm, maps.Keys(playerMap)); len(results) > 0 {
		// 	fmt.Printf("%#v\n", results)
		// }
	},
}

func init() {
	playerSearchCmd.Flags().StringP("search", "s", "", "Search term")
	playerCmd.AddCommand(playerSearchCmd)
	rootCmd.AddCommand(playerCmd)
}
