package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/client"
)

var savedViewsCmd = &cli.Command{
	Name:    "saved-view",
	Aliases: []string{"saved-views", "view", "views"},
	Usage:   "Create, replace, and drop views in your database.",
	Subcommands: []*cli.Command{
		{
			Name:    "create-or-replace",
			Usage:   "Create or replace a view.",
			Aliases: []string{"create", "replace", "new", "update"},
			Flags: []cli.Flag{
				orgFlag(),
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "The name of the view. Must be a database identifier (alphanumeric, spaces, and underscores).",
				},
				&cli.StringFlag{
					Name:  "sql",
					Usage: "SQL SELECT statement to run. Must include the 'CREATE VIEW' part of the query.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SavedViewCreateInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Name:          c.String("name"),
					Sql:           c.String("sql"),
				}
				out, err := client.SavedViewCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name:  "list",
			Usage: "List all saved views.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SavedViewListInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.SavedViewList(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
		{
			Name:  "delete",
			Usage: "Delete the saved view.",
			Flags: []cli.Flag{savedViewNameFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.SavedViewDelete(ctx, ac.Auth, client.SavedViewDeleteInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Name:          getSavedViewNameArgOrFlag(c),
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

func savedViewNameFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "name",
		Aliases: s1("n"),
		Usage:   usage("Name of the view. Run `webhookdb view list` to see a list of all your saved views."),
	}
}

func getSavedViewNameArgOrFlag(c *cli.Context) string {
	return requireFlagOrArg(c, "name", "Use `webhookdb view list` to see available saved views.")
}
