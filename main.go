package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/solarsailer/makeinvoice/converter"
	"github.com/solarsailer/makeinvoice/extensions"
	"github.com/solarsailer/makeinvoice/parser"
	"github.com/solarsailer/makeinvoice/table"
	"github.com/solarsailer/makeinvoice/template"
)

// -------------------------------------------------------
// Constants & Variables.
// -------------------------------------------------------

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
	app.Action = wrap

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
			Usage: "template file (Markdown or HTML)",
		},
	}

	app.Run(os.Args)
}

// Wrap the call and exit with error if needed.
func wrap(c *cli.Context) {
	if err := run(c); err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	if !c.Args().Present() {
		return errors.New("no arguments passed")
	}

	return process(c, c.Args().First())
}

func process(c *cli.Context, csvInput string) error {
	data, err := parser.ParseCSV(csvInput)
	if err != nil {
		return err
	}

	// Get the template filename from the flags.
	templateFilename := ""
	if c.GlobalIsSet(templateFlag) {
		templateFilename = c.GlobalString(templateFlag)
	}

	template, err := template.Parse(templateFilename)
	if err != nil {
		return err
	}

	// Format the data as a byte slice.
	content := table.Format(data)
	isMarkdown := filepath.Ext(templateFilename) == extensions.HTML

	// If the template is an HTML file, convert the data into HTML
	// before passing it to the template.
	if isMarkdown {
		html := converter.ConvertMarkdownToHTML([]byte(content))
		content = strings.TrimSpace(string(html))
	}

	// Create a buffer and execute the template with it.
	buffer := new(bytes.Buffer)
	template.Execute(buffer, Output{Table: content})

	// Export to markdown or pdf (depends on the extension).
	if c.GlobalIsSet(outputFlag) {
		return export(isMarkdown, buffer.Bytes(), c.GlobalString(outputFlag))
	}

	// No output? Just print on the stdout.
	fmt.Print(buffer)

	return nil
}

func export(fromMarkdown bool, data []byte, path string) error {
	if filepath.Ext(path) == extensions.PDF {
		// TODO css c.GlobalString(styleFlag)
		return converter.ExportPDF(fromMarkdown, data, path)
	}

	if filepath.Ext(path) == extensions.HTML {
		return converter.ExportHTML(fromMarkdown, data, path)
	}

	return converter.ExportMarkdown(data, path)
}
