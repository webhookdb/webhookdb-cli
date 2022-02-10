package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/urfave/cli/v2"
)

var fixturesCmd = &cli.Command{
	Name:  "fixtures",
	Flags: []cli.Flag{},
	Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
		serviceName, err := extractServiceName(0, c)
		out, err := client.GetFixtures(ctx, ac.Auth, client.GetFixturesInput{
			ServiceName: serviceName,
		})
		if err != nil {
			return err
		}
		fmt.Println(out.SchemaSql)
		return nil
	}),
}
