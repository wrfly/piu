package main

import (
	"os"

	"gopkg.in/urfave/cli.v2"
)

var appName = "template-app"

func main() {
	app := &cli.App{
		Name:        appName,
		Usage:       "Some template application",
		Authors:     author,
		Version:     simpleVersionInfo,
		HideVersion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "docker",
				Usage:   "docker socket path",
				Aliases: []string{"socket"},
				Value:   "/var/run/docker.sock",
			},
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "print the version",
			},
			&cli.BoolFlag{
				Name: "debug",
			},
		},
		CustomAppHelpTemplate: helpTemplate,
		Action:                run,
	}

	app.Run(os.Args)
}
