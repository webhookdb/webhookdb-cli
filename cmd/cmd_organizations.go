package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/urfave/cli/v2"
	"strings"
)

var organizationsCmd = &cli.Command{
	Name:  "org",
	Usage: "Create and activate an organization, invite new members, and change membership roles.",
	Subcommands: []*cli.Command{
		{
			Name:  "activate",
			Usage: "Change the default organization for any command you run.",
			Flags: []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.OrgGet(ctx, ac.Auth, client.OrgGetInput{
					OrgIdentifier: types.OrgIdentifierFromSlug(requireFlagOrArg(c, "org", "Run `webhookdb org list` to see available orgs.")),
				})
				if err != nil {
					return err
				}
				ac.Prefs = ac.Prefs.ChangeOrg(out.Org)
				if err := ac.SavePrefs(); err != nil {
					return err
				}
				fmt.Fprintln(c.App.Writer, fmt.Sprintf("%s is now your active organization. ", out.Org.DisplayString()))
				wasmUpdateAuthDisplay(ac.Prefs)
				return nil
			}),
		},
		{
			Name:  "changerole",
			Usage: "Change the role of members of your organization.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "usernames",
					Aliases: nil,
					Usage:   "Takes multiple emails.",
				},
				&cli.StringFlag{
					Name:    "role",
					Aliases: s1("r"),
					Usage:   "Role name, like 'member' or 'admin'.",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgChangeRolesInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Emails:        c.String("usernames"),
					RoleName:      c.String("role"),
				}
				out, err := client.OrgChangeRoles(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Fprintln(c.App.Writer, out.Message)
				return nil
			}),
		},
		{
			Name:  "close",
			Usage: "Close down this organization.",
			Flags: []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgCloseInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)}
				return stateMachineResponseRunner(ctx, ac.Auth)(client.OrgClose(ctx, ac.Auth, input))
			}),
		},
		{
			Name:  "create",
			Usage: "Create and activate an organization.",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "name", Aliases: s1("n"), Usage: "Name of the new organization. The unique key for the org is derived from this name."},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgCreateInput{
					OrgName: requireFlagOrArg(c, "name", ""),
				}
				out, err := client.OrgCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Fprintln(c.App.Writer, out.Message)
				ac.Prefs.CurrentOrg = out.Organization
				if err := ac.SavePrefs(); err != nil {
					return err
				}
				return nil
			}),
		},
		{
			Name:  "invite",
			Usage: "Invite a user to your organization.",
			Flags: []cli.Flag{orgFlag(), usernameFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgInviteInput{
					Email:         requireFlagOrArg(c, "username", ""),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.OrgInvite(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Fprintln(c.App.Writer, out.Message)
				return nil
			}),
		},
		{
			Name:  "join",
			Usage: "join an organization using a join code.",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "code", Aliases: s1("c"), Usage: "Invitation code sent to the new member. Has 'join-' prefix."},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgJoinInput{
					InvitationCode: requireFlagOrArg(c, "code", ""),
				}
				out, err := client.OrgJoin(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Fprintln(c.App.Writer, out.Message)
				return nil
			}),
		},
		{
			Name:  "list",
			Usage: "List all organizations that you are a member of.",
			Flags: []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.MeOrgMembershipsInput{
					ActiveOrgIdentifier: types.OrgIdentifierFromSlug(ac.Prefs.CurrentOrg.Key),
				}
				out, err := client.MeOrgMemberships(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				_, err = out.Blocks.WriteTo(c.App.Writer)
				return err
			}),
		},
		{
			Name:  "current",
			Usage: "Display the name and slug of the currently active organization.",
			Flags: []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				fmt.Fprintln(c.App.Writer, ac.Prefs.CurrentOrg.DisplayString())
				return nil
			}),
		},
		{
			Name:  "members",
			Usage: "List all members of the given organization.",
			Flags: []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.OrgMembers(ctx, ac.Auth, client.OrgMembersInput{OrgIdentifier: getOrgFlag(c, ac.Prefs)})
				if err != nil {
					return err
				}
				orgsLen := len(out.Data)
				members := make([]string, orgsLen)
				for i, value := range out.Data {
					if value.Status != "" {
						members[i] = value.CustomerEmail + " (" + value.Status + ")"
					} else {
						members[i] = value.CustomerEmail
					}
				}
				fmt.Fprintln(c.App.Writer, strings.Join(members, "\n"))
				return nil
			}),
		},
		{
			Name:  "remove",
			Usage: "Remove a member from an organization.",
			Flags: []cli.Flag{orgFlag(), usernameFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgRemoveInput{
					Email:         requireFlagOrArg(c, "username", ""),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.OrgRemove(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Fprintln(c.App.Writer, out.Message)
				return nil
			}),
		},
		{
			Name:  "rename",
			Usage: "Change the name of the organization. Does not change the org key, which is immutable.",
			Flags: []cli.Flag{
				orgFlag(),
				&cli.StringFlag{
					Name:    "name",
					Aliases: s1("n"),
					Usage:   "New name of the organization",
				},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgRenameInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					Name:          flagOrArg(c, "name"),
				}
				out, err := client.OrgRename(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Fprintln(c.App.Writer, out.Message)
				return nil
			}),
		},
	},
}
