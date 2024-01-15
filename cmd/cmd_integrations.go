package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/client"
	"github.com/webhookdb/webhookdb-cli/types"
)

var integrationsCmd = &cli.Command{
	Name:    "integrations",
	Aliases: []string{"integration"},
	Usage:   "Make sure that you're working on the correct organization when you create an integration.",
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
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
					Confirm:               c.String("confirm"),
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
			Name:  "info",
			Usage: "Display information about given integration.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
				&cli.StringFlag{
					Name:    "field",
					Aliases: s1("f"),
					Usage:   "The field that you want information about.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsInfoInput{
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
					Field:                 c.String("field"),
				}
				out, err := client.IntegrationsInfo(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				_, err = out.Blocks.WriteTo(c.App.Writer)
				return err
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
			Name:  "setup",
			Usage: "Ensure all the necessary fields are set to receive webhooks.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsSetupInput{
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
				}
				return stateMachineResponseRunner(ctx, ac.Auth)(client.IntegrationsSetup(ctx, ac.Auth, input))
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
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
				}
				return stateMachineResponseRunner(ctx, ac.Auth)(client.IntegrationsReset(ctx, ac.Auth, input))
			}),
		},
		{
			Name:  "roll-key",
			Usage: "Roll the API key used to access this integration. Only relevant for certain integrations.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				// This get some more work if it becomes more widely needed, but for now, it should be very rare.
				input := client.IntegrationsRollKeyInput{
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
				}
				out, err := client.IntegrationsRollKey(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Fprint(c.App.Writer, out.WebhookdbApiKey)
				return nil
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
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
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
