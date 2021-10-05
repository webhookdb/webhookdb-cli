package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.DbConnection(ctx, ac.Auth, client.DbConnectionInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.DbTables(ctx, ac.Auth, client.DbTablesInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				q, err := extractPositional(0, c, "You must enter a query string.")
				if err != nil {
					return err
				}
				out, err := client.DbSql(ctx, ac.Auth, client.DbSqlInput{OrgIdentifier: getOrgFlag(c, ac.Prefs), Query: q})
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.DbRollCredentials(ctx, ac.Auth, client.DbRollCredentialsInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
				if err != nil {
					return err
				}
				fmt.Print(out.ConnectionUrl)
				return nil
			}),
		},
	},
}
