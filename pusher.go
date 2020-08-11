package main

import (
	"log"
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/streadway/amqp"
	"github.com/spf13/viper"
)

type trigger struct {
	Source string
	Endpoint string
	Sink string
}

type config struct {
	Url string
	Concurrency int // not sure if we need this in go at all
	Triggers []trigger
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
}
  
func parseConfig() (cfg config) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	failOnError(viper.ReadInConfig(), "Config could not be loaded")
	failOnError(viper.Unmarshal(&cfg), "Config could not be unmarshaled")

	return
}

func declareQueue(name string, ch *amqp.Channel) {
	if name != "" {
		_, err := ch.QueueDeclare(
			name,    // name
			true,    // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare a queue")
		log.Println("Declared queue:", name)
	}
}

func process(t trigger, ch *amqp.Channel, msg []byte) {
	r := bytes.NewReader(msg)
	resp, err := http.Post(t.Endpoint, "application/json", r)
	failOnError(err, "Failed to send json to the endpoint")
	defer resp.Body.Close()
	// todo: check status code and nack the message
	// (resp.StatusCode == http.StatusOK)
	if t.Sink != "" {
		body, err := ioutil.ReadAll(resp.Body)
		failOnError(err, "Could not read response body")
		data := amqp.Publishing {
			ContentType: "text/plain",
			Body: body,
		}
		err = ch.Publish(
			"",     // exchange
			t.Sink, // routing key
			false,  // mandatory
			false,  // immediate
			data)
		failOnError(err, "Failed to publish response")
	}
}

func listen(t trigger, ch *amqp.Channel) {
	msgs, err := ch.Consume(
		t.Source, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	failOnError(err, "Failed to register a consumer")
	log.Println("Listening on:", t.Source)
	for msg := range msgs {
		go process(t, ch, msg.Body)
	} 
}

func main() {
	cfg := parseConfig()

	conn, err := amqp.Dial(cfg.Url)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	
	for _, t := range cfg.Triggers {
		declareQueue(t.Source, ch)
		declareQueue(t.Sink, ch)
		go listen(t, ch)
	}
	for {}
}
