package wechat

import (
	"fmt"

	"github.com/CharlesBases/common/databases/redis"
	"github.com/CharlesBases/common/http/web/request"
)

var (
	address = "192.168.1.88:6379"
)

const (
	WECHAT_APPID     = "WECHAT_APPID"
	WECHAT_APPSECRET = "WECHAT_APPSECRET"
	WECHAT_TOKEN     = "WECHAT_TOKEN"
	WECHAT_AESKEY    = "WECHAT_AESKEY"
)

func GetWeChat() *Wechat {
	config := Config{
		WECHAT_APPID,
		WECHAT_APPSECRET,
		WECHAT_TOKEN,
		WECHAT_AESKEY,
	}
	return &Wechat{conf: config}
}

func (wechat *Wechat) getAccessToken() (string, bool) {
	Redis := redis.GetRedis(address)
	var accessToken string
	if err := Redis.Get(redis.HASH, "authorizer_access_token", &accessToken); err == nil {
		if checkAccessToken(accessToken) {
			return accessToken, true
		}
	}
	if resp, err := request.Request("GET", getAccessTokenUrl(wechat.conf.AppID, wechat.conf.AppSecret)); err == nil {
		if resp.Body["errcode"] != nil {
			return fmt.Sprintf(`%v`, resp.Body["errcode"]), false
		} else {
			Redis.Set(redis.HASH, "authorizer_access_token", resp.Body["access_token"])
			return resp.Body["access_token"].(string), true
		}
	}
	return "error", false
}

func checkAccessToken(accessToken string) bool {
	if resp, err := request.Request("GET", genAccessTokenUrl(accessToken)); err == nil {
		if resp.Body["errcode"] != nil {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

func getAccessTokenUrl(appID string, appSecret string) string {
	return fmt.Sprintf(`https://api.weixin.qq.com/cgi-bin/token?appid=%s&secret=%s&grant_type=client_credential`, appID, appSecret)
}

func genAccessTokenUrl(accessToken string) string {
	return fmt.Sprintf(`https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s`, accessToken)
}
