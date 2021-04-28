package main

import (
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"rabbitmq_demo/mere"
	"strconv"
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
	_, err = ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments)
	)
	mere.FailOnError(err, "Failed to declare a queue")


	body := "Hello World!"

	for {
		number := rand.Intn(10)
		err = ch.Publish(
			"",     	// exchange
			"hello", 		// routing key
			false,  	// mandatory
			false,  	// immediate
			amqp.Publishing {
				ContentType: "text/plain",
				Body:        []byte(body+strconv.Itoa(number)),
			})
		mere.FailOnError(err, "Failed to publish a message")
		//time.Sleep(time.Second * 2)
		log.Printf("%s-----%s\n", body, strconv.Itoa(number))
	}
}

