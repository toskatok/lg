/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 20-11-2017
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

func main() {
	// Flags
	var rate = flag.Int64("rate", 1, "Sends one packet each ? millisecond")
	var broker = flag.String("broker", "127.0.0.1:1883", "MQTT Broker IP:Port address")
	var topic = flag.String("topic", "application/app/node/n/rx", "Publish on ? topic")
	flag.Parse()

	// Read message
	message, err := ioutil.ReadFile("message.bin")
	if err != nil {
		panic(err)
	}

	// Create an MQTT Client.
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			log.Println(err)
		},
	})

	// Connect to the MQTT Server.
	err = cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  *broker,
		ClientID: []byte(fmt.Sprintf("isrc-lg-%d", rand.Int63())),
	})
	if err != nil {
		panic(err)
	}

	// Tick Tick
	tick := time.Tick(time.Duration(*rate) * time.Millisecond)

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	for {
		select {
		case <-tick:
			err = cli.Publish(&client.PublishOptions{
				QoS:       mqtt.QoS0,
				TopicName: []byte(*topic),
				Message:   message,
			})
			if err != nil {
				log.Printf("MQTT Publish: %s", err)
			}
		case <-sigc:
			// Disconnect the Network Connection.
			cli.Disconnect()
			return
		}
	}
}
