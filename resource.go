/**

Copyright (C) 2017 ZhiQiang Huang, All Rights Reserved.

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
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func fetchResource(filename string) ([]byte, error) {
	if filename == "" {
		return nil, errors.New("ngx: empty resource name")
	}

	dir := filepath.Join(configDir, "rc")

	if err := mkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	fp := filepath.Join(dir, filename)

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		url := resourceURL + filename

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
	}

	bytes, err := ioutil.ReadFile(fp)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func writeResource(filename string) error {
	if filename == "" {
		return errors.New("ngx: empty resource name")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := sameDir(filename, 0700); err != nil {
			return err
		}

		bytes, err := fetchResource(filepath.Base(filename))

		if err != nil {
			return err
		}

		return ioutil.WriteFile(filename, bytes, 0600)

	}

	return nil
}
