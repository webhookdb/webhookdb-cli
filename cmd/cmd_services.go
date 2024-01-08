package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/client"
)

var servicesCmd = &cli.Command{
	Name:    "services",
	Aliases: []string{"service"},
	Usage:   "Work with available services that can be hooked up to reflect data to WebhookDB.",
	Subcommands: []*cli.Command{
		{
			Name:  "list",
			Usage: "List all available services.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.ServicesList(ctx, ac.Auth, client.ServicesListInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				})
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
	},
}
