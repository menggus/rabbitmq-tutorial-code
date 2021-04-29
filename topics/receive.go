package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"rabbitmq_demo/mere"
)

// To receive all the logs, you can:
// go run receive.go "#"

// TO receive all logs from the facility "kern", you can:
// go run receive.go "kern.*"

// if you want to hear only about "critical" logs, you can:
// gi run receive.go "*.critical"

// you can create multiple bindings
// go run receive.go "kern.*" "*.critical"

// emit a message with a routing key "kern.critical" type:
// go run publish.go "kern.critical" "A critical kernel error"

func main() {
	// 1. 建立与 rabbitmq 的连接
	connect, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	mere.FailOnError(err, mere.CONNECT_ERROR)
	defer connect.Close()

	// 2. 建立通道
	ch, err := connect.Channel()
	mere.FailOnError(err, mere.CHANNEL_ERROR)
	defer ch.Close()

	// 3. 声明 exchange 这里使用的类型是 direct
	err = ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	mere.FailOnError(err, mere.EXCHANGE_ERROR)

	// 4. 声明 消息队列
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	mere.FailOnError(err, mere.QUEUE_CREATE_ERROR)

	// 5. 判断参数的数量
	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	// 6. 消息队列 绑定 代理
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q, "logs_direct", s)

		err = ch.QueueBind(q.Name, s, "logs_topic", false, nil)
		mere.FailOnError(err, mere.QUEUE_BIND_ERROR)
	}

	// 7. 消费消息
	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	mere.FailOnError(err, mere.CONSUME_MSG_ERROR)

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Printf("[x] %s", d.Body)
		}
	}()

	<-forever
}
