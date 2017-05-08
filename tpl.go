package main

import (
	"html/template"
	"os"
)

func writeTpl(tpl *template.Template, fp string, data interface{}) (bool, error) {

	if _, err := os.Stat(fp); os.IsNotExist(err) {

		fn, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

		if err != nil {
			return false, err
		}

		defer fn.Close()

		if err := tpl.Execute(fn, data); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func editTpl(tpl *template.Template, fp string, data interface{}) error {

	fn, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return err
	}

	defer fn.Close()

	if err := tpl.Execute(fn, data); err != nil {
		return err
	}

	return nil
}
