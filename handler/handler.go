package handler

import (
	"agus/configs"
	"agus/lib"
	"strings"
	// "go.mau.fi/whatsmeow/proto/waE2E"
)

func Handler(conn *lib.IClient, m *lib.M) {
	// for _, v := range m.Mentions {
	// 	if v.User == conn.WA.Store.ID.User {
	// 		conn.SendText(*conn.WA.Store.ID, m.Text, &waE2E.ContextInfo{
	// 			QuotedMessage: m.Full.Message,
	// 		})
	// 	}
	// }

	if !configs.CONFIG.Public && !m.IsOwner {
		return
	}
	plugins := lib.GetList()
	for _, cmd := range plugins {
		if m.IsCmd && !cmd.NoPrefix && m.Command == cmd.Cmd[0] {
			if cmd.IsOwner && !m.IsOwner {
				// fmt.Println("User biasa")
				return
			}
			go cmd.Run(*conn, *m)
		} else if cmd.NoPrefix && strings.HasPrefix(m.Text, cmd.Cmd[0]) {
			if cmd.IsOwner && !m.IsOwner {
				return
			}
			go cmd.Run(*conn, *m)
		}
	}
}
