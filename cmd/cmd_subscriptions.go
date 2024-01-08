package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/client"
	"github.com/webhookdb/webhookdb-cli/whbrowser"
)

var subscriptionsCmd = &cli.Command{
	Name:    "subscription",
	Aliases: []string{"subscriptions", "sub", "subs"},
	Usage:   "Work with your WebhookDB subscription.",
	Subcommands: []*cli.Command{
		{
			Name:  "info",
			Usage: "Get information about an organization's software subscription.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.SubscriptionInfo(ctx, ac.Auth, client.SubscriptionInfoInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteSingle(c.App.Writer, out)
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
					Usage:    usage("Plan name, like 'yearly', 'monthly', etc. Use `webhookdb subscription plans` to see available plans."),
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
				fmt.Fprintln(c.App.Writer, "Your browser was opened redirected to the Stripe Billing Portal:", out.SessionUrl)
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
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
	},
}
