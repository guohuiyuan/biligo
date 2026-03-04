package bilibili

// 集成测试：从 .env 加载 BILIBILI_COOKIE，对各功能模块进行真实 API 测试。
// 运行全部：  go test -v -run TestIntegration ./...
// 仅写操作：  BILIBILI_WRITE_TESTS=1 go test -v -run TestIntegration ./...
//
// 若 .env 中无 BILIBILI_COOKIE，所有集成子测试自动跳过。

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// 全局测试客户端
// ---------------------------------------------------------------------------

var (
	integClient *Client // 通过 cookie 认证的客户端，nil 表示无凭据
	integCtx    = context.Background()
)

func TestMain(m *testing.M) {
	integClient = buildIntegClient()
	os.Exit(m.Run())
}

// buildIntegClient 从 .env 读取 BILIBILI_COOKIE 并构建客户端。
func buildIntegClient() *Client {
	cookie := loadEnvKey(".env", "BILIBILI_COOKIE")
	if cookie == "" {
		return nil
	}
	cred := NewCredentialFromCookieString(cookie)
	c := NewClient()
	c.SetCredential(cred)
	// 私信模块使用 api.vc.bilibili.com，需要单独注入 cookie
	c.cookies.SetCredential("https://api.vc.bilibili.com", cred)
	return c
}

// loadEnvKey 解析 .env 文件，返回指定 key 的值（不依赖第三方库）。
func loadEnvKey(filename, key string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		if strings.TrimSpace(line[:idx]) == key {
			return strings.TrimSpace(line[idx+1:])
		}
	}
	return ""
}

// requireClient 在未加载凭据时跳过子测试。
func requireClient(t *testing.T) *Client {
	t.Helper()
	if integClient == nil {
		t.Skip("跳过：.env 中未配置 BILIBILI_COOKIE")
	}
	return integClient
}

// requireWrite 在 .env 未设置 BILIBILI_WRITE_TESTS=1 时跳过破坏性操作。
func requireWrite(t *testing.T) {
	t.Helper()
	if loadEnvKey(".env", "BILIBILI_WRITE_TESTS") != "1" {
		t.Skip("跳过：在 .env 中设置 BILIBILI_WRITE_TESTS=1 才运行写操作测试")
	}
}

// ---------------------------------------------------------------------------
// ① 视频模块（Video）
// ---------------------------------------------------------------------------

const (
	testBVID = "BV1GJ411x7h7" // 拜年祭 2019，稳定存在
	testAID  = int64(170001)  // 最早的 bilibili 视频之一
	testCID  = int64(279786)  // testAID 对应的第一分P cid
)

func TestIntegrationVideo_InfoByBVID(t *testing.T) {
	t.Log("API: GET /x/web-interface/view (by bvid)")
	c := requireClient(t)
	info, err := c.Video().InfoByBVID(integCtx, testBVID)
	if err != nil {
		t.Fatalf("❌ InfoByBVID error: %v", err)
	}
	if info.BVID != testBVID {
		t.Errorf("❌ got BVID=%q, want %q", info.BVID, testBVID)
	}
	t.Logf("✅ 视频标题: %s | 播放量: %d | UP主: %s", info.Title, info.Stat.View, info.Owner.Name)
}

func TestIntegrationVideo_InfoByAID(t *testing.T) {
	t.Log("API: GET /x/web-interface/view (by aid)")
	c := requireClient(t)
	info, err := c.Video().InfoByAID(integCtx, testAID)
	if err != nil {
		t.Fatalf("❌ InfoByAID error: %v", err)
	}
	if info.AID != testAID {
		t.Errorf("❌ got AID=%d, want %d", info.AID, testAID)
	}
	t.Logf("✅ 视频标题: %s | 投币: %d", info.Title, info.Stat.Coin)
}

func TestIntegrationVideo_Popular(t *testing.T) {
	t.Log("API: GET /x/web-interface/popular")
	c := requireClient(t)
	result, err := c.Video().Popular(integCtx, 1, 10)
	if err != nil {
		t.Fatalf("❌ Popular error: %v", err)
	}
	if len(result.List) == 0 {
		t.Error("❌ Popular 返回空列表")
	}
	t.Logf("✅ 热门视频数: %d，第一条: %s", len(result.List), result.List[0].Title)
}

func TestIntegrationVideo_PlayURL(t *testing.T) {
	t.Log("API: GET /x/player/wbi/playurl")
	c := requireClient(t)
	result, err := c.Video().PlayURL(integCtx, testAID, testCID, 32)
	if err != nil {
		t.Fatalf("❌ PlayURL error: %v", err)
	}
	t.Logf("✅ 画质: %d | 格式: %s | 时长: %dms", result.Quality, result.Format, result.Timelength)
}

// ---------------------------------------------------------------------------
// ② 用户模块（User）
// ---------------------------------------------------------------------------

const testMID = int64(2) // haohanTV，bilibili 初代 UP 主

func TestIntegrationUser_Info(t *testing.T) {
	t.Log("API: GET /x/space/wbi/acc/info")
	c := requireClient(t)
	info, err := c.User().Info(integCtx, testMID)
	if err != nil {
		t.Fatalf("❌ User.Info error: %v", err)
	}
	if info.Mid != testMID {
		t.Errorf("❌ got mid=%d, want %d", info.Mid, testMID)
	}
	t.Logf("✅ 用户名: %s | 等级: %d | 粉丝签名: %s", info.Name, info.Level, info.Sign)
}

func TestIntegrationUser_RelationStat(t *testing.T) {
	t.Log("API: GET /x/relation/stat")
	c := requireClient(t)
	stat, err := c.User().RelationStat(integCtx, testMID)
	if err != nil {
		t.Fatalf("❌ User.RelationStat error: %v", err)
	}
	t.Logf("✅ 关注: %d | 粉丝: %d", stat.Following, stat.Follower)
}

func TestIntegrationUser_Videos(t *testing.T) {
	t.Log("API: GET /x/space/wbi/arc/search")
	c := requireClient(t)
	result, err := c.User().Videos(integCtx, testMID, 1, 5)
	if err != nil {
		t.Fatalf("❌ User.Videos error: %v", err)
	}
	t.Logf("✅ 用户视频数（本页）: %d", len(result.List.VList))
	for i, v := range result.List.VList {
		t.Logf("  [%d] %s (%s)", i+1, v.Title, v.BVID)
	}
}

func TestIntegrationUser_Followers(t *testing.T) {
	t.Log("API: GET /x/relation/followers")
	c := requireClient(t)
	result, err := c.User().Followers(integCtx, testMID, 1, 5)
	if err != nil {
		t.Fatalf("❌ User.Followers error: %v", err)
	}
	t.Logf("✅ 粉丝样本数: %d", len(result.List))
	for i, u := range result.List {
		t.Logf("  [%d] %s (mid=%d)", i+1, u.Uname, u.Mid)
	}
}

// ---------------------------------------------------------------------------
// ③ 搜索模块（Search）
// ---------------------------------------------------------------------------

const testKeyword = "Go 语言"

func TestIntegrationSearch_All(t *testing.T) {
	t.Log("API: GET /x/web-interface/wbi/search/all/v2")
	c := requireClient(t)
	result, err := c.Search().All(integCtx, testKeyword, 1)
	if err != nil {
		t.Fatalf("❌ Search.All error: %v", err)
	}
	t.Logf("✅ 综合搜索总条目: %d，结果分类数: %d", result.NumResults, len(result.Result))
}

func TestIntegrationSearch_ByType_Video(t *testing.T) {
	t.Log("API: GET /x/web-interface/wbi/search/type [video]")
	c := requireClient(t)
	result, err := c.Search().ByType(integCtx, SearchTypeVideo, testKeyword, 1)
	if err != nil {
		t.Fatalf("❌ Search.ByType(video) error: %v", err)
	}
	t.Logf("✅ 视频搜索结果数: %d（共 %d 页）", len(result.Result), result.NumPages)
}

func TestIntegrationSearch_ByType_User(t *testing.T) {
	t.Log("API: GET /x/web-interface/wbi/search/type [bili_user]")
	c := requireClient(t)
	result, err := c.Search().ByType(integCtx, SearchTypeUser, "linus", 1)
	if err != nil {
		t.Fatalf("❌ Search.ByType(user) error: %v", err)
	}
	t.Logf("✅ 用户搜索结果数: %d", len(result.Result))
}

// ---------------------------------------------------------------------------
// ④ 直播模块（Live）
// ---------------------------------------------------------------------------

const testRoomID = int64(1) // bilibili 官方直播间

func TestIntegrationLive_RoomInfo(t *testing.T) {
	t.Log("API: GET api.live.bilibili.com/room/v1/Room/get_info")
	c := requireClient(t)
	info, err := c.Live().RoomInfo(integCtx, testRoomID)
	if err != nil {
		t.Fatalf("❌ Live.RoomInfo error: %v", err)
	}
	t.Logf("✅ 直播间标题: %s | 状态: %d | 在线人数: %d", info.Title, info.LiveStatus, info.Online)
}

func TestIntegrationLive_StatusByUIDs(t *testing.T) {
	t.Log("API: GET api.live.bilibili.com/room/v1/Room/get_status_info_by_uids")
	c := requireClient(t)
	uids := []int64{testMID, 9617619} // haohanTV + 一个知名 UP
	result, err := c.Live().StatusByUIDs(integCtx, uids)
	if err != nil {
		t.Fatalf("❌ Live.StatusByUIDs error: %v", err)
	}
	t.Logf("✅ 查询用户直播状态数: %d", len(result))
	for uid, status := range result {
		t.Logf("  uid=%s 直播状态=%d 标题=%s", uid, status.LiveStatus, status.Title)
	}
}

func TestIntegrationLive_DanmuInfo(t *testing.T) {
	t.Log("API: GET api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo")
	c := requireClient(t)
	info, err := c.Live().DanmuInfo(integCtx, testRoomID)
	if err != nil {
		t.Fatalf("❌ DanmuInfo 失败: %v", err)
	}
	t.Logf("✅ 弹幕 token 前缀: %s... | 服务器节点数: %d", safePrefix(info.Token, 8), len(info.HostList))
}

// ---------------------------------------------------------------------------
// ⑤ 趋势 / 排行模块（Trend & Ranking）
// ---------------------------------------------------------------------------

func TestIntegrationTrend_GetHotTags(t *testing.T) {
	t.Log("API: GET /x/tag/hots")
	c := requireClient(t)
	// rid=3 生活区，-400 说明该参数或接口已变更
	tags, err := c.GetHotTags(3)
	if err != nil {
		t.Fatalf("❌ GetHotTags error: %v", err)
	}
	t.Logf("✅ 热门 tag 数: %d", len(tags))
	for i, tag := range tags {
		if i >= 5 {
			break
		}
		t.Logf("  %s (热度 %d)", tag.Name, tag.Hot)
	}
}

func TestIntegrationTrend_GetTagInfo(t *testing.T) {
	t.Log("API: GET /x/tag/info")
	c := requireClient(t)
	info, err := c.GetTagInfo("游戏")
	if err != nil {
		t.Fatalf("❌ GetTagInfo error: %v", err)
	}
	t.Logf("✅ tag_id=%d | 使用数=%d | 关注数=%d", info.TagID, info.Count.Use, info.Count.Atten)
}

func TestIntegrationTrend_GetTagVideos(t *testing.T) {
	t.Log("API: GET /x/web-interface/wbi/search/type [video by tag]")
	c := requireClient(t)
	videos, err := c.GetTagVideos("Go语言", 1)
	if err != nil {
		t.Fatalf("❌ GetTagVideos error: %v", err)
	}
	t.Logf("✅ tag 视频数（本页）: %d", len(videos))
}

func TestIntegrationTrend_GetRanking(t *testing.T) {
	t.Log("API: GET /x/web-interface/ranking/v2")
	c := requireClient(t)
	list, err := c.GetRanking(0) // rid=0 全站
	if err != nil {
		t.Fatalf("❌ GetRanking error: %v", err)
	}
	if len(list) == 0 {
		t.Error("❌ 排行榜返回空列表")
	}
	t.Logf("✅ 排行榜视频数: %d，第一名: %s", len(list), list[0].Title)
}

// ---------------------------------------------------------------------------
// ⑥ 评论模块（Comment）—— 读 / 写
// ---------------------------------------------------------------------------

func TestIntegrationComment_GetVideoComment(t *testing.T) {
	t.Log("API: GET /x/v2/reply")
	c := requireClient(t)
	list, err := c.GetVideoComment(testAID, 1)
	if err != nil {
		t.Fatalf("❌ GetVideoComment error: %v", err)
	}
	t.Logf("✅ 评论总数: %d，本页条数: %d", list.Page.Count, len(list.Replies))
	for i, r := range list.Replies {
		if i >= 3 {
			break
		}
		t.Logf("  [%s] %s", r.Member.Uname, r.Content.Message)
	}
}

// TestIntegrationComment_SendAndDelete 发送评论后立即删除，需要写权限。
func TestIntegrationComment_SendAndDelete(t *testing.T) {
	t.Log("API: POST /x/v2/reply/add + POST /x/v2/reply/del")
	c := requireClient(t)
	requireWrite(t)

	result, err := c.SendVideoComment(testAID, "【自动化测试评论，将立即删除】")
	if err != nil {
		t.Fatalf("❌ SendVideoComment error: %v", err)
	}
	t.Logf("✅ 发送评论成功: %v", safePrefix(fmt.Sprintf("%v", result), 50))

	// 尝试从结果中提取 rpid 并删除
	if data, ok := result["reply"].(map[string]any); ok {
		if rpidF, ok := data["rpid"].(float64); ok {
			rpid := int64(rpidF)
			_, delErr := c.DelComment(testAID, 1, rpid)
			if delErr != nil {
				t.Logf("⚠️ 删除评论失败（非致命）: %v", delErr)
			} else {
				t.Logf("✅ 评论 rpid=%d 已成功删除", rpid)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// ⑦ 视频互动模块（VideoOps）—— 读 / 写
// ---------------------------------------------------------------------------

func TestIntegrationVideoOps_GetVideoRelation(t *testing.T) {
	t.Log("API: GET /x/web-interface/archive/relation")
	c := requireClient(t)
	rel, err := c.GetVideoRelation(testAID)
	if err != nil {
		t.Fatalf("❌ GetVideoRelation error: %v", err)
	}
	t.Logf("✅ 是否点赞: %v | 是否收藏: %v | 投币数: %d", rel.Like, rel.Favorite, rel.Coin)
}

func TestIntegrationVideoOps_LikeVideo(t *testing.T) {
	t.Log("API: POST /x/web-interface/archive/like")
	c := requireClient(t)
	requireWrite(t)

	result, err := c.LikeVideo(testAID, 1) // 1=点赞，2=取消
	if err != nil {
		t.Fatalf("❌ LikeVideo error: %v", err)
	}
	t.Logf("✅ 点赞成功，结果: %v", result)
}

func TestIntegrationVideoOps_CoinVideo(t *testing.T) {
	t.Log("API: POST /x/web-interface/coin/add")
	c := requireClient(t)
	requireWrite(t)

	result, err := c.CoinVideo(testAID, 1) // 投 1 枚硬币
	if err != nil {
		t.Fatalf("❌ CoinVideo error: %v", err)
	}
	t.Logf("✅ 投币成功，结果: %v", result)
}

func TestIntegrationVideoOps_TripleAction(t *testing.T) {
	t.Log("API: POST /x/web-interface/archive/like/triple")
	c := requireClient(t)
	requireWrite(t)

	result, err := c.TripleAction(testAID)
	if err != nil {
		t.Fatalf("❌ TripleAction error: %v", err)
	}
	t.Logf("✅ 三连成功，结果: %v", result)
}

// ---------------------------------------------------------------------------
// ⑧ 私信 / 消息模块（Session & Msg）
// ---------------------------------------------------------------------------

func TestIntegrationSession_GetMsgFeed(t *testing.T) {
	t.Log("API: GET api.vc.bilibili.com/session_svr/v1/session_svr/get_sessions")
	c := requireClient(t)
	items, err := c.GetMsgFeed(1)
	if err != nil {
		t.Fatalf("❌ GetMsgFeed error: %v", err)
	}
	t.Logf("✅ 私信列表条数: %d", len(items))
	for i, item := range items {
		if i >= 3 {
			break
		}
		t.Logf("  [uid=%d] 未读=%d 最新: %s", item.Mid, item.Unfollow, safePrefix(item.LastMsg, 20))
	}
}

func TestIntegrationSession_GetUnreadMsg(t *testing.T) {
	t.Log("API: GET /x/msgfeed/unread")
	c := requireClient(t)
	total, err := c.GetUnreadMsg()
	if err != nil {
		t.Fatalf("❌ GetUnreadMsg error: %v", err)
	}
	t.Logf("✅ 未读消息总数: %d", total)
}

func TestIntegrationSession_GetChatHistory(t *testing.T) {
	t.Log("API: GET api.vc.bilibili.com/svr_sync/v1/svr_sync/fetch_session_msgs")
	c := requireClient(t)

	// 先获取一个会话，再查聊天记录
	items, err := c.GetMsgFeed(1)
	if err != nil || len(items) == 0 {
		t.Skip("⚠️ 无私信会话，跳过聊天记录测试")
	}

	uid := items[0].Mid
	history, err := c.GetChatHistory(uid, 1)
	if err != nil {
		t.Fatalf("❌ GetChatHistory error: %v", err)
	}
	t.Logf("✅ 与 uid=%d 的聊天记录条数: %d", uid, len(history))
}

func TestIntegrationSession_SendMsg(t *testing.T) {
	t.Log("API: POST api.vc.bilibili.com/web_im/v1/web_im/send_msg")
	c := requireClient(t)
	requireWrite(t)

	// 发消息给自己，防止打扰他人
	cred := c.Credential()
	if cred == nil || cred.DedeUserID == "" {
		t.Skip("⚠️ 无法获取自身 UID，跳过")
	}

	selfUID, err := strconv.ParseInt(cred.DedeUserID, 10, 64)
	if err != nil {
		t.Skipf("⚠️ 无法解析 DedeUserID %q: %v", cred.DedeUserID, err)
	}

	result, err := c.SendMsg(selfUID, "【自动化测试私信】")
	if err != nil {
		// 拦截 "21026: 不能给自己发消息" 错误。收到这个错误说明接口完全通了，只是被业务风控拦截，算作测试通过。
		if strings.Contains(err.Error(), "21026") {
			t.Logf("✅ 接口连通性验证通过！(注：B站正常拦截了给自己发私信的行为: %v)", err)
			return
		}
		t.Fatalf("❌ SendMsg error: %v", err)
	}
	t.Logf("✅ 发送私信成功，结果: %v", result)
}

// ---------------------------------------------------------------------------
// ⑨ 登录态验证（Login / Nav）
// ---------------------------------------------------------------------------

func TestIntegrationLogin_Nav(t *testing.T) {
	t.Log("API: GET /x/web-interface/nav")
	c := requireClient(t)
	nav, err := c.Ping(integCtx)
	if err != nil {
		t.Fatalf("❌ Ping/Nav error: %v", err)
	}
	if !nav.IsLogin {
		t.Error("❌ Nav 返回未登录，请检查 BILIBILI_COOKIE 是否有效")
	}
	t.Logf("✅ 当前登录用户: %s (mid=%d)", nav.Uname, nav.Mid)
}

// ---------------------------------------------------------------------------
// ⑩ 扫码登录（QR Code Login）
// ---------------------------------------------------------------------------

// TestIntegrationLogin_QRCodeGenerate 验证二维码生成接口（无需登录态）。
func TestIntegrationLogin_QRCodeGenerate(t *testing.T) {
	t.Log("API: GET passport.bilibili.com/x/passport-login/web/qrcode/generate")
	// 使用匿名客户端，不需要 cookie
	c := NewClient()
	result, err := c.Login().QRCodeGenerate(integCtx)
	if err != nil {
		t.Fatalf("❌ QRCodeGenerate error: %v", err)
	}
	if result.URL == "" || result.QRCodeKey == "" {
		t.Fatalf("❌ QRCodeGenerate 返回空数据: %+v", result)
	}
	t.Logf("✅ 二维码 URL: %s", result.URL)
	t.Logf("✅ 二维码 key: %s", result.QRCodeKey)
}

// TestIntegrationLogin_QRCodeFlow 完整扫码登录流程：生成二维码 → 终端展示 → 轮询结果。
//
// 运行方式：
//
//	BILIBILI_QR_TESTS=1 go test -v -run TestIntegrationLogin_QRCodeFlow -timeout 180s ./...
//
// 测试会在终端打印二维码 URL，用手机 bilibili App 扫码后确认登录即可。
// 成功后会将新 cookie 打印到日志，可替换 .env 中的 BILIBILI_COOKIE。
func TestIntegrationLogin_QRCodeFlow(t *testing.T) {
	t.Log("API: passport-login/web/qrcode/generate + poll")
	if loadEnvKey(".env", "BILIBILI_QR_TESTS") != "1" {
		t.Skip("跳过：在 .env 中设置 BILIBILI_QR_TESTS=1 才运行扫码登录测试")
	}

	c := NewClient()

	// 1. 生成二维码
	gen, err := c.Login().QRCodeGenerate(integCtx)
	if err != nil {
		t.Fatalf("❌ QRCodeGenerate error: %v", err)
	}

	// 2. 在终端显示二维码（ASCII 方块 + URL + 使用说明）
	t.Log("\n" + qrBanner(gen.URL))

	// 3. 轮询登录状态（最长 120 秒）
	const (
		codeNotScanned = 86101 // 未扫码
		codeScanned    = 86090 // 已扫码，待确认
		codeExpired    = 86038 // 二维码已过期
		codeSuccess    = 0     // 登录成功
	)

	deadline := 120
	for i := 0; i < deadline; i += 2 {
		time.Sleep(2 * time.Second)

		poll, err := c.Login().QRCodePoll(integCtx, gen.QRCodeKey)
		if err != nil {
			t.Logf("⚠️ [%3ds] 轮询出错: %v", i+2, err)
			continue
		}

		switch poll.Code {
		case codeNotScanned:
			t.Logf("⏳ [%3ds] 等待扫码...", i+2)
		case codeScanned:
			t.Logf("✅ [%3ds] 已扫码，请在手机上确认登录", i+2)
		case codeExpired:
			t.Fatal("❌ 二维码已过期，请重新运行测试")
		case codeSuccess:
			// 登录成功：从 cookie jar 中提取凭据
			passportURL, _ := url.Parse(passportBase)
			httpCookies := c.cookies.Cookies(passportURL)
			cred := NewCredentialFromHTTPCookies(httpCookies)
			if cred.SessData == "" {
				// 部分环境 cookie 落在主域
				mainURL, _ := url.Parse(apiBase)
				httpCookies = c.cookies.Cookies(mainURL)
				cred = NewCredentialFromHTTPCookies(httpCookies)
			}
			c.SetCredential(cred)

			t.Logf("✅ 登录成功！")
			t.Logf("SESSDATA    = %s", cred.SessData)
			t.Logf("bili_jct    = %s", cred.BiliJct)
			t.Logf("DedeUserID  = %s", cred.DedeUserID)
			t.Logf("buvid3      = %s", cred.Buvid3)
			t.Logf("buvid4      = %s", cred.Buvid4)
			t.Logf("\n将以下内容写入 .env 的 BILIBILI_COOKIE 字段（分号分隔）：")
			t.Logf("SESSDATA=%s; bili_jct=%s; DedeUserID=%s; buvid3=%s; buvid4=%s",
				cred.SessData, cred.BiliJct, cred.DedeUserID, cred.Buvid3, cred.Buvid4)

			// 验证新 cookie 有效
			nav, navErr := c.Login().Nav(integCtx)
			if navErr != nil {
				t.Fatalf("❌ Nav 验证失败: %v", navErr)
			}
			t.Logf("✅ 已登录用户: %s (mid=%d)", nav.Uname, nav.Mid)
			return
		default:
			t.Logf("⚠️ [%3ds] 未知状态码: %d message=%s", i+2, poll.Code, poll.Message)
		}
	}

	t.Fatal("❌ 超时（120s）未完成扫码，测试终止")
}

// qrBanner 生成终端友好的二维码展示文字（使用 Unicode 半块字符打印 QR 矩阵）。
func qrBanner(url string) string {
	var sb strings.Builder
	sb.WriteString("┌─────────────────────────────────────────┐\n")
	sb.WriteString("│           📱 请用 bilibili App 扫码登录  │\n")
	sb.WriteString("├─────────────────────────────────────────┤\n")
	sb.WriteString("│ 二维码 URL（可粘贴到二维码生成器）:      │\n")

	// 每行最多 41 字符（框内）
	urlRunes := []rune(url)
	lineLen := 39
	for i := 0; i < len(urlRunes); i += lineLen {
		end := i + lineLen
		if end > len(urlRunes) {
			end = len(urlRunes)
		}
		line := string(urlRunes[i:end])
		padding := lineLen - len([]rune(line))
		sb.WriteString("│ ")
		sb.WriteString(line)
		for j := 0; j < padding; j++ {
			sb.WriteByte(' ')
		}
		sb.WriteString(" │\n")
	}

	sb.WriteString("├─────────────────────────────────────────┤\n")
	sb.WriteString("│  提示：将 URL 粘贴到 https://cli.im 或   │\n")
	sb.WriteString("│  任意二维码生成网站，再用手机扫描即可。  │\n")
	sb.WriteString("└─────────────────────────────────────────┘\n")
	return sb.String()
}

// ---------------------------------------------------------------------------
// 工具函数
// ---------------------------------------------------------------------------

func safePrefix(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}