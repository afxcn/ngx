package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

const (
	settingFile   = "setting.json"
	siteConfFile  = "site.conf"
	siteIndexFile = "index.html"
	siteRawURL    = "https://raw.githubusercontent.com/afxcn/ngx-cli/master/"
)

var configDir string
var siteConfDir string
var siteRootDir string

func init() {
	configDir = os.Getenv("NGC_CONFIG")
	siteConfDir = os.Getenv("NGC_SITE_CONFIG")
	siteRootDir = os.Getenv("NGC_SITE_ROOT")

	if configDir == "" {

		if u, err := user.Current(); err == nil {
			configDir = filepath.Join(u.HomeDir, ".config", "ngc")
		}
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
	Server []serverConfig
}

type serverConfig struct {
	ServerName string
	IPAddress  string
	Username   string
	Password   string
}

func readConfig() (*userConfig, error) {
	b, err := ioutil.ReadFile(filepath.Join(configDir, settingFile))
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
	return ioutil.WriteFile(filepath.Join(configDir, settingFile), b, 0600)
}

func readRC(filename string) ([]byte, error) {
	rcDir := filepath.Join(configDir, "rc")

	if err := createDir(rcDir, 0700); err != nil {
		return nil, err
	}

	fp := filepath.Join(rcDir, filename)

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		url := siteRawURL + "rc/" + filename

		fn, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

		if err != nil {
			return nil, err
		}

		defer fn.Close()

		resp, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		_, err = io.Copy(fn, resp.Body)

		if err != nil {
			return nil, err
		}
	}

	bytes, err := ioutil.ReadFile(fp)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
