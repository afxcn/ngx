package main

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"path/filepath"

	"golang.org/x/crypto/acme"
)

var (
	certExpiry = 365 * 12 * time.Hour
	certBundle = true
)

func register(client *acme.Client) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	uc := &userConfig{}
	a, err := client.Register(ctx, &uc.Account, prompt)

	if err != nil {
		return err
	}

	uc.Account = *a

	if err := writeConfig(uc); err != nil {
		return err
	}

	return nil
}

func authz(ctx context.Context, client *acme.Client, domainPublic string, domain string) error {
	z, err := client.Authorize(ctx, domain)
	if err != nil {
		return err
	}
	if z.Status == acme.StatusValid {
		return nil
	}
	var chal *acme.Challenge
	for _, c := range z.Challenges {
		if c.Type == "http-01" {
			chal = c
			break
		}
	}
	if chal == nil {
		return errors.New("no supported challenge found")
	}

	val, err := client.HTTP01ChallengeResponse(chal.Token)

	if err != nil {
		return err
	}

	path := filepath.Join(domainPublic, client.HTTP01ChallengePath(chal.Token))

	if err := createFileDir(path, 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, []byte(val), 0644); err != nil {
		return err
	}

	if _, err := client.Accept(ctx, chal); err != nil {
		return fmt.Errorf("accept challenge: %v", err)
	}
	_, err = client.WaitAuthorization(ctx, z.URI)
	return err
}

func readKey(path string) (crypto.Signer, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d, _ := pem.Decode(b)
	if d == nil {
		return nil, fmt.Errorf("no block found in %q", path)
	}
	switch d.Type {
	case rsaPrivateKey:
		return x509.ParsePKCS1PrivateKey(d.Bytes)
	case ecPrivateKey:
		return x509.ParseECPrivateKey(d.Bytes)
	default:
		return nil, fmt.Errorf("%q is unsupported", d.Type)
	}
}

func writeKey(path string, key crypto.PrivateKey) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	var block *pem.Block

	switch privkey := key.(type) {
	case *ecdsa.PrivateKey:
		bytes, err := x509.MarshalECPrivateKey(privkey)
		if err != nil {
			return err
		}
		block = &pem.Block{Type: ecPrivateKey, Bytes: bytes}
	case *rsa.PrivateKey:
		bytes := x509.MarshalPKCS1PrivateKey(privkey)
		block = &pem.Block{Type: rsaPrivateKey, Bytes: bytes}
	default:
		return errors.New("unknown private key type")
	}

	if err := pem.Encode(f, block); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func anyKey(filename string) (crypto.Signer, error) {
	key, err := readKey(filename)
	if err == nil {
		return key, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}

	if strings.Contains(filename, ".ecdsa") {
		privkey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		return privkey, writeKey(filename, privkey)
	}

	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privkey, writeKey(filename, privkey)

}

func prompt(tos string) bool {
	return true
}
