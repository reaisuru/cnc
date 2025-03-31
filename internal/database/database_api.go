package database

import (
	"database/sql"
)

type ApiEntry struct {
	ID int

	ApiName string
	ApiLink string

	Method string
	Times  int
}

// Create creates the table
func (t *ApiTable) Create() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS apis (
		id SERIAL PRIMARY KEY,
		api_name TEXT NOT NULL,
		api_link TEXT NOT NULL,
		method TEXT NOT NULL,
		times INT NOT NULL
	);`

	_, err := Sql.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

// Insert will insert a user into the database.
func (t *ApiTable) Insert(apiName, apiLink, method string, times int) error {
	_, err := Sql.Exec(`INSERT INTO apis (api_name, api_link, method, times) VALUES ($1, $2, $3, $4);`, apiName, apiLink, method, times)
	return err
}

func (t *ApiTable) SelectAll() ([]*ApiEntry, error) {
	var targets []*ApiEntry

	// select all targets with all information
	rows, err := Sql.Query("SELECT * FROM apis")
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

func (t *ApiTable) SelectAllByMethod(method string) ([]*ApiEntry, error) {
	var targets []*ApiEntry

	// select all targets with all information
	rows, err := Sql.Query("SELECT * FROM apis WHERE method=$1", method)
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

func (t *ApiTable) SelectByName(name string) ([]*ApiEntry, error) {
	var targets []*ApiEntry

	// select all targets with all information
	rows, err := Sql.Query("SELECT * FROM apis WHERE api_name=$1", name)
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

// Drop will remove the blacklisted target from the database.
func (t *ApiEntry) Drop() error {
	_, err := Sql.Exec("DELETE FROM apis WHERE id=$1", t.ID)
	return err
}

// scan scans a blacklisted target into a pointer
func (t *ApiTable) scan(rows *sql.Rows) (*ApiEntry, error) {
	target := new(ApiEntry)
	return target, rows.Scan(
		&target.ID,
		&target.ApiName,
		&target.ApiLink,
		&target.Method,
		&target.Times,
	)
}
