package bilibili

import "context"

type HotTag struct {
	Name string `json:"tag_name"`
	Hot  int64  `json:"heat"`
}

type TagInfo struct {
	TagID int64 `json:"tag_id"`
	Hot   int64 `json:"atten"`
	Count int64 `json:"count"`
}

type RankingV2 struct {
	List []VideoInfo `json:"list"`
}

func (c *Client) GetHotTags(rid int32) ([]HotTag, error) {
	var out []HotTag
	err := c.NewRequest(endpointZoneHotTags).
		ParamInt("rid", int64(rid)).
		Do(context.Background(), &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetTagInfo(tagName string) (*TagInfo, error) {
	var out TagInfo
	err := c.NewRequest(endpointTagInfo).
		Param("tag_name", tagName).
		Do(context.Background(), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetTagVideos(tagName string, page int32) ([]VideoInfo, error) {
	result, err := c.Search().ByType(context.Background(), SearchTypeVideo, tagName, int(page))
	if err != nil {
		return nil, err
	}

	videos := make([]VideoInfo, 0, len(result.Result))
	for _, item := range result.Result {
		video := VideoInfo{}
		if bvid, ok := item["bvid"].(string); ok {
			video.BVID = bvid
		}
		if aid, ok := item["aid"].(float64); ok {
			video.AID = int64(aid)
		}
		if title, ok := item["title"].(string); ok {
			video.Title = title
		}
		if author, ok := item["author"].(string); ok {
			video.Owner.Name = author
		}
		if mid, ok := item["mid"].(float64); ok {
			video.Owner.Mid = int64(mid)
		}
		if pic, ok := item["pic"].(string); ok {
			video.Pic = pic
		}
		if play, ok := item["play"].(float64); ok {
			video.Stat.View = int64(play)
		}
		if like, ok := item["like"].(float64); ok {
			video.Stat.Like = int64(like)
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func (c *Client) GetRanking(rid int32) ([]VideoInfo, error) {
	var out RankingV2
	err := c.NewRequest(endpointZoneRanking).
		ParamInt("rid", int64(rid)).
		Param("type", "all").
		Param("web_location", "333.934").
		Do(context.Background(), &out)
	if err != nil {
		return nil, err
	}
	return out.List, nil
}

func (c *Client) GetFans(page, pageSize int32) ([]struct {
	Mid   int64
	Uname string
	Face  string
}, error) {
	uid, err := c.currentUserID(context.Background())
	if err != nil {
		return nil, err
	}

	resp, err := c.User().Followers(context.Background(), uid, int(page), int(pageSize))
	if err != nil {
		return nil, err
	}

	out := make([]struct {
		Mid   int64
		Uname string
		Face  string
	}, 0, len(resp.List))
	for _, item := range resp.List {
		out = append(out, struct {
			Mid   int64
			Uname string
			Face  string
		}{
			Mid:   item.Mid,
			Uname: item.Uname,
			Face:  item.Face,
		})
	}
	return out, nil
}

func (c *Client) GetUserVideos(mid int64, page, pageSize int) ([]VideoInfo, error) {
	resp, err := c.User().Videos(context.Background(), mid, page, pageSize)
	if err != nil {
		return nil, err
	}
	return resp.List.VList, nil
}

func (c *Client) SearchType(keyword, typ string, page int32) ([]struct {
	Title      string
	Play       int64
	VideoCount int64
}, error) {
	result, err := c.Search().ByType(context.Background(), SearchType(typ), keyword, int(page))
	if err != nil {
		return nil, err
	}

	out := make([]struct {
		Title      string
		Play       int64
		VideoCount int64
	}, 0, len(result.Result))
	for _, item := range result.Result {
		entry := struct {
			Title      string
			Play       int64
			VideoCount int64
		}{}
		if title, ok := item["title"].(string); ok {
			entry.Title = title
		}
		if play, ok := item["play"].(float64); ok {
			entry.Play = int64(play)
		}
		if count, ok := item["video_count"].(float64); ok {
			entry.VideoCount = int64(count)
		}
		out = append(out, entry)
	}
	return out, nil
}
