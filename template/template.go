package template

import (
	"errors"
	"io/ioutil"
	"text/template"
)

const defaultTemplate = "{{ range . }}{{ . }}\n\n{{ end }}"

// -------------------------------------------------------
// Parser.
// -------------------------------------------------------

// Parse a template (and select the appropriate engine based on the file's extension).
func Parse(templateFilename string) (*template.Template, error) {

	// If not available, use the default template.
	if templateFilename == "" {
		return useDefaultTemplate(), nil
	}

	content, err := ioutil.ReadFile(templateFilename)
	if err != nil {
		return nil, errors.New("cannot read the template")
	}

	return useTextTemplate(string(content))
}

func useTextTemplate(content string) (*template.Template, error) {
	t, err := template.New("").Parse(content)
	if err != nil {
		return nil, errors.New("invalid Text/Markdown template file")
	}

	return t, nil
}

func useDefaultTemplate() *template.Template {
	return template.Must(
		template.New("").Parse(defaultTemplate),
	)
}
