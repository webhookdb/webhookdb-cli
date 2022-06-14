package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var synctargetCmd = &cli.Command{
	Name:  "sync-target",
	Usage: "",
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Usage: "Create a sync target for the given organization.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
				connectionUrlFlag(),
				periodFlag(),
				schemaFlag(),
				tableFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetCreateInput{
					OrgIdentifier:       getOrgFlag(c, ac.Prefs),
					IntegrationOpaqueId: c.String("integration"),
					ConnectionUrl:       c.String("connection-url"),
					Period:              c.Int("period"),
					Schema:              c.String("schema"),
					Table:               c.String("table"),
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
			Usage: "Delete a sync target.",
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
			Usage: "Update a sync target",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag(),
				periodFlag(),
				schemaFlag(),
				tableFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.SyncTargetUpdateInput{
					OpaqueId:      getSyncTargetFlagOrArg(c),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
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
			Name:  "update-creds",
			Usage: "Update credentials for a sync target",
			Flags: []cli.Flag{
				orgFlag(),
				syncTargetFlag(),
				&cli.StringFlag{
					Name:    "user",
					Aliases: s1("u"),
					Usage:   "Database username",
				},
				passwordFlag(),
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
			Name:  "sync",
			Usage: "Sync a sync target",
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
