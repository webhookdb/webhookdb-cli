package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var integrationsCmd = &cli.Command{
	Name:        "integrations",
	Description: "Make sure that you're working on the correct organization when you create an integration.",
	Subcommands: []*cli.Command{
		{
			Name:        "create",
			Description: "create an integration for the given organization",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				serviceName, err := extractPositional(0, c, "Service name required. Use 'webhookdb services list' to view all available services.")
				if err != nil {
					return err
				}
				input := client.IntegrationsCreateInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					ServiceName:   serviceName,
				}
				step, err := client.IntegrationsCreate(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				if err := client.NewStateMachine().Run(ctx, ac.Auth, step); err != nil {
					return err
				}
				return nil
			}),
		},
		{
			Name:        "list",
			Description: "list all integrations for the given organization",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsListInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.IntegrationsList(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				if len(out.Data) == 0 {
					// TODO: fix the bug where the cli library's timestamp shows up if we return an `errors.New()`
					return errors.New("This organization doesn't have any integrations set up yet.")
				}
				// TODO: Get this spacing correct
				fmt.Println("service name \t\t\t\t table name \t\t\t\t id")
				for _, value := range out.Data {
					fmt.Println(value.ServiceName + " \t\t\t " + value.TableName + " \t\t\t " + value.OpaqueId)
				}
				return nil
			}),
		},
		{
			Name:        "reset",
			Description: "Reset the webhook secret for this integration.",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				opaqueId, err := extractIntegrationId(0, c)
				if err != nil {
					return err
				}
				input := client.IntegrationsResetInput{
					OpaqueId: opaqueId,
				}
				step, err := client.IntegrationsReset(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				if err := client.NewStateMachine().Run(ctx, ac.Auth, step); err != nil {
					return err
				}
				return nil
			}),
		},
	},
}
