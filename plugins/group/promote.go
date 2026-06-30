package plugins

import (
	"agus/lib"
	"encoding/json"

	"go.mau.fi/whatsmeow/types"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"promote"},
		Tags:    "group",
		IsOwner: true,
		Run: func(conn lib.IClient, m lib.M) {
			if !m.IsGroup {
				m.Reply("Perintah di Group")
				return
			}
			if m.IsGroup {
				if len(m.ContextInfo.MentionedJID) > 0 {
					jid, _ := types.ParseJID(m.ContextInfo.MentionedJID[0])
					if _, err := conn.PromoteParticipant(m.Chat, jid); err != nil {
						m.Reply("Terjadi Kesalahan")
						return
					}
					m.Reply("Berhasil Menjadikan Admin")
				} else if m.Quoted != nil {
					jid, _ := types.ParseJID(*m.ContextInfo.Participant)

					if g, err := conn.PromoteParticipant(m.Chat, jid); err != nil || g[0].Error != 0 {
						m.Reply("Terjadi Kesalahan")
						return
					}
					m.Reply("Berhasil Menjadikan Admin")
				} else if len(m.Args) > 1 {
					jid := types.NewJID(m.Args[1], types.DefaultUserServer)
					g, err := conn.PromoteParticipant(m.Chat, jid)
					if err != nil {
						m.Reply("Terjadi Kesalahan")
						return
					}
					if g[0].Error != 0 {
						res, _ := json.Marshal(g)
						m.Reply(string(res))
						return
					}
					m.Reply("Berhasil Menjadikan Admin")
				} else {
					m.Reply("Example: add 62xxx")
					return
				}
			}
		},
	})
}
