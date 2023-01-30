package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var dbsyncCmd = &cli.Command{
	Name:    "dbsync",
	Aliases: []string{"sync"},
	Usage:   "Replicate data in a WebhookDB table into another database.",
	Subcommands: []*cli.Command{
		{
			Name: "create",
			Usage: "Create a database sync target for the specified integration. Data will be automatically synced from " +
				"the integration's WebhookDB table into the database specified by the connection string. " +
				"PostgresQL and SnowflakeDB databases are supported.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
				&cli.StringFlag{
					Name:    "connection-url",
					Aliases: s1("u"),
					Usage:   "The connection string for the database that WebhookDB should write data to.",
				},
				syncPeriodFlag(),
				syncSchemaFlag(),
				syncTableFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetCreateInput{
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					ConnectionUrl:         c.String("connection-url"),
					Period:                c.Int("period"),
					Schema:                c.String("schema"),
					Table:                 c.String("table"),
					SyncTypePathSlug:      "db",
				}
				out, err := client.SyncTargetCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), false)
				return nil
			}),
		},
		{
			Name:  "delete",
			Usage: "Delete the database sync target and stop any further syncing. The table being synced to is not modified.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag("db"),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetDeleteInput{
					OpaqueId:         getSyncTargetFlagOrArg(c, "db"),
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					SyncTypePathSlug: "db",
				}
				out, err := client.SyncTargetDelete(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name:  "list",
			Usage: "List all database sync targets for the given organization.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetListInput{
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					SyncTypePathSlug: "db",
				}
				out, err := client.SyncTargetList(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
		{
			Name:  "update",
			Usage: "Update the database sync target. Use `webhookdb dbsync list` to see all database sync targets.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag("db"),
				syncPeriodFlag(),
				syncSchemaFlag(),
				syncTableFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetUpdateInput{
					OpaqueId:         getSyncTargetFlagOrArg(c, "db"),
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					Period:           c.Int("period"),
					Schema:           c.String("schema"),
					Table:            c.String("table"),
					SyncTypePathSlug: "db",
				}
				out, err := client.SyncTargetUpdate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return nil
			}),
		},
		{
			Name:  "update-creds",
			Usage: "Update the username and password used to connect to the database being synced to.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag("db"),
				&cli.StringFlag{
					Name:    "user",
					Aliases: s1("u"),
					Usage:   "Database username",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: s1("p"),
					Usage:   "Takes a password.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetUpdateCredsInput{
					OpaqueId:         getSyncTargetFlagOrArg(c, "db"),
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					Username:         c.String("user"),
					Password:         c.String("password"),
					SyncTypePathSlug: "db",
				}
				out, err := client.SyncTargetUpdateCreds(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return nil
			}),
		},
		{
			Name: "trigger",
			Usage: "Trigger a database sync to the sync target. The database sync will happen at the earliest possible moment since " +
				"the last sync, no matter how long the configured period is on the sync target.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag("db"),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetSyncInput{
					OpaqueId:         getSyncTargetFlagOrArg(c, "db"),
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					SyncTypePathSlug: "db",
				}
				out, err := client.SyncTargetSync(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return nil
			}),
		},
	},
}
var httpsyncCmd = &cli.Command{
	Name:  "httpsync",
	Usage: "Sync changes to a WebhookDB table to a specified URL.",
	Subcommands: []*cli.Command{
		{
			Name: "create",
			Usage: "Create a HTTP sync target for the specified integration. Data will be automatically synced from " +
				"the integration's WebhookDB table to the url specified in the connection string. " +
				"Only https urls are supported.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
				&cli.StringFlag{
					Name:    "connection-url",
					Aliases: s1("u"),
					Usage:   "The connection string for the URL that WebhookDB should send data to.",
				},
				syncPeriodFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetCreateInput{
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					ConnectionUrl:         c.String("connection-url"),
					Period:                c.Int("period"),
					SyncTypePathSlug:      "http",
				}
				out, err := client.SyncTargetCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), false)
				return nil
			}),
		},
		{
			Name:  "delete",
			Usage: "Delete the http sync target and stop any further syncing.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag("http"),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetDeleteInput{
					OpaqueId:         getSyncTargetFlagOrArg(c, "http"),
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					SyncTypePathSlug: "http",
				}
				out, err := client.SyncTargetDelete(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message, false)
				return nil
			}),
		},
		{
			Name:  "list",
			Usage: "List all http sync targets for the given organization.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetListInput{
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					SyncTypePathSlug: "http",
				}
				out, err := client.SyncTargetList(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteCollection(c.App.Writer, out)
			}),
		},
		{
			Name:  "update",
			Usage: "Update the http sync target. Use `webhookdb httpsync list` to see all http sync targets.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag("http"),
				syncPeriodFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetUpdateInput{
					OpaqueId:         getSyncTargetFlagOrArg(c, "http"),
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					Period:           c.Int("period"),
					SyncTypePathSlug: "http",
				}
				out, err := client.SyncTargetUpdate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return nil
			}),
		},
		{
			Name:  "update-creds",
			Usage: "Update username and password used to access the http endpoint.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag("http"),
				&cli.StringFlag{
					Name:    "user",
					Aliases: s1("u"),
					Usage:   "Database username",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: s1("p"),
					Usage:   "Takes a password.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetUpdateCredsInput{
					OpaqueId:         getSyncTargetFlagOrArg(c, "http"),
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					Username:         c.String("user"),
					Password:         c.String("password"),
					SyncTypePathSlug: "http",
				}
				out, err := client.SyncTargetUpdateCreds(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return nil
			}),
		},
		{
			Name: "trigger",
			Usage: "Trigger a sync to the http sync target. The http sync will happen at the earliest possible moment since " +
				"the last sync, no matter how long the configured period is on the sync target.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag("http"),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetSyncInput{
					OpaqueId:         getSyncTargetFlagOrArg(c, "http"),
					OrgIdentifier:    getOrgFlag(c, ac.Prefs),
					SyncTypePathSlug: "http",
				}
				out, err := client.SyncTargetSync(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				printlnif(c, out.Message(), true)
				return nil
			}),
		},
	},
}

func syncPeriodFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name:  "period",
		Usage: "Number of seconds between syncs.",
	}
}

func syncSchemaFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "schema",
		Usage: "Schema (or namespace) to write the table into. Default to no schema/namespace.",
	}
}

func syncTableFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "table",
		Usage: "Table to create and update. Default to match the table name of the service integration.",
	}
}

func runListMsg(syncTypePathSlug string) string {
	if syncTypePathSlug == "db" {
		return "Use `webhookdb dbsync list` to see a list of all your database sync targets."
	} else {
		return "Use `webhookdb httpsync list` to see a list of all your http sync targets."
	}
}

func syncTargetFlag(syncTypePathSlug string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "target",
		Aliases: s1("t"),
		Usage:   "Sync target opaque id. " + runListMsg(syncTypePathSlug),
	}
}

func getSyncTargetFlagOrArg(c *cli.Context, syncTypePathSlug string) string {
	return requireFlagOrArg(c, "target", runListMsg(syncTypePathSlug))
}
