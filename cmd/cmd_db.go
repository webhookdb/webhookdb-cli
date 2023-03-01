package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/go-aperitif/convext"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
	"strings"
)

var dbCmd = &cli.Command{
	Name:  "db",
	Usage: "Command namespace for interacting with your organization's database and tables.",
	Subcommands: []*cli.Command{
		{
			Name:  "connection",
			Usage: "Print the database connection url for an organization.",
			Flags: []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.DbConnection(ctx, ac.Auth, client.DbOrgIdentifierInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
				if err != nil {
					return err
				}
				fmt.Fprint(c.App.Writer, out.ConnectionUrl)
				return nil
			}),
		},
		{
			Name:  "tables",
			Usage: "List all tables in an organization's database.",
			Flags: []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.DbTables(ctx, ac.Auth, client.DbOrgIdentifierInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
				if err != nil {
					return err
				}
				printlnif(c, out.Message, true)
				printlnif(c, strings.Join(out.TableNames, "\n"), false)
				return nil
			}),
		},
		{
			Name:  "sql",
			Usage: "Execute query on organization's database.",
			Flags: []cli.Flag{
				orgFlag(),
				&cli.StringFlag{Name: "query", Aliases: s1("u"), Usage: "Query string to execute using your connection."},
				&cli.BoolFlag{Name: "color", Aliases: s1("c"), Usage: "Display colors. Default true if tty.", Value: IsTty},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				useColors := c.Bool("color")
				input := client.DbSqlInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Query:         extractPositional(0, c, "You must enter a query string."),
				}
				out, err := client.DbSql(ctx, ac.Auth, input)
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
			Name:  "roll-credentials",
			Usage: "Roll the credentials for an organization's database to something newly randomly generated.",
			Flags: []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.DbRollCredentials(ctx, ac.Auth, client.DbOrgIdentifierInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name: "fdw",
			Usage: "Write out commands that can be used to generate a FDW against your WebhookDB database and " +
				"import them into materialized views. See flags for further usage.",
			Flags: []cli.Flag{
				orgFlag(),
				&cli.BoolFlag{Name: "raw", Usage: "If given, print the raw SQL returned from the server. Useful if you want to pipe through jq or something similar."},
				&cli.BoolFlag{Name: "fdw", Usage: "Write the FDW SQL to stdout"},
				&cli.BoolFlag{Name: "views", Usage: "Write the SQL to create the materialized views to stdout"},
				&cli.BoolFlag{Name: "all", Usage: "Write a single SQL statement containing FDW and view creation code. Default if neither --fdw or --views are passed."},
				&cli.StringFlag{Name: "remote", Value: "webhookdb_remote", Usage: usage("The remote server name, used in the `CREATE SERVER <remote>` call")},
				&cli.StringFlag{Name: "fetch", Value: "50000", Usage: "fetch_size option used during server creation"},
				&cli.StringFlag{Name: "into-schema", Value: "webhookdb_remote", Usage: usage("Name of the schema to import the remote tables into (in `IMPORT FOREIGN SCHEMA public INTO <into schema>` call).")},
				&cli.StringFlag{Name: "views-schema", Value: "webhookdb", Usage: "Create materialized views in this schema. You can use 'public' if you do not want to qualify webhookdb tables."},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.DbFdwInput{
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					MessageFdw:       c.Bool("fdw"),
					MessageViews:     c.Bool("views"),
					MessageAll:       c.Bool("all"),
					RemoteServerName: c.String("remote"),
					FetchSize:        c.String("fetch"),
					LocalSchema:      c.String("into-schema"),
					ViewSchema:       c.String("views-schema"),
				}
				out, err := client.DbFdw(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				if c.Bool("raw") {
					fmt.Fprintln(c.App.Writer, convext.MustMarshal(out))
				} else {
					fmt.Fprintln(c.App.Writer, out["message"])
				}
				return nil
			}),
		},
		{
			Name:  "rename-table",
			Usage: "Rename the database table associated with the integration.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
				&cli.StringFlag{
					Name:    "new-name",
					Aliases: s1("n"),
					Usage:   "The new name of the table. " + tableNameRules,
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.DbRenameTableInput{
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
					NewName:               c.String("new-name"),
				}
				out, err := client.DbRenameTable(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), false)
				return nil
			}),
		},
		{
			Name:  "migrations",
			Usage: "Command namespace for interacting with your organizations database migrations.",
			Subcommands: []*cli.Command{
				{
					Name:  "start",
					Usage: "Enqueue a migration of all your organization's data to a new database.",
					Flags: []cli.Flag{
						orgFlag(),
						&cli.StringFlag{
							Name:     "admin-url",
							Aliases:  s1("a"),
							Required: false,
						},
						&cli.StringFlag{
							Name:     "readonly-url",
							Aliases:  s1("r"),
							Required: false,
						},
					},
					Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
						input := client.DbMigrationsStartInput{
							OrgIdentifier: getOrgFlag(c, ac.Prefs),
							AdminUrl:      c.String("admin-url"),
							ReadonlyUrl:   stringPtrFlag(c, "readonly-url"),
						}
						out, err := client.DbMigrationsStart(ctx, ac.Auth, input)
						if err != nil {
							return err
						}
						printlnif(c, out.Message, false)
						return nil
					}),
				},
				{
					Name:  "list",
					Usage: "List all database migrations.",
					Flags: []cli.Flag{
						orgFlag(),
						formatFlag(),
					},
					Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
						out, err := client.DbMigrationsList(ctx, ac.Auth, client.DbOrgIdentifierInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
						if err != nil {
							return err
						}
						printlnif(c, out.Message(), true)
						return getFormatFlag(c).WriteCollection(c.App.Writer, out)
					}),
				},
			},
		},
	},
}
