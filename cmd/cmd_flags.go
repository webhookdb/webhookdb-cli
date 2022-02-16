package cmd

import (
	"fmt"
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

func serviceFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "service",
		Aliases: s1("s"),
		Usage:   "Name of the service. Run `webhookdb services list` to see a list of all services available to your organization.",
	}
}

func getServiceFlagOrArg(c *cli.Context) string {
	return flagOrArg(c, "service", "Use `webhookdb services list` to see available integrations.")
}

func integrationFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "integration",
		Aliases: s1("i"),
		Usage:   "Integration opaque id, starting with 'svi_'. Run `webhookdb integrations list` to see a list of all your integrations.",
	}
}

func getIntegrationFlagOrArg(c *cli.Context) string {
	return flagOrArg(c, "integration", "Use `webhookdb integrations list` to see available integrations.")
}

func usernameFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "username",
		Aliases: s1("u"),
		Usage:   "Takes an email.",
	}
}

func extractPositional(idx int, c *cli.Context, msg string) string {
	a := c.Args().Get(idx)
	if a == "" {
		panic(CliError{Message: msg, Code: 1})
	}
	return a
}

func flagOrArg(c *cli.Context, param, extraMsg string) string {
	v := c.String(param)
	if v != "" {
		return v
	}
	v = c.Args().First()
	if v != "" {
		return v
	}
	msg := fmt.Sprintf("Please pass --%s or provide it as the first argument.", param)
	if extraMsg != "" {
		msg += " " + extraMsg
	}
	panic(CliError{Code: 1, Message: msg})
}
