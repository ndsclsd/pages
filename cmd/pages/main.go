package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"

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
	in := flag.String("i", "", "input file")
	out := flag.String("o", "", "output file")
	tmpl := flag.String("t", "pages.tmpl", "template file")

	flag.Parse()

	var (
		b bytes.Buffer
	)

	t := template.New(*tmpl)

	if err := tmplfunc.ParseFiles(t, *tmpl); err != nil {
		return err
	}

	tm := t.New(*in)

	if err := tmplfunc.ParseFiles(tm, *in); err != nil {
		return err
	}

	if err := tm.Execute(&b, nil); err != nil {
		return err
	}

	p := &markdown.Parser{
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

	d := p.Parse(b.String())

	m := markdown.ToHTML(d)

	var w io.Writer = os.Stdout
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	return t.Execute(w, map[string]any{"Content": template.HTML(m)})
}
