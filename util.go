package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var (
	confSslRegex = regexp.MustCompile(`(ssl_certificate|ssl_certificate_key|ssl_session_ticket_key|ssl_dhparam|ssl_trusted_certificate)\s+([a-z0-9_\-\.\/]+?);`)
)

func createDir(dir string, perm os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {

		if err := os.MkdirAll(dir, perm); err != nil {
			return err
		}
	}

	return nil
}

func createFileDir(filename string, perm os.FileMode) error {
	dir := filepath.Dir(filename)
	return createDir(dir, perm)
}

func parseSiteConfSSL(domain string, confFilename string) error {
	text, err := ioutil.ReadFile(confFilename)

	if err != nil {
		return err
	}

	matches := confSslRegex.FindAllStringSubmatch(string(text), -1)

	type cert struct {
		privkey   string
		fullchain string
	}

	var c *cert

	for _, match := range matches {

		if len(match) == 3 {
			key, value := match[1], match[2]

			switch key {
			case "ssl_certificate":
				if c == nil {
					c = &cert{
						fullchain: value,
					}
				} else {
					c.fullchain = value
					////
					c = nil
				}
			case "ssl_certificate_key":
				if c == nil {
					c = &cert{
						privkey: value,
					}
				} else {
					c.privkey = value
					////
					c = nil
				}
			case "ssl_session_ticket_key":

			case "ssl_dhparam":

			case "ssl_trusted_certificate":

			}
		}
	}

	return nil
}

func writeTpl(tpl *template.Template, fp string, data interface{}) error {

	if _, err := os.Stat(fp); os.IsNotExist(err) {

		fn, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

		if err != nil {
			return err
		}

		defer fn.Close()

		if err := tpl.Execute(fn, data); err != nil {
			return err
		}

	}

	return nil
}

func reloadNginx() error {
	cmd := exec.Command("/bin/sh", "-c", "sudo nginx -s reload")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
