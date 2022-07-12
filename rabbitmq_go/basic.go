package rabbitmq_go

import (
	"encoding/json"
	"github.com/Pangjiping/go-example/rabbitmq_go/lib"
	"github.com/streadway/amqp"
	"log"
)

/*
简单队列模式是RabbitMQ的常规用法，简单理解就是消息生产者发送消息给一个队列，然后消息的消费者从队列中读取消息
当多个消费者订阅同一个队列的时候，队列中的消息是平均分摊给多个消费者处理
*/

// =============================生产者=========================================//
type simpleDemo struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

func produce() {
	// 连接rabbitMQ服务器
	conn, err := RabbitMQConn()
	lib.ErrorHandle(err, lib.ErrConnectRabbit)
	defer conn.Close()

	// 新建一个channel
	ch, err := conn.Channel()
	lib.ErrorHandle(err, lib.ErrOpenChannel)
	defer ch.Close()

	// 声明或者创建一个队列来保存消息
	q, err := ch.QueueDeclare(
		"simple:queue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // argument
	)
	lib.ErrorHandle(err, lib.ErrDeclareQueue)
	data := simpleDemo{
		Name: "Tom",
		Addr: "Shanghai",
	}
	dataBytes, err := json.Marshal(data)
	lib.ErrorHandle(err, lib.ErrMarshalJSON)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        dataBytes,
		},
	)
	lib.ErrorHandle(err, lib.ErrPublishMsg)
	log.Printf(" [x] Sent %s", dataBytes)
}

//==================================消费者======================================//
func consume() {
	conn, err := RabbitMQConn()
	lib.ErrorHandle(err, lib.ErrConnectRabbit)
	defer conn.Close()

	ch, err := conn.Channel()
	lib.ErrorHandle(err, lib.ErrOpenChannel)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"simple:queue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // args
	)
	lib.ErrorHandle(err, lib.ErrDeclareQueue)
	// 定义一个消费者
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	lib.ErrorHandle(err, lib.ErrRegisterConsumer)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	select {}
}
