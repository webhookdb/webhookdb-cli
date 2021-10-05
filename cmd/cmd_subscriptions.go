package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				out, err := client.SubscriptionInfo(ctx, client.SubscriptionInfoInput{AuthCookie: p.AuthCookie, OrgIdentifier: getOrgFlag(c, p)})
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				out, err := client.SubscriptionEdit(ctx, client.SubscriptionEditInput{AuthCookie: p.AuthCookie, OrgIdentifier: getOrgFlag(c, p)})
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
