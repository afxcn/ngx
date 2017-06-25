/**

Copyright (C) 2017 ZhiQiang Huang (email: ngxpkg@gmail.com)

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

func register(client *acme.Client) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	uc := &userConfig{}

	account, err := client.Register(ctx, &uc.Account, prompt)

	if err != nil {
		return err
	}

	uc.Account = *account

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

	if err := sameDir(path, 0755); err != nil {
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

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)

	if block == nil {
		return nil, fmt.Errorf("no block found in %q", path)
	}

	switch block.Type {
	case rsaPrivateKey:
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case ecPrivateKey:
		return x509.ParseECPrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("%q is unsupported", block.Type)
	}
}

func writeKey(path string, key crypto.PrivateKey) error {

	fn, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

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

	if err := pem.Encode(fn, block); err != nil {
		fn.Close()
		return err
	}

	return fn.Close()
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

func memKey(filename string) (crypto.Signer, error) {

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

		return privkey, nil
	}

	privkey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, err
	}

	return privkey, nil
}

func prompt(tos string) bool {
	return true
}

func parseCertificate(path string) (*x509.Certificate, error) {

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)

	cert, err := x509.ParseCertificate(block.Bytes)

	if err != nil {
		return nil, err
	}

	return cert, nil
}
