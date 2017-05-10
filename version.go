package main

import "fmt"

var (
	version = "0.0.1"
	osarch  string // set by ldflags

	cmdVersion = &command{
		run:       runVersion,
		UsageLine: "version",
		Short:     "display ngx version",
		Long:      "display ngx version and build info.\n",
	}
)

func init() {
	commands = append([]*command{cmdVersion}, commands...)
}

func runVersion(args []string) {
	fmt.Println(version, osarch)
}
