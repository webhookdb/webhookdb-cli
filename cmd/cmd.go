package cmd

import (
	"github.com/urfave/cli/v2"
)

func BuildApp() *cli.App {
	app := &cli.App{
		Name:     "webhookdb",
		HelpName: helpName,
		Usage: `CLI for WebhookDB <https://github.com/webhookdb/webhookdb>.
WebhookDB replicates any API into a database,
so you have immediate, reliable access to all your data.

To create an account and get started, run:

    webhookdb auth login

The CLI will guide you from there.

The CLI also gives you quick access to the WebhookDB documentation:

    webhookdb docs html # Open a browser to https://docs.webhookdb.com
    webhookdb docs tui # Display the CLI Reference in the terminal

Source code: <https://github.com/webhookdb/webhookdb-cli>

Documentation: <https://docs.webhookdb.com/docs/cli-reference.md>
`,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Value: false},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: s1("q"),
				Value:   false,
				Usage:   "Do not print messages. Mostly used for collection endpoints, where you just want the returned data, not the help message."},
		},
		Commands: []*cli.Command{
			authCmd,
			backfillCmd,
			dbCmd,
			debugCmd,
			docsCmd,
			errorHandlersCmd,
			fixturesCmd,
			integrationsCmd,
			organizationsCmd,
			replayCmd,
			savedQueriesCmd,
			savedViewsCmd,
			servicesCmd,
			subscriptionsCmd,
			syncCmd(dbSyncType),
			syncCmd(httpSyncType),
			updateCmd,
			webhooksCmd,
			versionCmd,
		},
	}
	return app
}
