package cmd

import (
	"fmt"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/urfave/cli/v2"
)

const PASSWORD_RETRY_ATTEMPTS = 3

var authCmd = &cli.Command{
	Name:        "auth",
	Description: "These commands control the auth process.",
	Subcommands: []*cli.Command{
		{
			Name:        "login",
			Description: "logs a user in, sends them an otp",
			Flags:       []cli.Flag{usernameFlag()},
			Action: func(c *cli.Context) error {
				ctx := newCtx(newAppCtx(c))
				output, err := client.AuthLogin(ctx, client.AuthLoginInput{
					Username: c.String("username"),
				})
				if err != nil {
					return err
				}
				fmt.Println(output.Message)
				return nil
			},
		},
		{
			Name:        "otp",
			Description: "registers the user's otp",
			Flags:       []cli.Flag{usernameFlag(), tokenFlag()},
			Action: func(c *cli.Context) error {
				ctx := newCtx(newAppCtx(c))
				output, err := client.AuthOTP(ctx, client.AuthOTPInput{
					Username: c.String("username"),
					Token:    c.String("token"),
				})
				if err != nil {
					return err
				}
				p := prefs.Prefs{
					AuthCookie: output.AuthCookie,
					CurrentOrg: output.DefaultOrgKey,
				}
				prefs.Save(p)
				fmt.Println(output.Message)
				return nil
			},
		},
		{
			Name:        "logout",
			Description: "logs the current user out",
			Action: func(c *cli.Context) error {
				ctx := newCtx(newAppCtx(c))
				output, err := client.AuthLogout(ctx)
				if err != nil {
					return err
				}
				prefs.Delete()
				fmt.Println(output.Message)
				return nil
			},
		},
	},
}
