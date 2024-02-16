package cmd

import (
	"github.com/Wolechacho/ticketmaster-backend/waitingservice"
	"github.com/spf13/cobra"
)

func waitingServiceCommand() *cobra.Command {
	var waitingCmd = &cobra.Command{
		Use:   "waitservice",
		Short: "Waiting service for users",
		Long:  `Waiting service for users to book expired seats`,
		Run: func(cmd *cobra.Command, args []string) {
			waitingservice.Run()
		},
	}

	return waitingCmd
}
