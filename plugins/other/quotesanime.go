package plugins

import (
	"agus/lib"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

type Quote struct {
	Link     string `json:"link"`
	Gambar   string `json:"gambar"`
	Karakter string `json:"karakter"`
	Anime    string `json:"anime"`
	Episode  string `json:"episode"`
	UpAt     string `json:"up_at"`
	Quotes   string `json:"quotes"`
}

func quotesAnime() ([]Quote, error) {
	page := rand.Intn(184)
	url := fmt.Sprintf("https://otakotaku.com/quote/feed/%d", page)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var hasil []Quote
	doc.Find("div.kotodama-list").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a.kuroi").Attr("href")
		fmt.Println(link)
		gambar, _ := s.Find("img").Attr("data-src")
		karakter := s.Find("div.char-name").Text()
		anime := s.Find("div.anime-title").Text()
		episode := s.Find("div.meta").Text()
		upAt := s.Find("small.meta").Text()
		quotes := s.Find("div.quote").Text()

		hasil = append(hasil, Quote{
			Link:     link,
			Gambar:   gambar,
			Karakter: strings.TrimSpace(karakter),
			Anime:    strings.TrimSpace(anime),
			Episode:  strings.TrimSpace(episode),
			UpAt:     upAt,
			Quotes:   strings.TrimSpace(quotes),
		})
	})

	return hasil, nil
}

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:     []string{"quotesanime"},
		Tags:    "other",
		IsOwner: false,
		Run: func(conn lib.IClient, m lib.M) {
			quotes, err := quotesAnime()
			if err != nil {
				fmt.Println(err)
			}
			id := rand.Intn(len(quotes))
			hasil := quotes[id]

			str := fmt.Sprintf("```%s```\n\n_%s_", hasil.Quotes, strings.TrimSpace(hasil.Karakter))
			// conn.SendInteractive(m.Chat, str, &waE2E.ContextInfo{
			// 	MentionedJID: []string{m.Sender.String()},
			// 	ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
			// 		ShowAdAttribution: proto.Bool(true),
			// 		SourceURL:         &hasil.Link,
			// 	},
			// })
			conn.SendText(m.Chat, str, &waE2E.ContextInfo{
				MentionedJID: []string{m.Sender.String()},
				ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
					// Title:        &hasil.Anime,
					// Body:         &hasil.Karakter,
					SourceURL:         &hasil.Link,
					ShowAdAttribution: proto.Bool(true),
				},
			})
		},
	})

}
