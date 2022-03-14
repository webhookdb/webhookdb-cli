package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/whbrowser"
	"github.com/urfave/cli/v2"
)

var subscriptionsCmd = &cli.Command{
	Name:  "subscription",
	Usage: "Work with your WebhookDB subscription.",
	Subcommands: []*cli.Command{
		{
			Name:  "info",
			Usage: "Get information about an organization's software subscription.",
			Flags: []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.SubscriptionInfo(ctx, ac.Auth, client.SubscriptionInfoInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
				if err != nil {
					return err
				}
				printlnif(c, out.Message)
				return nil
			}),
		},
		{
			Name:  "edit",
			Usage: "Open stripe portal to edit subscription.",
			Flags: []cli.Flag{
				orgFlag(),
				&cli.StringFlag{
					Name:     "plan",
					Aliases:  s1("p"),
					Required: false,
					Usage:    "Plan name, like 'yearly', 'monthly', etc. Use `webhookdb subscription plans` to see available plans.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SubscriptionEditInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Plan:          c.String("plan"),
				}
				out, err := client.SubscriptionEdit(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				if err := whbrowser.OpenURL(out.SessionUrl); err != nil {
					return err
				}
				fmt.Fprintln(c.App.Writer, "You have been redirected to the Stripe Billing Portal:")
				fmt.Fprintln(c.App.Writer, out.SessionUrl)
				return nil
			}),
		},
		{
			Name:  "plans",
			Usage: "Print information about the WebhookDB pricing plans.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SubscriptionPlansInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.SubscriptionPlans(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message())
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
	},
}
