package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/mitchellh/mapstructure"
	"github.com/urfave/cli/v2"
	"strings"
)

var authCmd = &cli.Command{
	Name:        "auth",
	Description: "These commands control the auth process.",
	Subcommands: []*cli.Command{
		{
			Name:        "whoami",
			Description: "Print information about the current user",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				output, err := client.AuthGetMe(ctx, ac.Auth)
				if err != nil {
					return err
				}
				output.PrintTo(c.App.Writer)
				return nil
			}),
		},
		{
			Name:        "login",
			Aliases:     []string{"signup", "signin", "register"},
			Description: "Sign up or log in.",
			Flags: []cli.Flag{
				usernameFlag(),
				&cli.StringFlag{
					Name:    "token",
					Aliases: s1("t"),
					Usage:   "One-time-password token received in your email after running 'auth login'.",
				}},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.AuthLogin(ctx, client.AuthLoginInput{
					Username: c.String("username"),
					Token:    c.String("token"),
				})
				if err != nil {
					return err
				}
				step, err := client.NewStateMachine().RunWithOutput(ctx, ac.Auth, out)
				if err != nil {
					return err
				}

				// org information is coming in as a map[string]interface{}
				defaultOrg := types.Organization{}
				if err := mapstructure.Decode(step.Extras["current_customer"]["default_organization"], &defaultOrg); err != nil {
					return err
				}
				ac.Prefs.CurrentOrg = defaultOrg

				authTokHeader := out.RawResponse.Header().Get(client.AuthTokenHeader)
				ac.Prefs.AuthToken = types.AuthToken(strings.Split(authTokHeader, ";")[0])
				if err := ac.SavePrefs(); err != nil {
					return err
				}
				return nil
			}),
		},
		{
			Name:        "logout",
			Description: "Log out of your current session.",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "remove", Aliases: s1("r"), Usage: "If given, remove all WebhookDB preferences."},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				output, err := client.AuthLogout(ctx, ac.Auth)
				if err != nil {
					return err
				}
				if c.Bool("remove") {
					if err := prefs.DeleteAll(ac.FS); err != nil {
						return err
					}
				} else {
					ac.GlobalPrefs.ClearNS(ac.Config.PrefsNamespace)
					if err := prefs.Save(ac.FS, ac.GlobalPrefs); err != nil {
						return err
					}
				}
				fmt.Fprintln(c.App.Writer, output.Message)
				return nil
			}),
		},
	},
}
