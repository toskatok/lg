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
	"os"

	"github.com/abiosoft/ishell"
	"github.com/aiotrc/pm/client"
	"github.com/fatih/color"
)

type config struct {
	PM struct {
		URL string
	}

	Broker struct {
	}
}

func main() {
	pmClient := client.New("http://127.0.0.1:8080")

	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// display welcome info.
	shell.Println("MQTT Load Generator (parham.alvani@gmail.com)")

	shell.SetPrompt(fmt.Sprintf("%v %v ", color.RedString("mqttlg"), color.GreenString(">>>")))

	shell.AddCmd(&ishell.Cmd{
		Name: "about",
		Help: "about",
		Func: func(c *ishell.Context) {
			c.Println("18.20 is leaving us")
		},
	})

	projectCmd := &ishell.Cmd{
		Name: "projects",
	}
	shell.AddCmd(projectCmd)

	projectCmd.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "show [project_id]",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Err(fmt.Errorf("show [project_id]"))
				return
			}
			name := c.Args[0]

			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			project, err := pmClient.ProjectsShow(name)
			c.ProgressBar().Stop()

			if err != nil {
				c.Err(err)
				return
			}
			c.Printf("%+v\n", project)
		},
	})

	projectCmd.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: "create [project_id]",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Err(fmt.Errorf("create [project_id]"))
				return
			}
			name := c.Args[0]

			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			project, err := pmClient.ProjectsCreate(name)
			c.ProgressBar().Stop()

			if err != nil {
				c.Err(err)
				return
			}
			c.Printf("%+v\n", project)
		},
	})

	projectCmd.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: "delete [project_id]",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Err(fmt.Errorf("delete [project_id]"))
				return
			}
			name := c.Args[0]

			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			project, err := pmClient.ProjectsDelete(name)
			c.ProgressBar().Stop()

			if err != nil {
				c.Err(err)
				return
			}
			c.Printf("%+v\n", project)
		},
	})

	projectCmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "list",
		Func: func(c *ishell.Context) {
			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			projects, err := pmClient.ProjectsList()
			c.ProgressBar().Stop()

			if err != nil {
				c.Err(err)
				return
			}

			for _, project := range projects {
				c.Printf("%+v\n", project)
			}
		},
	})

	generatorCmd := &ishell.Cmd{
		Name: "generator",
	}
	shell.AddCmd(generatorCmd)

	generatorCmd.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: "create lora []",
	})

	// run shell
	if len(os.Args) > 1 && os.Args[1] == "--" {
		shell.Process(os.Args[2:]...)
	} else {
		// start shell
		shell.Run()
	}
}
