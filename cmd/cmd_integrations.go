package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var integrationsCmd = &cli.Command{
	Name: "integrations",
	Subcommands: []*cli.Command{
		{
			Name:        "create",
			Description: "TODO",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				if c.NArg() != 1 {
					return errors.New("service name required. Use 'webhookdb services list'.")
				}
				var orgKey string
				if c.String("org") != "" {
					orgKey = c.String("org")
				} else {
					orgKey = p.CurrentOrg
				}
				input := client.IntegrationsCreateInput{
					AuthCookie:  p.AuthCookie,
					OrgKey:      orgKey,
					ServiceName: c.Args().Get(0),
				}
				step, err := client.IntegrationsCreate(ctx, input)
				if err != nil {
					return err
				}
				if err := client.NewStateMachine().Run(ctx, p, step); err != nil {
					return err
				}
				return nil
			}),
		},
	},
}
