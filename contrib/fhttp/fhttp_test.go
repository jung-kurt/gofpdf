package fhttp_test

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/fhttp"
	"path/filepath"
)

const (
	cnGofpdfDir  = "./"
	cnExampleDir = cnGofpdfDir + "/pdf"
)

func exampleFilename(baseStr string) string {
	return filepath.Join(cnExampleDir, baseStr+".pdf")
}

func summary(err error, fileStr string) {
	if err == nil {
		fileStr = filepath.ToSlash(fileStr)
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		fmt.Println(err)
	}
}

func ExampleFpdf_AddHttpImage() {
	pdf := gofpdf.New("", "", "", "")
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetFillColor(200, 200, 220)
	pdf.AddPage()

	url := "https://github.com/jung-kurt/gofpdf/raw/master/image/logo_gofpdf.jpg?raw=true"
	fhttp.RegisterHttpImage(pdf, url, "")
	pdf.Image(url, 100, 100, 20, 20, false, "", 0, "")

	fileStr := exampleFilename("Fpdf_AddHttpImage")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_AddHttpImage.pdf
}
