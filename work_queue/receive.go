package main

import (
	"github.com/streadway/amqp"
	"log"
	"rabbitmq_demo/mere"
)


// InitRabbitmq 初始化 rabbitmq 连接
func main() {
	// 1. 创建连接
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	mere.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 2. 打开一个通道
	ch, err := conn.Channel()
	mere.FailOnError(err, "Failed to open a channel")

	// 3. 创建消息队列
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments)
	)
	mere.FailOnError(err, "Failed to declare a queue")

	// 4. 创建消费者
	cs, err := ch.Consume(q.Name, "", false, false,false, false, nil)
	mere.FailOnError(err, "Failed create Consume")

	forever := make(chan bool)  // 用于夯住进程，不会退出

	go func() {
		for d := range cs {
			log.Printf("Reveived a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<- forever // 用于夯住进程，不会退出
}