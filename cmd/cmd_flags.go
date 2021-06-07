package cmd

import "github.com/urfave/cli/v2"

func orgFlag() *cli.StringFlag {
	// takes the org key
	return &cli.StringFlag{
		Name:     "org",
		Aliases:  s1("o"),
		Required: false,
		Usage:    "TODO",
	}
}

func roleFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "role",
		Aliases:  s1("r"),
		Required: true,
		Usage:    "TODO",
	}
}

func tokenFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "token",
		Aliases:  s1("t"),
		Required: true,
		Usage:    "TODO",
	}
}

func usernameFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "username",
		Aliases:  s1("u"),
		Required: true,
		Usage:    "TODO",
	}
}

func usernamesFlag() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:     "usernames",
		Aliases:  nil,
		Required: true,
		Usage:    "TODO",
	}
}
