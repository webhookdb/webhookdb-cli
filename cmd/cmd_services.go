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
	Name:        "services",
	Description: "We use the term \"services\" to describe all of the platforms currently available for integration",
	Subcommands: []*cli.Command{
		{
			Name:        "list",
			Description: "list all available services",
			Flags:       []cli.Flag{},
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
				fmt.Println(strings.Join(names, "\n"))
				return nil
			}),
		},
	},
}
