package cmd

import (
	"github.com/spf13/cobra"
)

func WaitingServiceCommand() *cobra.Command {
	var waitingCmd = &cobra.Command{
		Use:   "waiting",
		Short: "Waiting service for users",
		Long:  `Waiting service for users to book expired seats`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return waitingCmd
}
