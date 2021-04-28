package main

import (
	"github.com/streadway/amqp"
	"log"
	"rabbitmq_demo/mere"
)

func main() {
	// 1. 建立连接
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	mere.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	// 2. 创建通道
	ch, err := conn.Channel()
	mere.FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	// 3. 声明代理，名 logs  类型为 fanout 扇出 以扇形方式发出
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	mere.FailOnError(err, "Failed to declare an exchange")

	// 4. 声明队列
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	mere.FailOnError(err, "Failed to declare a queue")

	// 5. 绑定队列
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	mere.FailOnError(err, "Failed to bind a queue")

	// 6. 消费消息
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	mere.FailOnError(err, "Failed to register a consumer")

	// 7. 通过无缓冲通道，完成 夯住进程，防止退出
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
