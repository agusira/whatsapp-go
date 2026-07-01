package plugins

import (
	"agus/lib"

	"github.com/purpshell/meowcaller"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"call"},
		Tags:    "owner",
		IsOwner: true,
		Run: func(conn lib.IClient, m lib.M) {
			if m.Query == "" {
				m.Reply("Bot Active")
				return
			}
			call, _ := conn.Call(m.Query)
			call.OnReady(func() {
				aud, err := meowcaller.MP3File("./call.mp3")
				if err != nil {
					call.Hangup()
				}
				player := call.Play(aud)
				player.OnFinish(func() {
					call.Hangup()
				})
				call.OnEnd(func(reason string) {
					player.Stop()
				})
			})
		},
	})
}
