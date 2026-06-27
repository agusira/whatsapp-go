package plugins

import (
	"agus/lib"
	"encoding/json"

	"github.com/robertkrimen/otto"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:      []string{">>"},
		Tags:     "owner",
		IsOwner:  true,
		NoPrefix: true,
		Run: func(conn lib.IClient, m lib.M) {
			vm := otto.New()
			vm.Set("M", m)
			vm.Set("Conn", conn)
			vm.Set("GetList", lib.GetList())

			h, err := vm.Run(m.Query)
			if err != nil {
				m.Reply(err.Error())
				return
			}

			if h.IsObject() {
				var data any
				h, _ := vm.Run("JSON.stringify(" + m.Query + ")")
				json.Unmarshal([]byte(h.String()), &data)
				pe, _ := json.MarshalIndent(data, "", "  ")
				m.Reply(string(pe))
			} else {
				m.Reply(h.String())
			}
		},
	})
}
