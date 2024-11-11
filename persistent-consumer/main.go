/*
This is a persistent simple example of a RabbitMQ consumer that listens to the changeset event.
The consumer listens to the univer-event-sync.changeset topic and prints the changeset event to the console.
*/
package main

import (
	"encoding/json"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	eventTypeChangeset = "changeset"
	topicPrefix        = "univer-event-sync."
	exchangeName       = "univer-event-sync"
)

func main() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/" // replace with your RabbitMQ server
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Dial error: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("conn error: %s", err)
	}
	defer ch.Close()

	// Univer server will declare the exchange after bootstrap, so you don't need to declare it here.
	// In certain situations, if you want to declare it in advance,
	// please ensure the parameters are consistent with the following code.
	// err = ch.ExchangeDeclare(
	// 	exchangeName, // name
	// 	"topic",      // type
	// 	true,         // durable
	// 	false,         // auto-deleted
	// 	false,        // internal
	// 	false,        // no-wait
	// 	nil,
	// )
	// if err != nil {
	// 	log.Fatalf("ExchangeDeclare error: %s", err)
	// }

	// You only need to declare it once.
	// If you modify the amqp.Table value and declare it again, an error will be returned.
	q, err := ch.QueueDeclare(
		exchangeName+"-"+eventTypeChangeset+"-"+"persistent", // Declare the queue with a fixed name. You can customize it
		true,  // need durable
		false, // close AD because msg need to persistent
		false,
		false,
		amqp.Table{
			// When the 100001th message arrives, the queue discards the oldest message according to the FIFO principle
			"x-max-length": 100000,
		},
	)
	if err != nil {
		log.Fatalf("QueueDeclare error: %s", err)
	}

	if err = ch.QueueBind(
		q.Name,
		topicPrefix+eventTypeChangeset, // now is univer-event-sync.changeset. you can use univer-event-sync.* to get all events
		exchangeName,
		false,
		nil,
	); err != nil {
		log.Fatalf("QueueBind fail: %s", err)
	}

	delivery, err := ch.Consume(
		q.Name,
		"myConsumer-persistent",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for msg := range delivery {
			var event Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Error: %s", err)
			}
			log.Printf("traceId:%s\n", msg.Headers["trace-id"])
			switch event.EventType {
			case eventTypeChangeset:
				// Do something with the changeset event
				log.Printf("Changeset event: %+v\n", event.CsAckEvent.Cs)
			default:
			}
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")

	select {}
}

type Mutation struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type Cs struct {
	UnitID    string     `json:"unitID"`
	Type      int        `json:"type"` // 1: doc, 2: sheet
	BaseRev   int        `json:"baseRev"`
	Revision  int        `json:"revision"`
	UserID    string     `json:"userID"`
	Mutations []Mutation `json:"mutations"`
	MemberID  string     `json:"memberID"`
}

type CsAckEvent struct {
	Cs Cs `json:"cs"`
}

type Event struct {
	EventID    string     `json:"eventId"`
	EventType  string     `json:"eventType"` // event type: changeset, more types in the future
	CsAckEvent CsAckEvent `json:"csAckEvent"`
}
