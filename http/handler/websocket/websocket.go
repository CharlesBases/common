package websocket

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// WebSocket websocket
type WebSocket interface {
	// Init init options
	Init(opts ...Option) error
	// Options return options
	Options() Options
	// Connection 建立连接
	Connection(w http.ResponseWriter, r *http.Request) error
	// Subscription 订阅
	Subscription()
	// Unsubscription 取消订阅
	Unsubscription()
	// Read read params from the request
	Read() error
	// Write write data to the response
	Write(v interface{}) error
	// Ping ping of the websocket
	Ping()
	// Close 关闭 websocket 连接
	Close()
}

// websocketRequest
type websocketRequest struct {
	Version string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	ID      *json.RawMessage `json:"id"`
	Params  *json.RawMessage `json:"params"`
}

// websocketResponse
type websocketResponse struct {
	Version string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	ID      *json.RawMessage `json:"id"`
	Data    interface{}      `json:"data"`
}

// opt websocket options
type opt struct {
	options Options
}

// NewHandler .
func NewHandler(opts ...Option) *opt {
	var opt = new(opt)
	opt.init(opts...)
	return opt
}

func (opt *opt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stream := opt.newConn(r)
	defer stream.Close()

	stream.Connection(w, r)
	return
}

// init .
func (opt *opt) init(opts ...Option) {
	for _, o := range opts {
		o(&opt.options)
	}
}

// newConn new connection for the websocket
func (opt *opt) newConn(r *http.Request) *stream {
	var capacity = 32
	if opt.options.Buffer != 0 {
		capacity = opt.options.Buffer
	}

	return &stream{
		id:            uuid.New().String(),
		metadata:      opt.parseHeaderFromRequest(r),
		request:       make(chan *websocketRequest, capacity),
		response:      make(chan *websocketResponse, capacity),
		broadcast:     make(chan struct{}, capacity),
		subscriptions: make([]string, 0),
		ctx:           r.Context(),
		options:       opt.options,
		ready:         true,
		close:         make(chan struct{}, 0),
	}
}

// parseHeaderFromRequest parse header from request
func (opt *opt) parseHeaderFromRequest(r *http.Request) metadata {
	var data = make(map[string]string, 0)
	for key, val := range r.Header {
		data[key] = strings.Join(val, ",")
	}
	return data
}

// writerErrorToResponse write error for response
func (opt *opt) writerErrorToResponse(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
	data, _ := json.Marshal(map[string]interface{}{
		"code":    statusCode,
		"message": http.StatusText(statusCode),
	})
	rw.Write(data)
}
