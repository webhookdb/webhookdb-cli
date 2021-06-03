package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/urfave/cli/v2"
	"strings"
)

var organizationsCmd = &cli.Command{
	Name: "org",
	Subcommands: []*cli.Command{
		{
			Name:        "invite",
			Description: "TODO",
			Flags:       []cli.Flag{orgFlag(), usernameFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				var orgKey string
				if c.String("org") != "" {
					orgKey = c.String("org")
				} else {
					orgKey = p.CurrentOrg
				}
				input := client.OrgInviteInput{
					AuthCookie: p.AuthCookie,
					Email: c.String("username"),
					OrgKey: orgKey,
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
			Description: "TODO",
			Flags:       []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				if c.NArg() != 1 {
					return errors.New("You must enter an invitation code")
				}
				input := client.OrgJoinInput{
					AuthCookie: p.AuthCookie,
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
			Description: "TODO",
			Flags:       []cli.Flag{},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				out, err := client.OrgList(ctx, client.OrgListInput{AuthCookie: p.AuthCookie})
				if err != nil {
					return err
				}
				orgsLen := len(out.Data)
				keySlugs := make([]string, orgsLen)
				for i, value := range out.Data {
					if value.Key == p.CurrentOrg {
						keySlugs[i] = (value.Key + " (active)")
					} else {
						keySlugs[i] = value.Key
					}
				}
				fmt.Println(strings.Join(keySlugs, "\n"))
				return nil
			}),
		},
		{
			Name:        "members",
			Description: "TODO",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				var orgKey string
				if c.String("org") != "" {
					orgKey = c.String("org")
				} else {
					orgKey = p.CurrentOrg
				}

				out, err := client.OrgMembers(ctx, client.OrgMembersInput{AuthCookie: p.AuthCookie, OrgKey: orgKey})
				if err != nil {
					return err
				}
				orgsLen := len(out.Data)
				members := make([]string, orgsLen)
				for i, value := range out.Data {
					if value.Status != "" {
						members[i] = (value.CustomerEmail + " (" + value.Status + ")")
					} else {
						members[i] = value.CustomerEmail
					}
				}
				fmt.Println(strings.Join(members, "\n"))
				return nil
			}),
		},
	},
}

func orgFlag() *cli.StringFlag {
	// takes the org ID
	return &cli.StringFlag{
		Name:     "org",
		Aliases:  s1("o"),
		Required: false,
		Usage:    "TODO",
	}
}