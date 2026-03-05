package bilibili

import "context"

type LoginService struct {
	client *Client
}

func (s *LoginService) QRCodeGenerate(ctx context.Context) (*QRCodeGenerate, error) {
	var out QRCodeGenerate
	err := s.client.NewRequest(endpointLoginQRCodeGenerate).Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *LoginService) QRCodePoll(ctx context.Context, key string) (*QRCodePoll, error) {
	var out QRCodePoll
	err := s.client.NewRequest(endpointLoginQRCodePoll).
		Param("qrcode_key", key).
		Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *LoginService) Nav(ctx context.Context) (*NavInfo, error) {
	var out NavInfo
	err := s.client.NewRequest(endpointNav).Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

type QRCodeGenerate struct {
	URL       string `json:"url"`
	QRCodeKey string `json:"qrcode_key"`
}

type QRCodePoll struct {
	URL          string `json:"url"`
	RefreshToken string `json:"refresh_token"`
	Timestamp    int64  `json:"timestamp"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
}

type NavInfo struct {
	IsLogin bool    `json:"isLogin"`
	Mid     int64   `json:"mid"`
	Uname   string  `json:"uname"`
	Money   float64 `json:"money"`
	WBIImg  struct {
		ImgURL string `json:"img_url"`
		SubURL string `json:"sub_url"`
	} `json:"wbi_img"`
}
