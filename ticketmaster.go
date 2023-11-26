package ticketmaster

import (
	"os"

	"github.com/Wolechacho/ticketmaster-backend/cmd"
	"github.com/spf13/cobra"
)

type TicketMaster struct {
	RootCmd *cobra.Command
}

func New() *TicketMaster {
	tc := &TicketMaster{
		RootCmd: &cobra.Command{
			Use:   "[command]",
			Short: "Ticketmaster CLI",
			// no need to provide the default cobra completion command
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}

	return tc
}

func (ticketmaster *TicketMaster) Start() error {
	ticketmaster.RootCmd.AddCommand(cmd.ServeApiCommand())
	ticketmaster.RootCmd.AddCommand(cmd.WaitingServiceCommand())
	ticketmaster.RootCmd.AddCommand(cmd.ActivationReservationCommand())

	return ticketmaster.execute()
}

func (ticketmaster *TicketMaster) execute() error {
	err := ticketmaster.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	return nil
}
