package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"rsc.io/markdown"
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

	flag.Parse()

	var (
		b   []byte
		err error
	)

	if *in != "" {
		b, err = os.ReadFile(os.Args[1])
	} else {
		b, err = io.ReadAll(os.Stdin)

	}
	if err != nil {
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

	d := p.Parse(string(b))

	m := markdown.ToHTML(d)

	if *out != "" {
		err = os.WriteFile(*out, []byte(m), 0666)
	} else {
		_, err = os.Stdout.Write([]byte(m))
	}

	return err
}
