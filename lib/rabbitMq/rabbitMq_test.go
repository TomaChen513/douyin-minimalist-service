package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/streadway/amqp"
)


func TestXxx(t *testing.T) {
	MQURL := "amqp://tiktok:tiktok@121.5.231.228:5672/"

	fmt.Println("2")
	Rmq := &RabbitMQ{
		mqurl: MQURL,
	}
	dial, err := amqp.Dial(Rmq.mqurl)
	Rmq.failOnErr(err, "创建连接失败")
	Rmq.conn = dial
}

