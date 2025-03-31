package database

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type FloodLog struct {
	// ID is an increasing ID for the DB, I guess
	ID int

	// AttackID is the attack id sent to the client
	AttackID int

	// UserID is the user id of the user that sent the attack
	UserID int

	// MethodID is the method id used on the target.
	MethodID int

	// Targets are the targets being attacked
	Targets []string

	// Duration is the attack duration wow who would've thought
	Duration int

	// Clients are the clients that the user has sent with.
	Clients int

	// Started is the time, where the attack/flood has been launched.
	Started time.Time

	// Ended is the end time of the attack.
	Ended time.Time

	IsAPI bool
}

func (t *LogsTable) handleDailyAttacks() {
	for {
		all, err := t.SelectAll()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// kek
		for _, log := range all {
			if log.Ended.Day() != time.Now().Day() {
				_ = t.DeleteID(log.ID)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

// Create simply creates the table.
func (t *LogsTable) Create() error {
	go t.handleDailyAttacks()

	createTableSQL := `CREATE TABLE IF NOT EXISTS logs (
		id SERIAL PRIMARY KEY, 
		attack_id INT,
		user_id INT,
		method_id INT,
		
		targets TEXT[],
		duration INT,
		clients INT,
		api boolean,
		
		time_started timestamptz,
		time_end timestamptz
	);`

	_, err := Sql.Exec(createTableSQL)
	return err
}

// Insert inserts an attack log into the database.
func (t *LogsTable) Insert(log *FloodLog) error {
	_, err := Sql.Exec(`INSERT INTO logs (
                   attack_id, 
                   user_id, 
                   method_id, 
                  
                   targets, 
                   duration, 
                   clients, 
                   api,
                  
                   time_started, 
                   time_end) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,

		log.AttackID,
		log.UserID,
		log.MethodID,

		pq.Array(log.Targets),
		log.Duration,
		log.Clients,
		log.IsAPI,

		log.Started,
		log.Ended,
	)

	return err
}

// SelectAll will get all logs from the database.
func (t *LogsTable) SelectAll() ([]*FloodLog, error) {
	// select all logs with all information
	rows, err := Sql.Query("SELECT * FROM logs")
	if err != nil {
		return nil, err
	}

	// scan all values and close rows after everything is done
	defer closeRows(rows)
	return t.scanAll(rows)
}

// DeleteAll deletes all logs
func (t *LogsTable) DeleteAll() error {
	_, err := Sql.Exec("DELETE FROM logs")
	return err
}

// DeleteID deletes log by ID
func (t *LogsTable) DeleteID(id int) error {
	_, err := Sql.Exec("DELETE FROM logs WHERE id=$1", id)
	return err
}

// SelectRunning gets the current running floods
func (t *LogsTable) SelectRunning(api bool) ([]*FloodLog, error) {
	// select all logs with all information
	rows, err := Sql.Query("SELECT * FROM logs WHERE time_end > NOW() AND api=$1", api)
	if err != nil {
		return nil, err
	}

	// scan all values and close rows after everything is done
	defer closeRows(rows)
	return t.scanAll(rows)
}

// RunningCount will get the current running flood count
func (t *LogsTable) RunningCount(api bool) (count int) {
	// select the count of running attacks
	row := Sql.QueryRow("SELECT COUNT(*) FROM logs WHERE time_end > NOW() AND api=$1", api)

	// scan in the count
	err := row.Scan(&count)
	if err != nil {
		return 0
	}

	return count
}

func (t *LogsTable) SelectByAttackID(id int) (*FloodLog, error) {
	// select log by attack id
	rows, err := Sql.Query("SELECT * FROM logs WHERE attack_id=$1", id)
	if err != nil {
		return nil, err
	}

	// close rows when everything done yippee
	defer closeRows(rows)

	// if no existo we throw error yes
	if !rows.Next() {
		return nil, errors.New("no such flood with that id")
	}

	// finally scan
	return t.scan(rows)
}

// Modify will modify the attack log.
func (t *FloodLog) Modify() error {
	query := `UPDATE logs SET attack_id=$1, user_id=$2, method_id=$3, targets=$4, duration=$5, clients=$6, time_started=$7, time_end=$8 WHERE id=$9;`
	_, err := Sql.Exec(query,
		t.AttackID,
		t.UserID,
		t.MethodID,

		pq.Array(t.Targets),
		t.Duration,
		t.Clients,

		t.Started,
		t.Ended,

		t.ID,
	)

	return err
}

func (t *LogsTable) LastByUserID(id int) (*FloodLog, error) {
	query := "SELECT * FROM logs WHERE id = (SELECT MAX(id) FROM logs WHERE user_id=$1)"

	// select last attack by user
	rows, err := Sql.Query(query, id)
	if err != nil {
		return nil, err
	}

	// close rows after everything is done
	defer closeRows(rows)

	// check if there is a next role
	if !rows.Next() {
		return nil, errors.New("no last attack recorded")
	}

	return t.scan(rows)
}

func (t *LogsTable) LastByUserID_2(id int, api bool) (*FloodLog, error) {
	query := "SELECT * FROM logs WHERE id = (SELECT MAX(id) FROM logs WHERE user_id=$1 AND api=$2)"

	// select last attack by user
	rows, err := Sql.Query(query, id, api)
	if err != nil {
		return nil, err
	}

	// close rows after everything is done
	defer closeRows(rows)

	// check if there is a next role
	if !rows.Next() {
		return nil, errors.New("no last attack recorded")
	}

	return t.scan(rows)
}

func (t *LogsTable) LastGlobalFlood(api bool) (*FloodLog, error) {
	query := "SELECT * FROM logs WHERE id = (SELECT MAX(id) FROM logs WHERE api=$1)"

	// select last attack
	rows, err := Sql.Query(query, api)
	if err != nil {
		return nil, err
	}

	// close rows after everything is done
	defer closeRows(rows)

	// check if there is a next role
	if !rows.Next() {
		return nil, nil
	}

	return t.scan(rows)
}

// scan scans an attack log, wow
func (t *LogsTable) scan(rows *sql.Rows) (*FloodLog, error) {
	log := new(FloodLog)
	return log, rows.Scan(
		&log.ID,
		&log.AttackID,
		&log.UserID,
		&log.MethodID,

		pq.Array(&log.Targets),
		&log.Duration,
		&log.Clients,
		&log.IsAPI,

		&log.Started,
		&log.Ended,
	)
}

// scanAll scans all possible rows
func (t *LogsTable) scanAll(rows *sql.Rows) ([]*FloodLog, error) {
	var logs = make([]*FloodLog, 0)

	for rows.Next() {
		log, err := t.scan(rows)

		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// Drop removes the attack from the logs.
func (t *FloodLog) Drop() error {
	_, err := Sql.Exec("DELETE FROM logs WHERE id=$1", t.ID)
	return err
}
