package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/config"
)

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "Print version and exit.",
	Action: func(c *cli.Context) error {
		shaPart := config.BuildSha
		if len(shaPart) >= 8 {
			shaPart = fmt.Sprintf(" (%s)", config.BuildSha[0:8])
		}
		fmt.Fprintf(c.App.Writer, "%s%s\n", config.Version, shaPart)
		return nil
	},
}
