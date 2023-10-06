package cmd

import (
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/spf13/cobra"
)

func NewServeCommand(app core.App) *cobra.Command {

	command := &cobra.Command{
		Use:   "serve",
		Short: "Starts the web server (default to :8185)",
		Run: func(cmd *cobra.Command, args []string) {
			e := app.GetEcho()
			if err := e.Start(":8185"); err != nil && err != http.ErrServerClosed {
				e.Logger.Fatal("shutting down the server")
			}
		},
	}
	return command
}
