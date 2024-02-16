package tools

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PDFService struct {
}

func (pdfService PDFService) GeneratePDF(data any) ([]byte, error) {
	var templ *template.Template
	var err error
	rootPath := "ticketmaster-backend"
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	index := strings.Index(currentWorkingDirectory, rootPath)
	if index == -1 {
		return nil, errors.New("App Root Folder Path not found")
	}

	file := filepath.Join(currentWorkingDirectory[:index], rootPath, "tools", "invoice.html")

	// use Go's default HTML template generation tools to generate your HTML
	if templ, err = template.ParseFiles(file); err != nil {
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
