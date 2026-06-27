package plugins

import (
	"agus/lib"
	"go.mau.fi/whatsmeow/types"
	// "strings"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"kick"},
		Tags:    "group",
		IsOwner: true,
		Run: func(conn lib.IClient, m lib.M) {
			if !m.IsGroup {
				m.Reply("Perintah di Group")
				return
			}
			if m.IsGroup {

				if len(m.ContextInfo.MentionedJID) > 0 {
					var jid []types.JID
					for _, j := range m.ContextInfo.MentionedJID {
						p, _ := types.ParseJID(j)
						jid = append(jid, p)
					}
					if _, err := conn.RemoveParticipant(m.Chat, jid); err != nil {
						m.Reply("Terjadi Kesalahan")
						return
					}
					m.Reply("Berhasil Mengeluarkan Member")
				} else if m.Quoted != nil {
					// jid, _ := types.ParseJID(*m.ContextInfo.Participant)
					_, err := conn.RemoveParticipant(m.Chat, []types.JID{m.Quoted.Participant})
					if err != nil {
						m.Reply("Terjadi Kesalahan")
						return
					}
					m.Reply("Berhasil Mengeluarkan Member")
					return
				} else if len(m.Args) > 1 {
					jid := types.NewJID(m.Args[1], types.DefaultUserServer)
					if _, err := conn.RemoveParticipant(m.Chat, []types.JID{jid}); err != nil {
						m.Reply("Terjadi Kesalahan")
						return
					}
					m.Reply("Berhasil Mengeluarkan Member")
				} else {
					m.Reply("Example: kick 62xxx or Reply Target")
					return
				}
			}
		},
	})
}
