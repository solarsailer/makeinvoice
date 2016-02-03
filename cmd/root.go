package cmd

import (
	"errors"

	"github.com/solarsailer/makeinvoice/flow"
	"github.com/solarsailer/makeinvoice/parser"
	"github.com/solarsailer/makeinvoice/template"
	"github.com/spf13/cobra"
)

// -------------------------------------------------------
// Constants & Variables.
// -------------------------------------------------------

const (
	outputFlag   = "output"
	cssFlag      = "css"
	templateFlag = "template"
)

// -------------------------------------------------------
// Command.
// -------------------------------------------------------

// RootCmd is the main entry point of the application.
var Root = &cobra.Command{
	Use:   "makeinvoice",
	Short: "Create an invoice populated with data from a CSV file.",
	Long:  `Create an invoice populated with data from a (or multiple) CSV file(s).`,
	Example: `  makeinvoice data.csv
  makeinvoice --output invoice.pdf data.csv
  makeinvoice --output invoice.html --template tpl.html data01.csv data02.csv`,
	RunE: run,
}

// -------------------------------------------------------
// Init.
// -------------------------------------------------------

func init() {
	Root.Flags().StringP(outputFlag, "o", "", "export to Markdown, HTML or PDF")
	Root.Flags().StringP(cssFlag, "c", "", "decorate the output with a CSS file (only for PDF)")
	Root.Flags().StringP(templateFlag, "t", "", "template file (Text, Markdown or HTML)")
}

// -------------------------------------------------------
// Functions.
// -------------------------------------------------------

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("no arguments passed")
	}

	// Parse the provided CSV files.
	data, err := parser.ParseCSVFiles(args)
	if err != nil {
		return err
	}

	// Get the template filename from the flags and parse it.
	templateFilename := cmd.Flag(templateFlag).Value.String()
	template, err := template.Parse(templateFilename)
	if err != nil {
		return err
	}

	return flow.Export(
		template,
		data,
		cmd.Flag(outputFlag).Value.String(),
		cmd.Flag(cssFlag).Value.String(),
	)
}
