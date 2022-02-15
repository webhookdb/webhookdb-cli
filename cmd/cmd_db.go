package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/go-aperitif/convext"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

var dbCmd = &cli.Command{
	Name:        "db",
	Description: "Command namespace for interacting with your organization's database and tables.",
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
			Flags: []cli.Flag{
				orgFlag(),
				&cli.StringFlag{Name: "query", Aliases: s1("u"), Usage: "Query string to execute using your connection."},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.DbSqlInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Query:         extractPositional(0, c, "You must enter a query string."),
				}
				out, err := client.DbSql(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				table := tablewriter.NewWriter(os.Stdout)

				table.SetHeader(out.Columns)
				headerCols := make([]tablewriter.Colors, len(out.Columns))
				for i := range out.Columns {
					headerCols[i] = tablewriter.Colors{tablewriter.FgHiGreenColor}
				}
				table.SetHeaderColor(headerCols...)

				for _, row := range out.Rows {
					rowStr := make([]string, len(row))
					colors := make([]tablewriter.Colors, len(row))
					for i, cell := range row {
						if cell == nil {
							rowStr[i] = "<null>"
							colors[i] = tablewriter.Colors{tablewriter.FgYellowColor}
						} else {
							rowStr[i] = fmt.Sprintf("%v", cell)
							colors[i] = tablewriter.Colors{}
						}
					}
					table.Rich(rowStr, colors)
				}
				if out.MaxRowsReached {
					table.SetCaption(true, "Results have been truncated.")
				}
				table.Render()
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
		{
			Name: "fdw",
			Description: "Write out commands that can be used to generate a FDW against your WebhookDB database and " +
				"import them into materialized views. See flags for further usage.",
			Flags: []cli.Flag{
				orgFlag(),
				&cli.BoolFlag{Name: "raw", Usage: "If given, print the raw SQL returned from the server. Useful if you want to pipe through jq or something similar."},
				&cli.BoolFlag{Name: "fdw", Usage: "Write the FDW SQL to stdout"},
				&cli.BoolFlag{Name: "views", Usage: "Write the SQL to create the materialized views to stdout"},
				&cli.BoolFlag{Name: "all", Usage: "Write a single SQL statement containing FDW and view creation code. Default if neither --fdw or --views are passed."},
				&cli.StringFlag{Name: "remote", Value: "webhookdb_remote", Usage: "The remote server name, used in the 'CREATE SERVER <remote>' call"},
				&cli.StringFlag{Name: "fetch", Value: "50000", Usage: "fetch_size option used during server creation"},
				&cli.StringFlag{Name: "into-schema", Value: "webhookdb_remote", Usage: "Name of the schema to import the remote tables into (IMPORT FOREIGN SCHEMA public INTO <into schema>."},
				&cli.StringFlag{Name: "views-schema", Value: "webhookdb", Usage: "Create materialized views in this schema. You can use 'public' if you do not want to qualify webhookdb tables."},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgFdwInput{
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					MessageFdw:       c.Bool("fdw"),
					MessageViews:     c.Bool("views"),
					MessageAll:       c.Bool("all"),
					RemoteServerName: c.String("remote"),
					FetchSize:        c.String("fetch"),
					LocalSchema:      c.String("into-schema"),
					ViewSchema:       c.String("views-schema"),
				}
				out, err := client.OrgFdw(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				if c.Bool("raw") {
					fmt.Println(convext.MustMarshal(out))
				} else {
					fmt.Println(out["message"])
				}
				return nil
			}),
		},
	},
}
