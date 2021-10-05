package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/ask"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/urfave/cli/v2"
	"strings"
)

var organizationsCmd = &cli.Command{
	Name:        "org",
	Description: "To set up integrations, you need to be part of an Organization. These commands will allow you to see and manipulate membership status for your organization.",
	Subcommands: []*cli.Command{
		{
			Name:        "activate",
			Description: "change the default organization for any command you run",
			Flags:       []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				orgSlug, err := extractPositional(0, c, "You must enter an organization key.")
				if err != nil {
					return err
				}
				out, err := client.OrgGet(ctx, ac.Auth, client.OrgGetInput{
					OrgIdentifier: types.OrgIdentifierFromSlug(orgSlug),
				})
				if err != nil {
					return err
				}
				ac.GlobalPrefs.SetNS(ac.Config.PrefsNamespace, ac.Prefs.ChangeOrg(out.Org))
				if err := prefs.Save(ac.GlobalPrefs); err != nil {
					return err
				}
				fmt.Println(fmt.Sprintf("%s is now your active organization. ", out.Org.DisplayString()))
				return nil
			}),
		},
		{
			Name:        "changerole",
			Description: "TODO",
			Flags:       []cli.Flag{roleFlag(), usernamesFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.OrgChangeRolesInput{
					Emails:        c.String("usernames"),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					RoleName:      c.String("role"),
				}
				out, err := client.OrgChangeRoles(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Println(out)
				return nil
			}),
		},
		{
			Name:        "create",
			Description: "create an organization",
			Flags:       []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				orgName, err := ask.Ask("What is your organization name? ")
				if err != nil {
					return err
				}
				input := client.OrgCreateInput{
					OrgName: orgName,
				}
				out, err := client.OrgCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				fmt.Println(out.Message)
				// Do we want to activate the org too?
				return nil
			}),
		},
		{
			Name:        "invite",
			Description: "invite a user to your organization",
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
			Flags:       []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				invCode, err := extractPositional(0, c, "You must enter an invitation code.")
				if err != nil {
					return err
				}
				input := client.OrgJoinInput{
					InvitationCode: invCode,
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
			Description: "list all organizations that you're a member of",
			Flags:       []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.OrgList(ctx, ac.Auth, client.OrgListInput{})
				if err != nil {
					return err
				}
				orgsLen := len(out.Items)
				keySlugs := make([]string, orgsLen)
				for i, value := range out.Items {
					if value.Id == ac.Prefs.CurrentOrg.Id {
						keySlugs[i] = value.Name + " (active)"
					} else {
						keySlugs[i] = value.Name
					}
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
			Description: "remove a member from an organization",
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
