package cmd

import (
	"fmt"
	"github.com/lithictech/webhookdb-cli/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func Execute() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Value: false},
		},
		Commands: []*cli.Command{
			authCmd,
			integrationsCmd,
			servicesCmd,
			{
				Name: "version",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(os.Stdout, config.BuildSha[0:8])
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
