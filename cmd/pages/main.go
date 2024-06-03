package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"rsc.io/markdown"
	"rsc.io/tmplfunc"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	inputName := flag.String("input", "", "input filename")
	outputName := flag.String("output", "", "output filename")
	templateName := flag.String("template", "pages.tmpl", "template filename")

	flag.Parse()

	var (
		buffer bytes.Buffer
	)

	tmpl := template.New(*templateName)

	if err := tmplfunc.ParseFiles(tmpl, *templateName); err != nil {
		return err
	}

	input, err := os.ReadFile(*inputName)
	if err != nil {
		return err
	}

	header, input := parseHeader(input)

	tmplContent := tmpl.New(*inputName)

	if err := tmplfunc.Parse(tmplContent, string(input)); err != nil {
		return err
	}

	if err := tmplContent.Execute(&buffer, map[string]any{"Header": header}); err != nil {
		return err
	}

	content := buffer.String()

	if path.Ext(*inputName) == ".md" {
		parser := &markdown.Parser{
			HeadingIDs:         true,
			Strikethrough:      true,
			TaskListItems:      true,
			AutoLinkText:       true,
			AutoLinkAssumeHTTP: true,
			Table:              true,
			Emoji:              true,
			SmartDot:           true,
			SmartDash:          true,
			SmartQuote:         true,
		}

		document := parser.Parse(content)

		content = markdown.ToHTML(document)
	}

	buffer.Reset()

	if err := tmpl.Execute(&buffer, map[string]any{"Header": header, "Content": template.HTML(content)}); err != nil {
		return err
	}

	return os.WriteFile(*outputName, buffer.Bytes(), 0666)
}

var (
	jsonStart = []byte("<!--{")
	jsonEnd   = []byte("}-->")
	yamlStart = []byte("---\n")
	yamlEnd   = []byte("\n---\n")
)

func parseHeader(b []byte) (map[string]any, []byte) {
	switch {
	case bytes.HasPrefix(b, jsonStart):
		header := make(map[string]any)
		end := bytes.Index(b, jsonEnd)

		if end < 0 {
			return header, b
		}

		json.Unmarshal(b[len(jsonEnd)-1:end+1], &header)

		return header, b[end+len(jsonEnd):]
	case bytes.HasPrefix(b, yamlStart):
		header := make(map[string]any)
		end := bytes.Index(b, yamlEnd)
		if end < 0 {
			return header, b
		}

		yaml.Unmarshal(b[len(yamlStart):end+1], &header)

		return header, b[end+len(yamlEnd):]
	}

	return nil, b
}
