package plugins

import (
	"agus/lib"
	"fmt"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"menu"},
		Tags:    "other",
		IsOwner: false,
		Run: func(conn lib.IClient, m lib.M) {
			str := fmt.Sprintf("Halo %s\n\n", m.PushName)
			str += "*INFO BOT :*\n* *Name* : `WhatsGO`\n* *Creator* : `Agus`\n\n"
			str += lib.GetMenu()
			m.Reply(str)
		},
	})
}
