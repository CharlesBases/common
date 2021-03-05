package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/CharlesBases/common/log"
	"github.com/gorilla/websocket"

	"charlesbases/http/handler/websocket/pb"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type metadata map[string]string

// stream websockt stream
type stream struct {
	id            string        // 连接 id
	metadata      metadata      // 随路数据
	subscriptions []string      // 订阅的消息列表
	active        bool          // 当前 websocket 是否活跃
	disconnect    chan struct{} // websocket 退出

	request   chan *pb.WebSocketRequest   // 请求
	response  chan *pb.WebSocketResponse  // 响应
	broadcast chan *pb.WebSocketBroadcast // 广播

	ctx  context.Context
	conn *websocket.Conn

	options Options
}

// handling websocket 请求处理
func (stream *stream) handling() {
	for {
		select {
		case request := <-stream.request:
			switch request.Method {
			// 消息订阅
			case pb.Method_subscription.String():
			// 取消订阅
			case pb.Method_unsubscription.String():
			// 断开连接
			case pb.Method_disconnect.String():
				stream.Disconnect()
			}
		case response := <-stream.response:
			stream.Write(response)
		case broadcast := <-stream.broadcast:
			stream.Write(broadcast)
		case <-stream.disconnect:
			stream.Disconnect()
			return
		}
	}
}

// Init .
func (stream *stream) Init(opts ...Option) error {
	for _, o := range opts {
		o(&stream.options)
	}
	return nil
}

// Options .
func (s *stream) Options() Options {
	return s.options
}

// Connect .
func (stream *stream) Connect(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("websocket upgrade error: ", err)
		return err
	}

	stream.conn = conn

	// 监听 websockt 请求
	go stream.Read()

	// ping
	stream.Ping()

	// 处理 websocket 请求
	stream.handling()
	return nil
}

// Subscription .
func (stream *stream) Subscription() {

}

// Unsubscription .
func (stream *stream) Unsubscription() {

}

// Read .
func (stream *stream) Read() error {
	for {
		request := new(pb.WebSocketRequest)
		if err := stream.conn.ReadJSON(request); err != nil {
			switch stream.isCloseError(err) {
			case websocket.CloseNoStatusReceived:
				log.Debugf("websocket disconnect [ID: %s]", stream.id)
			default:
				log.Error("received request error: ", err)
			}

			stream.Disconnect()
			break
		}

		stream.request <- request
	}
	return nil
}

// Write .
func (stream *stream) Write(v interface{}) error {
	switch v.(type) {
	// case proto.Message:
	// 	data, _ := proto.Marshal(v.(proto.Message))
	// 	return stream.conn.WriteMessage(websocket.BinaryMessage, data)
	default:
		return stream.conn.WriteJSON(v)
	}
}

// Ping .
func (stream *stream) Ping() {
	stream.Write(&pb.WebSocketResponse{
		Version: "test",
		Method:  "ping",
		Id:      nil,
		Data:    []byte(time.Now().Format("2006-01-02 15:04:05")),
	})
}

// Disconnect .
func (stream *stream) Disconnect() {
	if stream.active {
		stream.active = false
		stream.disconnect <- struct{}{}

		close(stream.request)
		close(stream.response)
		close(stream.broadcast)
		close(stream.disconnect)

		if stream.conn != nil {
			stream.conn.Close()
			stream.conn = nil
		}
	}
}

// isCloseError .
func (stream *stream) isCloseError(err error) int {
	if e, ok := err.(*websocket.CloseError); ok {
		return e.Code
	}
	return -1
}
