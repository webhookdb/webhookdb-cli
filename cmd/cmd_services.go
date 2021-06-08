package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/urfave/cli/v2"
	"strings"
)

var servicesCmd = &cli.Command{
	Name: "services",
	Description: "We use the term \"services\" to describe all of the platforms currently available for integration",
	Subcommands: []*cli.Command{
		{
			Name:        "list",
			Description: "list all available services",
			Flags:       []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				out, err := client.ServicesList(ctx, client.ServicesListInput{AuthCookie: p.AuthCookie})
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
