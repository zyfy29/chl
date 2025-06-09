package feishu

import (
	"chl/config"
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
		return nil, resp.Error().(*resty.ResponseError).Err
	}

	return resp.Result().(*AuthResponse), nil
}
