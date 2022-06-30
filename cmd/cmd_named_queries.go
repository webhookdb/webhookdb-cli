package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var namedQueriesCmd = &cli.Command{
	Name:  "named-query",
	Usage: "Make sure that you're working on the correct organization when you create a named query.",
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Usage: "Create a named query for the given organization.",
			Flags: []cli.Flag{
				orgFlag(),
				&cli.StringFlag{
					Name:    "name",
					Aliases: s1("n"),
					Usage:   "Name of the query",
				},
				&cli.StringFlag{
					Name:    "sql",
					Aliases: s1("s"),
					Usage:   "Sql to be executed",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.NamedQueryCreateInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Name:          c.String("name"),
					Sql:           c.String("sql"),
				}
				out, err := client.NamedQueryCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name:  "info",
			Usage: "Retrieve info about the query with the given ID.",
			Flags: []cli.Flag{
				orgFlag(),
				namedQueryFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.NamedQueryInfoInput{
					QueryIdentifier: getNamedQueryFlagOrArg(c),
					OrgIdentifier:   getOrgFlag(c, ac.Prefs),
				}
				out, err := client.NamedQueryInfo(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteSingle(c.App.Writer, out)
			}),
		},
		{
			Name:  "list",
			Usage: "List all named queries for the given organization.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.NamedQueryListInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.NamedQueryList(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
		{
			Name:  "run",
			Usage: "Run the named query.",
			Flags: []cli.Flag{
				orgFlag(),
				namedQueryFlag(),
				&cli.BoolFlag{Name: "color", Aliases: s1("c"), Usage: "Display colors. Default true if tty.", Value: IsTty},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				useColors := c.Bool("color")
				input := client.NamedQueryRunInput{
					QueryIdentifier: getNamedQueryFlagOrArg(c),
					OrgIdentifier:   getOrgFlag(c, ac.Prefs),
				}
				out, err := client.NamedQueryRun(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				err = printSqlOutput(c, out, useColors)
				if err != nil {
					return err
				}
				return nil
			}),
		},
		{
			Name:  "update",
			Usage: "Update the name or contents of the query",
			Flags: []cli.Flag{
				orgFlag(),
				namedQueryFlag(),
				fieldFlag(),
				valueFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.NamedQueryUpdateInput{
					OrgIdentifier:   getOrgFlag(c, ac.Prefs),
					QueryIdentifier: getNamedQueryFlagOrArg(c),
					Field:           c.String("field"),
					Value:           c.String("value"),
				}
				out, err := client.NamedQueryUpdate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message, true)
				return nil
			}),
		},
	},
}

func namedQueryFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "query",
		Aliases: s1("q"),
		Usage:   "Named query identifier. This can either be the short id or the name you've given the query. Run `webhookdb named-query list` to see a list of all your named queries.",
	}
}

func getNamedQueryFlagOrArg(c *cli.Context) string {
	return requireFlagOrArg(c, "query", "Use `webhookdb named-query list` to see available named queries.")
}
