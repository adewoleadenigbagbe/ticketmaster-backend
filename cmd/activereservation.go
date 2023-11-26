package cmd

import "github.com/spf13/cobra"

func ActivationReservationCommand() *cobra.Command {
	var reservationCmd = &cobra.Command{
		Use:   "reservation",
		Short: "Reserve Seat Service",
		Long:  `Service reserve seat for a cinema show`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return reservationCmd
}
