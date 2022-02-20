package cmd

import (
	"github.com/urfave/cli/v2"
)

func BuildApp() *cli.App {
	app := &cli.App{
		Usage: `CLI for the WebhookDB (https://webhookdb.com) application. WebhookDB allows you
to query any API in real-time with SQL.

To create an account and get started, run:

	webhookdb auth login

The CLI will guide you from there.

The CLI also gives you quick access to the WebhookDB documentation:

	webhookdb docs html
	webhookdb docs tui
`,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Value: false},
		},
		Commands: []*cli.Command{
			authCmd,
			backfillCmd,
			dbCmd,
			docsCmd,
			fixturesCmd,
			integrationsCmd,
			organizationsCmd,
			servicesCmd,
			subscriptionsCmd,
			updateCmd,
			webhooksCmd,
			versionCmd,
		},
	}
	return app
}
