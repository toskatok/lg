/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 26-09-2018
 * |
 * | File Name:     mqtt.go
 * +===============================================
 */

package core

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// MQTTTransport implements transport interface for mqtt protocol
type MQTTTransport struct {
	cli *client.Client
}

// Init creates and connect mqtt client
func (mt MQTTTransport) Init(url string) error {
	// Create an MQTT Client.
	mt.cli = client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			log.Println(err)
		},
	})

	// Connect to the MQTT Server.
	return mt.cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  url,
		ClientID: []byte(fmt.Sprintf("I1820-lg-%d", rand.Intn(1024))),
	})

	return nil
}

// Transmit sends data on given mqtt topic
func (mt MQTTTransport) Transmit(topic string, data []byte) error {
	return mt.cli.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte(topic),
		Message:   data,
	})
}
