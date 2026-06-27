package plugins

import (
	"agus/lib"
	"bytes"
	"os/exec"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:      []string{"$"},
		Tags:     "owner",
		IsOwner:  true,
		NoPrefix: true,
		Run: func(conn lib.IClient, m lib.M) {
			if len(m.Args) < 1 {
				m.Reply("Terjadi Kesalahan")
				return
			}
			cmd := exec.Command(m.Args[0], m.Args[1:]...)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				m.Reply("Terjadi Kesalahan")
			}
			m.Reply(out.String())
		},
	})
}
