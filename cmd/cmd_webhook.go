package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var webhooksCmd = &cli.Command{
	Name:    "webhook",
	Aliases: []string{"webhooks"},
	Usage:   "Manage webhooks that will be notified when WebhookDB data is updated.",
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Usage: "Create a new webhook that WebhookDB will call on every data update.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
				&cli.StringFlag{
					Name:  "url",
					Usage: "Full URL to the endpoint that will be POSTed to whenever this organization or integration is updated.",
				},
				&cli.StringFlag{
					Name:  "secret",
					Usage: "Random secure secret to use to sign webhooks coming from WebhookDB.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.WebhookCreateInput{
					Url:           c.String("url"),
					WebhookSecret: c.String("secret"),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				if c.String("integration") != "" {
					input.IntegrationIdentifier = c.String("integration")
				}
				out, err := client.WebhookCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name:  "list",
			Usage: "List all created webhooks.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.WebhookListInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.WebhookList(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
		{
			Name:  "test",
			Usage: "Send a test event to webhook subscription with the given ID.",
			Flags: []cli.Flag{webhookFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.WebhookTest(ctx, ac.Auth, client.WebhookOpaqueIdInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					OpaqueId:      getWebhookFlagOrArg(c),
				})
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name:  "delete",
			Usage: "Delete this webhook subscription, so no future events will be sent.",
			Flags: []cli.Flag{webhookFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.WebhookDelete(ctx, ac.Auth, client.WebhookOpaqueIdInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					OpaqueId:      getWebhookFlagOrArg(c),
				})
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
	},
}

func webhookFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "webhook",
		Aliases: s1("w"),
		Usage:   usage("Webhook opaque id. Run `webhookdb webhook list` to see a list of all your webhooks."),
	}
}

func getWebhookFlagOrArg(c *cli.Context) string {
	return requireFlagOrArg(c, "webhook", "Use `webhookdb webhook list` to see available webhooks.")
}
