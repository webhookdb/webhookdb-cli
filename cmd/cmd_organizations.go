package cmd

import (
	"fmt"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/urfave/cli/v2"
	"strings"
)

var organizationsCmd = &cli.Command{
	Name: "org",
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
				out, err := client.OrganizationsList(ctx, client.OrganizationsListInput{AuthCookie: p.AuthCookie})
				if err != nil {
					return err
				}
				orgsLen := len(out.Data)
				keySlugs := make([]string, orgsLen)
				for i, value := range out.Data {
					if value.Key == p.CurrentOrg {
						keySlugs[i] = (value.Key + " (active)")
					} else {
						keySlugs[i] = value.Key
					}
				}
				fmt.Println(strings.Join(keySlugs, "\n"))
				return nil
			},
		},
	},
}
