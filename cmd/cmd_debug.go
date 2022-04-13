package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lithictech/go-aperitif/convext"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/ask"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/urfave/cli/v2"
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
	},
}
