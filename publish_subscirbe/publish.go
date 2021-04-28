package main

import (
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"os"
	"rabbitmq_demo/mere"
	"strconv"
	"strings"
)

func main() {
	// 1. 建立连接
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	mere.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	// 2. 建立通道
	ch, err := conn.Channel()
	mere.FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	// 3. 声明一个代理 名为 logs
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

	// 4. 准备发布内容
	body := bodyFrom(os.Args)

	// 5. 发布消息
	for {
		err = ch.Publish(
			"logs", // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body + strconv.Itoa(rand.Intn(10))),
			})
		mere.FailOnError(err, "Failed to publish a message")
	}


	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}