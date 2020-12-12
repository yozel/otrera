package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yozel/otrera/internal/update"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update AWS data",
	Run: func(cmd *cobra.Command, args []string) {
		update.Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
