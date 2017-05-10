package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	usageTemplate = `ngx is a cli tool for nginx
Usage:
	ngx command [arguments]

The commands are:
{{range .}}{{if .Runnable}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "ngx help [command]" for more information about a command.

`
)

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, commands)
	bw.Flush()
}

type errWriter struct {
	w   io.Writer
	err error
}

func (w *errWriter) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	if err != nil {
		w.err = err
	}
	return n, err
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"trim":       strings.TrimSpace,
		"capitalize": capitalize,
	})
	template.Must(t.Parse(text))
	ew := &errWriter{w: w}
	err := t.Execute(ew, data)
	if ew.err != nil {
		if strings.Contains(ew.err.Error(), "pipe") {
			os.Exit(1)
		}
		fatalf("writing output: %v", ew.err)
	}
	if err != nil {
		panic(err)
	}
}

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}
	if len(args) != 1 {
		fatalf("usage: ngx help command\n\nToo many arguments given.\n")
	}

	arg := args[0]
	for _, cmd := range commands {
		if cmd.Name() == arg {
			if cmd.Runnable() {
				fmt.Fprintf(os.Stdout, "usage: ngx %s\n", cmd.UsageLine)
			}
			data := struct {
				ConfigDir    string
				AccountFile  string
				AccountKey   string
				DefaultDisco string
				DiscoAliases map[string]string
			}{
				ConfigDir: configDir,
			}
			tmpl(os.Stdout, cmd.Long, data)
			return
		}
	}

	fatalf("Unknown help topic %q. Run 'ngx help'.\n", arg)
}
