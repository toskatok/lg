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
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/2tvenom/cbor"
	"github.com/aiotrc/uplink/lora"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type devEUI int64

func (d *devEUI) Set(v string) (err error) {
	i, err := strconv.ParseInt(v, 16, 64)
	*d = devEUI(i)

	return
}

func (d *devEUI) String() string {
	return "10"
}

func main() {
	// Flags
	var rate = flag.Int64("rate", 1000, "Sends one packet each ? millisecond")
	var broker = flag.String("broker", "127.0.0.1:1883", "MQTT Broker IP:Port address")
	var devID devEUI
	flag.Var(&devID, "deveui", "Device EUI")
	flag.Parse()

	// DevEUI
	devEUI := fmt.Sprintf("%016X", int64(devID))
	fmt.Println(devEUI)

	// Read message
	data, err := ioutil.ReadFile("message.json")
	if err != nil {
		panic(err)
	}

	var info []map[string]int
	if err := json.Unmarshal(data, &info); err != nil {
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
	if err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  *broker,
		ClientID: []byte(fmt.Sprintf("isrc-lg-%d", rand.Int63())),
	}); err != nil {
		panic(err)
	}

	// Tick Tick
	sendTick := time.Tick(time.Duration(*rate) * time.Millisecond)
	printTick := time.Tick(1 * time.Second)

	// Packet counter
	packets := 0

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	for {
		select {
		case <-sendTick:
			var buffer bytes.Buffer
			encoder := cbor.NewEncoder(&buffer)
			if ok, err := encoder.Marshal(info[rand.Intn(len(info))]); !ok {
				log.Printf("CBor Encoding: %s", err)
				continue
			}

			message, err := json.Marshal(lora.RxMessage{
				ApplicationID:   "1",
				ApplicationName: "isrc-platform",
				DeviceName:      "isrc-sensor",
				DevEUI:          devEUI,
				FPort:           5,
				FCnt:            10,
				RxInfo: []lora.RxInfo{
					lora.RxInfo{
						Mac:     "b827ebffff633260",
						Name:    "isrc-gateway",
						Time:    time.Now(),
						RSSI:    -57,
						LoRaSNR: 10,
					},
				},
				TxInfo: lora.TxInfo{
					Frequency: 868100000,
					Adr:       true,
					CodeRate:  "4/6",
				},
				Data: buffer.Bytes(),
			})
			if err != nil {
				log.Printf("JSON Encoding: %s", err)
				continue
			}

			if err := cli.Publish(&client.PublishOptions{
				QoS:       mqtt.QoS0,
				TopicName: []byte("application/app/node/n/rx"),
				Message:   message,
			}); err != nil {
				log.Printf("MQTT Publish: %s", err)
			}
			packets++
		case <-printTick:
			fmt.Printf("Sends %d packets\n", packets)
		case <-sigc:
			// Disconnect the Network Connection.
			if err := cli.Disconnect(); err != nil {
				return
			}
			return
		}
	}
}
