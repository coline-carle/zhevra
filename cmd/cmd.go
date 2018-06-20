package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	setImportdbCmdFlags()
	rootCmd.AddCommand(importdbCmd)
}

// Execute is the command to be run by main
func Execute() {
	rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "zhevra",
	Short: "zhevra is an addon manager for World of Warcraft",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("please provide a command")
	},
}
