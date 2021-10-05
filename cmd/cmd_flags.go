package cmd

import (
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/urfave/cli/v2"
)

func orgFlag() *cli.StringFlag {
	// takes the org key
	return &cli.StringFlag{
		Name:     "org",
		Aliases:  s1("o"),
		Required: false,
		Usage:    "Takes an org key. Run `webhook org list` to see a list of all your org keys.",
	}
}

func getOrgFlag(c *cli.Context, p prefs.Prefs) types.OrgIdentifier {
	slug := c.String("org")
	if slug == "" {
		return types.OrgIdentifierFromId(p.CurrentOrg.Id)
	}
	return types.OrgIdentifierFromSlug(slug)
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

func usernamesFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "usernames",
		Aliases:  nil,
		Required: true,
		Usage:    "Takes multiple emails.",
	}
}

func extractPositional(idx int, c *cli.Context, msg string) (string, error) {
	a := c.Args().Get(idx)
	if a == "" {
		return "", CliError{Message: msg, Code: 1}
	}
	return a, nil
}

func extractIntegrationId(idx int, c *cli.Context) (string, error) {
	return extractPositional(idx, c, "Integration Id required. Use `webhookdb integrations list` to view all integrations.")
}
