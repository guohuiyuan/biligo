package bilibili

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
)

type cookieJar struct {
	mu  sync.RWMutex
	jar http.CookieJar
}

func newCookieJar() *cookieJar {
	jar, _ := cookiejar.New(nil)
	return &cookieJar{jar: jar}
}

func (j *cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.jar.SetCookies(u, cookies)
}

func (j *cookieJar) Cookies(u *url.URL) []*http.Cookie {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.jar.Cookies(u)
}

func (j *cookieJar) SetCredential(rawURL string, credential *Credential) {
	if credential == nil {
		return
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	j.SetCookies(parsed, credential.ToHTTPCookies())
}
