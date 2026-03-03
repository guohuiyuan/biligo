package bilibili

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
)

var (
	endpointSessionList = endpoint{
		name:      "session.list",
		path:      "https://api.vc.bilibili.com/session_svr/v1/session_svr/get_sessions",
		method:    "GET",
		needLogin: true,
		dataField: "data",
	}
	endpointSessionFetch = endpoint{
		name:      "session.fetch",
		path:      "https://api.vc.bilibili.com/svr_sync/v1/svr_sync/fetch_session_msgs",
		method:    "GET",
		withWBI:   true,
		needLogin: true,
		dataField: "data",
	}
	endpointSessionSend = endpoint{
		name:      "session.send",
		path:      "https://api.vc.bilibili.com/web_im/v1/web_im/send_msg",
		method:    "POST",
		withWBI:   true,
		needLogin: true,
		needCSRF:  true,
		dataField: "data",
	}
	endpointMsgUnread = endpoint{
		name:      "msg.unread",
		path:      "/x/msgfeed/unread",
		method:    "GET",
		needLogin: true,
		dataField: "data",
	}
)

type SessionList struct {
	SessionList []Session `json:"session_list"`
}

type Session struct {
	TalkerID    int64 `json:"talker_id"`
	SessionType int   `json:"session_type"`
	UnreadCount int   `json:"unread_count"`
	IsFollow    int   `json:"is_follow"`
	LastMsg     struct {
		Content   string `json:"content"`
		Timestamp int64  `json:"timestamp"`
		MsgType   int    `json:"msg_type"`
	} `json:"last_msg"`
	AccountInfo struct {
		UID  int64  `json:"uid"`
		Name string `json:"name"`
		Pic  string `json:"pic"`
	} `json:"account_info"`
}

type SessionMessages struct {
	Messages []SessionMessage `json:"messages"`
}

type SessionMessage struct {
	MsgKey     int64           `json:"msg_key"`
	MsgSeqno   int64           `json:"msg_seqno"`
	Timestamp  int64           `json:"timestamp"`
	SenderUID  int64           `json:"sender_uid"`
	ReceiverID int64           `json:"receiver_id"`
	MsgType    int             `json:"msg_type"`
	Content    json.RawMessage `json:"content"`
}

type MsgFeedItem struct {
	Mid      int64
	Uname    string
	Avatar   string
	LastMsg  string
	Unfollow int64
}

type ChatHistoryItem struct {
	MsgID      int64
	Content    string
	SenderUID  int64
	SenderName string
	Timestamp  int64
}

func (c *Client) GetMsgFeed(page int32) ([]MsgFeedItem, error) {
	var out SessionList
	err := c.NewRequest(endpointSessionList).
		ParamInt("session_type", 1).
		ParamInt("group_fold", 1).
		ParamInt("unfollow_fold", 0).
		ParamInt("sort_rule", 2).
		ParamInt("build", 0).
		Param("mobi_app", "web").
		Do(context.Background(), &out)
	if err != nil {
		return nil, err
	}

	items := make([]MsgFeedItem, 0, len(out.SessionList))
	for _, session := range out.SessionList {
		items = append(items, MsgFeedItem{
			Mid:      session.TalkerID,
			Uname:    session.AccountInfo.Name,
			Avatar:   session.AccountInfo.Pic,
			LastMsg:  session.LastMsg.Content,
			Unfollow: int64(session.UnreadCount),
		})
	}
	return items, nil
}

func (c *Client) GetChatHistory(userID int64, page int32) ([]ChatHistoryItem, error) {
	beginSeqno := int64(0)
	if page > 1 {
		beginSeqno = int64((page - 1) * 30)
	}

	var out SessionMessages
	err := c.NewRequest(endpointSessionFetch).
		ParamInt("talker_id", userID).
		ParamInt("session_type", 1).
		ParamInt("begin_seqno", beginSeqno).
		Do(context.Background(), &out)
	if err != nil {
		return nil, err
	}

	items := make([]ChatHistoryItem, 0, len(out.Messages))
	for _, msg := range out.Messages {
		items = append(items, ChatHistoryItem{
			MsgID:      msg.MsgKey,
			Content:    decodeMessageContent(msg.Content),
			SenderUID:  msg.SenderUID,
			SenderName: "",
			Timestamp:  msg.Timestamp,
		})
	}
	return items, nil
}

func (c *Client) SendMsg(userID int64, content string) (map[string]any, error) {
	senderUID, err := c.currentUserID(context.Background())
	if err != nil {
		return nil, err
	}

	contentJSON, _ := json.Marshal(map[string]string{"content": content})
	timestamp := time.Now().Unix()

	req := c.NewRequest(endpointSessionSend).
		ParamInt("w_sender_uid", senderUID).
		ParamInt("w_receiver_id", userID).
		FormInt("msg[sender_uid]", senderUID).
		FormInt("msg[receiver_id]", userID).
		FormInt("msg[receiver_type]", 1).
		FormInt("msg[msg_type]", 1).
		FormInt("msg[msg_status]", 0).
		Form("msg[content]", string(contentJSON)).
		Form("msg[dev_id]", "A6716E9A-7CE3-47AF-994B-F0B34178D28D").
		FormInt("msg[new_face_version]", 0).
		FormInt("msg[timestamp]", timestamp).
		FormInt("from_filework", 0).
		FormInt("build", 0).
		Form("mobi_app", "web")

	var out map[string]any
	err = req.Do(context.Background(), &out)
	return out, err
}

func (c *Client) ReadMsg(userID int64) (map[string]any, error) {
	return map[string]any{}, nil
}

func (c *Client) GetUnreadMsg() (int64, error) {
	var out map[string]int64
	err := c.NewRequest(endpointMsgUnread).Do(context.Background(), &out)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, key := range []string{"unfollow_unread", "follow_unread", "unread_dustbin", "biz_msg_follow_unread", "biz_msg_unfollow_unread"} {
		total += out[key]
	}
	return total, nil
}

func (c *Client) currentUserID(ctx context.Context) (int64, error) {
	cred := c.Credential()
	if cred != nil && cred.DedeUserID != "" {
		id, err := strconv.ParseInt(cred.DedeUserID, 10, 64)
		if err == nil && id > 0 {
			return id, nil
		}
	}

	nav, err := c.Login().Nav(ctx)
	if err != nil {
		return 0, err
	}
	return nav.Mid, nil
}

func decodeMessageContent(raw json.RawMessage) string {
	var text struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(raw, &text); err == nil && text.Content != "" {
		return text.Content
	}
	return string(raw)
}
