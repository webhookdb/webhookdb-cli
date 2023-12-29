package cmd

import (
	"context"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var fixturesCmd = &cli.Command{
	Name:    "fixtures",
	Aliases: []string{"fixture"},
	Usage: "Output the SQL DDL (CREATE TABLE command) to create a DB table that matches what is in WebhookDB. " +
		"This can be used to generate .sql files that can be run as part of test database fixturing.",
	Flags: []cli.Flag{serviceFlag()},
	Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
		out, err := client.GetFixtures(ctx, ac.Auth, client.GetFixturesInput{
			ServiceName: getServiceFlagOrArg(c),
		})
		if err != nil {
			return err
		}
		printlnif(c, out.Message, true)
		printlnif(c, out.SchemaSql, false)
		return nil
	}),
}
