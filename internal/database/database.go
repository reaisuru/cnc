package database

import (
	"cnc/pkg/logging"
	"database/sql"
	"fmt"
)
import _ "github.com/lib/pq"

const (
	Host     = "127.0.0.1"
	Port     = 5432
	Username = "postgres"
	Database = "botnet"
)

var (
	// Sql is the SQL database instance
	Sql *sql.DB

	// User is a table for users. | Declare tables here
	User      = new(UserTable)
	Blacklist = new(BlacklistTable)
	Logs      = new(LogsTable)
	API       = new(ApiTable)

	// tables is there to initialize all tables
	tables = []Table{
		User,
		Blacklist,
		Logs,
		API,
	}
)

func Serve() {
	var err error

	// Open Database
	Sql, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", Host, Port, Username, Database))
	if err != nil {
		logging.Global.Fatal().Err(err).Msg("Failed to open database")
	}

	// Health check
	err = Sql.Ping()
	if err != nil {
		logging.Global.Fatal().Err(err).Msg("Failed to ping database")
	}

	logging.Global.Info().Msg("Successfully connected to PostgreSQL database")

	// Initialize all tables
	for _, t := range tables {
		if err := t.Create(); err != nil {
			logging.Global.Fatal().Err(err).Msg("Failed to create table")
		}
	}
}
