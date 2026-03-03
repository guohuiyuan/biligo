package bilibili

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wbiMixinKeyEncTab = []int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
	33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
	61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
	36, 20, 34, 44, 52,
}

type wbiManager struct {
	mu        sync.RWMutex
	mixinKey  string
	expiresAt time.Time
}

func newWBIManager() *wbiManager {
	return &wbiManager{}
}

func (m *wbiManager) sign(ctx context.Context, client *Client, params url.Values) error {
	mixinKey, err := m.mixinKeyFor(ctx, client)
	if err != nil {
		return err
	}

	params.Set("wts", strconv.FormatInt(time.Now().Unix(), 10))
	if params.Get("web_location") == "" {
		params.Set("web_location", "1550101")
	}

	query := make(url.Values, len(params))
	for key, values := range params {
		query[key] = append([]string(nil), values...)
	}

	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for i, key := range keys {
		if i > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(query.Get(key))
	}
	sum := md5.Sum([]byte(builder.String() + mixinKey))
	params.Set("w_rid", hex.EncodeToString(sum[:]))
	return nil
}

func (m *wbiManager) mixinKeyFor(ctx context.Context, client *Client) (string, error) {
	m.mu.RLock()
	if m.mixinKey != "" && time.Now().Before(m.expiresAt) {
		key := m.mixinKey
		m.mu.RUnlock()
		return key, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.mixinKey != "" && time.Now().Before(m.expiresAt) {
		return m.mixinKey, nil
	}

	var nav NavInfo
	if err := newRequestBuilder(client, endpointNav).Do(ctx, &nav); err != nil {
		return "", err
	}

	m.mixinKey = buildMixinKey(extractWBIKey(nav.WBIImg.ImgURL) + extractWBIKey(nav.WBIImg.SubURL))
	m.expiresAt = time.Now().Add(12 * time.Hour)
	return m.mixinKey, nil
}

func extractWBIKey(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	filename := path.Base(parsed.Path)
	return strings.TrimSuffix(filename, path.Ext(filename))
}

func buildMixinKey(raw string) string {
	var builder strings.Builder
	for _, idx := range wbiMixinKeyEncTab {
		if idx >= len(raw) {
			continue
		}
		builder.WriteByte(raw[idx])
	}
	key := builder.String()
	if len(key) > 32 {
		return key[:32]
	}
	return key
}
