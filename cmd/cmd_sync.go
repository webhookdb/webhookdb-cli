package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var synctargetCmd = &cli.Command{
	Name:  "sync",
	Usage: "Replicate data in a WebhookDB table into another database.",
	Subcommands: []*cli.Command{
		{
			Name: "create",
			Usage: "Create a sync target for the specified integration. Data will be automatically synced from " +
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
				syncPageSizeFlag(),
				syncPeriodFlag(),
				syncSchemaFlag(),
				syncTableFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetCreateInput{
					OrgIdentifier:         getOrgFlag(c, ac.Prefs),
					IntegrationIdentifier: getIntegrationFlagOrArg(c),
					ConnectionUrl:         c.String("connection-url"),
					PageSize:              c.Int("pagesize"),
					Period:                c.Int("period"),
					Schema:                c.String("schema"),
					Table:                 c.String("table"),
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
			Usage: "Delete the sync target and stop any further syncing. The table being synced to is not modified.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetDeleteInput{
					OpaqueId:      getSyncTargetFlagOrArg(c),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
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
			Usage: "List all sync targets for the given organization.",
			Flags: []cli.Flag{
				orgFlag(),
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetListInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
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
			Usage: "Update the sync target. Use `webhookdb sync list` to see all sync targets.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag(),
				syncPageSizeFlag(),
				syncPeriodFlag(),
				syncSchemaFlag(),
				syncTableFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetUpdateInput{
					OpaqueId:      getSyncTargetFlagOrArg(c),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					PageSize:      c.Int("pagesize"),
					Period:        c.Int("period"),
					Schema:        c.String("schema"),
					Table:         c.String("table"),
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
			Name: "update-creds",
			Usage: "Update credentials for the sync target. If the database URL used to sync is changing, " +
				"you must use update-creds so WebhookDB can continue to write to it.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag(),
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
					OpaqueId:      getSyncTargetFlagOrArg(c),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Username:      c.String("user"),
					Password:      c.String("password"),
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
			Usage: "Trigger a sync to the sync target. The sync will happen at the earliest possible moment since " +
				"the last sync, no matter how long the configured period is on the sync target.",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetSyncInput{
					OpaqueId:      getSyncTargetFlagOrArg(c),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
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

func syncPageSizeFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name:  "pagesize",
		Usage: "Max number of rows to retrieve in a single call. This is only relevant for HTTPS sync targets.",
	}
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

func syncTargetFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "target",
		Aliases: s1("t"),
		Usage:   "Sync target opaque id. Run `webhookdb sync list` to see a list of all your sync targets.",
	}
}

func getSyncTargetFlagOrArg(c *cli.Context) string {
	return requireFlagOrArg(c, "target", "Use `webhookdb sync list` to see a list of all your sync targets.")
}
