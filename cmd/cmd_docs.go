package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/whbrowser"
	"io"
	"regexp"
	"strings"
)

var docsCmd = &cli.Command{
	Name:  "docs",
	Usage: "Work with the WebhookDB docs and guide.",
	Subcommands: []*cli.Command{
		{
			Name:  "guide",
			Usage: "Open a browser to the WebhookDB Getting Started guide.",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				return whbrowser.OpenURL("https://docs.webhookdb.com/docs/getting-started/")
			}),
		},
		{
			Name:    "manual",
			Aliases: []string{"html"},
			Usage:   "Open a browser to the WebhookDB CLI reference.",
			Action: cliAction(func(c *cli.Context, ac appcontext.AppContext, ctx context.Context) error {
				return whbrowser.OpenURL("https://docs.webhookdb.com/docs/cli-reference.html")
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
	// c.App only has the 'docs' commands and I cannot find a way to get the root
	app := BuildApp()
	md, err := app.ToMarkdown()
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
