package cmd

import (
	"github.com/Wolechacho/ticketmaster-backend/activereservation"
	"github.com/spf13/cobra"
)

func activationReservationCommand() *cobra.Command {
	var reservationCmd = &cobra.Command{
		Use:   "reservation",
		Short: "Reserve Seat Service",
		Long:  `Service reserve seat for a cinema show`,
		Run: func(cmd *cobra.Command, args []string) {
			activereservation.Run()
		},
	}

	return reservationCmd
}
