package cmd

import (
	"github.com/urfave/cli/v2"
)

var authCmd = &cli.Command{
	Name: "auth",
	Subcommands: []*cli.Command{
		{
			Name:        "register",
			Description: "TODO",
			Flags:       []cli.Flag{usernameFlag()},
			Action: func(c *cli.Context) error {
				ctx := newCtx(newAppCtx(c))
				panic(ctx)
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
		Name:    "username",
		Aliases: s1("u"),
		Usage:   "TODO",
	}
}
