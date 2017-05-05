package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	commands = []*command{
		cmdNew,
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
}

type command struct {
	run       func(args []string)
	flag      flag.FlagSet
	UsageLine string
	Short     string
	Long      string
}

func (c *command) Name() string {
	name := c.UsageLine
	i := strings.IndexRune(name, ' ')
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *command) Usage() {
	help([]string{c.Name()})
	os.Exit(2)
}

func (c *command) Runnable() bool {
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
	f.StringVar(&configDir, "c", configDir, "")
}
