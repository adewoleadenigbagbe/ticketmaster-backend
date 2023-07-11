package cmd

// func NewMigrateCommand() *cobra.Command {
// 	desc := `
// Supported arguments are:
// - up                   - runs all available migrations.
// - down [number]        - reverts the last [number] applied migrations.
// - create name [folder] - creates new migration template file.
// - collections [folder] - (Experimental) creates new migration file with the most recent local collections configuration.`

// 	command := cobra.Command{
// 		Use:       "migrate",
// 		Short:     "Executes DB migration scripts",
// 		ValidArgs: []string{"up", "down"},
// 		Long:      desc,
// 		Run: func(command *cobra.Command, args []string) {
// 			db.Open()
// 		},
// 	}
// 	return &command
// }
