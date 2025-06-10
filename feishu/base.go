package feishu

import (
	"fmt"
	"github.com/zyfy29/chl/config"
	"resty.dev/v3"
)

var Api Client

type Client struct {
	r *resty.Client
}

func init() {
	r := resty.New()
	r.SetBaseURL("https://open.feishu.cn/open-apis")
	r.SetHeader("Authorization", "Bearer "+config.Conf.Feishu.TenantAccessToken)
	Api = Client{r}
}

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type AuthResponse struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

func (c *Client) Auth(appId, appSecret string) (*AuthResponse, error) {
	resp, err := c.r.R().
		SetResult(&AuthResponse{}).
		SetBody(map[string]interface{}{
			"app_id":     appId,
			"app_secret": appSecret,
		}).
		Post("/auth/v3/tenant_access_token/internal")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error getting tenant_access_token\nstatus: %d\nbody: %s", resp.StatusCode(), resp.String())
	}

	return resp.Result().(*AuthResponse), nil
}

// Index2Range Converts 0-indexed i, j to A1 notation
func Index2Range(i, j int) string {
	col := string('A' + byte(j))
	row := i + 1
	return fmt.Sprintf("%s%d", col, row)
}
