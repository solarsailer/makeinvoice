package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/solarsailer/makeinvoice/parser"
	"github.com/solarsailer/makeinvoice/table"
)

// -------------------------------------------------------
// Constants & Variables.
// -------------------------------------------------------

// DefaultTemplate is a bare template to use as fallback.
const DefaultTemplate = `{{ .Table }}`

const (
	outputFlag   = "output"
	styleFlag    = "style"
	templateFlag = "template"
)

// -------------------------------------------------------
// Types.
// -------------------------------------------------------

// Output is the data sent to the text/template file.
type Output struct {
	Table string
}

// -------------------------------------------------------
// Main.
// -------------------------------------------------------

func main() {
	app := cli.NewApp()
	app.Name = "makeinvoice"
	app.Usage = "create an invoice populated with data from a CSV file"
	app.Version = "1.0.0"

	// Set main action.
	app.Action = func(c *cli.Context) {

		// Wrap the call and exit with error if needed.
		if err := run(c); err != nil {
			fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
			os.Exit(1)
		}
	}

	// Set flags.
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  outputFlag + ", o",
			Usage: "export to Markdown, HTML or PDF",
		},
		cli.StringFlag{
			Name:  styleFlag + ", s",
			Usage: "decorate the output (only for PDF)",
		},
		cli.StringFlag{
			Name:  templateFlag + ", t",
			Usage: "template file (in markdown)",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	if !c.Args().Present() {
		return errors.New("no arguments passed")
	}

	return execute(c, c.Args().First())
}

// -------------------------------------------------------
// Execute.
// -------------------------------------------------------

func execute(c *cli.Context, csvFilename string) error {
	data, err := parser.ParseCSV(csvFilename)
	if err != nil {
		return err
	}

	template, err := parseTemplate(c)
	if err != nil {
		return err
	}

	// Create a buffer and execute the template with it.
	buffer := bytes.NewBuffer([]byte{})
	template.Execute(buffer, Output{Table: table.Format(data)})

	// Export to markdown or pdf (depends on the extension).
	// Show on the stdout if no export.
	if c.GlobalIsSet(outputFlag) {
		path := c.GlobalString(outputFlag)

		if strings.Contains(path, ".pdf") {
			return createPDF(buffer, path, c.GlobalString(styleFlag))
		}

		return createMarkdown(buffer, path)
	}

	fmt.Println(buffer)

	return nil
}

func createPDF(buffer *bytes.Buffer, path string, cssPath string) error {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "temp")
	if err != nil {
		return errors.New("cannot export to PDF")
	}
	defer os.Remove(tmpFile.Name())

	// Write the markdown to the temp file.
	writeMarkdownFile(buffer, tmpFile)

	// Construct the arguments.
	args := []string{}

	// Parse html tags in the markdown file.
	args = append(args, "-m", `{"html": true}`)

	// Style if CSS file provided.
	if cssPath != "" {
		args = append(args, "-s", cssPath)
	}

	// Output to given path with the tmp file content.
	args = append(args, "-o", path)
	args = append(args, tmpFile.Name())

	if err := exec.Command("markdown-pdf", args...).Run(); err != nil {
		return errors.New("cannot generate the PDF")
	}

	return nil
}

func createMarkdown(buffer *bytes.Buffer, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.New("cannot create a file to " + path)
	}
	defer file.Close()

	writeMarkdownFile(buffer, file)

	return nil
}

func writeMarkdownFile(buffer *bytes.Buffer, file *os.File) error {
	_, err := file.WriteString(buffer.String())
	if err != nil {
		return errors.New("cannot write to the output path")
	}

	return nil
}

// -------------------------------------------------------
// Template.
// -------------------------------------------------------

func parseTemplate(c *cli.Context) (*template.Template, error) {
	// Get the template path form the global option `template`.
	if !c.GlobalIsSet(templateFlag) {
		// If not available, use the default template.
		return useDefaultTemplate(), nil
	}

	content, err := ioutil.ReadFile(c.GlobalString(templateFlag))
	if err != nil {
		return nil, errors.New("cannot read the template")
	}

	template, err := template.New("").Parse(string(content))
	if err != nil {
		return nil, errors.New("invalid template file")
	}

	return template, nil
}

func useDefaultTemplate() *template.Template {
	return template.Must(
		template.New("").Parse(DefaultTemplate),
	)
}
