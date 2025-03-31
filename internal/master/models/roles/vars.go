package roles

import (
	"cnc/internal"
	"log"
)

var (
	List = map[string]*Role{
		"admin": {
			DisplayName: "A",
			Color:       "#ffffff",
		},
	}
)

type Roles struct{}

func (r *Roles) Serve() {
	err := internal.Config.Unmarshal(&List, "roles")
	if err != nil {
		log.Fatal(err)
	}
}
