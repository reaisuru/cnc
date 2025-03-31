package commands

import (
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"cnc/pkg/pattern"
	"github.com/alexeyco/simpletable"
	"github.com/dustin/go-humanize"
	"sort"
	"strconv"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"logs"},
		Roles:       []string{"admin", "reseller"},
		Description: "View previous attack logs.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			// retrieve all logs from the database
			logs, err := database.Logs.SelectAll()

			if err != nil {
				return err
			}

			sort.Slice(logs, func(i, j int) bool {
				return logs[i].ID < logs[j].ID
			})

			// create table & set headers
			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "ID"},
					{Align: simpletable.AlignCenter, Text: "Username"},
					{Align: simpletable.AlignCenter, Text: "Targets"},
					{Align: simpletable.AlignCenter, Text: "Duration"},
					{Align: simpletable.AlignCenter, Text: "Clients"},
					{Align: simpletable.AlignCenter, Text: "Method ID"},
					{Align: simpletable.AlignCenter, Text: "Ends"},
				},
			}

			// go through all logs and add it to the cell
			for _, flood := range logs {
				user, err := database.User.SelectByID(flood.UserID)
				if err != nil {
					continue
				}

				cell := []*simpletable.Cell{
					{Align: simpletable.AlignRight, Text: strconv.Itoa(flood.ID)},
					{Align: simpletable.AlignLeft, Text: user.Name},
					{Align: simpletable.AlignLeft, Text: pattern.FormatArray(flood.Targets)},
					{Align: simpletable.AlignRight, Text: strconv.Itoa(flood.Duration)},
					{Align: simpletable.AlignRight, Text: strconv.Itoa(flood.Clients)},
					{Align: simpletable.AlignRight, Text: strconv.Itoa(flood.MethodID)},
					{Align: simpletable.AlignLeft, Text: humanize.Time(flood.Ended)},
				}

				table.Body.Cells = append(table.Body.Cells, cell)
			}

			// render table
			table.SetStyle(simpletable.StyleCompact)
			return session.Table(table, 1)
		},
	})
}
