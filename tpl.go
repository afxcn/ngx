package main

import (
	"errors"
	"html/template"
	"os"
)

func writeTpl(tpl *template.Template, fp string, data interface{}) error {

	if _, err := os.Stat(fp); os.IsNotExist(err) {

		fn, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

		if err != nil {
			return err
		}

		defer fn.Close()

		if err := tpl.Execute(fn, data); err != nil {
			return err
		}

		return nil
	}

	return errors.New("tpl: file exists, skip")
}

func editTpl(tpl *template.Template, fp string, data interface{}) error {

	fn, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}

	defer fn.Close()

	if err := tpl.Execute(fn, data); err != nil {
		return err
	}

	return nil
}
