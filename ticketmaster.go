package ticketmaster

// import (
// 	"github.com/Wolechacho/ticketmaster-backend/cmd"
// 	"github.com/spf13/cobra"
// )

// type TicketMaster struct {
// 	RootCmd *cobra.Command
// }

// func New() *TicketMaster {
// 	tc := &TicketMaster{
// 		RootCmd: &cobra.Command{
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

// func (ticketmaster *TicketMaster) Start() int {
// 	ticketmaster.RootCmd.AddCommand(cmd.NewMigrateCommand())
// 	return 1
// }

// func (ticketmaster *TicketMaster) Execute() error {
// 	return nil
// }
