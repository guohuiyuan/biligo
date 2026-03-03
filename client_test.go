package bilibili

import "testing"

func TestBuildMixinKey(t *testing.T) {
	got := buildMixinKey("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")
	if got == "" {
		t.Fatal("expected mixin key")
	}
	if len(got) != 32 {
		t.Fatalf("expected 32 chars, got %d", len(got))
	}
}

func TestCredentialParse(t *testing.T) {
	credential := NewCredentialFromCookieString("SESSDATA=a; bili_jct=b; DedeUserID=1; buvid3=c; buvid4=d")
	if credential.SessData != "a" || credential.BiliJct != "b" || credential.DedeUserID != "1" {
		t.Fatal("credential parse failed")
	}
}
