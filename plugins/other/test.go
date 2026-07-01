package plugins

import (
	"agus/lib"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"test"},
		Tags:    "other",
		IsOwner: false,
		Run: func(conn lib.IClient, m lib.M) {
			m.Reply("Bot Active")
		},
	})
}
