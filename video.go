package bilibili

import (
	"context"
	"strconv"
	"strings"
)

type VideoService struct {
	client *Client
}

func (s *VideoService) InfoByBVID(ctx context.Context, bvid string) (*VideoInfo, error) {
	var out VideoInfo
	err := s.client.NewRequest(endpointVideoInfo).
		Param("bvid", strings.TrimSpace(bvid)).
		Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *VideoService) InfoByAID(ctx context.Context, aid int64) (*VideoInfo, error) {
	var out VideoInfo
	err := s.client.NewRequest(endpointVideoInfo).
		ParamInt("aid", aid).
		Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *VideoService) Popular(ctx context.Context, page, pageSize int) (*PopularVideos, error) {
	var out PopularVideos
	req := s.client.NewRequest(endpointVideoPopular)
	if page > 0 {
		req.ParamInt("pn", int64(page))
	}
	if pageSize > 0 {
		req.ParamInt("ps", int64(pageSize))
	}
	if err := req.Do(ctx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *VideoService) PlayURL(ctx context.Context, aid, cid int64, quality int) (*VideoPlayURL, error) {
	var out VideoPlayURL
	req := s.client.NewRequest(endpointVideoPlayURL).
		Param("avid", strconv.FormatInt(aid, 10)).
		ParamInt("cid", cid).
		ParamInt("qn", int64(quality)).
		Param("fnval", "4048").
		ParamBool("fourk", true)
	if err := req.Do(ctx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type VideoInfo struct {
	BVID    string `json:"bvid"`
	AID     int64  `json:"aid"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Pic     string `json:"pic"`
	PubDate int64  `json:"pubdate"`
	Owner   struct {
		Mid  int64  `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"owner"`
	Stat struct {
		View     int64 `json:"view"`
		Danmaku  int64 `json:"danmaku"`
		Reply    int64 `json:"reply"`
		Favorite int64 `json:"favorite"`
		Coin     int64 `json:"coin"`
		Share    int64 `json:"share"`
		Like     int64 `json:"like"`
	} `json:"stat"`
	Pages []struct {
		CID      int64  `json:"cid"`
		Page     int    `json:"page"`
		Part     string `json:"part"`
		Duration int    `json:"duration"`
	} `json:"pages"`
}

type PopularVideos struct {
	List []VideoInfo `json:"list"`
}

type VideoPlayURL struct {
	Quality           int      `json:"quality"`
	Format            string   `json:"format"`
	Timelength        int64    `json:"timelength"`
	AcceptDescription []string `json:"accept_description"`
	AcceptQuality     []int    `json:"accept_quality"`
	DURL              []struct {
		Order     int      `json:"order"`
		Length    int64    `json:"length"`
		Size      int64    `json:"size"`
		URL       string   `json:"url"`
		BackupURL []string `json:"backup_url"`
	} `json:"durl"`
	Dash *struct {
		Duration int `json:"duration"`
		Video    []struct {
			ID        int      `json:"id"`
			BaseURL   string   `json:"base_url"`
			BackupURL []string `json:"backup_url"`
			Codecs    string   `json:"codecs"`
			Width     int      `json:"width"`
			Height    int      `json:"height"`
			MimeType  string   `json:"mime_type"`
		} `json:"video"`
		Audio []struct {
			ID        int      `json:"id"`
			BaseURL   string   `json:"base_url"`
			BackupURL []string `json:"backup_url"`
			Codecs    string   `json:"codecs"`
			MimeType  string   `json:"mime_type"`
		} `json:"audio"`
	} `json:"dash"`
}
