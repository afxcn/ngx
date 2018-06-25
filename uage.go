/**

Copyright (C) 2017-2018 ZhiQiang Huang, All Rights Reserved.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

**/

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
				ConfigDir string
			}{
				ConfigDir: configDir,
			}
			tmpl(os.Stdout, cmd.Long, data)
			return
		}
	}

	fatalf("Unknown help topic %q. Run 'ngx help'.\n", arg)
}
