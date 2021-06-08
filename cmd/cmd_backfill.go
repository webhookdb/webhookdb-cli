package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var backfillCmd = &cli.Command{
	Name:        "backfill",
	Description: "You can run this command to start a backfill of all the resources available to an integration.",
	Flags:       []cli.Flag{},
	Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
		if c.NArg() != 1 {
			return errors.New("Integration ID required. Use 'webhookdb integrations list'.")
		}
		input := client.BackfillInput{
			AuthCookie: p.AuthCookie,
			OpaqueId:   c.Args().Get(0),
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
