package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lithictech/go-aperitif/convext"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/ask"
	"github.com/webhookdb/webhookdb-cli/client"
	"github.com/webhookdb/webhookdb-cli/types"
	"time"
)

var debugCmd = &cli.Command{
	Name:   "debug",
	Hidden: true,
	Subcommands: []*cli.Command{
		{
			Name: "config",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				j, err := json.MarshalIndent(ac.Config, "", "  ")
				convext.Must(err)
				fmt.Fprintln(c.App.Writer, string(j))
				return nil
			}),
		},
		{
			Name: "prompt",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				a := ask.New()
				v, err := a.Ask("here is a prompt")
				convext.Must(err)
				a.Feedback("And here is the feedback:")
				a.Feedback(v)
				time.Sleep(time.Second * 1) // Give time for it to appear
				v2, err := a.Ask("here is another prompt")
				convext.Must(err)
				a.Feedback("here is on the same line: " + v2)
				return nil
			}),
		},
		{
			Name: "readprefs",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				j, err := json.MarshalIndent(ac.Prefs, "", "  ")
				convext.Must(err)
				fmt.Fprintln(c.App.Writer, string(j))
				return nil
			}),
		},
		{
			Name: "setauth",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				ac.Prefs.AuthToken = types.AuthToken(c.Args().Get(1))
				return ac.SavePrefs()
			}),
		},
		{
			Name:  "printargs",
			Usage: "Print out positional arguments and flags, to help debug command parsing.",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "mystr"},
				&cli.BoolFlag{Name: "mybool"},
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				for _, f := range c.FlagNames() {
					fmt.Fprintf(c.App.Writer, "Flag %s: %v\n", f, c.String(f))
				}
				for i, a := range c.Args().Slice() {
					fmt.Fprintf(c.App.Writer, "Arg %d: %s\n", i, a)
				}
				return nil
			}),
		},
		{
			Name:  "update-auth-display",
			Usage: "Used to update the WASM terminal on startup.",
			Action: cliAction(func(c *cli.Context, appContext appcontext.AppContext, ctx context.Context) error {
				wasmUpdateAuthDisplay(appContext.Prefs)
				return nil
			}),
		},
		{
			Name:  "panic",
			Usage: "Test out a panic",
			Action: cliAction(func(c *cli.Context, appContext appcontext.AppContext, ctx context.Context) error {
				panic("can you see this?")
			}),
		},
		{
			Name:  "statusz",
			Usage: "Check the server status.",
			Action: cliAction(func(c *cli.Context, appContext appcontext.AppContext, ctx context.Context) error {
				r := client.RestyFromContext(ctx)
				resp, err := r.R().Get("/statusz")
				if err != nil {
					return err
				}
				fmt.Fprintln(c.App.Writer, string(resp.Body()))
				return nil
			}),
		},
		{
			Name:  "fourohfour",
			Usage: "Make a 404 request to see how the CLI responds.",
			Action: cliAction(func(c *cli.Context, appContext appcontext.AppContext, ctx context.Context) error {
				r := client.RestyFromContext(ctx)
				resp, err := r.R().Get("/this-does-not-exist")
				if err != nil {
					return err
				}
				fmt.Fprintln(c.App.Writer, string(resp.Body()))
				return nil
			}),
		},
	},
}
