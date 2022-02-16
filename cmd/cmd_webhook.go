package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var webhooksCmd = &cli.Command{
	Name:  "webhook",
	Usage: "Manage webhooks that will be notified when WebhookDB data is updated.",
	Subcommands: []*cli.Command{
		{
			Name: "create",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
				&cli.StringFlag{
					Name:  "url",
					Usage: "Full URL to the endpoint that will be POSTed to whenever this organization or integration is updated.",
				},
				&cli.StringFlag{
					Name:  "secret",
					Usage: "Random secure secret to use to sign webhooks coming from WebhookDB",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				_, err := client.WebhookCreate(ctx, ac.Auth, client.WebhookCreateInput{
					Url:           c.String("url"),
					WebhookSecret: c.String("secret"),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					SintOpaqueId:  c.String("integration"),
				})
				if err != nil {
					return err
				}
				return nil
			}),
		},
		{
			Name:  "test",
			Usage: "Send a test event to all webhook subscriptions associated with this integration.",
			Flags: []cli.Flag{integrationFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				_, err := client.WebhookTest(ctx, ac.Auth, client.WebhookOpaqueIdInput{OpaqueId: flagOrArg(c, "integration", "Use `webhookdb integrations list` to see available integrations.")})
				if err != nil {
					return err
				}
				return nil
			}),
		},
		{
			Name:  "delete",
			Usage: "Delete this webhook subscription, so no future events will be sent.",
			Flags: []cli.Flag{integrationFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.WebhookOpaqueIdInput{
					OpaqueId: getIntegrationFlagOrArg(c),
				}
				_, err := client.WebhookDelete(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				return nil
			}),
		},
	},
}
