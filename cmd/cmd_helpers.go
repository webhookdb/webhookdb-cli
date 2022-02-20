package cmd

import (
	"context"
	"github.com/lithictech/go-aperitif/convext"
	"github.com/lithictech/go-aperitif/logctx"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/config"
	"github.com/urfave/cli/v2"
	"os"
)

func s1(s string) []string {
	return []string{s}
}

func newAppCtx(c *cli.Context) appcontext.AppContext {
	if c.Bool("debug") {
		convext.Must(os.Setenv("LOG_LEVEL", "debug"))
		convext.Must(os.Setenv("DEBUG", "true"))
	}
	appCtx, err := appcontext.New(c.Command.FullName(), config.LoadConfig())
	if err != nil {
		panic(err)
	}
	return appCtx
}

func newCtx(appCtx appcontext.AppContext) context.Context {
	c := context.Background()
	c = logctx.WithLogger(c, appCtx.Logger())
	c = logctx.WithTraceId(c, logctx.ProcessTraceIdKey)
	c = logctx.WithTracingLogger(c)
	c = client.RestyInContext(c, appCtx.Resty)
	return appcontext.InContext(c, appCtx)
}

type cliActionCallback func(*cli.Context, appcontext.AppContext, context.Context) error

func cliAction(cb cliActionCallback) cli.ActionFunc {
	return func(c *cli.Context) (returnErr error) {
		ac := newAppCtx(c)
		ctx := newCtx(ac)
		defer func() {
			if r := recover(); r == nil {
				return
			} else {
				if ce, ok := r.(CliError); ok {
					returnErr = cli.Exit(ce.Message, ce.Code)
				} else {
					panic(r)
				}
			}
		}()
		if err := cb(c, ac, ctx); err != nil {
			if eresp, ok := err.(client.ErrorResponse); ok {
				if eresp.Err.Status == 401 {
					return cli.Exit("You are not logged in. Use 'webhookdb auth login' to get started.", 2)
				}
				msg := eresp.Err.Message
				if msg == "" {
					msg = "Sorry, something went wrong. Please report it to support@webhookdb.com."
				}
				return cli.Exit(msg, 2)
			}
			if ce, ok := err.(CliError); ok {
				return cli.Exit(ce.Message, ce.Code)
			}
			return err
		}
		return nil
	}
}

type CliError struct {
	Message string
	Code    int
}

func (e CliError) Error() string {
	return e.Message
}

func stateMachineResponseRunner(ctx context.Context, auth client.Auth) func(client.Step, error) error {
	return func(st client.Step, e error) error {
		_, err := client.StateMachineResponseRunner(ctx, auth)(st, e)
		return err
	}
}
