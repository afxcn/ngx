package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/acme"
)

const (
	accountFile    = "account.json"
	accountKeyFile = "account.ecdsa.pem"
	siteConfFile   = "site.conf"
	siteOcspFile   = "ocsp.pem"
	siteIndexFile  = "index.html"

	rsaPrivateKey = "RSA PRIVATE KEY"
	ecPrivateKey  = "EC PRIVATE KEY"
)

var (
	configDir       string
	directoryURL    string
	siteResourceURL string
	siteConfDir     string
	siteRootDir     string
)

func init() {
	configDir = os.Getenv("NGX_CONFIG")
	directoryURL = os.Getenv("NGX_DIRECTORY_URL")
	siteResourceURL = os.Getenv("NGX_SITE_RESOURCE")
	siteConfDir = os.Getenv("NGX_SITE_CONFIG")
	siteRootDir = os.Getenv("NGX_SITE_ROOT")

	if configDir == "" {
		if u, err := user.Current(); err == nil {
			configDir = filepath.Join(u.HomeDir, ".config", "ngxpkg")
		}
	}

	if directoryURL == "" {
		// https://acme-v01.api.letsencrypt.org/directory
		// https://acme-staging.api.letsencrypt.org/directory
		directoryURL = "https://acme-v01.api.letsencrypt.org/directory"
	}

	if siteResourceURL == "" {
		siteResourceURL = "https://rc.ngxpkg.com/"
	}

	if siteConfDir == "" {
		siteConfDir = "/etc/nginx/conf.d"
	}

	if siteRootDir == "" {
		siteRootDir = "/opt/local/www"
	}
}

type userConfig struct {
	acme.Account
}

type serverConfig struct {
	ServerName string
	IPAddress  string
	Username   string
	Password   string
}

func readConfig() (*userConfig, error) {
	b, err := ioutil.ReadFile(filepath.Join(configDir, accountFile))
	if err != nil {
		return nil, err
	}
	uc := &userConfig{}
	if err := json.Unmarshal(b, uc); err != nil {
		return nil, err
	}

	return uc, nil
}

func writeConfig(uc *userConfig) error {
	b, err := json.MarshalIndent(uc, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(configDir, accountFile), b, 0600)
}
