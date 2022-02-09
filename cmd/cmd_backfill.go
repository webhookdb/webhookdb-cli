package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var backfillCmd = &cli.Command{
	Name:        "backfill",
	Description: "You can run this command to start a backfill of all the resources available to an integration.",
	Flags:       []cli.Flag{orgFlag()},
	Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
		opaqueId, err := extractIntegrationId(0, c)
		if err != nil {
			return err
		}
		input := client.BackfillInput{
			OpaqueId:      opaqueId,
			OrgIdentifier: getOrgFlag(c, ac.Prefs),
		}
		step, err := client.Backfill(ctx, ac.Auth, input)
		if err != nil {
			return err
		}
		if err := client.NewStateMachine().Run(ctx, ac.Auth, step); err != nil {
			return err
		}
		return nil
	}),
	Subcommands: []*cli.Command{
		{
			Name:        "reset",
			Description: "Reset any stored API keys and secrets associated with backfilling this integration.",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				opaqueId, err := extractIntegrationId(0, c)
				if err != nil {
					return err
				}
				input := client.BackfillResetInput{
					OpaqueId:      opaqueId,
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				step, err := client.BackfillReset(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				if err := client.NewStateMachine().Run(ctx, ac.Auth, step); err != nil {
					return err
				}
				return nil
			}),
		},
	},
}
