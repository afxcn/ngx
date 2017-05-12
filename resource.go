package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func siteResource(filename string) ([]byte, error) {
	if filename == "" {
		return nil, errors.New("ngx: invalid resource filename")
	}

	dir := filepath.Join(configDir, "rc")

	if err := createDir(dir, 0700); err != nil {
		return nil, err
	}

	fp := filepath.Join(dir, filename)

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		url := siteResourceURL + filename

		resp, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(fp, bytes, 0600); err != nil {
			return nil, err
		}

		return bytes, nil

	} else if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(fp)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func writeResource(filename string) error {
	if filename == "" {
		return errors.New("ngx: invalid filename")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := initDir(filename, 0700); err != nil {
			return err
		}

		data, err := siteResource(filepath.Base(filename))

		if err != nil {
			return err
		}

		return ioutil.WriteFile(filename, data, 0600)

	}

	return nil
}
