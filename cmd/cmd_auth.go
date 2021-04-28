package cmd

import (
	"fmt"
	"github.com/lithictech/webhookdb-cli/ask"
	"github.com/urfave/cli/v2"
	"os"
)

const PASSWORD_RETRY_ATTEMPTS = 3

var authCmd = &cli.Command{
	Name: "auth",
	Subcommands: []*cli.Command{
		{
			Name:        "register",
			Description: "TODO",
			Flags:       []cli.Flag{usernameFlag()},
			Action: func(c *cli.Context) error {
				ctx := newCtx(newAppCtx(c))
				password, err := ask.HiddenAsk(ask.HiddenPrompt("Enter your password (12 or more characters):"))
				if err != nil {
					return err
				}
				attempt := 0
				for {
					confirmed, err := ask.HiddenAsk(ask.HiddenPrompt("Please repeat your password:"))
					if err != nil {
						return err
					}
					if password == confirmed {
						break
					}
					attempt += 1
					if attempt == PASSWORD_RETRY_ATTEMPTS {
						fmt.Println("Sorry, those passwords don't match. Exiting.")
						os.Exit(1)
					}
					fmt.Println("Sorry, those passwords don't match, try again.")
				}
				fmt.Println(ctx, c.String("username"), password)
				return nil
			},
		},
		{
			Name:        "login",
			Description: "TODO",
			Flags:       []cli.Flag{usernameFlag()},
			Action: func(c *cli.Context) error {
				panic("TODO")
			},
		},
		{
			Name:        "logout",
			Description: "TODO",
			Flags:       []cli.Flag{usernameFlag()},
			Action: func(c *cli.Context) error {
				panic("TODO")
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
