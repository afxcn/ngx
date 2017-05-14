package main

var (
	cmdReNew = &command{
		run:       runReNew,
		UsageLine: "renew [domain ...]",
		Short:     "renew ssl certificates base on domain conf",
		Long: `
Parse domain conf and renew ssl certificates.
If not domain input, will parse all domain conf at NGX_SITE_CONFIG dir.
`,
	}
)

func runReNew(args []string) {
	if len(args) == 0 {
		fatalf("no domain specified")
	}
}
