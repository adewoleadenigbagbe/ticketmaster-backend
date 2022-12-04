package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func AddNewMigrationCommand() {
	desc := `
Supported arguments are:
- up                   - runs all available migrations.
- down [number]        - reverts the last [number] applied migrations.
- create name [folder] - creates new migration template file.
- collections [folder] - (Experimental) creates new migration file with the most recent local collections configuration.`

	migrationCommand := cobra.Command{
		Use:       "migrate",
		Short:     "Executes DB migration scripts",
		ValidArgs: []string{"up", "down"},
		Long:      desc,
		Run: func(command *cobra.Command, args []string) {
			
		},
	}
}
