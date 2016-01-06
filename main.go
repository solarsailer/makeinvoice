package main

import (
	"bytes"
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
	app.Action = func(c *cli.Context) {
		if !c.Args().Present() {
			exit("No arguments passed.")
		}

		execute(c, c.Args().First())
	}
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

func exit(message string) {
	fmt.Println(color.RedString(message))
	os.Exit(1)
}

// -------------------------------------------------------
// Execute.
// -------------------------------------------------------

func execute(c *cli.Context, csvFilename string) {
	data, err := parser.ParseCSV(csvFilename)
	if err != nil {
		exit("Cannot open the provided CSV file.")
	}

	buffer := bytes.NewBuffer([]byte{})

	template := parseTemplate(c)
	template.Execute(buffer, Output{Table: table.Format(data)})

	// Export to markdown or pdf (depends on the extension).
	// Show on the stdout if no export.
	if c.GlobalIsSet(outputFlag) {
		path := c.GlobalString(outputFlag)

		if strings.Contains(path, ".pdf") {
			createPDF(buffer, path, c.GlobalString(styleFlag))
		} else {
			createMarkdown(buffer, path)
		}
	} else {
		fmt.Println(buffer)
	}
}

func createPDF(buffer *bytes.Buffer, path string, cssPath string) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "temp")
	if err != nil {
		exit("Cannot export to PDF.")
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
		exit("Cannot generate the PDF.")
	}
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
	if !c.GlobalIsSet(templateFlag) {
		// If not available, use the default template.
		return useDefaultTemplate()
	}

	content, err := ioutil.ReadFile(c.GlobalString(templateFlag))
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
