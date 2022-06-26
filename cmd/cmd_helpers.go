package cmd

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lithictech/go-aperitif/convext"
	"github.com/lithictech/go-aperitif/logctx"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/lithictech/webhookdb-cli/config"
	"github.com/urfave/cli/v2"
	"os"
	"regexp"
)

func s1(s string) []string {
	if len(s) != 1 {
		panic(fmt.Sprintf("s1 strings must be one char, got: %s", s))
	}
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
	appCtx.Resty.OnAfterResponse(func(rc *resty.Client, r *resty.Response) error {
		if r.StatusCode() == 404 || r.StatusCode() >= 500 {
			return onServerError(c, appCtx, r)
		}
		return nil
	})
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
		cb = guardInvalidArgs(cb)
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

func guardInvalidArgs(cb cliActionCallback) cliActionCallback {
	return func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
		if ac.Config.SkipArgFlagCheck || c.Args().Len() == 0 {
			return cb(c, ac, ctx)
		}
		for _, arg := range c.Args().Tail() {
			if isFlagArg.MatchString(arg) {
				return CliError{
					Message: fmt.Sprintf("Positional arguments must follow flags, but '%s' looks like a flag. "+
						"Please re-run the command, putting it before positional arguments (which start with '%s')."+
						"\nIf this placement is intentional, re-run this command with %s=1.",
						arg, c.Args().Get(0), config.SkipArgFlagCheckEnv,
					),
					Code: 1,
				}
			}
		}
		return cb(c, ac, ctx)
	}
}

var isFlagArg = regexp.MustCompile("^--?[a-z-]+=?$")

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

// Print msg if not quiet mode, and it's not empty.
// If linebr, print a newline after the message.
// Usually 'true' when printing a collection after the message.
func printlnif(c *cli.Context, msg string, linebr bool) {
	if c.Bool("quiet") {
		return
	}
	if len(msg) == 0 {
		return
	}
	fmt.Fprintln(c.App.Writer, msg)
	if linebr {
		fmt.Fprintln(c.App.Writer)
	}
}

const tableNameRules = "Valid table names must adhere to the following rules: " +
	"must begin with an ASCII letter, contain only ASCII letters, numbers, underscores, dashes, and spaces, " +
	"can begin and end with double quotes, and must otherwise be a valid Postgres identifier."
