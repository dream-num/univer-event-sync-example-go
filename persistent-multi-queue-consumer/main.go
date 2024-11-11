/*
This is a simple example of a RabbitMQ multi-consumer listening to changeset events.
Multiple consumer bindings on multiple queues.
Each queue can consume the full amount of data in the exchange
*/
package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	eventTypeChangeset = "changeset"
	topicPrefix        = "univer-event-sync."
	exchangeName       = "univer-event-sync"
	numQueues          = 3
)

func main() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
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

	var wg sync.WaitGroup
	for i := 0; i < numQueues; i++ {
		wg.Add(1)
		go func(queueID int) {
			defer wg.Done()

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
			queueName := exchangeName + "-" + eventTypeChangeset + "-" + strconv.Itoa(queueID)
			q, err := ch.QueueDeclare(
				queueName, // Declare the queue with a fixed name. You can customize it
				true,      // need durable
				false,     // close AD because msg need to persistent
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
				topicPrefix+eventTypeChangeset,
				exchangeName,
				false,
				nil,
			); err != nil {
				log.Fatalf("QueueBind fail: %s", err)
			}

			customerName := "myConsumer" + strconv.Itoa(queueID)
			delivery, err := ch.Consume(
				q.Name,
				customerName,
				true,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				log.Fatal(err)
			}

			for msg := range delivery {
				var event Event
				if err := json.Unmarshal(msg.Body, &event); err != nil {
					log.Printf("Error: %s", err)
				}
				log.Printf("Consume %d - traceId:%s\n", queueID, msg.Headers["trace-id"])
				switch event.EventType {
				case eventTypeChangeset:
					log.Printf("Consume %d - Changeset event: %+v\n", queueID, event.CsAckEvent.Cs)
				default:
				}
			}
		}(i)
	}

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	wg.Wait()
}

type Mutation struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type Cs struct {
	UnitID    string     `json:"unitID"`
	Type      int        `json:"type"`
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
	EventType  string     `json:"eventType"`
	CsAckEvent CsAckEvent `json:"csAckEvent"`
}
