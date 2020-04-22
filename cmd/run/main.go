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
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/toskatok/lg/instance"
	"gopkg.in/yaml.v2"
)

const (
	flagDestination = "destination"
	flagRate        = "rate"
	flagConfig      = "config"
)

// ReportDuration is a duration to print results
const ReportDuration = 1 * time.Second

// action runs a test instance
func main(rate time.Duration, destination string, config string) error {
	// cfg variable contains current user configuration
	var cfg instance.Config

	f, err := os.Open(config)
	if err != nil {
		return err
	}

	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return err
	}

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
			destination, err := cmd.Flags().GetString(flagDestination)
			if err != nil {
				return err
			}
			rate, err := cmd.Flags().GetDuration(flagRate)
			if err != nil {
				return err
			}
			config, err := cmd.Flags().GetString(flagConfig)
			if err != nil {
				return err
			}

			return main(rate, destination, config)
		},
	}

	cmd.Flags().String(flagDestination, "mqtt://127.0.0.1:1883", "scheme://(host or host:port)")
	cmd.Flags().Duration(flagRate, time.Second, "send interval")
	cmd.Flags().String(flagConfig, "config.yml", "load instance configuration")

	root.AddCommand(cmd)
}
