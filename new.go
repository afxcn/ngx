package main

import (
	"html/template"
	"path/filepath"
)

var (
	cmdNew = &command{
		run:       runNew,
		UsageLine: "new [-f force] domain [domain ...]",
		Short:     "Create an new empty Site or reinitialize an existing site",
		Long: `
Create an new empty Site or reinitialize an existing site.

`,
	}
)

func init() {

}

func runNew(args []string) {
	if len(args) == 0 {
		fatalf("no domain specified")
	}

	siteConfData, err := siteRC(siteConfFile)

	if err != nil {
		fatalf("read site conf failure: %v", err)
	}

	siteIndexData, err := siteRC(siteIndexFile)

	if err != nil {
		fatalf("read site index failure: %v", err)
	}

	siteConfTpl, err := template.New("siteConf").Parse(string(siteConfData))

	if err != nil {
		fatalf("parse site conf template failure: %v", err)
	}

	siteIndexTpl, err := template.New("siteIndex").Parse(string(siteIndexData))

	if err != nil {
		fatalf("parse site index template failure: %v", err)
	}

	for _, domain := range args {
		domainConfPath := filepath.Join(siteConfDir, domain+".conf")
		domainRootDir := filepath.Join(siteRootDir, domain)
		domainPublicDir := filepath.Join(domainRootDir, "public")
		domainIndexPath := filepath.Join(domainPublicDir, siteIndexFile)

		if err := createDir(domainRootDir, 0755); err != nil {
			fatalf("create domain root dir failure: %v", err)
		}

		if err := createDir(domainPublicDir, 0755); err != nil {
			fatalf("create domain public dir failure: %v", err)
		}

		data := struct {
			SiteRoot string
			Domain   string
			WithSSL  bool
		}{
			SiteRoot: siteRootDir,
			Domain:   domain,
			WithSSL:  false,
		}

		if err := writeTpl(siteConfTpl, domainConfPath, data); err != nil {
			fatalf("create domain conf failure: %v", err)
		}

		if err := writeTpl(siteIndexTpl, domainIndexPath, data); err != nil {
			fatalf("create domain index failure: %v", err)
		}

		if err := reloadNginx(); err != nil {
			fatalf("reload nginx failure: %v", err)
		}

		if err := parseSiteConfSSL(domain, domainConfPath); err != nil {
			fatalf("parse site conf ssl failure: %v", err)
		}
	}
}
