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
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/toskatok/lg/models"
	"github.com/urfave/cli"
)

// config loads configuration from file and return it
func config() models.Config {
	var config models.Config

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // looking for config in the working directory

	if err := viper.ReadInConfig(); err != nil { // find and read the config file
		log.Fatalf("config file: %s \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil { // parse configuration into Config structure
		log.Fatal(err)
	}

	return config
}

// ReportDuration is a duration to print results
const ReportDuration = 1 * time.Second

// action runs a test instance
func action(c *cli.Context) error {
	// cfg variable contains current user configuration
	cfg := config()

	i, err := models.NewInstance(cfg, c.Duration("rate"), c.String("destination"))
	if err != nil {
		return err
	}

	// prints generator information
	color.Blue(">>> Generator")
	color.Yellow("%+v\n", i.G)
	color.Blue(">>>")

	// runs the instance
	i.Run()

	// prints report in 1 second intervals
	go func() {
		for range time.Tick(ReportDuration) {
			fmt.Print(color.CyanString("%s", cfg.Generator.Name))
			fmt.Print(color.GreenString(" -> generates"))
			fmt.Print(color.CyanString(" %d ", i.R.Count()))
			fmt.Print(color.GreenString("number of packets\n"))
		}
	}()

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	<-sigc
	i.Stop()
	color.HiRed("Total packets send to %s: %d\n", cfg.Generator.Name, i.R.Count())

	return nil
}

func main() {
	app := &cli.App{
		Name:        "LG",
		Description: "Load Generator",
		Authors: []cli.Author{
			{
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
				Name:     "rate",
				Usage:    "Send interval",
				Required: true,
			},
		},
		Action: action,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
