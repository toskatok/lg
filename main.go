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
	"strings"
	"time"

	"github.com/I1820/lg/core"
	"github.com/I1820/lg/generators"
	"github.com/I1820/lg/receivers"
	"github.com/alecthomas/template"
	"github.com/urfave/cli"
)

type generator int

func (g *generator) Set(v string) error {
	switch v {
	case "aolab":
		*g = 1
	case "isrc":
		*g = 0
	case "atrovan":
		*g = 2
	case "fanco":
		*g = 3
	default:
		return fmt.Errorf("the %s generator hasn't been implemented yet", v)
	}
	return nil
}

func (g *generator) String() string {
	return string(*g)
}

// message is used to populate templates in the given message file.
var message struct {
	Count int64
}

func main() {
	// dirname
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:        "MQTT-LG",
		Description: "MQTT based Load Generator",
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
				Value: 1 * time.Millisecond,
				Usage: "Send interval",
			},
			&cli.StringFlag{
				Name:  "message",
				Value: path.Join(dir, "message.json"),
				Usage: "Raw message file path (relative to lg directory)",
			},
			&cli.StringFlag{
				Name:  "deveui",
				Value: "",
				Usage: "DevEUI [Device Unique Identifier]",
			},
			&cli.GenericFlag{
				Name:  "generator",
				Value: new(generator),
				Usage: "Generator [isrc, aolab, atrovan]",
			},
			&cli.StringFlag{
				Name:  "i1820",
				Value: "",
				Usage: "I1820 project name [Use this only with the I1820 platform otherwise leave it empty]",
			},
		},
		Action: func(c *cli.Context) error {
			// DevEUI
			devEUI := c.String("deveui")
			fmt.Println(">>> Device")
			fmt.Println(devEUI)
			fmt.Println(">>>")

			// Read message.json
			file, err := ioutil.ReadFile(c.String("message"))
			if err != nil {
				return err
			}

			// Create parent template with some useful function
			tmpl := template.New("lg").Funcs(template.FuncMap{
				"now":   time.Now,
				"randn": rand.Intn,
			})

			// Read data from the given message file, and then prase template strings if exist.
			var data []map[string]interface{}
			if err := json.Unmarshal(file, &data); err != nil {
				return err
			}
			fmt.Println(">>> Data")
			for i, d := range data {
				fmt.Printf("%v\n", d)
				for k, v := range d {
					if s, ok := v.(string); ok {
						p, err := tmpl.New(fmt.Sprintf("lg-%d-%s", i, k)).Parse(s)
						if err != nil {
							return err
						}
						d[k] = p
					}
				}
			}
			fmt.Println(">>>")

			// Generator selection
			var g generators.Generator
			switch *c.Generic("generator").(*generator) {
			case 0: // ISRC
				g = generators.ISRCGenerator{
					DevEUI:          devEUI,
					ApplicationName: "fake-application",
					ApplicationID:   "13731372",
					GatewayMac:      "b827ebffff633260",
					DeviceName:      "fake-device",
				}
			case 1: // Aolab
				g = generators.AolabGenerator{
					DevEUI: devEUI,
				}
			case 2: // Atrovan
				g = generators.AtrovanGenerator{}
			case 3: // Fanco
				g = generators.FancoGenerator{
					ThingID: devEUI,
				}
			}

			// I1820 mode
			var rs []receivers.Receiver
			var counter int
			if project := c.String("i1820"); project != "" {
				rs = append(rs, receivers.Receiver{
					Topic: []byte(fmt.Sprintf("i1820/project/%s/data", project)),
					Handler: func(topicName, message []byte) {
						counter++
					},
				})
			}

			// Runner
			r, err := core.NewRunner(
				g,
				c.Duration("rate"),
				func() interface{} {
					message.Count++
					d := make(map[string]interface{})
					for k, v := range data[rand.Intn(len(data))] {
						if tmpl, ok := v.(*template.Template); ok {
							var b strings.Builder
							if err := tmpl.Execute(&b, message); err != nil {
								continue
							}
							d[k] = b.String()
						} else {
							d[k] = v
						}
					}
					return d
				},
				c.String("broker"),
				rs...,
			)
			if err != nil {
				return err
			}
			r.Run()

			// Report
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
			r.Stop()
			fmt.Printf("Total packets send to %s: %d\n", devEUI, r.Count())
			fmt.Printf("Total packets receive from I1820: %d\n", counter)

			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
