package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/go-aperitif/convext"
	"github.com/lithictech/go-aperitif/logctx"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func Execute() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Value: false},
		},
		Commands: []*cli.Command{
			authCmd,
			{
				Name: "version",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(os.Stdout, config.BuildSha[0:8])
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func newAppCtx(c *cli.Context) appcontext.AppContext {
	if c.Bool("debug") {
		convext.Must(os.Setenv("LOG_LEVEL", "debug"))
	}
	appCtx, err := appcontext.New(c.Command.FullName(), config.LoadConfig())
	if err != nil {
		log.Fatal(err)
	}
	return appCtx
}

func newCtx(appCtx appcontext.AppContext) context.Context {
	c := context.Background()
	c = logctx.WithLogger(c, appCtx.Logger())
	c = logctx.WithTraceId(c, logctx.ProcessTraceIdKey)
	c = logctx.WithTracingLogger(c)
	return appcontext.InContext(c, appCtx)
}
