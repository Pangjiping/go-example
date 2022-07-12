package rabbitmq_go

import (
	"fmt"
	"github.com/streadway/amqp"
)

var (
	user = "admin"
	pwd  = "admin"
	host = "xxx.xxx.xxx.xxx"
	port = "5672"
)

func RabbitMQConn() (conn *amqp.Connection, err error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pwd, host, port)
	conn, err = amqp.Dial(url)
	return
}
