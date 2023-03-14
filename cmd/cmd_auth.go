package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/mitchellh/mapstructure"
	"github.com/urfave/cli/v2"
	"strings"
)

var authCmd = &cli.Command{
	Name:  "auth",
	Usage: "These commands control the auth process.",
	Subcommands: []*cli.Command{
		{
			Name:  "whoami",
			Usage: "Print information about the current user.",
			Flags: []cli.Flag{
				formatFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.AuthGetMe(ctx, ac.Auth)
				if err != nil {
					return err
				}
				wasmUpdateAuthDisplay(ac.Prefs)
				printlnif(c, out.Message(), true)
				return getFormatFlag(c).WriteSingle(c.App.Writer, out)
			}),
		},
		{
			Name:    "login",
			Aliases: []string{"signup", "signin", "register"},
			Usage:   "Sign up or log in.",
			Flags: []cli.Flag{
				usernameFlag(),
				&cli.StringFlag{
					Name:    "token",
					Aliases: s1("t"),
					Usage:   "One-time-password token received in your email after running 'auth login'.",
				}},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				authOut, err := client.AuthLogin(ctx, client.AuthLoginInput{
					Username: flagOrArg(c, "username"),
					Token:    c.String("token"),
				})
				if err != nil {
					return err
				}

				var result client.Step
				if authOut.Complete {
					result = authOut
				} else {
					result, err = client.NewStateMachine().Run(ctx, ac.Auth, authOut)
					if err != nil {
						return err
					}
				}

				// If the state machine finished running, it was probably successful, but maybe not.
				// If somehow we didn't get real current_customer info, assume it failed, and exit.
				if result.Extras["current_customer"] == nil || result.Extras["current_customer"]["email"] == nil {
					return nil
				}

				// Org information is coming in as a map[string]interface{}
				defaultOrg := types.Organization{}
				if err := mapstructure.Decode(result.Extras["current_customer"]["default_organization"], &defaultOrg); err != nil {
					return err
				}
				ac.Prefs.CurrentOrg = defaultOrg
				authTokHeader := result.RawResponse.Header().Get(client.AuthTokenHeader)
				ac.Prefs.AuthToken = types.AuthToken(strings.Split(authTokHeader, ";")[0])
				ac.Prefs.Email = result.Extras["current_customer"]["email"].(string)
				if err := ac.SavePrefs(); err != nil {
					return err
				}
				wasmUpdateAuthDisplay(ac.Prefs)
				return nil
			}),
		},
		{
			Name:  "logout",
			Usage: "Log out of your current session.",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "remove", Aliases: s1("r"), Usage: "If given, remove all WebhookDB preferences."},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.AuthLogout(ctx, ac.Auth)
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
				printlnif(c, out.Message, false)
				wasmUpdateAuthDisplay(prefs.Prefs{})
				return nil
			}),
		},
	},
}
