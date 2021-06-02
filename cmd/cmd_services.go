package cmd

import (
	"fmt"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/urfave/cli/v2"
	"strings"
)

var servicesCmd = &cli.Command{
	Name: "services",
	Subcommands: []*cli.Command{
		{
			Name:        "list",
			Description: "TODO",
			Flags:       []cli.Flag{},
			Action: func(c *cli.Context) error {
				ac := newAppCtx(c)
				ctx := newCtx(ac)
				p, err := prefs.Load()
				if err != nil {
					return err
				}
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
			},
		},
	},
}
