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
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"text/template"
	"time"

	"github.com/I1820/lg/core"
	"github.com/I1820/lg/generators"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// Config is a configuration structure for I1820/lg
type Config struct {
	Generator struct {
		Name string
		Info interface{} // this structure is not used in config, it is passed to generators
	}
	Token    string
	Messages []map[string]interface{}
}

// config variable contains current user configuration
var config Config

// message is used to populate templates in the given message file.
var message struct {
	Count int64
}

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // looking for config in the working directory

	if err := viper.ReadInConfig(); err != nil { // find and read the config file
		log.Fatalf("config file: %s \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil { // parse configuration into Config structure
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
				Name:  "destination",
				Value: "mqtt://127.0.0.1:1883",
				Usage: "scheme://(host or host:port)",
			},
			&cli.DurationFlag{
				Name:  "rate",
				Value: 1 * time.Millisecond,
				Usage: "Send interval",
			},
		},
		Action: func(c *cli.Context) error {
			// Create parent template with some useful function
			tmpl := template.New("lg").Funcs(template.FuncMap{
				"now":   time.Now,
				"randn": rand.Intn,
			})

			// Read data from the given message file, and then prase template strings if exist.
			fmt.Println(">>> Data")
			for i, d := range config.Messages {
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

			// Generator selection and configuration
			var g generators.Generator
			switch config.Generator.Name {
			case "isrc": // generators/isrc.go
				var isrc generators.ISRCGenerator
				// load genrator information from configuration file
				if err := viper.UnmarshalKey("generator.info", &isrc); err != nil {
					return err
				}
				g = isrc
			case "aolab": // generators/aolab.go
				var aolab generators.AolabGenerator
				// load genrator information from configuration file
				if err := viper.UnmarshalKey("generator.info", &aolab); err != nil {
					return err
				}
				g = aolab
			case "atrovan": // generators/atrovan.go
				var atrovan generators.AtrovanGenerator
				// load genrator information from configuration file
				if err := viper.UnmarshalKey("generator.info", &atrovan); err != nil {
					return err
				}
				g = atrovan
			case "fanco": // generators/fanco.go
				var fanco generators.FancoGenerator
				// load genrator information from configuration file
				if err := viper.UnmarshalKey("generator.info", &fanco); err != nil {
					return err
				}
				g = fanco
			case "ttn": // generators/ttn.go
				var ttn generators.TTNGenerator
				// load genrator information from configuration file
				if err := viper.UnmarshalKey("generator.info", &ttn); err != nil {
					return err
				}
				g = ttn
			default:
				return fmt.Errorf("Generator %s is not supported yet", "")
			}
			fmt.Println(">>> Generator")
			fmt.Printf("%+v\n", g)
			fmt.Println(">>>")

			// Runner
			r, err := core.NewRunner(core.RunnerConfig{
				Generator: g,
				Duration:  c.Duration("rate"),
				Pick: func() interface{} {
					message.Count++
					d := make(map[string]interface{})
					for k, v := range config.Messages[rand.Intn(len(config.Messages))] {
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
				URL:   c.String("destination"),
				Token: config.Token,
			})
			if err != nil {
				return err
			}
			r.Run()

			// Report
			go func() {
				for {
					time.Sleep(1 * time.Second)
					fmt.Printf("%s -> %d\n", config.Generator.Name, r.Count())
				}
			}()

			// Set up channel on which to send signal notifications.
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, os.Interrupt, os.Kill)

			<-sigc
			r.Stop()
			fmt.Printf("Total packets send to %s: %d\n", config.Generator.Name, r.Count())

			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
