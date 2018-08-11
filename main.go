/*
 *
 * In The Name of God
 *
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/I1820/lg/generators"
	"github.com/urfave/cli"
	"github.com/yosssi/gmq/mqtt/client"
)

type devEUI struct {
	v int64
}

func (d *devEUI) Set(v string) (err error) {
	i, err := strconv.ParseInt(v, 16, 64)
	d.v = i

	return
}

func (d *devEUI) String() string {
	return fmt.Sprintf("%016X", d.v)
}

func main() {
	// dirname
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:        "MQTT-LG",
		Description: "MQTT-LoRaServer.io Load Generator",
		Authors: []cli.Author{
			cli.Author{
				Name:  "Parham Alvani",
				Email: "parham.alvani@gmail.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "broker",
				Value: "127.0.0.1:1883",
				Usage: "MQTT `Broker` IP:Port address",
			},
			&cli.DurationFlag{
				Name:  "rate",
				Value: 1000,
				Usage: "Sends one packet each `?` ns",
			},
			&cli.StringFlag{
				Name:  "message",
				Value: path.Join(dir, "message.json"),
				Usage: "raw information file path",
			},
			&cli.GenericFlag{
				Name:  "deveui",
				Value: &devEUI{0x73},
				Usage: "DevEUI",
			},
		},
		Action: func(c *cli.Context) error {
			// DevEUI
			devEUI := fmt.Sprintf("%016X", c.Generic("deveui").(*devEUI).v)
			fmt.Println(">>> Device")
			fmt.Println(devEUI)
			fmt.Println(">>>")

			// Read message
			file, err := ioutil.ReadFile(c.String("message"))
			if err != nil {
				return err
			}

			var data []map[string]int
			if err := json.Unmarshal(file, &data); err != nil {
				return err
			}
			fmt.Println(">>> Data")
			for _, d := range data {
				fmt.Printf("%v\n", d)
			}
			fmt.Println(">>>")

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
				Address:  c.String("broker"),
				ClientID: []byte(fmt.Sprintf("I1820-lg-%s", devEUI)),
			}); err != nil {
				return err
			}

			r := generators.NewRunner(
				generators.LoRaApplicationGenerator{
					DevEUI:          devEUI,
					ApplicationName: "fake-application",
					ApplicationID:   "13731372",
					GatewayMac:      "b827ebffff633260",
					DeviceName:      "fake-device",
				},
				c.Duration("rate"),
				func() interface{} {
					return data[rand.Intn(len(data))]
				},
				cli,
			)
			r.Run()

			go func() {
				for {
					time.Sleep(1 * time.Second)
					fmt.Printf("%s -> %d\n", devEUI, r.Count())
				}
			}()

			// Set up channel on which to send signal notifications.
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, os.Interrupt, os.Kill)

			<-sigc

			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
