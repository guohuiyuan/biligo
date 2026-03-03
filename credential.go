package bilibili

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Credential stores cookies required by authenticated APIs.
type Credential struct {
	SessData    string `json:"sessdata"`
	BiliJct     string `json:"bili_jct"`
	DedeUserID  string `json:"dede_user_id"`
	Buvid3      string `json:"buvid3"`
	Buvid4      string `json:"buvid4"`
	AcTimeValue string `json:"ac_time_value"`
}

func NewCredentialFromCookieString(raw string) *Credential {
	credential := &Credential{}
	for _, part := range strings.Split(raw, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch strings.TrimSpace(kv[0]) {
		case "SESSDATA":
			credential.SessData = strings.TrimSpace(kv[1])
		case "bili_jct":
			credential.BiliJct = strings.TrimSpace(kv[1])
		case "DedeUserID":
			credential.DedeUserID = strings.TrimSpace(kv[1])
		case "buvid3":
			credential.Buvid3 = strings.TrimSpace(kv[1])
		case "buvid4":
			credential.Buvid4 = strings.TrimSpace(kv[1])
		case "ac_time_value":
			credential.AcTimeValue = strings.TrimSpace(kv[1])
		}
	}
	return credential
}

func NewCredentialFromHTTPCookies(cookies []*http.Cookie) *Credential {
	credential := &Credential{}
	for _, cookie := range cookies {
		switch cookie.Name {
		case "SESSDATA":
			credential.SessData = cookie.Value
		case "bili_jct":
			credential.BiliJct = cookie.Value
		case "DedeUserID":
			credential.DedeUserID = cookie.Value
		case "buvid3":
			credential.Buvid3 = cookie.Value
		case "buvid4":
			credential.Buvid4 = cookie.Value
		case "ac_time_value":
			credential.AcTimeValue = cookie.Value
		}
	}
	return credential
}

func (c *Credential) Clone() *Credential {
	if c == nil {
		return nil
	}
	clone := *c
	return &clone
}

func (c *Credential) Cookies() map[string]string {
	if c == nil {
		return nil
	}
	result := make(map[string]string)
	if c.SessData != "" {
		result["SESSDATA"] = c.SessData
	}
	if c.BiliJct != "" {
		result["bili_jct"] = c.BiliJct
	}
	if c.DedeUserID != "" {
		result["DedeUserID"] = c.DedeUserID
	}
	if c.Buvid3 != "" {
		result["buvid3"] = c.Buvid3
	}
	if c.Buvid4 != "" {
		result["buvid4"] = c.Buvid4
	}
	if c.AcTimeValue != "" {
		result["ac_time_value"] = c.AcTimeValue
	}
	return result
}

func (c *Credential) ToHTTPCookies() []*http.Cookie {
	cookies := c.Cookies()
	result := make([]*http.Cookie, 0, len(cookies))
	for name, value := range cookies {
		result = append(result, &http.Cookie{Name: name, Value: value})
	}
	return result
}

func (c *Credential) EnsureSessData() error {
	if c == nil || c.SessData == "" {
		return ErrMissingSessData
	}
	return nil
}

func (c *Credential) EnsureBiliJct() error {
	if c == nil || c.BiliJct == "" {
		return ErrMissingBiliJct
	}
	return nil
}

func (c *Credential) MarshalJSON() ([]byte, error) {
	type credential Credential
	return json.Marshal((*credential)(c))
}
