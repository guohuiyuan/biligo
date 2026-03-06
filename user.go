package bilibili

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

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

func (s *UserService) Videos(ctx context.Context, mid int64, page, pageSize int) (*UserVideoSearchResult, error) {
	var out UserVideoSearchResult
	req := s.client.NewRequest(endpointUserVideos).
		ParamInt("mid", mid).
		ParamInt("pn", int64(page)).
		ParamInt("ps", int64(pageSize)).
		ParamInt("tid", 0)
	err := req.Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *UserService) Followers(ctx context.Context, mid int64, page, pageSize int) (*UserFollowersResult, error) {
	var out UserFollowersResult
	req := s.client.NewRequest(endpointUserFollowers).
		ParamInt("vmid", mid).
		ParamInt("pn", int64(page)).
		ParamInt("ps", int64(pageSize)).
		Param("order", "desc")
	err := req.Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *UserService) Fans(ctx context.Context, mid int64, page, pageSize int) (*UserFollowersResult, error) {
	var out UserFollowersResult
	req := s.client.NewRequest(endpointUserFans).
		ParamInt("vmid", mid).
		ParamInt("pn", int64(page)).
		ParamInt("ps", int64(pageSize))
	err := req.Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *UserService) FollowersUnreadCount(ctx context.Context) (*FollowerUnreadCount, error) {
	var out FollowerUnreadCount
	err := s.client.NewRequest(endpointUserFollowersUnreadCount).Do(ctx, &out)
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

type UserVideoSearchResult struct {
	List struct {
		VList []UserVideoItem `json:"vlist"`
	} `json:"list"`
}

type FlexInt64 int64

func (f *FlexInt64) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		*f = 0
		return nil
	}
	if len(trimmed) >= 2 && trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"' {
		s := strings.Trim(trimmed, "\"")
		if s == "" {
			*f = 0
			return nil
		}
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			*f = FlexInt64(n)
			return nil
		}
		if fl, err := strconv.ParseFloat(s, 64); err == nil {
			*f = FlexInt64(int64(fl))
			return nil
		}
		*f = 0
		return nil
	}
	var num json.Number
	if err := json.Unmarshal(data, &num); err != nil {
		*f = 0
		return nil
	}
	if n, err := num.Int64(); err == nil {
		*f = FlexInt64(n)
		return nil
	}
	if fl, err := num.Float64(); err == nil {
		*f = FlexInt64(int64(fl))
		return nil
	}
	*f = 0
	return nil
}

type UserVideoItem struct {
	BVID    string    `json:"bvid"`
	Title   string    `json:"title"`
	Pic     string    `json:"pic"`
	Length  string    `json:"length"`
	Play    FlexInt64 `json:"play"`
	Comment FlexInt64 `json:"comment"`
	Created int64     `json:"created"`
	PubDate int64     `json:"pubdate"`
}

type UserFollowersResult struct {
	List []struct {
		Mid   int64  `json:"mid"`
		Uname string `json:"uname"`
		Face  string `json:"face"`
		MTime int64  `json:"mtime"`
	} `json:"list"`
}

type FollowerUnreadCount struct {
	Count int64 `json:"count"`
	Time  int64 `json:"time"`
}
