package main

import "fmt"

var (
	version = "1.1.1"
	osarch  string // set by ldflags

	cmdVersion = &command{
		run:       runVersion,
		UsageLine: "version",
		Short:     "display ngc version",
		Long:      "display ngc version and build info.\n",
	}
)

func init() {
	commands = append([]*command{cmdVersion}, commands...)
}

func runVersion(args []string) {
	fmt.Println(version, osarch)
}
