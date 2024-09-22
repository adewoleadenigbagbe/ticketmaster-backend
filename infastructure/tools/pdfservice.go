package tools

import (
	"bytes"
	"html/template"
	"path/filepath"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

const (
	templatePath = "../infastructure/tools/invoice.html"
)

type PDFService struct {
}

func (pdfService PDFService) GeneratePDF(data any) ([]byte, error) {
	var err error
	var templ *template.Template
	filePath, err := filepath.Abs(templatePath)
	if err != nil {
		return nil, err
	}

	// use Go's default HTML template generation tools to generate your HTML
	if templ, err = template.ParseFiles(filePath); err != nil {
		return nil, err
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	if err = templ.Execute(&body, data); err != nil {
		return nil, err
	}

	// initalize a wkhtmltopdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))

	// enable this if the HTML file contains local references such as images, CSS, etc.
	page.EnableLocalFileAccess.Set(true)

	// add the page to your generator
	pdfg.AddPage(page)

	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)

	// magic
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
