package bilibili

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"
)

// TestInteractiveQRCodeLogin 完整扫码登录流程：生成二维码 → 终端展示 → 轮询结果。
//
// 独立运行方式（不干扰其他自动化测试）：
// 	go test -v -run TestInteractiveQRCodeLogin -timeout 180s ./...
//
// 测试会在终端打印二维码 URL，用手机 bilibili App 扫码后确认登录即可。
// 成功后会将新 cookie 打印到日志，用于替换 .env 中的 BILIBILI_COOKIE。
func TestInteractiveQRCodeLogin(t *testing.T) {
	t.Log("=== 开始交互式扫码登录测试 ===")
	t.Log("API: passport-login/web/qrcode/generate + poll")

	c := NewClient()
	ctx := context.Background()

	// 1. 生成二维码
	gen, err := c.Login().QRCodeGenerate(ctx)
	if err != nil {
		t.Fatalf("❌ QRCodeGenerate error: %v", err)
	}

	// 2. 在终端显示二维码提取链接（整行输出，方便复制）
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

		poll, err := c.Login().QRCodePoll(ctx, gen.QRCodeKey)
		if err != nil {
			t.Logf("⚠️ [%3ds] 轮询出错 (可能是风控拦截，请确保 endpoint 加了 dataField: \"data\"): %v", i+2, err)
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
			passportURL, _ := url.Parse(passportBase) // 假设 passportBase 在 client.go 中已定义
			httpCookies := c.cookies.Cookies(passportURL)
			cred := NewCredentialFromHTTPCookies(httpCookies)
			if cred.SessData == "" {
				mainURL, _ := url.Parse(apiBase) // 假设 apiBase 在 client.go 中已定义
				httpCookies = c.cookies.Cookies(mainURL)
				cred = NewCredentialFromHTTPCookies(httpCookies)
			}
			c.SetCredential(cred)

			t.Logf("🎉 登录成功！")
			t.Logf("SESSDATA    = %s", cred.SessData)
			t.Logf("bili_jct    = %s", cred.BiliJct)
			t.Logf("DedeUserID  = %s", cred.DedeUserID)
			t.Logf("buvid3      = %s", cred.Buvid3)
			t.Logf("buvid4      = %s", cred.Buvid4)
			t.Logf("\n👇 请将以下内容写入 .env 的 BILIBILI_COOKIE 字段（分号分隔）👇")
			t.Logf("\nSESSDATA=%s; bili_jct=%s; DedeUserID=%s; buvid3=%s; buvid4=%s\n",
				cred.SessData, cred.BiliJct, cred.DedeUserID, cred.Buvid3, cred.Buvid4)

			// 验证新 cookie 有效
			nav, navErr := c.Login().Nav(ctx)
			if navErr != nil {
				t.Fatalf("❌ Nav 验证失败: %v", navErr)
			}
			t.Logf("✅ 已验证新凭据有效，当前用户: %s (mid=%d)", nav.Uname, nav.Mid)
			return
		default:
			t.Logf("⚠️ [%3ds] 未知状态码: %d message=%s", i+2, poll.Code, poll.Message)
		}
	}

	t.Fatal("❌ 超时（120s）未完成扫码，测试终止")
}

// qrBanner 生成终端友好的提示文字，链接保持单行方便复制。
func qrBanner(url string) string {
	var sb strings.Builder
	sb.WriteString("========================================================\n")
	sb.WriteString("📱 请用 bilibili App 扫码登录\n")
	sb.WriteString("========================================================\n")
	sb.WriteString("👉 请复制下方整行链接，粘贴到二维码生成器（如 https://cli.im ）中生成二维码后扫码：\n\n")
	sb.WriteString(url)
	sb.WriteString("\n\n========================================================")
	return sb.String()
}