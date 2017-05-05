package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const (
	settingFile      = "setting.json"
	siteTemplateFile = "site.conf"
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
			configDir = filepath.Join(u.HomeDir, ".config", "nginx")
		}
	}

	if siteConfDir == "" {
		siteConfDir = "/etc/nginx/conf.d"
	}

	if siteRootDir == "" {
		siteRootDir = "/opt/local/www"
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
