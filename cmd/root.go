package cmd

// import (
// 	"github.com/spf13/cobra"
// )

// type TicketMaster struct {
// 	rootCmd *cobra.Command
// }

// func NewTicketMaster() *TicketMaster {
// 	tc := &TicketMaster{
// 		rootCmd: &cobra.Command{
// 			Use:   "ticketmaster",
// 			Short: "Ticketmaster CLI",
// 			// no need to provide the default cobra completion command
// 			CompletionOptions: cobra.CompletionOptions{
// 				DisableDefaultCmd: true,
// 			},
// 		},
// 	}

// 	return tc
// }

// func (ticketmaster *TicketMaster) Start() error {
// 	ticketmaster.rootCmd.AddCommand(serveApiCommand())
// 	ticketmaster.rootCmd.AddCommand(waitingServiceCommand())
// 	ticketmaster.rootCmd.AddCommand(activationReservationCommand())

// 	return ticketmaster.execute()
// }

// func (ticketmaster *TicketMaster) execute() error {
// 	err := ticketmaster.rootCmd.Execute()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
