package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/CharlesBases/common/log"
	"github.com/gorilla/websocket"

	"charlesbases/http/handler/websocket/pb"
)

// upgrader websocker upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: time.Second * 3,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

type metadata map[string]string

// stream websockt stream
type stream struct {
	id            string          // 连接 id
	metadata      metadata        // 随路数据
	subscriptions map[string]bool // 订阅的消息列表
	active        bool            // 当前 websocket 是否活跃
	disconnect    chan struct{}   // websocket 退出

	request   chan *WebSocketRequest   // 请求
	response  chan *WebSocketResponse  // 响应
	broadcast chan *WebSocketBroadcast // 广播

	ctx  context.Context
	conn *websocket.Conn

	lock sync.RWMutex

	options Options
}

// isCloseError .
func (stream *stream) isCloseError(err error) int {
	if e, ok := err.(*websocket.CloseError); ok {
		return e.Code
	}
	return -1
}

// handling websocket 请求处理
func (stream *stream) handling() {
	for {
		select {
		case request := <-stream.request:
			switch request.Method {
			// 消息订阅
			case pb.Method_subscription.String():
				stream.eventSubscription(request)
			// 取消订阅
			case pb.Method_unsubscription.String():
				stream.eventUnsubscription(request)
			// 断开连接
			case pb.Method_disconnect.String():
				stream.Disconnect()
			}
		case response := <-stream.response:
			stream.Write(response)
		case broadcast := <-stream.broadcast:
			stream.eventBroadcast(broadcast)
		case <-stream.disconnect:
			stream.Disconnect()
			return
		}
	}
}

// eventBroadcast 广播
func (stream *stream) eventBroadcast(broadcast *WebSocketBroadcast) {
	var isBroadcast bool

	stream.lock.Lock()
	if _, ok := stream.subscriptions[broadcast.Topic]; ok {
		isBroadcast = true
	} else {
		for topic := range stream.subscriptions {
			// 订阅的 topic 是否支持前缀匹配
			if strings.HasSuffix(topic, "*") {
				if strings.HasPrefix(broadcast.Topic, strings.TrimSuffix(topic, "*")) {
					isBroadcast = true
					break
				}
			}
		}
	}
	stream.lock.Unlock()

	if isBroadcast {
		stream.Write(&WebSocketResponse{ID: stream.id, Method: pb.Method_broadcast.String(), Data: broadcast})
	}
}

// eventSubscription 广播订阅
func (stream *stream) eventSubscription(request *WebSocketRequest) {
	var topics = make([]string, 0)

	if err := stream.Unmarshal(*request.Params, &topics); err != nil {
		log.Error("[WebSocketID: %s] broadcast subscription error: invalid params format", stream.id)
		stream.Disconnect()
		return
	}

	stream.lock.Lock()
	for _, topic := range topics {
		stream.subscriptions[topic] = true
	}
	stream.lock.Unlock()
}

// eventUnsubscription 取消订阅
func (stream *stream) eventUnsubscription(request *WebSocketRequest) {
	var topics = make([]string, 0)

	if err := stream.Unmarshal(*request.Params, &topics); err != nil {
		log.Error("[WebSocketID: %s] broadcast unsubscription error: invalid params format", stream.id)
		stream.Disconnect()
		return
	}

	stream.lock.Lock()
	for _, topic := range topics {
		delete(stream.subscriptions, topic)
	}
	stream.lock.Unlock()
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

	log.Debugf("[WebSocketID: %s] connect", stream.id)

	// ping
	stream.Ping()

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
		request := new(WebSocketRequest)
		if err := stream.conn.ReadJSON(request); err != nil {
			switch stream.isCloseError(err) {
			case websocket.CloseNoStatusReceived:
			default:
				log.Errorf("[WebSocketID: %s] read message error: %v", stream.id, err)
			}

			stream.Disconnect()
			break
		}

		if request.Params == nil {
			log.Errorf("[WebSocketID: %s] read message error: params must be not nil", stream.id)
			stream.Disconnect()
			break
		}

		stream.request <- request
	}
	return nil
}

// Write .
func (stream *stream) Write(v *WebSocketResponse) error {
	if stream.active {
		return stream.conn.WriteJSON(v)
	}
	return nil
}

// Marshal .
func (stream *stream) Marshal(v interface{}) ([]byte, error) {
	// json
	return json.Marshal(v)
}

// Unmarshal .
func (stream *stream) Unmarshal(data []byte, v interface{}) error {
	// json
	return json.Unmarshal(data, v)
}

// Ping .
func (stream *stream) Ping() {
	stream.Write(&WebSocketResponse{
		ID:     stream.id,
		Method: "ping",
		Data:   "OK",
	})
}

// Disconnect .
func (stream *stream) Disconnect() {
	if stream.active {
		log.Debugf("[WebSocketID: %s] disconnect", stream.id)

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
