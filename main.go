/**

Copyright (C) 2017 ZhiQiang Huang, All Rights Reserved.

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
	"flag"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	commands = []*cmd{
		cmdNew,
		cmdRenew,
		cmdVersion,
	}
	exitMu     sync.Mutex
	exitStatus = 0
)

func main() {
	flag.Usage = usage
	flag.Parse() // catch -h argument
	log.SetFlags(0)

	args := flag.Args()

	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Runnable() {
			addFlags(&cmd.flag)
			cmd.flag.Usage = func() { cmd.Usage() }
			cmd.flag.Parse(args[1:])
			cmd.run(cmd.flag.Args())
			exit()
			return
		}
	}

	fatalf("Unknown subcommand %q.\nRun 'ngx help' for usage.\n", args[0])
}

type cmd struct {
	run       func(args []string)
	flag      flag.FlagSet
	UsageLine string
	Short     string
	Long      string
}

func (c *cmd) Name() string {
	name := c.UsageLine
	i := strings.IndexRune(name, ' ')
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *cmd) Usage() {
	help([]string{c.Name()})
	os.Exit(2)
}

func (c *cmd) Runnable() bool {
	return c.run != nil
}

var logf = log.Printf

func errorf(format string, args ...interface{}) {
	logf(format, args...)
	setExitStatus(1)
}

func fatalf(format string, args ...interface{}) {
	errorf(format, args...)
	exit()
}

func setExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

func exit() {
	os.Exit(exitStatus)
}

func addFlags(f *flag.FlagSet) {
	f.StringVar(&configDir, "configDir", configDir, "")
	f.StringVar(&directoryURL, "directoryURL", directoryURL, "")
	f.StringVar(&resourceURL, "resourceURL", resourceURL, "")
	f.StringVar(&siteConfDir, "siteConfDir", siteConfDir, "")
	f.StringVar(&siteRootDir, "siteRootDir", siteRootDir, "")
	f.IntVar(&allowRenewDays, "allowRenewDays", allowRenewDays, "")
}
