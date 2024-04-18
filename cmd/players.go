package cmd

import (
	"fmt"
	"strings"

	"github.com/ksc98/nbacli/nag"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mitchellh/mapstructure"
	"github.com/mitranim/repr"
	"github.com/spf13/cobra"
)

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
				continue
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

		var playerInfoResult map[string]interface{}
		mapstructure.Decode(p, &playerInfoResult)

		repr.Println(playerInfoResult)

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
