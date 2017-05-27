/**

Copyright (C) 2017  ZhiQiang Huang (email: ngxpkg@gmail.com)

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
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/acme"
)

var (
	cmdRenew = &command{
		run:       runRenew,
		UsageLine: "renew [domain ...]",
		Short:     "renew ssl certificates base on domain conf",
		Long: `
Parse domain conf and renew ssl certificates.
If not domain input, will parse all domain conf at NGX_SITE_CONFIG dir.
`,
	}
)

func runRenew(args []string) {

	var client *acme.Client
	certExpiry := 365 * 24 * time.Hour
	certBundle := true

	for _, domain := range args {
		domainConfPath := filepath.Join(siteConfDir, domain+".conf")

		if _, err := os.Stat(domainConfPath); os.IsNotExist(err) {
			logf("%s conf: %v", domain, err)
			continue
		}

		conf, err := parseSiteConf(domainConfPath)

		if err != nil {
			fatalf("%s conf: %v", domain, err)
		}

		for _, cert := range conf.Certificates {

			c, err := parseCertificate(cert.fullchain)

			if err != nil {
				fatalf("%s cert: %v", domain, err)
			}

			if !strings.Contains(c.Issuer.CommonName, "Let's Encrypt") {
				logf("%s Issuer '%s' not support acme, skip.", filepath.Base(cert.fullchain), c.Issuer.CommonName)
				continue
			}

			days := int(c.NotAfter.Sub(time.Now()).Hours() / 24)

			if days > allowRenewDays {
				logf("%s still %d days valid, skip.", filepath.Base(cert.fullchain), days)
				continue
			}

			if client == nil {

				accountKey, err := anyKey(filepath.Join(configDir, accountKeyFile))

				if err != nil {
					fatalf("account key: %v", err)
				}

				client = &acme.Client{
					Key:          accountKey,
					DirectoryURL: directoryURL,
				}

				if _, err := readConfig(); os.IsNotExist(err) {
					if err := register(client); err != nil {
						fatalf("register: %v", err)
					}
				}
			}

			req := &x509.CertificateRequest{
				Subject:  pkix.Name{CommonName: c.Subject.CommonName},
				DNSNames: c.DNSNames,
			}

			privkey, err := memKey(cert.privkey)

			if err != nil {
				fatalf("privkey: %v", err)
			}

			csr, err := x509.CreateCertificateRequest(rand.Reader, req, privkey)

			if err != nil {
				fatalf("csr: %v", err)
			}

			for _, dnsName := range c.DNSNames {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)

				if err := authz(ctx, client, conf.DomainPublicDir, dnsName); err != nil {
					fatalf("authz %s: %v", dnsName, err)
				}
				cancel()
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()
			certs, _, err := client.CreateCert(ctx, csr, certExpiry, certBundle)

			if err != nil {
				fatalf("cert: %v", err)
			}

			var pemcert []byte
			for _, b := range certs {
				b = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: b})
				pemcert = append(pemcert, b...)
			}

			if err := writeKey(cert.privkey, privkey); err != nil {
				fatalf("privkey: %v", err)
			}

			if err := ioutil.WriteFile(cert.fullchain, pemcert, 0644); err != nil {
				fatalf("cert: %v", err)
			}
		}
	}

	if client != nil {
		if err := nginxReload(); err != nil {
			fatalf("nginx: %v", err)
		}
	}
}
