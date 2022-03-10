package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	manual, err := os.Open("MANUAL.md")
	if err != nil {
		log.Fatal(err)
	}
	out, err := os.Create("../webhookdb-api/webhookdb-website/src/docs/manual.md")
	if err != nil {
		log.Fatal(err)
	}
	out.WriteString("---\ntitle: Manual\npath: /docs/manual\norder: 3\n---\nAll commands for the WebhookDB CLI.\n\n```toc\n```\n\n")
	scanner := bufio.NewScanner(manual)
	var prevLine string
	for scanner.Scan() {
		origLine := scanner.Text()
		line := origLine
		if strings.HasPrefix(line, "#") {
			// Indent all headings once
			line = "#" + line
			// If line has an uppercase, it's not a command, so title case it.
			// We must lower it first, since Title doesn't work right otherwise ('NAME' stays 'NAME').
			if strings.ToUpper(line) == line {
				line = strings.Title(strings.ToLower(line))
			}
			// Split out the ### prefix from the title
			headingLevelSep := strings.Index(line, " ")
			title := line[headingLevelSep+1:]
			// Create anchors with IDs so they can go into the TOC
			slug := ToKebabCase(title)
			anchor := fmt.Sprintf("<a id=\"%s\"></a>", slug)
			line = fmt.Sprintf("%s\n\n%s [%s](#%s)", anchor, line[:headingLevelSep], title, slug)
		} else if strings.HasPrefix(line, "\t") {
			if prevLine == "" {
				// Open a markdown ``` block with language if indented
				line = "```arff\n" + strings.TrimSpace(line)
			} else {
				// Keep going with the block
				line = strings.TrimSpace(line)
			}
		} else if origLine == "" && strings.HasPrefix(prevLine, "\t") {
			// End of indented block, finish the ```
			line = "```"
		} else if origLine == "```" && prevLine == "" {
			// Start of a ``` block with no lang, so add it
			line = "```arff"
		}
		out.WriteString(line)
		out.WriteString("\n")
		prevLine = origLine
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var matchNonAlphanum = regexp.MustCompile("[^a-zA-Z0-9]+")

func ToKebabCase(str string) string {
	keb := str
	keb = matchNonAlphanum.ReplaceAllString(keb, "")
	keb = matchFirstCap.ReplaceAllString(keb, "${1}_${2}")
	keb = matchAllCap.ReplaceAllString(keb, "${1}_${2}")
	keb = strings.ToLower(str)
	keb = strings.Replace(keb, " ", "-", -1)
	keb = strings.Replace(keb, "-_", "-", -1)
	return keb
}
