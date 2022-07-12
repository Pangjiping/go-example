package rabbitmq_go

import (
	"github.com/Pangjiping/go-example/rabbitmq_go/lib"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

/*
	工作队列也被称为任务队列
	任务队列是为了避免等待执行一些耗时的任务，而是将需要执行的任务封装为消息发送给工作队列，后台运行的工作进程将任务消息取出来并执行相关任务
	多个后台工作进程同时进行，他们之间共享任务（抢占）
*/

//===========================任务生产者=================================//
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "no task"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func taskGen() {
	conn, err := RabbitMQConn()
	lib.ErrorHandle(err, lib.ErrConnectRabbit)
	defer conn.Close()

	ch, err := conn.Channel()
	lib.ErrorHandle(err, lib.ErrOpenChannel)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task:queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // args
	)
	lib.ErrorHandle(err, lib.ErrDeclareQueue)

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         []byte(body),
		},
	)
	lib.ErrorHandle(err, lib.ErrPublishMsg)
	log.Printf("sent %s", body)
}

//===========================工人======================================//
func worker() {
	conn, err := RabbitMQConn()
	lib.ErrorHandle(err, lib.ErrConnectRabbit)
	defer conn.Close()

	ch, err := conn.Channel()
	lib.ErrorHandle(err, lib.ErrOpenChannel)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task:queue",
		false,
		false,
		false,
		false,
		nil)
	lib.ErrorHandle(err, lib.ErrDeclareQueue)

	// 将预取计数器设置为1
	// 在并行处理中将消息分配给不同的工作进程
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	lib.ErrorHandle(err, lib.ErrSetQoS)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	lib.ErrorHandle(err, lib.ErrRegisterConsumer)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
