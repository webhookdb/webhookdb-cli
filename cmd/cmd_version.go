package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/config"
)

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "Print version and exit.",
	Flags: []cli.Flag{&cli.BoolFlag{Name: "time", Usage: "Print the build time as well."}},
	Action: func(c *cli.Context) error {
		shaPart := fmt.Sprintf(" (%s)", config.BuildShaShort)
		fmt.Fprintf(c.App.Writer, "%s%s\n", config.Version, shaPart)
		if c.Bool("time") {
			fmt.Fprintf(c.App.Writer, "%s\n", config.BuildTime)
		}
		return nil
	},
}
