package handler

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var allowedHosts = []string{
	"cdn.citilink.ru",
	"static.citilink.ru",
}

func ProxyImage(w http.ResponseWriter, r *http.Request) {
	rawURL := r.URL.Query().Get("url")
	if rawURL == "" {
		http.Error(w, "параметр url обязателен", http.StatusBadRequest)
		return
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || !parsed.IsAbs() {
		http.Error(w, "некорректный url", http.StatusBadRequest)
		return
	}

	if !isAllowedHost(parsed.Host) {
		http.Error(w, "хост не разрешён", http.StatusForbidden)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, rawURL, nil)
	if err != nil {
		http.Error(w, "ошибка создания запроса", http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://www.citilink.ru/")
	req.Header.Set("Accept", "image/webp,image/avif,image/*,*/*;q=0.8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("image proxy: ошибка запроса к %s: %v", rawURL, err)
		http.Error(w, "не удалось загрузить изображение", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	w.Header().Set("Cache-Control", "public, max-age=604800, immutable")

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func isAllowedHost(host string) bool {
	for _, allowed := range allowedHosts {
		if strings.EqualFold(host, allowed) {
			return true
		}
	}
	return false
}
