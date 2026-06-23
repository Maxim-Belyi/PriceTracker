package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsAllowedHost(t *testing.T) {
	h := NewImageProxyHandler()

	testCases := []struct {
		name     string
		host     string
		expected bool
	}{
		{
			name:     "разрешённый хост cdn.citilink.ru",
			host:     "cdn.citilink.ru",
			expected: true,
		},
		{
			name:     "разрешённый хост static.citilink.ru",
			host:     "static.citilink.ru",
			expected: true,
		},
		{
			name:     "регистр не важен — CDN.CITILINK.RU тоже разрешён",
			host:     "CDN.CITILINK.RU",
			expected: true,
		},
		{
			name:     "чужой хост запрещён",
			host:     "evil.com",
			expected: false,
		},
		{
			name:     "пустая строка запрещена",
			host:     "",
			expected: false,
		},
		{
			name:     "похожий, но не тот хост",
			host:     "cdn.citilink.ru.evil.com",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := h.isAllowedHost(tc.host)
			if result != tc.expected {
				t.Errorf("isAllowedHost(%q) = %v, хотели %v", tc.host, result, tc.expected)
			}
		})
	}
}

func TestProxyImage_Validation(t *testing.T) {
	fakeCDN := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Referer") != "https://www.citilink.ru/" {
			t.Errorf("ожидали Referer citilink.ru, получили %q", r.Header.Get("Referer"))
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("fake image bytes"))
	}))
	defer fakeCDN.Close() 

	fakeHost := fakeCDN.Listener.Addr().String()

	h := &ImageProxyHandler{
		client:       fakeCDN.Client(), 
		allowedHosts: []string{fakeHost},
	}

	testCases := []struct {
		name           string
		queryURL       string
		expectedStatus int
	}{
		{
			name:           "нет параметра url → 400",
			queryURL:       "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "некорректный url (без схемы) → 400",
			queryURL:       "not-an-absolute-url",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "относительный url → 400",
			queryURL:       "/relative/path",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "запрещённый хост → 403",
			queryURL:       "https://evil.com/image.jpg",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "разрешённый хост возвращает 200 с картинкой",
			queryURL:       "http://" + fakeHost + "/image.jpg",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			target := "/image-proxy"
			if tc.queryURL != "" {
				target += "?url=" + tc.queryURL
			}

			req := httptest.NewRequest(http.MethodGet, target, nil)
			w := httptest.NewRecorder()

			h.Handle(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("ожидали статус %d, получили %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestProxyImage_PassesContentType(t *testing.T) {
	fakeCDN := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.WriteHeader(http.StatusOK)
	}))
	defer fakeCDN.Close()

	fakeHost := fakeCDN.Listener.Addr().String()
	h := &ImageProxyHandler{
		client:       fakeCDN.Client(),
		allowedHosts: []string{fakeHost},
	}

	req := httptest.NewRequest(http.MethodGet, "/image-proxy?url=http://"+fakeHost+"/img.webp", nil)
	w := httptest.NewRecorder()

	h.Handle(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "image/webp" {
		t.Errorf("ожидали Content-Type image/webp, получили %q", ct)
	}

	cacheControl := w.Header().Get("Cache-Control")
	if cacheControl == "" {
		t.Error("ожидали заголовок Cache-Control, его нет")
	}
}
