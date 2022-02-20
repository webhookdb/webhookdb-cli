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
	Name: "subscription",
	Subcommands: []*cli.Command{
		{
			Name:        "info",
			Description: "Get information about an organization's software subscription.",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.SubscriptionInfo(ctx, ac.Auth, client.SubscriptionInfoInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
				if err != nil {
					return err
				}
				out.PrintTo(c.App.Writer)
				return nil
			}),
		},
		{
			Name:        "edit",
			Description: "Open stripe portal to edit subscription.",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.SubscriptionEdit(ctx, ac.Auth, client.SubscriptionEditInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
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
	},
}
