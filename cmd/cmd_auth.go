package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/urfave/cli/v2"
	"os"
)

var authCmd = &cli.Command{
	Name:        "auth",
	Description: "These commands control the auth process.",
	Subcommands: []*cli.Command{
		{
			Name:        "whoami",
			Description: "Print information about the current user",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				output, err := client.AuthGetMe(ctx)
				if err != nil {
					return err
				}
				output.PrintTo(os.Stdout)
				return nil
			}),
		},
		{
			Name:        "login",
			Description: "Sign up or log in.",
			Flags:       []cli.Flag{usernameFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				output, err := client.AuthLogin(ctx, client.AuthLoginInput{
					Username: c.String("username"),
				})
				if err != nil {
					return err
				}
				fmt.Println(output.Message)
				return nil
			}),
		},
		{
			Name:        "otp",
			Description: "Finish sign up or login in using the given One Time Password.",
			Flags:       []cli.Flag{usernameFlag(), tokenFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				output, err := client.AuthOTP(ctx, client.AuthOTPInput{
					Username: c.String("username"),
					Token:    c.String("token"),
				})
				if err != nil {
					return err
				}
				p.AuthCookie = output.AuthCookie
				p.CurrentOrg = output.CurrentCustomer.DefaultOrganization
				ac.GlobalPrefs.SetNS(ac.Config.PrefsNamespace, p)
				if err := prefs.Save(ac.GlobalPrefs); err != nil {
					return err
				}
				fmt.Println(output.Message)
				return nil
			}),
		},
		{
			Name:        "logout",
			Description: "Log out of your current session.",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context, p prefs.Prefs) error {
				output, err := client.AuthLogout(ctx)
				if err != nil {
					return err
				}
				ac.GlobalPrefs.ClearNS(ac.Config.PrefsNamespace)
				if err := prefs.Save(ac.GlobalPrefs); err != nil {
					return err
				}
				fmt.Println(output.Message)
				return nil
			}),
		},
	},
}
