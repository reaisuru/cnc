package database

import (
	"cnc/pkg/logging"
	"crypto/sha256"
	"database/sql"
	"fmt"
)

type Table interface {
	Create() error
}

type UserTable struct{}
type BlacklistTable struct{}
type LogsTable struct{}
type ApiTable struct{}

// We'll insert util methods here, don't mind it

// closeRows will close the rows without an error.
func closeRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		logging.Global.Fatal().Err(err).Msg("Failed to close rows")
		return
	}
}

// Hash will hash a value and return a hexadecimal string.
func Hash(v []byte) string {
	var sha = sha256.New()
	sha.Write(v)
	return fmt.Sprintf("%x", sha.Sum(nil))
}
