package mere

import "log"

const (
	CONNECT_ERROR = "Failed create connect"
	CHANNEL_ERROR = "Failed create channel"
	EXCHANGE_ERROR = "Failed create exchange"
	PUBLISH_ERROR = "Failed create publish"
	QUEUE_CREATE_ERROR = "Failed create QUEUE"
	QUEUE_BIND_ERROR = "Failed BIND QUEUE"
	CONSUME_MSG_ERROR = "Failed consume msg error"
)


//FailOnError 错误处理
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
