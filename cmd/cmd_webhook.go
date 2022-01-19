package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var webhooksCmd = &cli.Command{
	Name: "webhook",
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Flags: []cli.Flag{serviceIntegrationFlag(), orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				url, err := extractWebhookUrl(0, c)
				if err != nil {
					return err
				}
				_, err = client.WebhookCreate(ctx, ac.Auth, client.WebhookCreateInput{
					Url:           url,
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
			Flags: []cli.Flag{serviceIntegrationFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				_, err := client.WebhookTest(ctx, ac.Auth, client.WebhookOpaqueIdInput{OpaqueId: c.String("integration")})
				if err != nil {
					return err
				}
				return nil
			}),
		},
		{
			Name:  "delete",
			Flags: []cli.Flag{serviceIntegrationFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				_, err := client.WebhookDelete(ctx, ac.Auth, client.WebhookOpaqueIdInput{OpaqueId: c.String("integration")})
				if err != nil {
					return err
				}
				return nil
			}),
		},
	},
}
