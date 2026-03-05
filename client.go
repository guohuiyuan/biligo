package bilibili

import (
	"context"
	"net/http"
	"sync"
)

// HTTPDoer abstracts the transport layer for testing and custom clients.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the root entry point for all module APIs.
type Client struct {
	config     Config
	httpClient HTTPDoer
	cookies    *cookieJar

	mu         sync.RWMutex
	credential *Credential

	wbi *wbiManager

	videoOnce  sync.Once
	video      *VideoService
	userOnce   sync.Once
	user       *UserService
	searchOnce sync.Once
	search     *SearchService
	liveOnce   sync.Once
	live       *LiveService
	loginOnce  sync.Once
	login      *LoginService
}

// NewClient creates a new Bilibili client.
func NewClient(opts ...Option) *Client {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	jar := newCookieJar()
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.Timeout,
			Jar:     jar,
		}
	}

	return &Client{
		config:     cfg,
		httpClient: httpClient,
		cookies:    jar,
		wbi:        newWBIManager(),
	}
}

// Config returns a copy of the client config.
func (c *Client) Config() Config {
	return c.config
}

// SetCredential updates the active credential and syncs cookies into the jar.
func (c *Client) SetCredential(credential *Credential) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.credential = credential.Clone()
	if c.credential == nil {
		return
	}

	for _, base := range []string{c.config.BaseURL, c.config.PassportBaseURL, c.config.LiveBaseURL} {
		c.cookies.SetCredential(base, c.credential)
	}
	// 私信/会话接口走 api.vc.bilibili.com，需单独注入凭证
	c.cookies.SetCredential("https://api.vc.bilibili.com", c.credential)
}

// Credential returns a defensive copy of the active credential.
func (c *Client) Credential() *Credential {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.credential.Clone()
}

// Video returns the video service.
func (c *Client) Video() *VideoService {
	c.videoOnce.Do(func() {
		c.video = &VideoService{client: c}
	})
	return c.video
}

// User returns the user service.
func (c *Client) User() *UserService {
	c.userOnce.Do(func() {
		c.user = &UserService{client: c}
	})
	return c.user
}

// Search returns the search service.
func (c *Client) Search() *SearchService {
	c.searchOnce.Do(func() {
		c.search = &SearchService{client: c}
	})
	return c.search
}

// Live returns the live service.
func (c *Client) Live() *LiveService {
	c.liveOnce.Do(func() {
		c.live = &LiveService{client: c}
	})
	return c.live
}

// Login returns the login service.
func (c *Client) Login() *LoginService {
	c.loginOnce.Do(func() {
		c.login = &LoginService{client: c}
	})
	return c.login
}

func (c *Client) credentialOrDefault() *Credential {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.credential == nil {
		return &Credential{}
	}
	return c.credential.Clone()
}

func (c *Client) NewRequest(endpoint endpoint) *RequestBuilder {
	return newRequestBuilder(c, endpoint)
}

// Ping verifies transport and cookie state by requesting nav.
func (c *Client) Ping(ctx context.Context) (*NavInfo, error) {
	return c.Login().Nav(ctx)
}
