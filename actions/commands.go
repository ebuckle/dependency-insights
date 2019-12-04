package actions

import (
	"os"

	"github.com/urfave/cli"
)

// Commands creates the CLI commands
func Commands() {
	app := cli.NewApp()
	app.Name = "dependency-insights"
	app.Version = "dev"
	app.Usage = "Analyse and report on project dependencies"

	app.Commands = []cli.Command{
		{
			Name:  "local",
			Usage: "analyse a locally saved project",

			Flags: []cli.Flag{
				cli.StringFlag{Name: "path, p", Usage: "The path to the project", Required: true},
				cli.StringFlag{Name: "language, l", Usage: "The project language", Required: true},
			},
			Action: func(c *cli.Context) error {
				InsightsLocalProject(c)
				return nil
			},
		},
		{
			Name:  "docker",
			Usage: "analyse a containerized project",

			Flags: []cli.Flag{
				cli.StringFlag{Name: "conid", Usage: "The container id of the project", Required: true},
				cli.StringFlag{Name: "language, l", Usage: "The project language", Required: true},
			},
			Action: func(c *cli.Context) error {
				InsightsDockerProject(c)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	println(err)
}
