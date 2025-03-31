package database

import (
	"database/sql"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

type BlacklistTarget struct {
	ID int

	Prefix       string
	Netmask      uint8
	CreationDate time.Time
}

// Create creates the table
func (t *BlacklistTable) Create() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS blacklist (
		id SERIAL PRIMARY KEY,
		prefix TEXT NOT NULL,
		netmask INT NOT NULL,
		creation_date timestamptz NOT NULL
	);`

	_, err := Sql.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

// Insert will insert a user into the database.
func (t *BlacklistTable) Insert(prefix string, netmask int) error {
	_, err := Sql.Exec(`INSERT INTO blacklist (prefix, netmask, creation_date) VALUES ($1, $2, $3);`, prefix, netmask, time.Now())
	return err
}

func (t *BlacklistTable) SelectAll() ([]*BlacklistTarget, error) {
	var targets []*BlacklistTarget

	// select all targets with all information
	rows, err := Sql.Query("SELECT * FROM blacklist")
	if err != nil {
		return nil, err
	}

	defer closeRows(rows)

	// iterate through all rows and add targets into a slice
	for rows.Next() {
		target, err := t.scan(rows)
		if err != nil {
			return nil, err
		}

		targets = append(targets, target)
	}

	return targets, nil
}

// Select a blacklisted target by prefix and netmask
func (t *BlacklistTable) Select(prefix string, netmask int) (*BlacklistTarget, error) {
	// select blacklisted target from database
	rows, err := Sql.Query("SELECT * FROM blacklist WHERE prefix=$1 AND netmask=$2", prefix, netmask)
	if err != nil {
		return nil, err
	}

	// close rows after everything is done
	defer closeRows(rows)

	// if there's no row the user doesn't exist, so we return an error
	if !rows.Next() {
		return nil, errors.New("blacklisted target doesn't exist")
	}

	// now finally scan user profile
	return t.scan(rows)
}

// Drop will remove the blacklisted target from the database.
func (t *BlacklistTarget) Drop() error {
	_, err := Sql.Exec("DELETE FROM blacklist WHERE id=$1", t.ID)
	return err
}

// scan scans a blacklisted target into a pointer
func (t *BlacklistTable) scan(rows *sql.Rows) (*BlacklistTarget, error) {
	target := new(BlacklistTarget)
	return target, rows.Scan(
		&target.ID,
		&target.Prefix,
		&target.Netmask,
		&target.CreationDate,
	)
}

// Is checks if a target is blacklisted
func (t *BlacklistTable) Is(target uint32, netmask uint8) bool {
	targets, err := Blacklist.SelectAll()
	if err != nil {
		return false
	}

	for _, entry := range targets {
		prefix := binary.BigEndian.Uint32(net.ParseIP(entry.Prefix).To4())

		switch {
		case netmask > entry.Netmask:
			if netshift(prefix, entry.Netmask) == netshift(target, entry.Netmask) {
				return true
			}
		case netmask < entry.Netmask:
			if netshift(target, netmask) == netshift(prefix, netmask) {
				return true
			}
		default:
			if prefix == target {
				return true
			}
		}
	}

	return false
}

func netshift(prefix uint32, netmask uint8) uint32 {
	return prefix >> (32 - netmask)
}
