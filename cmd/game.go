/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	"github.com/ksc98/nbacli/ui"

	"github.com/spf13/cobra"
)

// args
var date = ""

// var gameID = ""

var hasYesterday = false
var hasTomorrow = false

// gameCmd represents the game command
var gameCmd = &cobra.Command{
	Use:   "games",
	Short: "Get the NBA schedule for a specific date",
	Run: func(cmd *cobra.Command, args []string) {

		// no date then get today's games
		dateArg := time.Now()

		if hasYesterday {
			dateArg = time.Now().AddDate(0, 0, -1)
		}
		if hasTomorrow {
			dateArg = time.Now().AddDate(0, 0, 1)
		}
		if date != "" {
			dateArg, _ = time.Parse("20060102", date)
		}

		// start the tui
		ui.StartTea(dateArg)

	},
}

var StandingCmd = &cobra.Command{
	Use:   "standings",
	Short: "Get the NBA standings for the current season",
	Run: func(cmd *cobra.Command, args []string) {
		// start the tui
		ui.StartStanding()
	},
}

func init() {
	rootCmd.AddCommand(gameCmd)
	rootCmd.PersistentFlags().StringVarP(&date, "date", "d", "", "Date to get the schedule for (YYYYMMDD)")
	rootCmd.PersistentFlags().BoolVarP(&hasYesterday, "yesterday", "y", false, "Get yesterday's games")
	rootCmd.PersistentFlags().BoolVarP(&hasTomorrow, "tomorrow", "t", false, "Get tomorrow's games")
	rootCmd.MarkFlagsMutuallyExclusive("yesterday", "tomorrow", "date")

	rootCmd.AddCommand(StandingCmd)
}
