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
	Description: "TODO",
	Subcommands: []*cli.Command{
		{
			Name:        "login",
			Description: "TODO",
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
			Description: "TODO",
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
			Description: "TODO",
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

func usernameFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "username",
		Aliases:  s1("u"),
		Required: true,
		Usage:    "TODO",
	}
}

func tokenFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "token",
		Aliases:  s1("t"),
		Required: true,
		Usage:    "TODO",
	}
}
