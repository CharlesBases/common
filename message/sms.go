package message

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CharlesBases/daddylab/DaddyLab.MS.Common/baseservice/sms/sdk-submail/submail"
	"github.com/CharlesBases/daddylab/DaddyLab.MS.Common/web"

	"github.com/CharlesBases/common/log"
)

var (
	// 助通短信
	ZTSMS_URL      = "http://api.zthysms.com/sendSms.do"
	ZTSMS_USERNAME = "lbpc888hy"
	ZTSMS_PASSWORD = "DaddyLab520"
)

func SendVerificationCodeByZTSms(code, mobile string) (errmsg string, errcode int) {
	tkey := time.Now().Format("20060102150405")
	password := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%x", md5.Sum([]byte(ZTSMS_PASSWORD)))+tkey)))
	resp, err := web.Get(ZTSMS_URL +
		"?content=【老爸评测】用户注册验证码:验证码" + code + "，您正在注册成为老爸评测用户，感谢您的支持！" +
		"&mobile=" + mobile +
		"&username=" + ZTSMS_USERNAME +
		"&tkey=" + tkey +
		"&password=" + password)
	if err == nil {
		if resp.Raw.StatusCode == 200 {
			data := []byte(resp.Body)
			result := strings.Split(string(data), ",")
			if result[0] == "1" {
				return "短信发送成功", 200
			}
			return result[1], 401
		}
	}
	return "获取验证码请求失败", 402
}

type KVPair struct {
	Key   string
	Value string
}

/**

请求成功
{
      "status":"success"
      "send_id":"093c0a7df143c087d6cba9cdf0cf3738"
      "fee":1,
      "sms_credits":14197
}

请求失败
{
      "status":"error",
      "code":"1xx",
      "msg":"error message"
}

https://www.mysubmail.com/chs/documents/developer/t2f1J2
*/
type ResultSendSms struct {
	Status                    string `json:"status"` //"success"/"error"
	Send_id                   string `json:"send_id"`
	Fee                       int    `json:"fee"`
	Sms_credits               string `json:"sms_credits"`
	Code                      string `json:"code"`
	Msg                       string `json:"msg"`
	Transactional_sms_credits string `json:"transactional_sms_credits"`
}

/**
使用Mysubmail 发送短信
*/
func SendSmsByMysubmail(appid, appkey, sendto, project string, vars ...KVPair) ResultSendSms {
	messageconfig := make(map[string]string)
	messageconfig["appid"] = appid   //"24364"
	messageconfig["appkey"] = appkey //"5f5d902f5abda665288b784801e9f950"
	messageconfig["signtype"] = signtype
	//messagexsend
	messagexsend := submail.CreateMessageXSend()
	submail.MessageXSendAddTo(messagexsend, sendto)       //"13067863233")
	submail.MessageXSendSetProject(messagexsend, project) //"kJSad3")
	for _, v := range vars {
		submail.MessageXSendAddVar(messagexsend, v.Key, v.Value)
	}
	//{"status":"success","send_id":"7247b703a91f2721ff1f4c497a5e0369","fee":1,"sms_credits":"500031","transactional_sms_credits":"0"}
	requestBody := submail.MessageXSendBuildRequest(messagexsend)
	resultMessage := submail.MessageXSendRun(requestBody, messageconfig)
	log.Info("MessageXSend ", resultMessage)
	result := new(ResultSendSms)
	e := json.Unmarshal([]byte(resultMessage), result)
	if e != nil {
		result.Status = "error"
		result.Code = "001" //submail未返回有效数据
		result.Msg = resultMessage
	}
	return *result
}

var (
	appid                   = "24364"
	appkey                  = "f78f865741ae065b1f26faec06407177"
	signtype                = "md5"
	verificationCodeProject = "kJSad3"
	noticeMsgProject        = "5thzo2"
)

/**
发送验证码短信
*/
func SendVerificationCodeByMysubmail(sendto, code string) ResultSendSms {
	vars := KVPair{Key: "code", Value: code}
	return SendSmsByMysubmail(appid, appkey, sendto, verificationCodeProject, vars)
}

/**
发送通知
*/
func SendNoticeMsgByMysubmail(sendto, msg string) ResultSendSms {
	vars := KVPair{Key: "msg", Value: msg}
	return SendSmsByMysubmail(appid, appkey, sendto, noticeMsgProject, vars)
}

//:7470 必须使用 POST
//
//SubHook回调接口
//
//resquest			//发送请求被接收
//
//{
//"events":"request",
//"address":"138xxxxxxx",
//"send_id":"093c0a7df143c087d6cba9cdf0cf3738",
//"app":xxxxx,
//"timestamp":1415014855,
//"token":"067ef7e2f286a9a56eabb07dc9657852",
//"signature":"a70d09a9345adfdd353d34a505dac4ca"
//}
//
//delivered			//发送成功
//
//{
//"events":"delivered",
//"address":"138xxxxxxxx",
//"send_id":"093c0a7df143c087d6cba9cdf0cf3738",
//"app":xxxxxx,
//"timestamp":1415014855,
//"token":"067ef7e2f286a9a56eabb07dc9657852",
//"signature":"a70d09a9345adfdd353d34a505dac4ca"
//}
//
//dropped			//发送失败
//
//{
//"events":"dropped",
//"address":"138xxxxxxxx",
//"send_id":"093c0a7df143c087d6cba9cdf0cf3738",
//"report":" UNDELIV",
//"app":xxxxxx,
//"timestamp":1415014855,
//"token":"067ef7e2f286a9a56eabb07dc9657852",
//"signature":"a70d09a9345adfdd353d34a505dac4ca"
//}
//
//sending			//正在发送
//
//{
//"events":"sending",
//"address":"138xxxxxxxx",
//"send_id":"093c0a7df143c087d6cba9cdf0cf3738",
//"app":xxxxxx,
//"timestamp":1415014855,
//"token":"067ef7e2f286a9a56eabb07dc9657852",
//"signature":"a70d09a9345adfdd353d34a505dac4ca"
//}
//
//mo				//短信上行（指用户回复和上行）
//
//{
//"events":"mo",
//"address":"138xxxxxxxx",
//"app":xxxxxx,
//"content":"xxxxxx",
//"timestamp":1415014855,
//"token":"067ef7e2f286a9a56eabb07dc9657852",
//"signature":"a70d09a9345adfdd353d34a505dac4ca"
//}
//
//unkown			//未知状态（指无法从网关获取短信回执状态）
//
//{
//"events":"unkown",
//"address":"138xxxxxxx",
//"send_id":"093c0a7df143c087d6cba9cdf0cf3738",
//"app":xxxxx,
//"timestamp":1415014855,
//"token":"067ef7e2f286a9a56eabb07dc9657852",
//"signature":"a70d09a9345adfdd353d34a505dac4ca"
//}
//
//template_accept	//短信模板审核通过
//
//{
//"events":"template_accept",
//"template_id":"H5OSN4",
//"timestamp":1415014855,
//"token":"067ef7e2f286a9a56eabb07dc9657852",
//"signature":"a70d09a9345adfdd353d34a505dac4ca"
//}
//
//template_reject	//短信模板审核未通过
//
//{
//"events":"template_reject",
//"reason":"签名不正确",
//"template_id":"H5OSN4",
//"timestamp":1415014855,
//"token":"067ef7e2f286a9a56eabb07dc9657852",
//"signature":"a70d09a9345adfdd353d34a505dac4ca"
//}
//
//https://www.mysubmail.com/chs/documents/developer/Uau792
func SubHookHandler(resp http.ResponseWriter, req *http.Request) {
	defer log.Flush()
	events := req.FormValue("events")
	address := req.FormValue("address")
	send_id := req.FormValue("send_id")
	app := req.FormValue("app")
	timestamp := req.FormValue("timestamp")
	token := req.FormValue("token")
	signature := req.FormValue("signature")
	content := req.FormValue("content")
	template_id := req.FormValue("template_id")
	reason := req.FormValue("reason")
	report := req.FormValue("report")

	timestampint64, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Error(err)
	}
	subhookdata := SubHookData{
		Events:      events,
		Address:     address,
		Send_id:     send_id,
		App:         app,
		Timestamp:   timestampint64,
		Token:       token,
		Signature:   signature,
		Content:     content,
		Template_id: template_id,
		Reason:      reason,
		Report:      report,
		EventTime:   time.Unix(timestampint64, 0).UTC(),
	}
	var shdh = SubHookDataHandler
	for shdh != nil {
		shdh = shdh(&subhookdata)
	}
}

//实现此方法用于回调subHook的回调，返回nil结束回调
var SubHookDataHandler subHookDataHandler

type subHookDataHandler func(subhookdata *SubHookData) subHookDataHandler

type SubHookData struct {
	Events      string    `json:"events"`
	Address     string    `json:"address"`
	Send_id     string    `json:"send_id"`
	App         string    `json:"app"`
	Timestamp   int64     `json:"timestamp"`
	Token       string    `json:"token"`
	Signature   string    `json:"signature"`
	Content     string    `json:"content"`
	Template_id string    `json:"template_id"`
	Reason      string    `json:"reason"`
	Report      string    `json:"report"`
	EventTime   time.Time `json:"event_time"`
}

//推送数据验证
//
// 	key : SUBHOOK 密匙
func (shd SubHookData) checkSignature(key string) bool {
	tokenjkey := shd.Token + key
	hash := md5.New()
	hash.Write([]byte(tokenjkey))
	sum := hash.Sum(nil)
	md555 := hex.EncodeToString(sum)
	return md555 == shd.Signature
}
