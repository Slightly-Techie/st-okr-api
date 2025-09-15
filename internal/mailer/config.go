package mailer

import (
	"bytes"
	"fmt"
	"html/template"
)

func LoadTemplate(templateName string, data any) (string, error) {

	t, err := template.ParseFiles(fmt.Sprintf("templates/%s.html", templateName))

	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	err = t.Execute(&body, data)
	if err != nil {
		return "", err
	}
	return body.String(), nil

}
