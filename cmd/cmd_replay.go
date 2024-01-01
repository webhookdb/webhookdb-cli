package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var replayCmd = &cli.Command{
	Name:    "replay",
	Aliases: []string{"replays"},
	Usage: "Replay webhooks received by WebhookDB. Useful if webhooks were rejected due to " +
		"invalid secrets/verification and need to be re-processed.\n\n" +
		"Note that the maximum window that can be replayed is managed by the server, and is normally 4 hours. " +
		"To replay a larger window, issue multiple `webhookdb replay` calls.",
	Flags: []cli.Flag{
		orgFlag(),
		(func() cli.Flag {
			f := integrationFlag()
			f.Usage = "Integration identifier. This can either be the service name, the table name, " +
				"or the opaque id, which is a unique code that starts with 'svi_'. " +
				"If not given, replay webhooks for all integrations in the current organization."
			return f
		})(),
		&cli.IntFlag{
			Name:    "hours",
			Aliases: s1("r"),
			Usage:   "Number of hours before now (or the --before value) to replay webhooks.",
		},
		&cli.StringFlag{
			Name:    "before",
			Aliases: s1("b"),
			Usage:   "Replay webhooks received before this time. Timestamp can be anything time-like, preferably ISO8601 (ie, 2012-11-22T06:00:00Z). Defaults to now.",
		},
		&cli.StringFlag{
			Name:    "after",
			Aliases: s1("a"),
			Usage: "Replay webhooks received after this time. Takes precedence over --hours. " +
				"If neither are given, use 1 hour.",
		},
	},
	Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
		input := client.ReplayInput{
			OrgIdentifier:         getOrgFlag(c, ac.Prefs),
			IntegrationIdentifier: c.String("integration"),
			Hours:                 c.Int("hours"),
			Before:                c.String("before"),
			After:                 c.String("after"),
		}
		out, err := client.Replay(ctx, ac.Auth, input)
		if err != nil {
			return err
		}
		printlnif(c, out.Message, true)
		return nil
	}),
}
