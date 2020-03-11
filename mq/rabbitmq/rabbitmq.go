package rabbitmq

import (
	"sync"

	mq "github.com/streadway/amqp"

	"charlesbases/common/log"
)

const (
	rabbiturl = "amqp://guest:guest@localhost:5672/"
)

var (
	mqconn *mq.Connection
	mqchan *mq.Channel
)

type Producer interface {
	MsgContent() string
}

type Consumer interface {
	Consumer([]byte) error
}

type RabbitMQ struct {
	connection   *mq.Connection
	channel      *mq.Channel
	queueName    string // 队列名称
	routingKey   string // key 名称
	exchangeName string // 交换机名称
	exchangeType string // 交换机类型
	producers    []Producer
	consumers    []Consumer
	mutex        sync.RWMutex
}

type QueueExchange struct {
	QueueName    string
	RoutingKey   string
	ExchangeName string
	ExchangeType string
}

func New(quene *QueueExchange) *RabbitMQ {
	return &RabbitMQ{
		queueName:    quene.QueueName,
		routingKey:   quene.RoutingKey,
		exchangeName: quene.ExchangeName,
		exchangeType: quene.ExchangeType,
	}
}

func (r *RabbitMQ) Start() {

}

func (r *RabbitMQ) connect() {
	conn, err := mq.Dial(rabbiturl)
	if err != nil {
		log.Error("RabbitMQ Connect Error: ", err)
		return
	}
	if conn == nil {
		log.Error("RabbitMQ Connect Error: conn is nil")
		return
	}
	r.connection = conn
	channel, err := conn.Channel()
	if err != nil {
		log.Error("RabbitMQ Channel Error: ", err)
		return
	}
	if channel == nil {
		log.Error("RabbitMQ Channel Error: channel is nil")
		return
	}
	r.channel = channel
}

func (r *RabbitMQ) close() {
	if r.channel != nil {
		err := r.channel.Close()
		if err != nil {
			log.Error("RabbitMQ Channel Close Error: ", err)
			return
		}
	}
	if r.connection != nil {
		err := r.connection.Close()
		if err != nil {
			log.Error("RabbitMQ Connection Close Error: ", err)
			return
		}
	}
}
