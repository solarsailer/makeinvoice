package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/solarsailer/makeinvoice/converter"
	"github.com/solarsailer/makeinvoice/extensions"
	"github.com/solarsailer/makeinvoice/parser"
	"github.com/solarsailer/makeinvoice/table"
	"github.com/solarsailer/makeinvoice/template"
	"github.com/spf13/cobra"
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
// Command.
// -------------------------------------------------------

// RootCmd is the main entry point of the application.
var Root = &cobra.Command{
	Use:     "makeinvoice",
	Short:   "Create an invoice populated with data from a CSV file.",
	Long:    `Create an invoice populated with data from a (or multiple) CSV file(s).`,
	Example: `  makeinvoice --output invoice42.pdf data.csv`,
	RunE:    run,
}

// -------------------------------------------------------
// Init.
// -------------------------------------------------------

func init() {
	Root.Flags().StringP(outputFlag, "o", "", "export to Markdown, HTML or PDF")
	Root.Flags().StringP(styleFlag, "s", "", "decorate the output (only for PDF)")
	Root.Flags().StringP(templateFlag, "t", "", "template file (Markdown or HTML)")
}

// -------------------------------------------------------
// Functions.
// -------------------------------------------------------

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("no arguments passed")
	}

	// Parse the provided CSV file.
	data, err := parser.ParseCSV(args[0])
	if err != nil {
		return err
	}

	// Get the template filename from the flags and parse it.
	templateFilename := cmd.Flag(templateFlag).Value.String()
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
	template.Execute(buffer, struct{ Table string }{Table: content})

	// Export to markdown or pdf (depends on the extension).
	outputPath := cmd.Flag(outputFlag).Value.String()
	if outputPath != "" {
		return export(isMarkdown, buffer.Bytes(), outputPath)
	}

	// No output? Just print on the stdout.
	fmt.Print(buffer)

	return nil
}

func export(fromMarkdown bool, data []byte, path string) error {
	if filepath.Ext(path) == extensions.PDF {
		// TODO: Add "--user-style-sheet path/to/css" to wkhtmltopdf command.
		return converter.ExportPDF(fromMarkdown, data, path)
	}

	if filepath.Ext(path) == extensions.HTML {
		return converter.ExportHTML(fromMarkdown, data, path)
	}

	return converter.ExportMarkdown(data, path)
}
