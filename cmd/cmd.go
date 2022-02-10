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
			backfillCmd,
			dbCmd,
			fixturesCmd,
			integrationsCmd,
			organizationsCmd,
			servicesCmd,
			subscriptionsCmd,
			updateCmd,
			{
				Name: "version",
				Action: func(c *cli.Context) error {
					shaPart := config.BuildSha
					if len(shaPart) >= 8 {
						shaPart = fmt.Sprintf(" (%s)", config.BuildSha[0:8])
					}
					fmt.Fprintf(os.Stdout, "%s%s\n", config.Version, shaPart)
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
