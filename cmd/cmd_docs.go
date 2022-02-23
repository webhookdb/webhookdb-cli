package cmd

import (
	"context"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/lithictech/webhookdb-cli/appcontext"
	"github.com/lithictech/webhookdb-cli/whbrowser"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var docsCmd = &cli.Command{
	Name:  "docs",
	Usage: "Work with the WebhookDB docs and guide.",
	Subcommands: []*cli.Command{
		{
			Name:  "html",
			Usage: "Open a browser to the WebhookDB HTML guide.",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				return whbrowser.OpenURL("https://webhookdb.com/docs/cli")
			}),
		},
		{
			Name:  "tui",
			Usage: "Render the WebhookDB guide into a local Markdown viewer.",
			Flags: []cli.Flag{orgFlag()},
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				resp, err := ac.Resty.R().Get(fmt.Sprintf("%s/docs/cli.md", ac.Config.WebsiteHost))
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
					cm := exec.Command(pa[0], pa[1:]...)
					cm.Stdin = strings.NewReader(string(result))
					cm.Stdout = c.App.Writer
					return cm.Run()
				}
				fmt.Fprint(c.App.Writer, string(result))
				return nil
			}),
		},
		{
			Name:  "build",
			Usage: "Build the docs for the app.",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "format", Usage: "One of: markdown, man"},
			},
		},
	},
}

func docsBuildFunc(c *cli.Context) error {
	app := BuildApp()
	app.Setup()
	writestr := func(s string, i ...interface{}) {
		app.Writer.Write([]byte(fmt.Sprintf(s, i...)))
	}
	if c.String("format") == "man" {
		writestr(app.ToMan())
	} else {
		writestr(app.ToMarkdown())
	}
	return nil
}
func init() {
	// We must do it this way to avoid a cyclical compile error
	for _, subcmd := range docsCmd.Subcommands {
		if subcmd.Name == "build" {
			subcmd.Action = docsBuildFunc
		}
	}
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
