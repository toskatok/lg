package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/toskatok/lg/cmd/run"
	"github.com/toskatok/lg/cmd/server"

	"github.com/spf13/cobra"
)

// ExitFailure status code
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var root = &cobra.Command{
		Use:   "lg",
		Short: "Fully customizable load generator for MQTT, HTTP, etc.",
	}

	root.Println("13 Feb 2020, Best Day Ever")

	run.Register(root)
	server.Register(root)

	if err := root.Execute(); err != nil {
		logrus.Errorf("failed to execute root command: %s", err.Error())
		os.Exit(ExitFailure)
	}
}
