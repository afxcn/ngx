package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func siteResource(filename string) ([]byte, error) {
	dir := filepath.Join(configDir, "rc")

	if err := createDir(dir, 0700); err != nil {
		return nil, err
	}

	fp := filepath.Join(dir, filename)

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		url := siteResourceURL + filename

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
			fn.Close()
			os.Remove(fp)
			return nil, err
		}
	}

	bytes, err := ioutil.ReadFile(fp)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
