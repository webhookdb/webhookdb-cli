package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/client"
)

var savedQueriesCmd = &cli.Command{
	Name:    "saved-query",
	Aliases: []string{"saved-queries", "custom-query", "custom-queries"},
	Usage:   "Manage your library of saved queries, including ones that can be accessed publicly.",
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Usage: "Create a new saved query.",
			Flags: []cli.Flag{
				orgFlag(),
				&cli.StringFlag{
					Name:    "description",
					Aliases: []string{"d", "desc"},
					Usage:   "Explain what this query is used for.",
				},
				&cli.StringFlag{
					Name:  "sql",
					Usage: "SQL statement to run.",
				},
				&cli.BoolFlag{
					Name: "public",
					Usage: "If true, the query can be accessed publicly, without authentication. " +
						"Allows a saved query to be used on public dashboards or websites, " +
						"without exposing a database connection string or " +
						"allowing public access to a database.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SavedQueryCreateInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Description:   c.String("description"),
					Sql:           c.String("sql"),
					Public:        false,
				}
				out, err := client.SavedQueryCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name:  "list",
			Usage: "List all saved queries.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SavedQueryListInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.SavedQueryList(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
		{
			Name:    "update",
			Usage:   "Update a new saved query.",
			Aliases: []string{"edit", "modify"},
			Flags: []cli.Flag{
				orgFlag(),
				savedQueryFlag(),
				fieldFlag(),
				valueFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SavedQueryUpdateInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Identifier:    getSavedQueryArgOrFlag(c),
					Field:         c.String("field"),
					Value:         c.String("value"),
				}
				out, err := client.SavedQueryUpdate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name:  "info",
			Usage: "Display information about given saved query.",
			Flags: []cli.Flag{
				orgFlag(),
				savedQueryFlag(),
				&cli.StringFlag{
					Name:    "field",
					Aliases: s1("f"),
					Usage:   "The field that you want information about.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SavedQueryInfoInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Identifier:    getSavedQueryArgOrFlag(c),
					Field:         c.String("field"),
				}
				out, err := client.SavedQueryInfo(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				_, err = out.Blocks.WriteTo(c.App.Writer)
				return err
			}),
		},
		{
			Name:  "run",
			Usage: "Run the query.",
			Flags: []cli.Flag{
				savedQueryFlag(),
				colorFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				useColors := c.Bool("color")
				out, err := client.SavedQueryRun(ctx, ac.Auth, client.SavedQueryIdentifierInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Identifier:    getSavedQueryArgOrFlag(c),
				})
				if err != nil {
					return err
				}
				err = printSqlOutput(c, out, useColors)
				return nil
			}),
		},
		{
			Name:  "delete",
			Usage: "Delete this saved query.",
			Flags: []cli.Flag{savedQueryFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.SavedQueryDelete(ctx, ac.Auth, client.SavedQueryIdentifierInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Identifier:    getSavedQueryArgOrFlag(c),
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

func savedQueryFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "saved-query",
		Aliases: s1("q"),
		Usage:   usage("Saved query opaque id. Run `webhookdb saved-query list` to see a list of all your saved queries."),
	}
}

func getSavedQueryArgOrFlag(c *cli.Context) string {
	return requireFlagOrArg(c, "saved-query", "Use `webhookdb saved-query list` to see available saved queries.")
}
