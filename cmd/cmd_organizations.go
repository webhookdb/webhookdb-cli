package cmd

import (
	"context"
	"errors"
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				if c.NArg() != 1 {
					return errors.New("You must enter an organization key.")
				}
				orgSlug := c.Args().Get(0)
				out, err := client.OrgGet(ctx, client.OrgGetInput{
					AuthCookie:    p.AuthCookie,
					OrgIdentifier: types.OrgIdentifierFromSlug(orgSlug),
				})
				if err != nil {
					return err
				}
				if err := prefs.Save(p.ChangeOrg(out.Org)); err != nil {
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				input := client.OrgChangeRolesInput{
					AuthCookie: p.AuthCookie,
					Emails: c.StringSlice("usernames"),
					OrgIdentifier: getOrgFlag(c, p),
					RoleName: c.String("role"),
				}
				out, err := client.OrgChangeRoles(ctx, input)
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				orgName, err := ask.Ask("What is your organization name? ")
				if err != nil {
					return err
				}
				input := client.OrgCreateInput{
					AuthCookie: p.AuthCookie,
					OrgName:    orgName,
				}
				out, err := client.OrgCreate(ctx, input)
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				input := client.OrgInviteInput{
					AuthCookie:    p.AuthCookie,
					Email:         c.String("username"),
					OrgIdentifier: getOrgFlag(c, p),
				}
				out, err := client.OrgInvite(ctx, input)
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				if c.NArg() != 1 {
					return errors.New("You must enter an invitation code.")
				}
				input := client.OrgJoinInput{
					AuthCookie:     p.AuthCookie,
					InvitationCode: c.Args().Get(0),
				}
				out, err := client.OrgJoin(ctx, input)
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				out, err := client.OrgList(ctx, client.OrgListInput{AuthCookie: p.AuthCookie})
				if err != nil {
					return err
				}
				orgsLen := len(out.Items)
				keySlugs := make([]string, orgsLen)
				for i, value := range out.Items {
					if value.Id == p.CurrentOrg.Id {
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				fmt.Println(p.CurrentOrg.DisplayString())
				return nil
			}),
		},
		{
			Name:        "members",
			Description: "list all members of the given organization",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				out, err := client.OrgMembers(ctx, client.OrgMembersInput{AuthCookie: p.AuthCookie, OrgIdentifier: getOrgFlag(c, p)})
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				input := client.OrgRemoveInput{
					AuthCookie:    p.AuthCookie,
					Email:         c.String("username"),
					OrgIdentifier: getOrgFlag(c, p),
				}
				out, err := client.OrgRemove(ctx, input)
				if err != nil {
					return err
				}
				fmt.Println(out.Message)
				return nil
			}),
		},
	},
}
