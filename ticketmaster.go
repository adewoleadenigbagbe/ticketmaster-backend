package ticketmaster

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/Wolechacho/ticketmaster-backend/cmd"
	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/spf13/cobra"
)

type appWrapper struct {
	core.App
}

type TicketMaster struct {
	// RootCmd is the main console command
	RootCmd *cobra.Command

	*appWrapper
}

func New() *TicketMaster {
	tc := &TicketMaster{
		RootCmd: &cobra.Command{
			Use:   "ticketmaster",
			Short: "Ticketmaster CLI",
			// no need to provide the default cobra completion command
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}

	tc.appWrapper = &appWrapper{
		core.NewBaseApp(),
	}

	return tc
}

func (ticketmaster *TicketMaster) Start() error {
	ticketmaster.RootCmd.AddCommand(cmd.NewServeCommand(ticketmaster.appWrapper.App))

	return ticketmaster.Execute()
}

func (ticketmaster *TicketMaster) Execute() error {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
	}()

	// execute the root command
	go func() {
		defer wg.Done()
		if err := ticketmaster.RootCmd.Execute(); err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

	return nil
}
