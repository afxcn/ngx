package main

import (
	"io/ioutil"
	"regexp"
)

var (
	confSslRegex = regexp.MustCompile(`(ssl_certificate|ssl_certificate_key|ssl_session_ticket_key|ssl_dhparam|ssl_trusted_certificate)\s+([a-z0-9_\-\.\/]+?);`)
)

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
