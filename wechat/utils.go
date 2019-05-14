package wechat

type Wechat struct {
	conf Config
	msg  Msg
}

type Msg struct {
	ToUserName   string  // 消息接收方
	FromUserName string  // 消息发送方
	CreateTime   int     // 消息创建时间
	MsgType      string  // 消息类型
	Content      string  // 文本消息内容
	PicUrl       string  // 图片链接
	MediaId      string  // 多媒体文件id
	Format       string  // 语音格式，如amr，speex等
	Recognition  string  // 语音识别结果，UTF8编码（需要公众号开通语音识别）
	ThumbMediaId string  // 视频消息缩略图的媒体id
	LocationX    float64 // 地理位置维度
	LocationY    float64 // 地理位置经度
	Scale        int     // 地图缩放大小
	Label        string  // 地理位置信息
	Title        string  // 消息标题
	Description  string  // 消息描述
	Url          string  // 消息链接
	MsgId        int64   // 消息ID
	Event        string  // 事件类型
	EventKey     string  // 事件KEY值，qrscene_为前缀，后面为二维码的参数值
	Ticket       string  // 二维码的ticket，可用来换取二维码图片
	Latitude     float64 // 地理位置纬度
	Longitude    float64 // 地理位置经度
	Precision    float64 // 地理位置精度
}

type Config struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
}
