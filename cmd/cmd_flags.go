package cmd

import (
	"fmt"
	"github.com/lithictech/webhookdb-cli/formatting"
	"github.com/lithictech/webhookdb-cli/prefs"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/urfave/cli/v2"
	"strings"
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
	return requireFlagOrArg(c, "service", "Use `webhookdb services list` to see available integrations.")
}

func integrationFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "integration",
		Aliases: s1("i"),
		Usage:   "Integration opaque id, starting with 'svi_'. Run `webhookdb integrations list` to see a list of all your integrations.",
	}
}

func getIntegrationFlagOrArg(c *cli.Context) string {
	return requireFlagOrArg(c, "integration", "Use `webhookdb integrations list` to see available integrations.")
}

func usernameFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "username",
		Aliases: []string{"u", "email"},
		Usage:   "Takes an email.",
	}
}

func formatFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "format",
		Aliases: s1("f"),
		Value:   formatting.Table.FlagValue,
		Usage:   "Format of the output. One of: " + strings.Join(formatting.FormatFlagValues(), ", "),
	}
}
func getFormatFlag(c *cli.Context) formatting.Format {
	f, ok := formatting.LookupByFlag(c.String("format"))
	if !ok {
		panic(CliError{
			Message: fmt.Sprintf("Invalid --format flag value: %s. Must be one of: %s", c.String("format"), strings.Join(formatting.FormatFlagValues(), ", ")),
			Code:    1,
		})
	}
	return f
}

func webhookFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "webhook",
		Aliases: s1("w"),
		Usage:   "Webhook opaque id. Run `webhookdb webhook list` to see a list of all your webhooks.",
	}
}

func getWebhookFlagOrArg(c *cli.Context) string {
	return requireFlagOrArg(c, "webhook", "Use `webhookdb webhook list` to see available webhooks.")
}

func extractPositional(idx int, c *cli.Context, msg string) string {
	a := c.Args().Get(idx)
	if a == "" {
		panic(CliError{Message: msg, Code: 1})
	}
	return a
}

func flagOrArg(c *cli.Context, param string) string {
	v := c.String(param)
	if v != "" {
		return v
	}
	return c.Args().First()
}

func requireFlagOrArg(c *cli.Context, param, extraMsg string) string {
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

func stringPtrFlag(c *cli.Context, key string) *string {
	if !c.IsSet(key) {
		return nil
	}
	s := c.String(key)
	return &s
}
