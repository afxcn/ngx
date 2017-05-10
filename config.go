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
	siteResourceURL string
	siteConfDir     string
	siteRootDir     string
)

func init() {
	configDir = os.Getenv("NGX_CONFIG")
	siteResourceURL = os.Getenv("NGX_SITE_RESOURCE")
	siteConfDir = os.Getenv("NGX_SITE_CONFIG")
	siteRootDir = os.Getenv("NGX_SITE_ROOT")

	if configDir == "" {
		if u, err := user.Current(); err == nil {
			configDir = filepath.Join(u.HomeDir, ".config", "ngxpkg")
		}
	}

	if siteResourceURL == "" {
		siteResourceURL = "https://raw.githubusercontent.com/afxcn/ngxpkg/master/rc/"
	}

	if siteConfDir == "" {
		siteConfDir = "/etc/nginx/conf.d"
	}

	if siteRootDir == "" {
		siteRootDir = "/opt/local/www"
	}

	if err := createDir(configDir, 0700); err != nil {
		fatalf("create configDir failure: %v", err)
	}

	if err := createDir(siteConfDir, 0700); err != nil {
		fatalf("create siteConfDir failure: %v", err)
	}

	if err := createDir(siteRootDir, 0755); err != nil {
		fatalf("create siteRootDir failure: %v", err)
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
