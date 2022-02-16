package cmd

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/client"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"os"
)

var integrationsCmd = &cli.Command{
	Name:        "integrations",
	Description: "Make sure that you're working on the correct organization when you create an integration.",
	Subcommands: []*cli.Command{
		{
			Name:        "create",
			Description: "Create an integration for the given organization",
			Flags: []cli.Flag{
				orgFlag(),
				serviceFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsCreateInput{
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
					ServiceName:   flagOrArg(c, "service", "Use `webhookdb services list` to see available integrations."),
				}
				return stateMachineResponseRunner(ctx, ac.Auth)(client.IntegrationsCreate(ctx, ac.Auth, input))
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
					fmt.Println("This organization doesn't have any integrations set up yet.")
					fmt.Println("Use `webhookdb services list` and `webhookdb integrations create` to set one up.")
					return nil
				}
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Service", "Table", "Id"})
				configTableWriter(table)
				for _, value := range out.Data {
					table.Append([]string{value.ServiceName, value.TableName, value.OpaqueId})
				}
				table.Render()
				return nil
			}),
		},
		{
			Name:        "reset",
			Description: "Reset the webhook secret for this integration.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsResetInput{
					OpaqueId:      getIntegrationFlagOrArg(c),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				return stateMachineResponseRunner(ctx, ac.Auth)(client.IntegrationsReset(ctx, ac.Auth, input))
			}),
		},
		{
			Name:        "status",
			Description: "Get statistics about webhooks for this integration.",
			Flags: []cli.Flag{
				orgFlag(),
				integrationFlag(),
			},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				input := client.IntegrationsStatusInput{
					OpaqueId:      getIntegrationFlagOrArg(c),
					OrgIdentifier: getOrgFlag(c, ac.Prefs),
				}
				out, err := client.IntegrationsStatus(ctx, ac.Auth, input)
				if err != nil {
					return err
				}
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader(out.Header)
				configTableWriter(table)
				for _, row := range out.Rows {
					table.Append(row)
				}
				table.Render()
				return nil
			}),
		},
	},
}
