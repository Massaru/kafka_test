package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	start := time.Now()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "kafka",
		"group.id":           "myGroup",
		"enable.auto.commit": false,
		"auto.offset.reset":  "earliest",
	})

	if err != nil {
		panic(err)
	}

	duration := time.Since(start)
	fmt.Println("NewConsumer Duração:", duration)
	start = time.Now()

	c.SubscribeTopics([]string{"^DT0*"}, nil)

	duration = time.Since(start)
	fmt.Println("SubscribeTopics Duração:", duration)
	start = time.Now()

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%%%s Message on %s:\n%s\n",
					time.Now(), e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
				duration = time.Since(start)
				fmt.Println("ReadMessage Duração:", duration)
				start = time.Now()
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()

	// for {
	// 	c.Poll(500)
	// 	msg, err := c.ReadMessage(-1)
	// 	if err == nil {
	// 		duration = time.Since(start)
	// 		fmt.Println("ReadMessage Duração:", duration)
	// 		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
	// 		start = time.Now()
	// 	} else {
	// 		// The client will automatically try to recover from all errors.
	// 		fmt.Printf("Consumer error: %v (%v)\n", err, msg)
	// 	}
	// }

	// // c.Close()
}
