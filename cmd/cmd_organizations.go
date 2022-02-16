package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/urfave/cli/v2"
	"strings"
)

var organizationsCmd = &cli.Command{
	Name:        "org",
	Description: "Create and activate an organization, invite new members, and change membership roles.",
	Subcommands: []*cli.Command{
		{
			Name:        "activate",
			Description: "Change the default organization for any command you run",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.OrgGet(ctx, ac.Auth, client.OrgGetInput{
					OrgIdentifier: types.OrgIdentifierFromSlug(flagOrArg(c, "org", "Run `webhookdb org list` to see available orgs.")),
				})
				if err != nil {
					return err
				}
				if err := prefs.SetNSAndSave(ac.GlobalPrefs, ac.Config.PrefsNamespace, ac.Prefs.ChangeOrg(out.Org)); err != nil {
					return err
				}
				fmt.Println(fmt.Sprintf("%s is now your active organization. ", out.Org.DisplayString()))
				return nil
			}),
		},
		{
			Name:        "changerole",
			Description: "Change the role of members of your organization.",
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
				fmt.Println(out.Message)
				return nil
			}),
		},
		{
			Name:        "create",
			Description: "Create and activate an organization.",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "name", Aliases: s1("n"), Usage: "Name of the new organization. The unique key for the org is derived from this name."},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgCreateInput{
					OrgName: c.String("name"),
				}
				out, err := client.OrgCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Println(out.Message)
				ac.Prefs.CurrentOrg = out.Organization
				if err := prefs.SetNSAndSave(ac.GlobalPrefs, ac.Config.PrefsNamespace, ac.Prefs); err != nil {
					return err
				}
				return nil
			}),
		},
		{
			Name:        "invite",
			Description: "Invite a user to your organization",
			Flags:       []cli.Flag{orgFlag(), usernameFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgInviteInput{
					Email:         c.String("username"),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.OrgInvite(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Println(out.Message)
				return nil
			}),
		},
		{
			Name:        "join",
			Description: "join an organization using a join code",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "code", Aliases: s1("c"), Usage: "Invitation code sent to the new member. Has 'join-' prefix."},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgJoinInput{
					InvitationCode: c.String("code"),
				}
				out, err := client.OrgJoin(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Println(out.Message)
				return nil
			}),
		},
		{
			Name:        "list",
			Description: "List all organizations that you are a member of.",
			Flags:       []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.OrgList(ctx, ac.Auth, client.OrgListInput{})
				if err != nil {
					return err
				}
				orgsLen := len(out.Items)
				keySlugs := make([]string, orgsLen)
				for i, value := range out.Items {
					line := fmt.Sprintf("%s\t%s", value.Name, value.Key)
					if value.Id == ac.Prefs.CurrentOrg.Id {
						line += " (active)"
					}
					keySlugs[i] = line
				}
				fmt.Println(strings.Join(keySlugs, "\n"))
				return nil
			}),
		},
		{
			Name:        "current",
			Description: "display the name and slug of the currently active org",
			Flags:       []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				fmt.Println(ac.Prefs.CurrentOrg.DisplayString())
				return nil
			}),
		},
		{
			Name:        "members",
			Description: "list all members of the given organization",
			Flags:       []cli.Flag{orgFlag()},
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
				fmt.Println(strings.Join(members, "\n"))
				return nil
			}),
		},
		{
			Name:        "remove",
			Description: "Remove a member from an organization",
			Flags:       []cli.Flag{orgFlag(), usernameFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgRemoveInput{
					Email:         c.String("username"),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.OrgRemove(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Println(out.Message)
				return nil
			}),
		},
	},
}
