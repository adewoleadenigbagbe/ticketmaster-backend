package migrate

const migrationTableName = "_migration"

// type Runner struct {
// 	db        *dbx.DB
// 	tableName string
// }

// // builds the runner struct and call createMigrationsTable function
// func NewRunner(db *dbx.DB) (*Runner, error) {
// 	runner := &Runner{
// 		db:        db,
// 		tableName: migrationTableName,
// 	}

// 	err := runner.createMigrationsTable()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return runner, nil
// }

// // create a migration table if does not exist , builds up the query to do that
// func (r *Runner) createMigrationsTable() error {
// 	rawQuery := fmt.Sprintf(
// 		"CREATE TABLE IF NOT EXISTS %v (file VARCHAR(255) PRIMARY KEY NOT NULL, applied INTEGER NOT NULL)",
// 		r.db.QuoteTableName(r.tableName),
// 	)

// 	_, err := r.db.NewQuery(rawQuery).Execute()

// 	return err
// }

// func (r *Runner) Run(args ...string) {

// }
