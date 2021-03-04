package websocket

import (
	"context"
	"net/http"

	"github.com/CharlesBases/common/log"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type metadata map[string]string

// stream websockt stream
type stream struct {
	id            string        // 连接 id
	metadata      metadata      // 随路数据
	subscriptions []string      // 订阅的消息列表
	ready         bool          // 当前 websocket 是否就绪
	close         chan struct{} // websocket 退出

	request   chan *websocketRequest  // 请求
	response  chan *websocketResponse // 响应
	broadcast chan struct{}           // 广播

	ctx  context.Context
	conn *websocket.Conn

	options Options
}

// handling websocket 请求处理
func (stream *stream) handling() {
	for {
		select {
		case <-stream.request:
		case response := <-stream.response:
			stream.Write(response)
		case broadcast := <-stream.broadcast:
			stream.Write(broadcast)
		case <-stream.close:
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

// Connection .
func (stream *stream) Connection(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("websocket upgrade error: ", err)
		return err
	}

	stream.conn = conn
	defer stream.Close()

	// 监听 websockt 请求
	go stream.Read()

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
		request := new(websocketRequest)
		if err := stream.conn.ReadJSON(request); err != nil {
			log.Error("received request error: ", err)
			break
		}

		stream.request <- request
	}
	return nil
}

// Write .
func (stream *stream) Write(v interface{}) error {
	return stream.conn.WriteJSON(v)
}

// Ping .
func (stream *stream) Ping() {

}

// Close .
func (stream *stream) Close() {
	close(stream.request)
	close(stream.response)
	close(stream.broadcast)

	stream.close <- struct{}{}
	close(stream.close)

	if stream.conn != nil {
		stream.conn.Close()
		stream.conn = nil
	}

}
