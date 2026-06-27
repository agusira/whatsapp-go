/*
   tiktok.go — TikTok video downloader (no watermark).

   @creator: Lune

   Description:
     Resolve a TikTok URL via tikwm.com, fallback to tmate.cc, save MP4.

   Usage:
     1. Edit `tiktokURL` constant below with your target video URL.
     2. Run: go run ./downloader/tiktok/
     3. Output saved to ./out/<author>_<id>.mp4

   No install needed — only stdlib. Requires Go 1.22+.
*/

package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	tiktokURL = "https://www.tiktok.com/@sadampermanawiyana/video/7646767031202663701"
	outDir    = "./tmp"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type video struct {
	id, title, author, url, referer, source string
	size                                    int64
}

func Tiktokdl(link string) ([]byte, error) {
	log.SetFlags(log.Ltime)
	ctx := context.Background()
	client := newClient()

	v, err := fetchTikwm(ctx, client, link)
	if err != nil {
		log.Printf("tikwm: %v", err)
		v, err = fetchTmate(ctx, client, link)
		if err != nil {
			log.Fatalf("tmate: %v", err)
		}
	}

	hasil, err := download(ctx, client, v, outDir)
	if err != nil {
		log.Fatalf("download: %v", err)
	}
	return hasil, err
}

const tikwmEndpoint = "https://www.tikwm.com/api/"

type tikwmResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Play   string `json:"play"`
		HDPlay string `json:"hdplay"`
		Size   int64  `json:"size"`
		HDSize int64  `json:"hd_size"`
		Author struct {
			UniqueID string `json:"unique_id"`
		} `json:"author"`
	} `json:"data"`
}

func fetchTikwm(ctx context.Context, c *http.Client, tiktok string) (*video, error) {
	form := url.Values{"url": {tiktok}, "hd": {"1"}}
	resp, err := doRetry(ctx, c, func(ctx context.Context) (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, tikwmEndpoint, strings.NewReader(form.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")
		return req, nil
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "json") {
		return nil, fmt.Errorf("non-json response %q (possibly geo-blocked)", ct)
	}
	var r tikwmResp
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&r); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("code=%d msg=%s", r.Code, r.Msg)
	}

	mp4, size := r.Data.HDPlay, r.Data.HDSize
	if mp4 == "" {
		mp4, size = r.Data.Play, r.Data.Size
	}
	if mp4 == "" {
		return nil, fmt.Errorf("empty play and hdplay")
	}
	if strings.HasPrefix(mp4, "/") {
		mp4 = "https://www.tikwm.com" + mp4
	}
	return &video{
		id:      r.Data.ID,
		title:   r.Data.Title,
		author:  r.Data.Author.UniqueID,
		url:     mp4,
		size:    size,
		source:  "tikwm",
		referer: "https://www.tikwm.com/",
	}, nil
}

const (
	tmateHome   = "https://tmate.cc/"
	tmateAction = "https://tmate.cc/action"
)

var (
	tmateTokenRe = regexp.MustCompile(`name="token"\s+type="hidden"\s+value="([a-f0-9]{32})"`)
	tmateHrefRe  = regexp.MustCompile(`href="([^"]+)"[^>]*class="abutton is-success[^"]*"[^>]*>\s*<span><span>Download without Watermark`)
)

func fetchTmate(ctx context.Context, c *http.Client, tiktok string) (*video, error) {
	resp, err := doRetry(ctx, c, func(ctx context.Context) (*http.Request, error) {
		return http.NewRequestWithContext(ctx, http.MethodGet, tmateHome, nil)
	})
	if err != nil {
		return nil, fmt.Errorf("home: %w", err)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("home read: %w", err)
	}
	m := tmateTokenRe.FindSubmatch(body)
	if m == nil {
		return nil, fmt.Errorf("token not found in homepage")
	}
	token := string(m[1])

	form := url.Values{"url": {tiktok}, "token": {token}}
	resp, err = doRetry(ctx, c, func(ctx context.Context) (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, tmateAction, strings.NewReader(form.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Origin", "https://tmate.cc")
		req.Header.Set("Referer", tmateHome)
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("Accept", "*/*")
		return req, nil
	})
	if err != nil {
		return nil, fmt.Errorf("action: %w", err)
	}
	defer resp.Body.Close()

	var r struct {
		Error bool   `json:"error"`
		Data  string `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&r); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if r.Error {
		return nil, fmt.Errorf("tmate reported error")
	}
	hm := tmateHrefRe.FindStringSubmatch(r.Data)
	if hm == nil {
		return nil, fmt.Errorf("no-watermark anchor not found")
	}
	return &video{
		url:     hm[1],
		source:  "tmate",
		referer: tmateHome,
	}, nil
}

func download(ctx context.Context, c *http.Client, v *video, dir string) ([]byte, error) {
	var body []byte
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return body, err
	}
	dst := filepath.Join(dir, filename(v))
	if _, err := os.Stat(dst); err == nil {
		log.Printf("[%s] skip (exists): %s", v.source, dst)
		return body, err
	}

	resp, err := doRetry(ctx, c, func(ctx context.Context) (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.url, nil)
		if err != nil {
			return nil, err
		}
		if v.referer != "" {
			req.Header.Set("Referer", v.referer)
		}
		return req, nil
	})
	if err != nil {
		return body, fmt.Errorf("mp4 get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("mp4 status %d", resp.StatusCode)
	}
	body, err = io.ReadAll(resp.Body)

	// tmp, err := os.CreateTemp(dir, filepath.Base(dst)+".*.partial")
	// if err != nil {
	// 	return body, err
	// }
	// tmpName := tmp.Name()
	// defer os.Remove(tmpName)

	// n, err := io.Copy(tmp, resp.Body)
	// fmt.Println(resp.Body)
	// if cerr := tmp.Close(); err == nil {
	// 	err = cerr
	// }
	// if err != nil {
	// 	return body, fmt.Errorf("save: %w", err)
	// }
	// if err := os.Rename(tmpName, dst); err != nil {
	// 	return body, err
	// }
	// log.Printf("[%s] saved %s (%s)", v.source, dst, humanBytes(n))
	return body, nil
}

func filename(v *video) string {
	var parts []string
	if v.author != "" {
		parts = append(parts, safeName(v.author))
	}
	if v.id != "" {
		parts = append(parts, v.id)
	}
	if len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("tiktok_%d", time.Now().Unix()))
	}
	return strings.Join(parts, "_") + ".mp4"
}

func newClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{Timeout: 60 * time.Second, Jar: jar}
}

func doRetry(ctx context.Context, c *http.Client, build func(context.Context) (*http.Request, error)) (*http.Response, error) {
	const maxRetry = 3
	const base = 500 * time.Millisecond
	var lastErr error
	for attempt := 0; attempt <= maxRetry; attempt++ {
		if attempt > 0 {
			d := base << (attempt - 1)
			d += time.Duration(rand.Int63n(int64(d)/2 + 1))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(d):
			}
		}
		req, err := build(ctx)
		if err != nil {
			return nil, err
		}
		if req.Header.Get("User-Agent") == "" {
			req.Header.Set("User-Agent", userAgent)
		}
		resp, err := c.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("http %d", resp.StatusCode)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("after %d attempts: %w", maxRetry+1, lastErr)
}

var illegalChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func safeName(s string) string {
	s = strings.TrimSpace(s)
	s = illegalChars.ReplaceAllString(s, "_")
	s = strings.Trim(s, "._-")
	if len(s) > 100 {
		s = s[:100]
	}
	if s == "" {
		return "untitled"
	}
	return s
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(n)/float64(div), "KMGT"[exp])
}
