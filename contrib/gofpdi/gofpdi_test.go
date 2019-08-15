package gofpdi

import (
	"bytes"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/internal/example"
	"io"
)

func ExampleGofpdiImporter() {
	// create new pdf
	pdf := gofpdf.New("P", "pt", "A4", "")

	// for testing purposes, get an arbitrary template pdf as stream
	rs, _ := getTemplatePdf()

	// create a new Importer instance
	imp := NewImporter()

	// import first page and determine page sizes
	tpl := imp.ImportPageFromStream(pdf, &rs, 1, "/MediaBox")
	pageSizes := imp.GetPageSizes()
	nrPages := len(imp.GetPageSizes())

	// add all pages from template pdf
	for i := 1; i <= nrPages; i++ {
		pdf.AddPage()
		if i > 1 {
			tpl = imp.ImportPageFromStream(pdf, &rs, i, "/MediaBox")
		}
		imp.UseImportedTemplate(pdf, tpl, 0, 0, pageSizes[i]["/MediaBox"]["w"], pageSizes[i]["/MediaBox"]["h"])
	}

	// output
	fileStr := example.Filename("contrib_gofpdi_Importer")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_gofpdi_Importer.pdf
}

func getTemplatePdf() (io.ReadSeeker, error) {
	tpdf := gofpdf.New("P", "pt", "A4", "")
	tpdf.AddPage()
	tpdf.SetFont("Arial", "", 12)
	tpdf.Text(20, 20, "Example Page 1")
	tpdf.AddPage()
	tpdf.Text(20, 20, "Example Page 2")
	tbuf := bytes.Buffer{}
	err := tpdf.Output(&tbuf)
	return bytes.NewReader(tbuf.Bytes()), err
}
