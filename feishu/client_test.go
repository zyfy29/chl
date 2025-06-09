package feishu

import (
	"chl/config"
	"testing"
)

func TestAuth(t *testing.T) {
	authResponse, err := Api.Auth(config.Conf.Feishu.AppID, config.Conf.Feishu.AppSecret)
	if err != nil {
		t.Fatalf("Auth failed: %v", err)
	}
	if authResponse.Code != 0 {
		t.Errorf("Expected code 0, got %d: %s", authResponse.Code, authResponse.Msg)
	} else {
		t.Logf("Auth successful: %s", authResponse.TenantAccessToken)
	}
}
