package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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

func createMockSSL(domain string, confFilename string) error {
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

					if err := createMockCert(c.privkey, c.fullchain, domain); err != nil {
						return err
					}

					c = nil
				}
			case "ssl_certificate_key":
				if c == nil {
					c = &cert{
						privkey: value,
					}
				} else {
					c.privkey = value

					if err := createMockCert(c.privkey, c.fullchain, domain); err != nil {
						return err
					}

					c = nil
				}
			case "ssl_session_ticket_key":
				if err := createSessionTicketKey(value); err != nil {
					return err
				}
			case "ssl_dhparam":
				if err := createDHParam(value); err != nil {
					return err
				}
			case "ssl_trusted_certificate":
				if err := createTrustedCertificate(value); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func createMockCert(privkey string, fullchain string, domain string) error {
	if err := createFileDir(privkey, 0700); err != nil {
		return err
	}

	if err := createFileDir(fullchain, 0700); err != nil {
		return err
	}

	if strings.Contains(privkey, "ecdsa") {
		privData, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

		if err != nil {
			return err
		}

		writeKey(privkey, privData)
		createSelfSignCert(fullchain, privData, domain)
	} else {
		privData, err := rsa.GenerateKey(rand.Reader, 2048)

		if err != nil {
			return err
		}

		writeKey(privkey, privData)
		createSelfSignCert(fullchain, privData, domain)
	}

	return nil
}

func createSessionTicketKey(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := createFileDir(filename, 0700); err != nil {
			return err
		}

		bytes := make([]byte, 80)
		_, err := rand.Read(bytes)

		if err != nil {
			return err
		}

		return ioutil.WriteFile(filename, bytes, 0600)
	}

	return nil
}

func createDHParam(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := createFileDir(filename, 0700); err != nil {
			return err
		}

		cmd := exec.Command("/bin/sh", "-c", "openssl dhparam -out "+filename+" 2048")
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return nil
}

func createTrustedCertificate(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {

		if err := createFileDir(filename, 0700); err != nil {
			return err
		}

		siteOcspData, err := siteRC(siteOcspFile)

		if err != nil {
			return err
		}

		return ioutil.WriteFile(filename, siteOcspData, 0644)
	}

	return nil
}

func writeKey(path string, key crypto.PrivateKey) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	var b *pem.Block

	switch privkey := key.(type) {
	case *rsa.PrivateKey:
		bytes := x509.MarshalPKCS1PrivateKey(privkey)
		b = &pem.Block{Type: rsaPrivateKey, Bytes: bytes}
	case *ecdsa.PrivateKey:
		bytes, err := x509.MarshalECPrivateKey(privkey)
		if err != nil {
			f.Close()
			return err
		}
		b = &pem.Block{Type: ecPrivateKey, Bytes: bytes}
	default:
		f.Close()
		return errors.New("unknown private key type")
	}

	if err := pem.Encode(f, b); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}

func createSelfSignCert(path string, priv crypto.PrivateKey, domain string) error {
	var notBefore = time.Now()
	notAfter := notBefore.AddDate(0, 0, 2)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   domain,
			Organization: []string{"Mock Cert"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := []string{
		domain,
		"www." + domain,
	}

	for _, h := range hosts {
		template.DNSNames = append(template.DNSNames, h)
	}

	template.IsCA = false
	template.KeyUsage |= x509.KeyUsageCertSign

	pubkey, err := publicKey(priv)

	if err != nil {
		return err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, pubkey, priv)

	if err != nil {
		return err
	}

	certOut, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}

	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}

	return nil
}

func publicKey(priv crypto.PrivateKey) (crypto.PublicKey, error) {
	switch key := priv.(type) {
	case *rsa.PrivateKey:
		return &key.PublicKey, nil
	case *ecdsa.PrivateKey:
		return &key.PublicKey, nil
	default:
		return nil, errors.New("unknown private key type")
	}
}

func reloadNginx() error {
	cmd := exec.Command("/bin/sh", "-c", "sudo nginx -s reload")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
