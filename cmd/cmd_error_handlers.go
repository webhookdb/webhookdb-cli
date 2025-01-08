package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/client"
)

var errorHandlersCmd = &cli.Command{
	Name:    "error-handler",
	Aliases: []string{"error-handlers"},
	Usage:   "Automatically alert an endpoint when a replicator error occurs, rather than sending an email. See https://docs.webhookdb.com/docs/integrating/error-handlers.html for help.",
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Usage: "Create a new error handler.",
			Flags: []cli.Flag{
				orgFlag(),
				&cli.StringFlag{
					Name:    "url",
					Aliases: []string{"u"},
					Usage:   "URL to notify about errors.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.ErrorHandlerCreateInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Url:           c.String("url"),
				}
				out, err := client.ErrorHandlerCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name:  "list",
			Usage: "List all error handlers.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.ErrorHandlerListInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.ErrorHandlerList(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
		{
			Name:  "delete",
			Usage: "Delete an error handler.",
			Flags: []cli.Flag{errorHandlerFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.ErrorHandlerDelete(ctx, ac.Auth, client.ErrorHandlerIdentifierInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Identifier:    getErrorHandlerArgOrFlag(c),
				})
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
	},
}

func errorHandlerFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "error-handler",
		Aliases: s1("e"),
		Usage:   usage("Error handler id. Run `webhookdb error-handler list` to see a list of all your error handlers."),
	}
}

func getErrorHandlerArgOrFlag(c *cli.Context) string {
	return requireFlagOrArg(c, "error-handler", "Use `webhookdb error-handler list` to see available error handlers.")
}
