package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

type syncType struct {
	Slug                 string
	Aliases              []string
	FullName             string
	Destination          string
	SupportedProtocolMsg string
	UniqueFieldFlags     []cli.Flag
}

func (st syncType) Cmd() string {
	return fmt.Sprintf("%ssync", st.Slug)
}

var dbSyncType = syncType{
	Slug:                 "db",
	Aliases:              []string{"sync"},
	FullName:             "database",
	Destination:          "database",
	SupportedProtocolMsg: "PostgresQL and SnowflakeDB databases are supported.",
	UniqueFieldFlags:     []cli.Flag{syncSchemaFlag(), syncTableFlag()},
}

var httpSyncType = syncType{
	Slug:                 "http",
	Aliases:              []string{},
	FullName:             "http",
	Destination:          "http endpoint",
	SupportedProtocolMsg: "Only https urls are supported.",
	UniqueFieldFlags:     []cli.Flag{syncPageSizeFlag()},
}

func syncCmd(st syncType) *cli.Command {
	return &cli.Command{
		Name:    st.Cmd(),
		Aliases: st.Aliases,
		Usage:   fmt.Sprintf("Replicate data in a WebhookDB table to another %s.", st.Destination),
		Subcommands: []*cli.Command{
			{
				Name: "create",
				Usage: fmt.Sprintf("Create a %s sync target for the specified integration. Data will be "+
					"automatically synced from the integration's WebhookDB table into the %s specified by the "+
					"connection string. %s",
					st.FullName, st.Destination, st.SupportedProtocolMsg),
				Flags: append(
					[]cli.Flag{
						orgFlag(),
						integrationFlag(),
						&cli.StringFlag{
							Name:    "connection-url",
							Aliases: s1("u"),
							Usage: fmt.Sprintf(
								"The connection string for the %s that WebhookDB should write data to.", st.Destination),
						},
						syncPeriodFlag(),
					},
					st.UniqueFieldFlags...,
				),
				Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
					input := client.SyncTargetCreateInput{
						OrgIdentifier:         getOrgFlag(c, ac.Prefs),
						IntegrationIdentifier: getIntegrationFlagOrArg(c),
						ConnectionUrl:         c.String("connection-url"),
						PageSize:              c.Int("pagesize"),
						Period:                c.Int("period"),
						Schema:                c.String("schema"),
						Table:                 c.String("table"),
						SyncTypeSlug:          st.Slug,
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
				Name: "delete",
				Usage: fmt.Sprintf(
					"Delete the %s sync target and stop any further syncing. The %s being synced to is not modified.",
					st.FullName, st.Destination),
				Flags: []cli.Flag{
					orgFlag(),
					syncTargetFlag(st),
				},
				Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
					input := client.SyncTargetDeleteInput{
						OpaqueId:      getSyncTargetFlagOrArg(c, st),
						OrgIdentifier: getOrgFlag(c, ac.Prefs),
						SyncTypeSlug:  st.Slug,
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
				Usage: fmt.Sprintf("List all %s sync targets for the given organization.", st.FullName),
				Flags: []cli.Flag{
					orgFlag(),
					formatFlag(),
				},
				Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
					input := client.SyncTargetListInput{
						OrgIdentifier: getOrgFlag(c, ac.Prefs),
						SyncTypeSlug:  st.Slug,
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
				Name: "update",
				Usage: fmt.Sprintf(
					"Update the %s sync target. Use `webhookdb %s list` to see all %s sync targets.",
					st.FullName, st.Cmd(), st.FullName),
				Flags: append(
					[]cli.Flag{
						orgFlag(),
						syncTargetFlag(st),
						syncPeriodFlag(),
					},
					st.UniqueFieldFlags...),
				Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
					input := client.SyncTargetUpdateInput{
						OpaqueId:      getSyncTargetFlagOrArg(c, st),
						OrgIdentifier: getOrgFlag(c, ac.Prefs),
						PageSize:      c.Int("pagesize"),
						Period:        c.Int("period"),
						Schema:        c.String("schema"),
						Table:         c.String("table"),
						SyncTypeSlug:  st.Slug,
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
				Usage: fmt.Sprintf(
					"Update the username and password used to connect to the %s being synced to.",
					st.Destination),
				Flags: []cli.Flag{
					orgFlag(),
					syncTargetFlag(st),
					&cli.StringFlag{
						Name:    "user",
						Aliases: s1("u"),
						Usage:   "Takes a username.",
					},
					&cli.StringFlag{
						Name:    "password",
						Aliases: s1("p"),
						Usage:   "Takes a password.",
					},
				},
				Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
					input := client.SyncTargetUpdateCredsInput{
						OpaqueId:      getSyncTargetFlagOrArg(c, st),
						OrgIdentifier: getOrgFlag(c, ac.Prefs),
						Username:      c.String("user"),
						Password:      c.String("password"),
						SyncTypeSlug:  st.Slug,
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
				Usage: fmt.Sprintf("Trigger a %s sync to the sync target. The %s sync will happen at the "+
					"earliest possible moment since the last sync, no matter how long the configured period is on the "+
					"sync target.", st.FullName, st.FullName),
				Flags: []cli.Flag{
					orgFlag(),
					syncTargetFlag(st),
				},
				Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
					input := client.SyncTargetSyncInput{
						OpaqueId:      getSyncTargetFlagOrArg(c, st),
						OrgIdentifier: getOrgFlag(c, ac.Prefs),
						SyncTypeSlug:  st.Slug,
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
}

func syncPageSizeFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name:        "pagesize",
		Usage:       "Max number of rows WebhookDB sends the sync target in each call.",
		DefaultText: "unset",
	}
}

func syncPeriodFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name:        "period",
		Usage:       "Number of seconds between syncs.",
		DefaultText: "unset",
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

func syncTargetFlag(input syncType) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "target",
		Aliases: s1("t"),
		Usage: fmt.Sprintf(
			"Sync target opaque id. Use `webhookdb %s list` to see a list of all your %s sync targets.",
			input.Cmd(), input.FullName),
	}
}

func getSyncTargetFlagOrArg(c *cli.Context, input syncType) string {
	return requireFlagOrArg(
		c,
		"target",
		fmt.Sprintf(
			"Use `webhookdb %s list` to see a list of all your %s sync targets.",
			input.Cmd(), input.FullName),
	)
}
