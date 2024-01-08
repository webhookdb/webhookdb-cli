package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/client"
)

var webhooksCmd = &cli.Command{
	Name:    "notification",
	Aliases: []string{"notifications", "webhook", "webhooks"},
	Usage:   "Manage endpoints that will be notified when WebhookDB data is updated.",
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Usage: "Create a new notification that WebhookDB will call on every data update.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
				&cli.StringFlag{
					Name:  "url",
					Usage: "Full URL to the endpoint that will be POSTed to whenever this organization or integration is updated.",
				},
				&cli.StringFlag{
					Name:  "secret",
					Usage: "Random secure secret to use to sign notifications coming from WebhookDB.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.WebhookCreateInput{
					Url:                   c.String("url"),
					WebhookSecret:         c.String("secret"),
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
					IntegrationIdentifier: flagOrArg(c, "integration"),
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
			Usage: "List all created notifications.",
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
			Usage: "Send a test event to the notification subscription with the given ID.",
			Flags: []cli.Flag{notificationFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.WebhookTest(ctx, ac.Auth, client.WebhookOpaqueIdInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					OpaqueId:      getNotificationArgOrFlag(c),
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
			Usage: "Delete this notification subscription, so no future events will be sent.",
			Flags: []cli.Flag{notificationFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.WebhookDelete(ctx, ac.Auth, client.WebhookOpaqueIdInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					OpaqueId:      getNotificationArgOrFlag(c),
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

func notificationFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "notification",
		Aliases: s1("n"),
		Usage:   usage("Notification opaque id. Run `webhookdb notification list` to see a list of all your notification subscriptions."),
	}
}

func getNotificationArgOrFlag(c *cli.Context) string {
	return requireFlagOrArg(c, "notification", "Use `webhookdb notification list` to see available notifications.")
}
