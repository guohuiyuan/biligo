package bilibili

import (
	"context"
	"strconv"
)

var (
	endpointCommentGet = endpoint{
		name:      "comment.get",
		path:      "/x/v2/reply",
		method:    "GET",
		dataField: "data",
	}
	endpointCommentSend = endpoint{
		name:      "comment.send",
		path:      "/x/v2/reply/add",
		method:    "POST",
		withWBI:   true,
		needLogin: true,
		needCSRF:  true,
		dataField: "data",
	}
	endpointCommentDelete = endpoint{
		name:      "comment.del",
		path:      "/x/v2/reply/del",
		method:    "POST",
		needLogin: true,
		needCSRF:  true,
		dataField: "data",
	}
	endpointCommentLike = endpoint{
		name:      "comment.like",
		path:      "/x/v2/reply/action",
		method:    "POST",
		needLogin: true,
		needCSRF:  true,
		dataField: "data",
	}
)

type CommentPage struct {
	Count int64 `json:"count"`
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
}

type VideoComment struct {
	Rpid    int64 `json:"rpid"`
	Ctime   int64 `json:"ctime"`
	Like    int   `json:"like"`
	Content struct {
		Message string `json:"message"`
	} `json:"content"`
	Member struct {
		Mid   int64  `json:"mid,string"` // API 返回字符串形式的 mid
		Uname string `json:"uname"`
	} `json:"member"`
}

type VideoCommentList struct {
	Page    CommentPage    `json:"page"`
	Replies []VideoComment `json:"replies"`
}

func (c *Client) GetVideoComment(aid int64, page int32) (*VideoCommentList, error) {
	var out VideoCommentList
	err := c.NewRequest(endpointCommentGet).
		ParamInt("oid", aid).
		ParamInt("pn", int64(page)).
		ParamInt("type", 1).
		ParamInt("sort", 2).
		Do(context.Background(), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SendVideoComment(aid int64, content string) (map[string]any, error) {
	return c.sendComment(context.Background(), aid, 1, content, 0, 0)
}

func (c *Client) ReplyComment(oid string, typ int, content string, root string, parent string) (map[string]any, error) {
	oidInt, err := strconv.ParseInt(oid, 10, 64)
	if err != nil {
		return nil, err
	}

	var rootInt, parentInt int64
	if root != "" {
		rootInt, err = strconv.ParseInt(root, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if parent != "" {
		parentInt, err = strconv.ParseInt(parent, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return c.sendComment(context.Background(), oidInt, typ, content, rootInt, parentInt)
}

func (c *Client) DelComment(oid int64, typ int, rpid int64) (map[string]any, error) {
	var out map[string]any
	err := c.NewRequest(endpointCommentDelete).
		FormInt("oid", oid).
		FormInt("type", int64(typ)).
		FormInt("rpid", rpid).
		Do(context.Background(), &out)
	return out, err
}

func (c *Client) LikeComment(oid int64, typ int, rpid int64, action int) (map[string]any, error) {
	var out map[string]any
	err := c.NewRequest(endpointCommentLike).
		FormInt("oid", oid).
		FormInt("type", int64(typ)).
		FormInt("rpid", rpid).
		FormInt("action", int64(action)).
		Do(context.Background(), &out)
	return out, err
}

func (c *Client) sendComment(ctx context.Context, oid int64, typ int, content string, root int64, parent int64) (map[string]any, error) {
	req := c.NewRequest(endpointCommentSend).
		FormInt("oid", oid).
		FormInt("type", int64(typ)).
		Form("message", content).
		FormInt("plat", 1)

	if root > 0 && parent == 0 {
		req.FormInt("root", root)
		req.FormInt("parent", root)
	} else if root > 0 && parent > 0 {
		req.FormInt("root", root)
		req.FormInt("parent", parent)
	}

	var out map[string]any
	err := req.Do(ctx, &out)
	return out, err
}
