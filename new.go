package main

import "fmt"

var (
	cmdNew = &command{
		run:       runNew,
		UsageLine: "new [-f force] domain [domain ...]",
		Short:     "Create an new empty Site or reinitialize an existing site",
		Long: `
Create an new empty Site or reinitialize an existing site.

`,
	}

	force bool
)

func init() {
	cmdNew.flag.BoolVar(&force, "f", false, "domain of site to create.")
}

func runNew(args []string) {
	if len(args) == 0 {
		fatalf("no domain specified")
	}

	for _, domain := range args {
		fmt.Println(domain)
	}
}
