package plugins

import (
	"agus/lib"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"go.mau.fi/whatsmeow/proto/waE2E"
)

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"ghibli"},
		Tags:    "other",
		IsOwner: true,
		Run: func(conn lib.IClient, m lib.M) {
			pict, err := conn.WA.Download(context.Background(), m.Full.Message.ImageMessage)
			if err != nil {
				m.Reply("Terjadi Kesalahan")
				return
			}

			var b bytes.Buffer
			w := multipart.NewWriter(&b)
			fw, err := w.CreateFormFile("image", "image.jpg")
			if err != nil {
				m.Reply("terjadi kesalahan")
				fmt.Println(err)
				return
			}
			fw.Write(pict)

			w.Close()
			url := "https://api.siputzx.my.id/api/image2ghibli"
			resp, err := http.Post(url, w.FormDataContentType(), &b)
			if err != nil {
				return
			}
			if resp.StatusCode != http.StatusOK {
				m.Reply(resp.Status)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				m.Reply("Terjadi Kesalahan")
				return
			}
			conn.SendImage(m.Chat, body, "Done.", &waE2E.ContextInfo{
				QuotedMessage: m.Full.Message,
			})
		},
	})
}
