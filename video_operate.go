package bilibili

import (
	"context"
	"fmt"
)

var (
	endpointVideoRelation = endpoint{
		name:      "video.relation",
		path:      "/x/web-interface/archive/relation",
		method:    "GET",
		needLogin: true,
		dataField: "data",
	}
	endpointVideoLike = endpoint{
		name:      "video.like",
		path:      "/x/web-interface/archive/like",
		method:    "POST",
		needLogin: true,
		needCSRF:  true,
		dataField: "data",
	}
	endpointVideoCoin = endpoint{
		name:      "video.coin",
		path:      "/x/web-interface/coin/add",
		method:    "POST",
		needLogin: true,
		needCSRF:  true,
		dataField: "data",
	}
	endpointVideoTriple = endpoint{
		name:      "video.triple",
		path:      "/x/web-interface/archive/like/triple",
		method:    "POST",
		needLogin: true,
		needCSRF:  true,
		dataField: "data",
	}
	endpointVideoFavorite = endpoint{
		name:      "video.favorite",
		path:      "/x/v3/fav/resource/deal",
		method:    "POST",
		needLogin: true,
		needCSRF:  true,
		dataField: "data",
	}
)

type VideoRelation struct {
	Attention bool `json:"attention"`
	Favorite  bool `json:"favorite"`
	Coin      int  `json:"coin"`
	Like      bool `json:"like"`
	Dislike   bool `json:"dislike"`
}

func (c *Client) LikeVideo(aid int64, like int) (map[string]any, error) {
	var out map[string]any
	err := c.NewRequest(endpointVideoLike).
		FormInt("aid", aid).
		FormInt("like", int64(like)).
		Do(context.Background(), &out)
	return out, err
}

func (c *Client) CoinVideo(aid int64, multiply int32) (map[string]any, error) {
	var out map[string]any
	err := c.NewRequest(endpointVideoCoin).
		FormInt("aid", aid).
		FormInt("multiply", int64(multiply)).
		FormInt("select_like", 0).
		Do(context.Background(), &out)
	return out, err
}

func (c *Client) TripleAction(aid int64) (map[string]any, error) {
	var out map[string]any
	err := c.NewRequest(endpointVideoTriple).
		FormInt("aid", aid).
		Do(context.Background(), &out)
	return out, err
}

func (c *Client) GetVideoRelation(aid int64) (*VideoRelation, error) {
	var out VideoRelation
	err := c.NewRequest(endpointVideoRelation).
		ParamInt("aid", aid).
		Do(context.Background(), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) FavVideo(aid int, mediaID int) (map[string]any, error) {
	var out map[string]any
	err := c.NewRequest(endpointVideoFavorite).
		FormInt("rid", int64(aid)).
		FormInt("type", 2).
		Form("add_media_ids", contextlessComma(mediaID)).
		Form("del_media_ids", "").
		Do(context.Background(), &out)
	return out, err
}

func contextlessComma(v int) string {
	return fmt.Sprintf("%d", v)
}
