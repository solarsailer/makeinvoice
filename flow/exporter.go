package flow

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/solarsailer/makeinvoice/common/extensions"
)

const (
	pdfConverter = "wkhtmltopdf"
)

// -------------------------------------------------------
// Public.
// -------------------------------------------------------

// Export the data to path.
func Export(data []byte, path string) error {
	if filepath.Ext(path) == extensions.PDF {
		// TODO: Add "--user-style-sheet path/to/css" to wkhtmltopdf command.
		return ExportPDF(data, path)
	}

	if filepath.Ext(path) == extensions.HTML {
		return ExportHTML(data, path)
	}

	return ExportMarkdown(data, path)
}

// ExportMarkdown creates a Markdown file.
func ExportMarkdown(data []byte, filename string) error {
	filename = extensions.Force(filename, extensions.Markdown)

	return createTextFile(data, filename)
}

// ExportHTML creates an HTML file for a given markdown data.
func ExportHTML(data []byte, filename string) error {
	data = ConvertMarkdownToHTML(data)
	filename = extensions.Force(filename, extensions.HTML)

	return createTextFile(data, filename)
}

// ExportPDF creates a PDF file for a given markdown data.
func ExportPDF(data []byte, filename string) error {
	data = ConvertMarkdownToHTML(data)
	filename = extensions.Force(filename, extensions.PDF)

	return createPDF(data, filename)
}

// -------------------------------------------------------
// Private.
// -------------------------------------------------------

func createTextFile(data []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.New("cannot create " + filename)
	}
	defer f.Close()

	// Write data to the file.
	_, err = f.Write(data)
	if err != nil {
		return errors.New("cannot write data to " + filename)
	}

	return nil
}

func createPDF(html HTML, filename string) error {

	// Create a tmp file.
	tmpFile, err := ioutil.TempFile(os.TempDir(), "mkinv_")
	if err != nil {
		return errors.New("cannot create a temporary file")
	}
	defer os.Remove(tmpFile.Name())

	// Fill the temp file.
	_, err = tmpFile.Write(html)
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
