package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/urfave/cli/v2"
)

var authCmd = &cli.Command{
	Name:        "auth",
	Description: "These commands control the auth process.",
	Subcommands: []*cli.Command{
		{
			Name:        "login",
			Description: "logs a user in, sends them an otp",
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
			Description: "registers the user's otp",
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
				p.CurrentOrg = output.DefaultOrg
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
			Description: "logs the current user out",
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
