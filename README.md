# biligo

Go 1.18+ Bilibili 非官方 SDK，原生类型化客户端，无反射、无代码生成。

```
go get github.com/guohuiyuan/biligo
```

---

## 目录

- [快速开始](#快速开始)
- [认证](#认证)
  - [Cookie 登录](#cookie-登录)
  - [扫码登录](#扫码登录)
- [功能模块](#功能模块)
  - [视频 Video](#视频-video)
  - [用户 User](#用户-user)
  - [搜索 Search](#搜索-search)
  - [直播 Live](#直播-live)
  - [评论 Comment](#评论-comment)
  - [视频互动 VideoOps](#视频互动-videoops)
  - [私信消息 Session](#私信消息-session)
  - [趋势排行 Trend](#趋势排行-trend)
- [配置项](#配置项)
- [测试](#测试)

---

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"

    bilibili "github.com/guohuiyuan/biligo"
)

func main() {
    client := bilibili.NewClient()

    info, err := client.Video().InfoByBVID(context.Background(), "BV1GJ411x7h7")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s — 播放量 %d\n", info.Title, info.Stat.View)
}
```

---

## 认证

### Cookie 登录

将浏览器 Cookie 字符串（`document.cookie`）传入：

```go
cred := bilibili.NewCredentialFromCookieString(
    "SESSDATA=xxx; bili_jct=yyy; DedeUserID=123; buvid3=zzz; buvid4=www",
)
client := bilibili.NewClient()
client.SetCredential(cred)
```

也可以从 `[]*http.Cookie` 构建：

```go
cred := bilibili.NewCredentialFromHTTPCookies(resp.Cookies())
```

#### 从 .env 加载（推荐用于本地开发）

在项目根目录创建 `.env`：

```dotenv
BILIBILI_COOKIE=SESSDATA=xxx; bili_jct=yyy; DedeUserID=123; buvid3=zzz; buvid4=www
```

---

### 扫码登录

```go
client := bilibili.NewClient()
ctx := context.Background()

// 1. 生成二维码
gen, err := client.Login().QRCodeGenerate(ctx)
// gen.URL       → 二维码内容（用扫码 App 扫描）
// gen.QRCodeKey → 用于轮询

// 2. 将 gen.URL 渲染成二维码，或粘贴到 https://cli.im 生成图片

// 3. 轮询
for {
    time.Sleep(2 * time.Second)
    poll, err := client.Login().QRCodePoll(ctx, gen.QRCodeKey)
    if err != nil { ... }
    switch poll.Code {
    case 86101:
        fmt.Println("等待扫码…")
    case 86090:
        fmt.Println("已扫码，请确认…")
    case 86038:
        log.Fatal("二维码已过期")
    case 0:
        fmt.Println("登录成功")
        // Cookie 已自动写入 client 的 jar，可直接调用需登录的 API
        nav, _ := client.Login().Nav(ctx)
        fmt.Println("用户:", nav.Uname)
        return
    }
}
```

> **集成测试** 中提供了可交互的扫码测试：
> ```bash
> BILIBILI_QR_TESTS=1 go test -v -run TestIntegrationLogin_QRCodeFlow -timeout 180s ./...
> ```

---

## 功能模块

### 视频 Video

```go
v := client.Video()
ctx := context.Background()

// 按 BVID 查询
info, _ := v.InfoByBVID(ctx, "BV1GJ411x7h7")
fmt.Println(info.Title, info.AID, info.Stat.Like)

// 按 AID 查询
info, _ = v.InfoByAID(ctx, 170001)

// 热门视频（page/pageSize）
popular, _ := v.Popular(ctx, 1, 20)

// 播放地址（需登录获取高清）
play, _ := v.PlayURL(ctx, aid, cid, 116 /*4K*/)
fmt.Println(play.Dash.Video[0].BaseURL)
```

**VideoInfo 关键字段**

| 字段 | 含义 |
|---|---|
| `BVID` / `AID` | 视频 ID |
| `Title` / `Desc` | 标题 / 简介 |
| `Owner.Name` | UP 主名称 |
| `Stat.View` / `Like` / `Coin` | 播放 / 点赞 / 投币 |
| `Pages[].CID` | 分P的 cid（PlayURL 需要） |

---

### 用户 User

```go
u := client.User()

info, _  := u.Info(ctx, 2)           // 用户信息
stat, _  := u.RelationStat(ctx, 2)   // 关注 / 粉丝数
videos,_ := u.Videos(ctx, 2, 1, 20)  // 投稿视频列表
fans, _  := u.Followers(ctx, 2, 1, 20) // 粉丝列表（需登录）
```

---

### 搜索 Search

```go
s := client.Search()

// 综合搜索
all, _ := s.All(ctx, "Go 语言", 1)

// 分类型搜索：SearchTypeVideo / SearchTypeUser / SearchTypeBangumi ...
result, _ := s.ByType(ctx, bilibili.SearchTypeVideo, "golang", 1)
```

搜索类型常量：

| 常量 | 含义 |
|---|---|
| `SearchTypeVideo` | 视频 |
| `SearchTypeBangumi` | 番剧 |
| `SearchTypeFT` | 影视 |
| `SearchTypeUser` | 用户 |
| `SearchTypeLiveRoom` | 直播间 |
| `SearchTypeTopic` | 话题 |

---

### 直播 Live

```go
l := client.Live()

room, _ := l.RoomInfo(ctx, 1)                  // 直播间信息
status, _ := l.StatusByUIDs(ctx, []int64{2, 3}) // 批量查询直播状态
danmu, _ := l.DanmuInfo(ctx, 1)                // 弹幕连接信息（需鉴权）
```

---

### 评论 Comment

```go
// 获取视频评论（aid, page）
list, _ := client.GetVideoComment(170001, 1)
for _, r := range list.Replies {
    fmt.Printf("[%s] %s\n", r.Member.Uname, r.Content.Message)
}

// 发送评论（需登录）
client.SendVideoComment(170001, "这个视频很棒！")

// 回复评论（需登录）
client.ReplyComment("170001", 1, "同意！", rootRpid, parentRpid)

// 删除评论（需登录）
client.DelComment(170001, 1, rpid)

// 点赞评论（action=1点赞，0取消）
client.LikeComment(170001, 1, rpid, 1)
```

---

### 视频互动 VideoOps

```go
// 查询与视频的关系（是否点赞/收藏/投币）
rel, _ := client.GetVideoRelation(170001)
fmt.Println(rel.Like, rel.Favorite, rel.Coin)

// 点赞（like=1点赞，2取消）
client.LikeVideo(170001, 1)

// 投币（multiply=1投1枚，2投2枚）
client.CoinVideo(170001, 1)

// 三连（点赞+投币+收藏）
client.TripleAction(170001)

// 收藏到收藏夹
client.FavVideo(170001, mediaID)
```

---

### 私信消息 Session

```go
// 私信会话列表（需登录）
feed, _ := client.GetMsgFeed(1)
for _, item := range feed {
    fmt.Println(item.Uname, item.LastMsg)
}

// 聊天记录（page=1最新30条）
history, _ := client.GetChatHistory(userID, 1)

// 发送私信（需登录）
client.SendMsg(targetUID, "你好！")

// 未读消息总数
total, _ := client.GetUnreadMsg()
```

---

### 趋势排行 Trend

```go
// 全站排行榜（rid=0全站，可指定分区 rid）
list, _ := client.GetRanking(0)

// 热门 tag（rid 指定分区）
tags, _ := client.GetHotTags(3)

// tag 信息
info, _ := client.GetTagInfo("游戏")
fmt.Printf("使用数=%d 关注数=%d\n", info.Count.Use, info.Count.Atten)

// tag 相关视频（内部调用搜索）
videos, _ := client.GetTagVideos("golang", 1)

// 粉丝列表（当前登录用户，需登录）
fans, _ := client.GetFans(1, 20)

// 用户投稿视频
videos, _ = client.GetUserVideos(mid, 1, 20)
```

---

## 配置项

```go
client := bilibili.NewClient(
    bilibili.WithTimeout(15 * time.Second),   // 请求超时（默认 30s）
    bilibili.WithUserAgent("MyBot/1.0"),       // 自定义 UA
    bilibili.WithDebug(true),                  // 打印请求/响应日志
    bilibili.WithWBIRetryTimes(5),             // WBI 签名失败重试次数（默认 3）
    bilibili.WithHTTPClient(myHTTPClient),     // 自定义 http.Client
)
```

---

## 测试

### 配置 .env

所有开关均通过 `.env` 文件控制，无需设置环境变量：

```dotenv
# 必填：登录态 cookie
BILIBILI_COOKIE=SESSDATA=xxx; bili_jct=yyy; DedeUserID=123; buvid3=zzz; buvid4=www

# 可选：开启写操作测试（点赞 / 投币 / 发评论 / 私信）
# BILIBILI_WRITE_TESTS=1

# 可选：开启扫码登录交互测试
# BILIBILI_QR_TESTS=1
```

### 运行测试

```bash
# 单元测试（无需 cookie）
go test -v -run ^TestBuild ./...

# 集成测试（读操作，需要 cookie）
go test -v -run TestIntegration -timeout 60s ./...

# 含写操作：在 .env 中取消注释 BILIBILI_WRITE_TESTS=1，然后：
go test -v -run TestIntegration -timeout 60s ./...

# 扫码登录测试：在 .env 中取消注释 BILIBILI_QR_TESTS=1，然后：
go test -v -run TestIntegrationLogin_QRCodeFlow -timeout 180s ./...
```

扫码测试成功后会在日志中打印新的 Cookie 字符串，复制后更新 `.env` 即可。

---

## 架构说明

```
Client
├── Config          — 超时、UA、BaseURL 等
├── Credential      — SESSDATA / bili_jct / DedeUserID 等 cookie
├── cookieJar       — 线程安全的 http.CookieJar 封装
├── wbiManager      — WBI 签名密钥的懒加载与自动刷新
├── RequestBuilder  — 参数组装 / WBI 签名 / CSRF 注入 / 响应解码
└── *Service        — 面向业务域的稳定 API
    ├── VideoService
    ├── UserService
    ├── SearchService
    ├── LiveService
    └── LoginService
```

`endpoint` 结构通过声明式元数据描述每个接口的路径、HTTP 方法、是否需要 WBI 签名、CSRF、登录态等，`RequestBuilder` 统一消费这些声明，业务层无需关心底层细节。

---

## License

MIT
