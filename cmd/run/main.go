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

package run

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toskatok/lg/instance"
)

const (
	flagDestination = "destination"
	flagRate        = "rate"
)

// config loads configuration from file and return it
func config() instance.Config {
	var config instance.Config

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
func main(rate time.Duration, destination string) error {
	// cfg variable contains current user configuration
	cfg := config()

	i, err := instance.New(cfg, rate, destination)
	if err != nil {
		return err
	}

	// prints generator information
	color.Blue(">>> Generator")
	color.Yellow("%+v\n", i.R.Generator)
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

func Register(root *cobra.Command) {
	cmd := &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			destination, err := cmd.Flags().GetString("destination")
			if err != nil {
				return err
			}
			rate, err := cmd.Flags().GetDuration("rate")
			if err != nil {
				return err
			}

			return main(rate, destination)
		},
	}

	cmd.Flags().String(flagDestination, "mqtt://127.0.0.1:1883", "scheme://(host or host:port)")
	cmd.Flags().Duration(flagRate, time.Second, "send interval")

	root.AddCommand(cmd)
}
