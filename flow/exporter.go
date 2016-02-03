package flow

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/solarsailer/makeinvoice/common/extensions"
)

const (
	pdfConverter = "wkhtmltopdf"
)

// -------------------------------------------------------
// Public.
// -------------------------------------------------------

// Export the data to path.
func Export(template *template.Template, data map[string]string, path string) error {
	buffer := new(bytes.Buffer)

	// No output? Just print on the stdout.
	if path == "" {
		template.Execute(buffer, data)
		return ExportStdout(buffer)
	}

	// We check the extension to know if we need to conver the markdown.
	ext := filepath.Ext(path)

	if ext == extensions.HTML || ext == extensions.PDF {
		// Export to PDF or HTML: we need to convert the markdown to HTML.
		template.Execute(buffer, ConvertAllMarkdownToHTML(data))

		if ext == extensions.PDF {
			// TODO: Add "--user-style-sheet path/to/css" to wkhtmltopdf command.
			return ExportPDF(buffer, path)
		}

		if ext == extensions.HTML {
			return ExportHTML(buffer, path)
		}
	}

	// Export to markdown: the data is already ready.
	template.Execute(buffer, data)
	return ExportMarkdown(buffer, path)
}

// ExportStdout redirects its output to the stdout.
func ExportStdout(buffer *bytes.Buffer) error {
	fmt.Print(buffer.String())
	return nil
}

// ExportMarkdown creates a Markdown file.
func ExportMarkdown(buffer *bytes.Buffer, filename string) error {
	filename = extensions.Force(filename, extensions.Markdown)
	return createTextFile(buffer, filename)
}

// ExportHTML creates an HTML file for a given markdown data.
func ExportHTML(buffer *bytes.Buffer, filename string) error {
	filename = extensions.Force(filename, extensions.HTML)
	return createTextFile(buffer, filename)
}

// ExportPDF creates a PDF file for a given markdown data.
func ExportPDF(buffer *bytes.Buffer, filename string) error {
	filename = extensions.Force(filename, extensions.PDF)
	return createPDF(buffer, filename)
}

// -------------------------------------------------------
// Private.
// -------------------------------------------------------

func createTextFile(buffer *bytes.Buffer, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.New("cannot create " + filename)
	}
	defer f.Close()

	// Write data to the file.
	_, err = f.Write(buffer.Bytes())
	if err != nil {
		return errors.New("cannot write data to " + filename)
	}

	return nil
}

func createPDF(buffer *bytes.Buffer, filename string) error {

	// Create a tmp file.
	tmpFile, err := ioutil.TempFile(os.TempDir(), "mkinv_")
	if err != nil {
		return errors.New("cannot create a temporary file")
	}
	defer os.Remove(tmpFile.Name())

	// Fill the temp file.
	_, err = tmpFile.Write(buffer.Bytes())
	if err != nil {
		return errors.New("cannot write data to a temporary file")
	}

	// Rename to ".html".
	// `wkhtmltopdf` **needs** an extension to determine the type of the file.
	htmlPath := tmpFile.Name() + extensions.HTML

	if err := os.Rename(tmpFile.Name(), htmlPath); err != nil {
		return errors.New("cannot create a temporary html file")
	}
	defer os.Remove(htmlPath)

	// Check for the converter availability.
	if _, err := exec.LookPath(pdfConverter); err != nil {
		return errors.New("impossible to call `" + pdfConverter + "`: install it and run this command again")
	}

	return exec.Command(pdfConverter, htmlPath, filename).Run()
}
