package plugins

import (
	"agus/lib"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"rvo"},
		Tags:    "other",
		IsOwner: false,
		Run: func(conn lib.IClient, m lib.M) {
			if m.Quoted == nil {
				m.Reply("Mana Foto/Vid sekali lihatnya?")
				return
			}
			qmsg := m.Quoted.Message
			if qmsg.ImageMessage != nil && *qmsg.ImageMessage.ViewOnce {
				conn.CopyNForward(m.Chat, qmsg)
				return
			} else if qmsg.VideoMessage != nil && *qmsg.VideoMessage.ViewOnce {
				conn.CopyNForward(m.Chat, qmsg)
				return
			} else if qmsg.AudioMessage != nil && *qmsg.AudioMessage.ViewOnce {
				conn.CopyNForward(m.Chat, qmsg)
				return
			}
		},
	})
}
