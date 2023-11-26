package cmd

import (
	"github.com/Wolechacho/ticketmaster-backend/apis"
	"github.com/spf13/cobra"
)

func ServeApiCommand() *cobra.Command {
	var apiCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve the API on the Specified host",
		Run: func(cmd *cobra.Command, args []string) {
			apis.Serve()
		},
	}

	return apiCmd
}
