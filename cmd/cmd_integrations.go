package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/urfave/cli/v2"
)

var integrationsCmd = &cli.Command{
	Name:  "integrations",
	Usage: "Make sure that you're working on the correct organization when you create an integration.",
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Usage: "Create an integration for the given organization.",
			Flags: []cli.Flag{
				orgFlag(),
				serviceFlag(),
				&cli.BoolFlag{
					Name:    "confirm",
					Aliases: s1("c"),
					Usage: "If there is already an integration for this service, " +
						"you will be prompted to confirm you want to create a new integration. " +
						"Pass --confirm to automatically accept this prompt and always create a new integration.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsCreateInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					ServiceName:   requireFlagOrArg(c, "service", "Use `webhookdb services list` to see available integrations."),
				}
				if c.IsSet("confirm") {
					input.GuardConfirm = types.SPtr("y")
				}
				return stateMachineResponseRunner(ctx, ac.Auth)(client.IntegrationsCreate(ctx, ac.Auth, input))
			}),
		},
		{
			Name:  "delete",
			Usage: "Delete an integration and its table.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
				&cli.StringFlag{
					Name:    "confirm",
					Aliases: s1("c"),
					Usage: "Confirm this action by providing a value of the integration's table name. " +
						"Will be prompted if not provided.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsDeleteInput{
					OpaqueId:      getIntegrationFlagOrArg(c),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Confirm:       c.String("confirm"),
				}
				out, err := client.IntegrationsDelete(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name:  "list",
			Usage: "List all integrations for the given organization.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsListInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.IntegrationsList(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
		{
			Name:  "reset",
			Usage: "Reset the webhook secret for this integration.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsResetInput{
					OpaqueId:      getIntegrationFlagOrArg(c),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				return stateMachineResponseRunner(ctx, ac.Auth)(client.IntegrationsReset(ctx, ac.Auth, input))
			}),
		},
		{
			Name:  "stats",
			Usage: "Get statistics about webhooks for this integration.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsStatsInput{
					OpaqueId:      getIntegrationFlagOrArg(c),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.IntegrationsStats(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteSingle(c.App.Writer, out)
			}),
		},
	},
}
