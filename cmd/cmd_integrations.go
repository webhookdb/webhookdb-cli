package cmd

import (
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var integrationsCmd = &cli.Command{
	Name: "integrations",
	Subcommands: []*cli.Command{
		{
			Name:        "create",
			Description: "TODO",
			Flags:       []cli.Flag{},
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					return errors.New("service name required. Use 'webhookdb services list'.")
				}
				ac := newAppCtx(c)
				ctx := newCtx(ac)
				out, err := client.IntegrationsCreate(ctx, client.IntegrationsCreateInput{ServiceName: c.Args().Get(0)})
				if err != nil {
					return err
				}
				if err := ac.StateMachine.Run(ctx, out.Step); err != nil {
					return err
				}
				return nil
			},
		},
	},
}
