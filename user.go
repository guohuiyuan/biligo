package bilibili

import "context"

type UserService struct {
	client *Client
}

func (s *UserService) Info(ctx context.Context, mid int64) (*UserInfo, error) {
	var out UserInfo
	err := s.client.NewRequest(endpointUserInfo).
		ParamInt("mid", mid).
		Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *UserService) RelationStat(ctx context.Context, mid int64) (*UserRelationStat, error) {
	var out UserRelationStat
	err := s.client.NewRequest(endpointUserStat).
		ParamInt("vmid", mid).
		Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

type UserInfo struct {
	Mid      int64  `json:"mid"`
	Name     string `json:"name"`
	Face     string `json:"face"`
	Sign     string `json:"sign"`
	Level    int    `json:"level"`
	Sex      string `json:"sex"`
	Official struct {
		Role  int    `json:"role"`
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Type  int    `json:"type"`
	} `json:"official"`
	VIP struct {
		Type   int `json:"type"`
		Status int `json:"status"`
	} `json:"vip"`
}

type UserRelationStat struct {
	Mid       int64 `json:"mid"`
	Following int64 `json:"following"`
	Whisper   int64 `json:"whisper"`
	Black     int64 `json:"black"`
	Follower  int64 `json:"follower"`
}
