package bilibili

import "net/http"

type endpoint struct {
	name       string
	baseURL    string
	path       string
	method     string
	withWBI    bool
	needLogin  bool
	needCSRF   bool
	dataField  string
	ignoreCode bool
}

func (e endpoint) url(c *Client) string {
	if len(e.path) >= 4 && e.path[:4] == "http" {
		return e.path
	}
	base := c.config.BaseURL
	switch e.baseURL {
	case passportBase:
		base = c.config.PassportBaseURL
	case liveBase:
		base = c.config.LiveBaseURL
	}
	return base + e.path
}

var (
	endpointNav = endpoint{
		name:      "nav",
		path:      "/x/web-interface/nav",
		method:    http.MethodGet,
		dataField: "data",
	}
	endpointVideoInfo = endpoint{
		name:      "video.info",
		path:      "/x/web-interface/view",
		method:    http.MethodGet,
		dataField: "data",
	}
	endpointVideoPlayURL = endpoint{
		name:      "video.playurl",
		path:      "/x/player/wbi/playurl",
		method:    http.MethodGet,
		withWBI:   true,
		dataField: "data",
	}
	endpointVideoPopular = endpoint{
		name:      "video.popular",
		path:      "/x/web-interface/popular",
		method:    http.MethodGet,
		dataField: "data",
	}
	endpointUserInfo = endpoint{
		name:      "user.info",
		path:      "/x/space/wbi/acc/info",
		method:    http.MethodGet,
		withWBI:   true,
		dataField: "data",
	}
	endpointUserStat = endpoint{
		name:      "user.stat",
		path:      "/x/relation/stat",
		method:    http.MethodGet,
		withWBI:   true,
		dataField: "data",
	}
	endpointUserVideos = endpoint{
		name:      "user.videos",
		path:      "/x/space/wbi/arc/search",
		method:    http.MethodGet,
		withWBI:   true,
		dataField: "data",
	}
	endpointUserFollowers = endpoint{
		name:      "user.followers",
		path:      "/x/relation/followers",
		method:    http.MethodGet,
		needLogin: true,
		dataField: "data",
	}
	endpointUserFans = endpoint{
		name:      "user.fans",
		path:      "/x/relation/fans",
		method:    http.MethodGet,
		needLogin: true,
		dataField: "data",
	}
	endpointUserFollowersUnreadCount = endpoint{
		name:      "user.followers.unread_count",
		path:      "/x/relation/followers/unread/count",
		method:    http.MethodGet,
		needLogin: true,
		dataField: "data",
	}
	endpointSearchAll = endpoint{
		name:      "search.all",
		path:      "/x/web-interface/wbi/search/all/v2",
		method:    http.MethodGet,
		withWBI:   true,
		dataField: "data",
	}
	endpointSearchType = endpoint{
		name:      "search.type",
		path:      "/x/web-interface/wbi/search/type",
		method:    http.MethodGet,
		withWBI:   true,
		dataField: "data",
	}
	endpointSearchSuggest = endpoint{
		name:      "search.suggest",
		path:      "/x/web-interface/search/suggest",
		method:    http.MethodGet,
		dataField: "result",
	}
	endpointZoneHotTags = endpoint{
		name:      "zone.hot.tags",
		path:      "/x/tag/hots",
		method:    http.MethodGet,
		dataField: "data",
	}
	endpointTagInfo = endpoint{
		name:      "tag.info",
		path:      "/x/tag/info",
		method:    http.MethodGet,
		dataField: "data",
	}
	endpointZoneRanking = endpoint{
		name:      "zone.ranking",
		path:      "/x/web-interface/ranking/v2",
		method:    http.MethodGet,
		withWBI:   true,
		dataField: "data",
	}
	endpointPGCRanking = endpoint{
		name:      "pgc.ranking",
		path:      "/pgc/web/rank/list",
		method:    http.MethodGet,
		dataField: "result",
	}
	endpointLiveRoomInfo = endpoint{
		name:      "live.room.info",
		baseURL:   liveBase,
		path:      "/room/v1/Room/get_info",
		method:    http.MethodGet,
		dataField: "data",
	}
	endpointLiveStatusByUIDs = endpoint{
		name:      "live.status.by_uids",
		baseURL:   liveBase,
		path:      "/room/v1/Room/get_status_info_by_uids",
		method:    http.MethodGet,
		dataField: "data",
	}
	endpointLiveDanmuInfo = endpoint{
		name:      "live.danmu.info",
		baseURL:   liveBase,
		path:      "/xlive/web-room/v1/index/getDanmuInfo",
		method:    http.MethodGet,
		withWBI:   true,
		dataField: "data",
	}
	endpointLoginQRCodeGenerate = endpoint{
		name:      "login.qrcode.generate",
		baseURL:   passportBase,
		path:      "/x/passport-login/web/qrcode/generate",
		method:    http.MethodGet,
		dataField: "data",
	}
	endpointLoginQRCodePoll = endpoint{
		name:      "login.qrcode.poll",
		baseURL:   passportBase,
		path:      "/x/passport-login/web/qrcode/poll",
		method:    http.MethodGet,
		dataField: "data",
	}
)
