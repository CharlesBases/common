package wechat

import (
	"fmt"

	"github.com/CharlesBases/common/db"
	"github.com/CharlesBases/common/web/request"
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
	Redis := db.GetRedis(address)
	if accessToken, err := Redis.Get("authorizer_access_token"); err == nil {
		if checkAccessToken(accessToken) {
			return accessToken, true
		}
	}

	if resp, err := request.Request("GET", "https://api.weixin.qq.com/cgi-bin/token?appid="+wechat.conf.AppID+"&secret="+wechat.conf.AppSecret+"&grant_type=client_credential"); err == nil {
		if resp.Body["errcode"] != nil {
			return fmt.Sprintf(`%v`, resp.Body["errcode"]), false
		} else {
			Redis.Set("authorizer_access_token", resp.Body["access_token"])
			return resp.Body["access_token"].(string), true
		}
	}
	return "error", false
}

func checkAccessToken(accessToken string) bool {
	if resp, err := request.Request("GET", "https://api.weixin.qq.com/cgi-bin/menu/get?access_token="+accessToken); err == nil {
		if resp.Body["errcode"] != nil {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}
