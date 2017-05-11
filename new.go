package main

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/acme"
)

var (
	cmdNew = &command{
		run:       runNew,
		UsageLine: "new domain [domain ...]",
		Short:     "Create an new empty Site or reinitialize an existing site",
		Long: `
Create an new empty Site or reinitialize an existing site.

`,
	}
)

func runNew(args []string) {
	if len(args) == 0 {
		fatalf("no domain specified")
	}

	siteConfData, err := siteResource(siteConfFile)

	if err != nil {
		fatalf("read conf: %v", err)
	}

	siteIndexData, err := siteResource(siteIndexFile)

	if err != nil {
		fatalf("read index: %v", err)
	}

	siteConfTpl, err := template.New("siteConf").Parse(string(siteConfData))

	if err != nil {
		fatalf("parse conf: %v", err)
	}

	siteIndexTpl, err := template.New("siteIndex").Parse(string(siteIndexData))

	if err != nil {
		fatalf("parse index: %v", err)
	}

	for _, domain := range args {
		domainConfPath := filepath.Join(siteConfDir, domain+".conf")
		domainRootDir := filepath.Join(siteRootDir, domain)
		domainPublicDir := filepath.Join(domainRootDir, "public")
		domainIndexPath := filepath.Join(domainPublicDir, siteIndexFile)

		if err := createDir(domainRootDir, 0755); err != nil {
			fatalf("create root dir: %v", err)
		}

		if err := createDir(domainPublicDir, 0755); err != nil {
			fatalf("create public dir: %v", err)
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

		created, err := writeTpl(siteConfTpl, domainConfPath, data)

		if err != nil {
			fatalf("create conf: %v", err)
		}

		if _, err := writeTpl(siteIndexTpl, domainIndexPath, data); err != nil {
			fatalf("create index: %v", err)
		}

		if err := nginxReload(); err != nil {
			fatalf("nginx: %v", err)
		}

		if created {

			data.WithSSL = true

			if err := editTpl(siteConfTpl, domainConfPath, data); err != nil {
				fatalf("edit conf: %v", err)
			}

			conf, err := parseSiteConf(domain, domainConfPath)

			if err != nil {
				fatalf("parse conf: %v", err)
			}

			if _, err := os.Stat(conf.SslDHParam); os.IsNotExist(err) {
				if err := createFileDir(conf.SslDHParam, 0700); err != nil {
					fatalf("dhparam dir: %v", err)
				}
				data, err := siteResource(filepath.Base(conf.SslDHParam))
				if err != nil {
					fatalf("dhparam: %v", err)
				}

				ioutil.WriteFile(conf.SslDHParam, data, 0600)
			}

			if _, err := os.Stat(conf.SslSessionTicketKey); os.IsNotExist(err) {
				if err := createFileDir(conf.SslSessionTicketKey, 0700); err != nil {
					fatalf("ticket Key dir: %v", err)
				}
				data, err := siteResource(filepath.Base(conf.SslSessionTicketKey))
				if err != nil {
					fatalf("ticket Key: %v", err)
				}

				ioutil.WriteFile(conf.SslSessionTicketKey, data, 0600)
			}

			if _, err := os.Stat(conf.SslTrustedCertificate); os.IsNotExist(err) {
				if err := createFileDir(conf.SslTrustedCertificate, 0700); err != nil {
					fatalf("create ocsp dir failure: %v", err)
				}
				data, err := siteResource(filepath.Base(conf.SslTrustedCertificate))
				if err != nil {
					fatalf("read ocsp failure: %v", err)
				}

				ioutil.WriteFile(conf.SslTrustedCertificate, data, 0600)
			}

			accountKey, err := anyKey(filepath.Join(configDir, accountKeyFile))

			if err != nil {
				fatalf("account key: %v", err)
			}

			client := &acme.Client{
				Key:          accountKey,
				DirectoryURL: directoryURL,
			}

			if _, err := readConfig(); os.IsNotExist(err) {
				if err := register(client); err != nil {
					fatalf("register: %v", err)
				}
			} else if err != nil {
				fatalf("user config: %v", err)
			}

			req := &x509.CertificateRequest{
				Subject: pkix.Name{CommonName: domain},
			}

			dnsNames := []string{
				domain,
				"www." + domain,
			}

			for _, cert := range conf.Certificates {

				if err := createFileDir(cert.privkey, 0700); err != nil {
					fatalf("create privkey dir: %v", err)
				}

				privkey, err := anyKey(cert.privkey)

				if err != nil {
					fatalf("privkey: %v", err)
				}

				req.DNSNames = dnsNames

				csr, err := x509.CreateCertificateRequest(rand.Reader, req, privkey)

				if err != nil {
					fatalf("csr: %v", err)
				}

				for _, dnsName := range dnsNames {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)

					if err := authz(ctx, client, domainPublicDir, dnsName); err != nil {
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

				if err := ioutil.WriteFile(cert.fullchain, pemcert, 0644); err != nil {
					fatalf("write cert: %v", err)
				}

			}

			if err := nginxReload(); err != nil {
				fatalf("nginx: %v", err)
			}

		}
	}
}
