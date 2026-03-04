package bilibili

import "context"

type LiveService struct {
	client *Client
}

func (s *LiveService) RoomInfo(ctx context.Context, roomID int64) (*LiveRoomInfo, error) {
	var out LiveRoomInfo
	err := s.client.NewRequest(endpointLiveRoomInfo).
		ParamInt("room_id", roomID).
		Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *LiveService) StatusByUIDs(ctx context.Context, uids []int64) (map[string]LiveStatus, error) {
	var out map[string]LiveStatus
	req := s.client.NewRequest(endpointLiveStatusByUIDs)
	for _, uid := range uids {
		req.ParamInt("uids[]", uid)
	}
	if err := req.Do(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *LiveService) DanmuInfo(ctx context.Context, roomID int64) (*LiveDanmuInfo, error) {
	var out LiveDanmuInfo
	err := s.client.NewRequest(endpointLiveDanmuInfo).
		ParamInt("id", roomID).
		Param("type", "0").
		Do(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

type LiveStatus struct {
	UID        int64  `json:"uid"`
	UName      string `json:"uname"`
	Title      string `json:"title"`
	RoomID     int64  `json:"room_id"`
	LiveStatus int    `json:"live_status"`
	Online     int64  `json:"online"`
	AreaName   string `json:"area_name"`
	Cover      string `json:"cover_from_user"`
}

type LiveRoomInfo struct {
	UID         int64  `json:"uid"`
	RoomID      int64  `json:"room_id"`
	ShortID     int64  `json:"short_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Attention   int64  `json:"attention"`
	Online      int64  `json:"online"`
	LiveStatus  int    `json:"live_status"`
	UserCover   string `json:"user_cover"`
	AreaName    string `json:"area_name"`
	LiveTime    string `json:"live_time"`
}

type LiveDanmuInfo struct {
	Group      string `json:"group"`
	BusinessID int64  `json:"business_id"` // API 返回整数 ID
	Token      string `json:"token"`
	HostList   []struct {
		Host    string `json:"host"`
		Port    int    `json:"port"`
		WSPort  int    `json:"ws_port"`
		WSSPort int    `json:"wss_port"`
	} `json:"host_list"`
}
