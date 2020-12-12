package cmd

import (
	"github.com/yozel/otrera/internal/log"

	"github.com/spf13/cobra"
)

// handlErrors wraps a cobra function, and will print and exit on error
//
// We don't use RunE because if that returns an error, Cobra will print usage.
// That behavior can be disabled with SilenceUsage option, but then Cobra arg/flag errors don't display usage. (sigh)
func handleErrors(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := f(cmd, args); err != nil {
			log.Log().Fatal().Err(err).Send()
		}
	}
}
