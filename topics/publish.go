package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"rabbitmq_demo/mere"
	"strings"
)

func main() {
	// 1. 建立与 rabbitmq 的连接
	connect, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	mere.FailOnError(err, mere.CONNECT_ERROR)
	defer connect.Close()

	// 2. 建立通道
	ch, err := connect.Channel()
	mere.FailOnError(err, mere.CHANNEL_ERROR)
	defer ch.Close()

	// 3. 声明代理 exchange 这次使用的类型是 direct
	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable  持久
		false,        // auto-deleted 自动删除
		false,        // internal 内部
		false,        // no-wait
		nil,          // arguments
	)
	mere.FailOnError(err, mere.EXCHANGE_ERROR)

	// 4. 发布消息
	body := bodyForm(os.Args)
	err = ch.Publish(
		"logs_topic",          // exchange name
		severityForm(os.Args), // route key
		false,                 // mandatory 强制性
		false,                 // immediate  立即
		amqp.Publishing{ // 发布消息
			ContentType: "text/plain", // 格式
			Body:        []byte(body), // 内容
		})
	mere.FailOnError(err, mere.QUEUE_CREATE_ERROR)

	log.Printf(" [x] Sent %s start ...... \n", body)
}

// bodyForm 构建发布的 消息体
func bodyForm(args []string) string {
	var s string

	if len(args) < 3 || args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

// severityForm 生成发布消息的 严重程度 warning  info  error等
func severityForm(args []string) string {
	var s string
	if len(args) < 2 || args[1] == "" {
		s = "anonymous.info"
	} else {
		s = args[1]
	}
	return s
}
