package cmd

import (
	"context"
	"fmt"
	"github.com/blang/semver"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/config"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var updateCmd = &cli.Command{
	Name: "update",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "version", Usage: "Use a specific version rather than latest. Can be used to downgrade."},
	},
	Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
		latest, found, err := selfupdate.DetectVersion(config.Repo, c.String("version"))
		if err != nil {
			return CliError{Message: fmt.Sprintf("Could not find latest release: %s", err)}
		}
		forceUpdate := c.String("version") != ""
		if !forceUpdate {
			v := semver.MustParse(config.Version)
			if !found || latest.Version.LTE(v) {
				fmt.Println("Already up-to-date.")
				return nil
			}
		}
		exe, err := os.Executable()
		if err != nil {
			return CliError{Message: "Could not locate executable path of this process"}
		}
		fmt.Printf("Updating from %s to %s\n", config.Version, latest.Version.String())
		if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
			return CliError{Message: fmt.Sprintf("Error occurred while updating binary: %s", err)}
		}
		log.Printf("Successfully updated to %s\n", latest.Version.String())
		return nil
	}),
}
