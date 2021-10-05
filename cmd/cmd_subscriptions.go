package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"os"
)

var subscriptionsCmd = &cli.Command{
	Name: "subscription",
	Subcommands: []*cli.Command{
		{
			Name:        "info",
			Description: "get information about an organization's subscription",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.SubscriptionInfo(ctx, ac.Auth, client.SubscriptionInfoInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
				if err != nil {
					return err
				}
				out.PrintTo(os.Stdout)
				return nil
			}),
		},
		{
			Name:        "edit",
			Description: "open stripe portal to edit subscription",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.SubscriptionEdit(ctx, ac.Auth, client.SubscriptionEditInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
				if err != nil {
					return err
				}
				if err := browser.OpenURL(out.SessionUrl); err != nil {
					return err
				}
				fmt.Println("You have been redirected to the Stripe Billing Portal:")
				fmt.Println(out.SessionUrl)
				return nil
			}),
		},
	},
}
