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

func Test_index2Range(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 0",
			args: args{i: 0, j: 0},
			want: "A1",
		},
		{
			name: "case 1",
			args: args{i: 0, j: 3},
			want: "D1",
		},
		{
			name: "case 2",
			args: args{i: 3, j: 4},
			want: "E4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Index2Range(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("Index2Range() = %v, want %v", got, tt.want)
			}
		})
	}
}
