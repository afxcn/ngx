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
	"io/ioutil"
	"regexp"
)

var (
	confSslRegex = regexp.MustCompile(`(root|ssl_certificate|ssl_certificate_key|ssl_session_ticket_key|ssl_dhparam|ssl_trusted_certificate)\s+([a-z0-9_\-\.\/]+?);`)
)

type ngxCertificate struct {
	privkey   string
	fullchain string
}

type ngxSiteConf struct {
	Certificates          []ngxCertificate
	SslSessionTicketKey   string
	SslDHParam            string
	SslTrustedCertificate string
	DomainPublicDir       string
}

func parseSiteConf(confFilename string) (*ngxSiteConf, error) {
	text, err := ioutil.ReadFile(confFilename)

	if err != nil {
		return nil, err
	}

	matches := confSslRegex.FindAllStringSubmatch(string(text), -1)

	var conf = &ngxSiteConf{}
	var cert *ngxCertificate

	for _, match := range matches {

		if len(match) == 3 {
			key, value := match[1], match[2]

			switch key {
			case "ssl_certificate":
				if cert == nil {
					cert = &ngxCertificate{
						fullchain: value,
					}
				} else {
					cert.fullchain = value
					conf.Certificates = append(conf.Certificates, *cert)
					cert = nil
				}
			case "ssl_certificate_key":
				if cert == nil {
					cert = &ngxCertificate{
						privkey: value,
					}
				} else {
					cert.privkey = value
					conf.Certificates = append(conf.Certificates, *cert)
					cert = nil
				}
			case "ssl_session_ticket_key":
				conf.SslSessionTicketKey = value
			case "ssl_dhparam":
				conf.SslDHParam = value
			case "ssl_trusted_certificate":
				conf.SslTrustedCertificate = value
			case "root":
				if conf.DomainPublicDir == "" {
					conf.DomainPublicDir = value
				}
			}
		}

	}

	return conf, nil
}
