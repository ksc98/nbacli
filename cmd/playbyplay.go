package cmd

import (
	"github.com/spf13/cobra"
)

var (
// id string
)

func init() {
	rootCmd.AddCommand(PlayByPlayCmd)
}

var PlayByPlayCmd = &cobra.Command{
	Use:     "playbyplay [GameID]",
	Aliases: []string{"pbp", "play"},
	Short:   "List games today",
	Long:    `Given a game id as the argument, display a UI for the game`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		playbyplay(id)
	},
	Args: cobra.MinimumNArgs(1),
}

func playbyplay(gameID string) {
	// pbp := play.New(gameID)
	// err := pbp.Get()
	// if err != nil {
	// 	log.Fatalf("Error looking up game %s", gameID)
	// }

	// pbp.UI()
}
