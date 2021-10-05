package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/urfave/cli/v2"
	"strings"
)

var dbCmd = &cli.Command{
	Name:        "db",
	Description: "",
	Flags:       []cli.Flag{orgFlag()},
	Subcommands: []*cli.Command{
		{
			Name:        "connection",
			Description: "Print the database connection url for an organization.",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				out, err := client.DbConnection(ctx, client.DbConnectionInput{AuthCookie: p.AuthCookie, OrgIdentifier: getOrgFlag(c, p)})
				if err != nil {
					return err
				}
				fmt.Print(out.ConnectionUrl)
				return nil
			}),
		},
		{
			Name:        "tables",
			Description: "List all tables in an organization's database.",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				out, err := client.DbTables(ctx, client.DbTablesInput{AuthCookie: p.AuthCookie, OrgIdentifier: getOrgFlag(c, p)})
				if err != nil {
					return err
				}
				fmt.Println(strings.Join(out.TableNames, "\n"))
				return nil
			}),
		},
		{
			Name:        "sql",
			Description: "Execute query on organization's database.",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				if c.NArg() != 1 {
					return errors.New("You must enter a query string.")
				}
				out, err := client.DbSql(ctx, client.DbSqlInput{AuthCookie: p.AuthCookie, OrgIdentifier: getOrgFlag(c, p), Query: c.Args().Get(0)})
				if err != nil {
					return err
				}
				fmt.Println(strings.Join(out.Columns, "t"))
				fmt.Println(strings.Join(out.Rows, "\n"))
				if out.MaxRowsReached {
					fmt.Println("Results have been truncated.")
				}
				return nil
			}),
		},
		{
			Name:        "roll-credentials",
			Description: "Roll the credentials for an organization's database to something newly randomly generated.",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				out, err := client.DbRollCredentials(ctx, client.DbRollCredentialsInput{AuthCookie: p.AuthCookie, OrgIdentifier: getOrgFlag(c, p)})
				if err != nil {
					return err
				}
				fmt.Print(out.ConnectionUrl)
				return nil
			}),
		},
	},
}
