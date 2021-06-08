package cmd

import "github.com/urfave/cli/v2"

func orgFlag() *cli.StringFlag {
	// takes the org key
	return &cli.StringFlag{
		Name:     "org",
		Aliases:  s1("o"),
		Required: false,
		Usage:    "Takes an org key. Run `webhook org list` to see a list of all your org keys.",
	}
}

func roleFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "role",
		Aliases:  s1("r"),
		Required: true,
		Usage:    "Takes a role name.",
	}
}

func tokenFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "token",
		Aliases:  s1("t"),
		Required: true,
		Usage:    "Takes a one time password—only used during auth.",
	}
}

func usernameFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "username",
		Aliases:  s1("u"),
		Required: true,
		Usage:    "Takes an email.",
	}
}

func usernamesFlag() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:     "usernames",
		Aliases:  nil,
		Required: true,
		Usage:    "Takes multiple emails.",
	}
}
