package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/urfave/cli/v2"
)

var backfillCmd = &cli.Command{
	Name:        "backfill",
	Description: "You can run this command to start a backfill of all the resources available to an integration.",
	Flags:       []cli.Flag{},
	Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
		opaqueId, err := extractIntegrationId(0, c)
		if err != nil {
			return err
		}
		input := client.BackfillInput{
			AuthCookie: p.AuthCookie,
			OpaqueId:   opaqueId,
		}
		step, err := client.Backfill(ctx, input)
		if err != nil {
			return err
		}
		if err := client.NewStateMachine().Run(ctx, p, step); err != nil {
			return err
		}
		return nil
	}),
}
