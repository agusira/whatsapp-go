package plugins

import (
	"agus/lib"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Data struct {
	Status bool   `json:"status"`
	Data   string `json:"data"`
	Error  string `json:"error"`
}

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"ai"},
		Tags:    "other",
		IsOwner: true,
		Run: func(conn lib.IClient, m lib.M) {
			if m.Query == "" {
				m.Query = "hello"
			}
			url := fmt.Sprintf("https://api.siputzx.my.id/api/ai/luminai?content=%s", url.QueryEscape(m.Query))
			resp, err := http.Get(url)
			if err != nil {
				m.Reply("Terjadi Kesalahan")
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			var data Data
			json.Unmarshal(body, &data)
			fmt.Println(url)
			if !data.Status {
				m.Reply("Terjadi Kesalahan")
				return
			}
			m.Reply(data.Data)
		},
	})

}
