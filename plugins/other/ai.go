package plugins

import (
	"agus/lib"
	"bytes"
	"encoding/json"
	"io"

	"net/http"
	"os"
)

// 1. Struct buat request body
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// 2. Struct buat response. Ambil yg dibutuhin aja
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func init() {
	lib.AddPlugins(&lib.Plugins{
		Cmd:      []string{"ai"},
		Tags:     "other",
		NoPrefix: false,
		IsOwner:  false,
		Run: func(conn lib.IClient, m lib.M) {
			if m.Query == "" {
				m.Query = "hello"
			}
			res := request(m.Query)
			m.Reply(res)
		},
	})
}

func request(query string) string {
	apiKey := os.Getenv("OPEN_ROUTER") // ganti punya kamu
	siteURL := ""                      // optional
	siteName := "Agus"                 // optional

	// 3. Bikin body JSON nya
	reqBody := ChatRequest{
		Model: "openrouter/free",
		Messages: []Message{
			{Role: "user", Content: query},
		},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err.Error()
	}

	// 4. Bikin request POST. Pakai NewRequest biar bisa set header custom
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return string(err.Error())
	}

	// 5. Set semua header, sama kayak di Python
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", siteURL)        // optional
	req.Header.Set("X-OpenRouter-Title", siteName) // optional

	// 6. Kirim
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	// 7. Baca & parse response
	bodyBytes, _ := io.ReadAll(resp.Body)

	var chatResp ChatResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return "Gagal parse JSON: " + string(bodyBytes)
		// log.Fatal("Gagal parse JSON:", err, "\nBody:", string(bodyBytes))
	}

	// 8. Ambil hasil teksnya
	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content
	} else {
		return ""
	}
}
