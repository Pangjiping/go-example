package lib

import "log"

const (
	ErrConnectRabbit    = "Failed to connect to RabbitMQ"
	ErrOpenChannel      = "Failed to open a channel"
	ErrDeclareQueue     = "Failed to declare a queue"
	ErrPublishMsg       = "Failed to publish a message"
	ErrRegisterConsumer = "Failed to register a consumer"
	ErrMarshalJSON      = "Failed to marshal json"
	ErrSetQoS           = "Failed to set QoS"
)

// ErrorHandle 错误处理函数
func ErrorHandle(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}
