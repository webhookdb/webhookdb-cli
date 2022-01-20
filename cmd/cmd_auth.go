package cmd

import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/lithictech/webhookdb-cli/types"
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
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				output, err := client.AuthGetMe(ctx, ac.Auth)
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
			Flags:       []cli.Flag{usernameFlag(), tokenFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				out, err := client.AuthLogin(ctx, client.AuthLoginInput{
					Username: c.String("username"),
					Token:    c.String("token"),
				})
				step, err := client.NewStateMachine().RunWithOutput(ctx, ac.Auth, out.OutputStep)
				if err != nil {
					return err
				}

				// org information is coming in as a map[string]interface{}
				defaultOrg := types.Organization{}
				defaultOrgMap := step.Extras["current_customer"]["default_organization"]
				jsonString, err := json.Marshal(defaultOrgMap)
				if err != nil {
					return err
				}
				json.Unmarshal(jsonString, &defaultOrg)
				ac.Prefs.CurrentOrg = defaultOrg

				ac.Prefs.AuthCookie = out.AuthCookie
				ac.GlobalPrefs.SetNS(ac.Config.PrefsNamespace, ac.Prefs)
				if err := prefs.Save(ac.GlobalPrefs); err != nil {
					return err
				}
				return nil
			}),
		},
		{
			Name:        "logout",
			Description: "Log out of your current session.",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				output, err := client.AuthLogout(ctx, ac.Auth)
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
