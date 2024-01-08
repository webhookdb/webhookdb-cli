package cmd

import (
	"context"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/config"
	"github.com/webhookdb/webhookdb-cli/whbrowser"
	"io"
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
				version := config.Version
				if version == config.UnknownVersion {
					version = "main"
				}
				manualUrl := fmt.Sprintf("https://raw.githubusercontent.com/lithictech/webhookdb-cli/%s/MANUAL.md", version)
				fmt.Println(manualUrl)
				resp, err := ac.Resty.R().Get(manualUrl)
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
				&cli.BoolFlag{Name: "docsite", Usage: "If given, format this for use on docs.webhookdb.com"},
			},
		},
	},
}

func docsBuildFunc(c *cli.Context) error {
	app := c.App
	if c.String("format") == "man" {
		docstr, err := app.ToMan()
		if err != nil {
			return err
		}
		app.Writer.Write([]byte(docstr))
	} else {
		return toMarkdown(app.Writer, c)
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

func toMarkdown(w io.Writer, c *cli.Context) error {
	writestr := func(s string, i ...interface{}) {
		w.Write([]byte(fmt.Sprintf(s, i...)))
	}
	writeln := func(s string, i ...interface{}) {
		writestr(s, i...)
		w.Write([]byte("\n"))
	}
	docsite := c.Bool("docsite")
	if docsite {
		writeln("---")
		writeln("title: CLI Reference")
		writeln("layout: home")
		writeln("nav_order: 600")
		writeln("---")
		writeln("")
	}
	md, err := c.App.ToMarkdown()
	if err != nil {
		return err
	}
	if docsite {
		// Add TOC above GLOBAL OPTIONS
		md = strings.Replace(md, "# GLOBAL OPTIONS", "## Table of contents\n{: .no_toc }\n\n1. TOC\n{:toc}\n\n# GLOBAL OPTIONS", 1)
	}
	md = strings.ReplaceAll(md, `**="":`, "=\"\"`:")
	md = strings.ReplaceAll(md, `**:`, "`:")
	md = strings.ReplaceAll(md, `**--`, "`--")

	// Get rid of a level of headers
	multiheaderRe := regexp.MustCompile("(?m)^##")
	md = multiheaderRe.ReplaceAllString(md, "#")

	// Wrap all the command names in backticks
	// I'm certain there's a better way to do this, but it's not worth figuring out for docbuilding.
	headerline1Re := regexp.MustCompile("(?m)^(#+) ([a-z0-9-]+)$")
	headerline2Re := regexp.MustCompile("(?m)^(#+) ([a-z0-9-]+), ([a-z0-9-]+)$")
	headerline3Re := regexp.MustCompile("(?m)^(#+) ([a-z0-9-]+), ([a-z0-9-]+), ([a-z0-9-]+)$")
	headerline4Re := regexp.MustCompile("(?m)^(#+) ([a-z0-9-]+), ([a-z0-9-]+), ([a-z0-9-]+), ([a-z0-9-]+)$")

	md = headerline1Re.ReplaceAllString(md, "$1 `$2`")
	md = headerline2Re.ReplaceAllString(md, "$1 `$2`, `$3`")
	md = headerline3Re.ReplaceAllString(md, "$1 `$2`, `$3`, `$4`")
	md = headerline4Re.ReplaceAllString(md, "$1 `$2`, `$3`, `$4`, `$5`")

	md = strings.ReplaceAll(md, "**Usage`:", "**Usage**:") // Unbreak this
	md = strings.ReplaceAll(md, "# NAME", "# `webhookdb`")
	if docsite {
		md = strings.ReplaceAll(md, "# COMMANDS", "# COMMANDS\n{: .no_toc }")
	}

	// Remove synopsis section
	const tbt = "```"
	md = strings.ReplaceAll(md, fmt.Sprintf(`# SYNOPSIS

webhookdb

%s
[--debug]
[--help|-h]
[--quiet|-q]
%s
`, tbt, tbt), "")

	writeln(md)
	return nil
}
