package cmd

import (
	"context"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var docsCmd = &cli.Command{
	Name:        "docs",
	Description: "Work with the WebhookDB docs and guide.",
	Subcommands: []*cli.Command{
		{
			Name:        "html",
			Description: "Open a browser to the WebhookDB HTML guide.",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				return browser.OpenURL("https://webhookdb.com/docs/cli")
			}),
		},
		{
			Name:        "tui",
			Description: "Render the WebhookDB guide into a local Markdown viewer.",
			Flags:       []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				resp, err := ac.Resty.R().Get("https://webhookdb.com/docs/cli.md")
				if err != nil {
					return err
				}
				if resp.StatusCode() >= 400 {
					return CliError{Message: "Sorry, could not fetch the guide Markdown: " + resp.String(), Code: 2}
				}
				md := resp.String()
				md = regexp.MustCompile("\\A---(.|\n)*?---").ReplaceAllString(md, "")
				md = regexp.MustCompile("<a id=\".*\"></a>").ReplaceAllString(md, "")
				result := markdown.Render(md, 80, 0)
				if pager := getPager(); pager != "" {
					pa := strings.Split(pager, " ")
					c := exec.Command(pa[0], pa[1:]...)
					c.Stdin = strings.NewReader(string(result))
					c.Stdout = os.Stdout
					return c.Run()
				}
				fmt.Print(string(result))
				return nil
			}),
		},
	},
}

func getPager() string {
	pager := os.Getenv("PAGER")
	if pager != "" {
		return pager
	}
	if _, err := exec.LookPath("less"); err == nil {
		return "less -r"
	}
	return ""
}
