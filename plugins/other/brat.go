package plugins

import (
	"agus/lib"
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	// "io"
	"net/http"
	"net/url"

	"github.com/HugoSmits86/nativewebp"
	// "github.com/chai2010/webp"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"brat"},
		Tags:    "other",
		IsOwner: true,
		Run: func(conn lib.IClient, m lib.M) {
			if m.Query == "" {
				m.Query = "hello"
			}
			url := fmt.Sprintf("https://api.siputzx.my.id/api/m/brat?text=%s&delay=500", url.QueryEscape(m.Query))
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				m.Reply("Terjadi Kesalahan")
				return
			}
			defer resp.Body.Close()
			img, _, err := image.Decode(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			var sticker bytes.Buffer
			if err := nativewebp.Encode(&sticker, img, nil); err != nil {
				fmt.Println(err)
			}

			conn.SendSticker(m.Chat, sticker.Bytes(), &waE2E.ContextInfo{})

		},
	})

}
