package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var backfillCmd = &cli.Command{
	Name:  "backfill",
	Usage: "Start backfilling all the resources available to the service integration.",
	Flags: []cli.Flag{
		orgFlag(),
		integrationFlag(),
	},
	Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
		input := client.BackfillInput{
			IntegrationIdentifier: getIntegrationFlagOrArg(c),
			OrgIdentifier:         getOrgFlag(c, ac.Prefs),
		}
		return stateMachineResponseRunner(ctx, ac.Auth)(client.Backfill(ctx, ac.Auth, input))
	}),
	Subcommands: []*cli.Command{
		{
			Name:  "reset",
			Usage: "Reset any stored API keys and secrets associated with backfilling this integration.",
			Flags: []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.BackfillResetInput{
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
				}
				return stateMachineResponseRunner(ctx, ac.Auth)(client.BackfillReset(ctx, ac.Auth, input))
			}),
		},
	},
}
