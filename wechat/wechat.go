package wechat

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CharlesBases/common/db"
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
	Redis := db.GetRedis()
	accessToken, err := Redis.Get("wx5f3252c4af7c1805_authorizer_access_token").Result()
	if err == nil {
		//if wechat.CheckAccessToken(access_token) {
		// 检查AccessToken是否有用
		return accessToken, true
		//}
		// 失效 - 抢新的
		//log.Info("抢 AccessToken")
	}

	resp, err := Get("https://api.weixin.qq.com/cgi-bin/token?appid=" + wechat.conf.AppID +
		"&secret=" + wechat.conf.AppSecret +
		"&grant_type=client_credential")
	if err == nil {
		if resp.IsOk() {
			data := []byte(resp.Body)
			var f map[string]interface{}
			json.Unmarshal(data, &f)
			errcode := f["errcode"]
			if errcode != nil {
				return fmt.Sprintf("%f", f["errcode"].(float64)), false
			} else {
				// TODO access_token存redis
				_ = Redis.Set("access_token", f["access_token"].(string), 7000*time.Second).Err()
				return f["access_token"].(string), true
			}
		}
		return "error", false
	}
	return "error", false
}
