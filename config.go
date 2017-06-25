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
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

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
	configDir      string
	directoryURL   string
	allowRenewDays int
	resourceURL    string
	siteConfDir    string
	siteRootDir    string
)

func init() {
	configDir = os.Getenv("NGX_CONFIG")
	directoryURL = os.Getenv("NGX_DIRECTORY_URL")
	resourceURL = os.Getenv("NGX_RESOURCE")
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

	if resourceURL == "" {
		resourceURL = "https://rc.ngxpkg.com/"
	}

	if siteConfDir == "" {
		siteConfDir = "/etc/nginx/conf.d"
	}

	if siteRootDir == "" {
		siteRootDir = "/opt/local/www"
	}

	allowRenewDays, err := strconv.Atoi(os.Getenv("NGX_ALLOW_RENEW_DAYS"))

	if err != nil {
		allowRenewDays = 30
	}

	if allowRenewDays < 7 {
		allowRenewDays = 7
	}

	if allowRenewDays > 30 {
		allowRenewDays = 30
	}
}

type userConfig struct {
	acme.Account
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
