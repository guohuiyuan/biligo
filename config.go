package bilibili

import "time"

const (
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36"

	apiBase      = "https://api.bilibili.com"
	passportBase = "https://passport.bilibili.com"
	liveBase     = "https://api.live.bilibili.com"

	bilibiliReferer = "https://www.bilibili.com/"
)

var defaultHeaders = map[string]string{
	"Accept":          "application/json, text/plain, */*",
	"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
	"Origin":          "https://www.bilibili.com",
	"Referer":         bilibiliReferer,
}

// Config configures the client transport and anti-spider behavior.
type Config struct {
	BaseURL         string
	PassportBaseURL string
	LiveBaseURL     string
	UserAgent       string
	Timeout         time.Duration
	HTTPClient      HTTPDoer
	EnableDebug     bool
	WBIRetryTimes   int
}

// DefaultConfig returns a safe baseline config.
func DefaultConfig() Config {
	return Config{
		BaseURL:         apiBase,
		PassportBaseURL: passportBase,
		LiveBaseURL:     liveBase,
		UserAgent:       defaultUserAgent,
		Timeout:         30 * time.Second,
		WBIRetryTimes:   3,
	}
}

// Option mutates Config.
type Option func(*Config)

func WithTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.Timeout = timeout
	}
}

func WithUserAgent(userAgent string) Option {
	return func(cfg *Config) {
		cfg.UserAgent = userAgent
	}
}

func WithHTTPClient(httpClient HTTPDoer) Option {
	return func(cfg *Config) {
		cfg.HTTPClient = httpClient
	}
}

func WithDebug(enabled bool) Option {
	return func(cfg *Config) {
		cfg.EnableDebug = enabled
	}
}

func WithWBIRetryTimes(times int) Option {
	return func(cfg *Config) {
		cfg.WBIRetryTimes = times
	}
}
