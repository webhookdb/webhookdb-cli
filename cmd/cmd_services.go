package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
	"strings"
)

var servicesCmd = &cli.Command{
	Name:  "services",
	Usage: "Work with available services that can be hooked up to reflect data to WebhookDB.",
	Subcommands: []*cli.Command{
		{
			Name:  "list",
			Usage: "List all available services.",
			Flags: []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.ServicesList(ctx, ac.Auth, client.ServicesListInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				})
				if err != nil {
					return err
				}
				servicesLen := len(out.Data)
				names := make([]string, servicesLen)
				for i, value := range out.Data {
					names[i] = value.Name
				}
				fmt.Fprintln(c.App.Writer, strings.Join(names, "\n"))
				return nil
			}),
		},
	},
}
