package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

// -------------------------------------------------------
// Constants & Variables.
// -------------------------------------------------------

// DefaultTemplate is a bare template to use as fallback.
const DefaultTemplate = `{{ .Table }}`

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
	app.Usage = "create invoice populated with data from a CSV file"
	app.Version = "1.0.0"
	app.Action = func(c *cli.Context) {
		if !c.Args().Present() {
			exit("No arguments passed.")
		}

		execute(c, c.Args().First())
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Usage: "export to markdown or PDF",
		},
		cli.StringFlag{
			Name:  "css, s",
			Usage: "decorate the output (only for PDF)",
		},
		cli.StringFlag{
			Name:  "template, t",
			Usage: "template file (in markdown)",
		},
	}
	app.Run(os.Args)
}

func exit(message string) {
	fmt.Println(color.RedString(message))
	os.Exit(1)
}

// -------------------------------------------------------
// Execute.
// -------------------------------------------------------

func execute(c *cli.Context, csvFilename string) {
	data, err := parseCSV(csvFilename)
	if err != nil {
		exit("Cannot open the provided CSV file.")
	}

	buffer := bytes.NewBuffer([]byte{})

	template := parseTemplate(c)
	template.Execute(buffer, Output{Table: format(data)})

	// Export to markdown or pdf (depends on the extension).
	// Show on the stdout if no export.
	if c.GlobalIsSet("output") {
		path := c.GlobalString("output")

		if strings.Contains(path, ".pdf") {
			createPDF(buffer, path)
		} else {
			createMarkdown(buffer, path)
		}
	} else {
		fmt.Println(buffer)
	}
}

func createPDF(buffer *bytes.Buffer, path string) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "temp")
	if err != nil {
		exit("Cannot export to PDF.")
	}
	defer os.Remove(tmpFile.Name())

	// Write the markdown to the temp file.
	writeMarkdownFile(buffer, tmpFile)

	// And convert.
	exec.Command(
		"markdown-pdf", "-o", path, tmpFile.Name(),
	).Run()
}

func createMarkdown(buffer *bytes.Buffer, path string) {
	file, err := os.Create(path)
	if err != nil {
		exit("Cannot create a file to " + path + ".")
	}
	defer file.Close()

	writeMarkdownFile(buffer, file)
}

func writeMarkdownFile(buffer *bytes.Buffer, file *os.File) {
	_, err := file.WriteString(buffer.String())
	if err != nil {
		exit("Cannot write to the output path.")
	}
}

// -------------------------------------------------------
// Template.
// -------------------------------------------------------

func parseTemplate(c *cli.Context) *template.Template {
	// Get the template path form the global option `template`.
	if !c.GlobalIsSet("template") {
		// If not available, use the default template.
		return useDefaultTemplate()
	}

	content, err := ioutil.ReadFile(c.GlobalString("template"))
	if err != nil {
		exit("Cannot read the template.")
	}

	template, err := template.New("").Parse(string(content))
	if err != nil {
		exit("Invalid template file.")
	}

	return template
}

func useDefaultTemplate() *template.Template {
	return template.Must(
		template.New("").Parse(DefaultTemplate),
	)
}

// -------------------------------------------------------
// CSV.
// -------------------------------------------------------

func parseCSV(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// -------------------------------------------------------
// Table formatting.
// -------------------------------------------------------

// Format a two-dimensional array into a markdown table.
//
// Example:
//
//   [["head1", "head2"], ["val1", "val2"]]
//   =>
//   "
//     head1|head2
//     -|-
//     val1|val2
//   "
func format(table [][]string) string {
	result := ""

	for i, row := range table {
		result += formatRow(row)

		if i == 0 {
			result += formatHeader(len(row))
		}
	}

	return strings.TrimSpace(result)
}

// Format the header line of a markdown table.
// For each col, put a "-" and separate each col with a "|".
//
// Example:
//
//   ["a", "b"] => "-|-\n"
func formatHeader(length int) string {
	result := "\n"

	for i := 0; i < length; i++ {
		if i != 0 {
			result += "|"
		}

		result += "-"
	}

	return result
}

// Format a row - separate each value by a "|".
//
// Example:
//
//   ["a", "b"] => "a|b\n"
func formatRow(row []string) string {
	result := "\n"

	for i, col := range row {
		if i != 0 {
			result += "|"
		}

		result += col
	}

	return result
}
