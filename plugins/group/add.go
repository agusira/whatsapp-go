package plugins

import (
	"agus/lib"
	"encoding/json"

	"go.mau.fi/whatsmeow/types"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"add"},
		Tags:    "group",
		IsOwner: true,
		Run: func(conn lib.IClient, m lib.M) {
			if !m.IsGroup {
				m.Reply("Perintah di Group")
				return
			}
			if m.IsGroup {
				if m.Quoted != nil {
					jid, _ := types.ParseJID(*m.ContextInfo.Participant)

					if _, err := conn.AddParticipant(m.Chat, jid); err != nil {
						m.Reply("Terjadi Kesalahan")
						return
					}
					m.Reply("Berhasil Menambahkan member bau")
				} else if len(m.Args) > 1 {
					jid := types.NewJID(m.Args[1], types.DefaultUserServer)
					g, err := conn.AddParticipant(m.Chat, jid)
					if err != nil {
						m.Reply("Terjadi Kesalahan")
						return
					}
					if g[0].Error != 0 {
						res, _ := json.Marshal(g)
						m.Reply(string(res))
						return
					}
					m.Reply("Berhasil Menambahkan Member Baru")
				} else {
					m.Reply("Example: add 62xxx")
				}
			}
			return
		},
	})
}
