package cmd

import (
	"context"
	"fmt"
	"github.com/blang/semver"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/config"
	"github.com/webhookdb/webhookdb-cli/whselfupdate"
	"os"
)

var updateCmd = &cli.Command{
	Name:    "update",
	Aliases: []string{"upgrade"},
	Usage:   "Update the version of the CLI in-place.",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "version", Usage: "Use a specific version rather than latest. Can be used to downgrade."},
		&cli.PathFlag{Name: "path", Usage: "Download the new version to the given path. Default to the current executable."},
	},
	Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
		latest, found, err := whselfupdate.DetectVersion(config.Repo, c.String("version"))
		if err == whselfupdate.ErrUnsupported {
			return err
		} else if err != nil {
			return CliError{Message: fmt.Sprintf("Could not find latest release: %s", err)}
		}
		forceUpdate := c.String("version") != ""
		if !forceUpdate {
			v := semver.MustParse(config.Version)
			if !found || latest.Version().LTE(v) {
				fmt.Fprintln(c.App.Writer, "Already up-to-date.")
				return nil
			}
		}
		path := c.Path("path")
		if path == "" {
			exe, err := os.Executable()
			if err != nil {
				return CliError{Message: "Could not locate executable path of this process"}
			}
			path = exe
		}
		fmt.Fprintf(c.App.Writer, "Updating from %s to %s\n", config.Version, latest.Version().String())
		if err := whselfupdate.UpdateTo(latest.AssetURL(), path); err != nil {
			return CliError{Message: fmt.Sprintf("Error occurred while updating binary: %s", err)}
		}
		fmt.Fprintf(c.App.Writer, "Successfully updated to %s\n", latest.Version().String())
		return nil
	}),
}
