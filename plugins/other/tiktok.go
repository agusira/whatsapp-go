package plugins

import (
	"agus/lib"
	"regexp"

	"go.mau.fi/whatsmeow/proto/waE2E"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"tiktok", "tt"},
		Tags:    "other",
		IsOwner: false,
		Run: func(conn lib.IClient, m lib.M) {
			vtRegex := regexp.MustCompile("https://vt.tiktok.com/[A-Za-z0-9]{9,11}/?$")

			if vtRegex.MatchString(m.Args[0]) {
				hasil, err := lib.Tiktokdl(m.Args[0])
				if err != nil {
					m.Reply("Terjadi kesalahan")
					return
				}
				conn.SendVideo(m.Chat, hasil, "Ini Video yg anda minta", &waE2E.ContextInfo{})
				return
			}
		},
	})
}
