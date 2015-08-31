package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/moul/dockerpatch"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = "Manfred Touron"
	app.Email = "https://github.com/moul/dockerpatch"
	app.Version = "0.1.0"
	app.Usage = "Read, write, manipulate, convert & apply filters to Dockerfiles"

	app.Commands = []cli.Command{
		{
			Name:        "patch",
			Usage:       "Patch a Dockerfile with filters",
			Description: "patch [--to-arm]",
			Action:      CmdPatch,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "to-arm",
					Usage: "Convert Dockerfile to armhf architecture",
				},
				cli.BoolFlag{
					Name:  "disable-network",
					Usage: "Remove network rules",
				},
			},
		},
	}
	app.Run(os.Args)
}

func CmdPatch(c *cli.Context) {
	if len(c.Args()) == 0 {
		cli.ShowSubcommandHelp(c)
		os.Exit(1)
	}

	var input io.Reader
	var err error
	path := c.Args()[0]
	switch path {
	case "-":
		input = os.Stdin
	default:
		input, err = os.Open(path)
		if err != nil {
			logrus.Fatalf("os.Open failed: %v", err)
		}
	}
	dockerfile, err := dockerpatch.DockerfileRead(input)
	if err != nil {
		logrus.Fatalf("dockerpatch.DockerfileRead failed: %v", err)
	}

	if c.Bool("to-arm") {
		if dockerfile.FilterToArm("armhf") != nil {
			logrus.Fatalf("dockerfile.FilterToArm failed: %v", err)
		}
	}

	if c.Bool("disable-network") {
		if dockerfile.FilterDisableNetwork() != nil {
			logrus.Fatalf("dockerfile.FilterDisableNetwork failed: %v", err)
		}
	}

	fmt.Println(dockerfile.String())
}
