package barcode_test

import (
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/qr"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/barcode"
	"github.com/jung-kurt/gofpdf/internal/example"
)

func createPdf() (pdf *gofpdf.Fpdf) {
	pdf = gofpdf.New("L", "mm", "A4", "")
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetFillColor(200, 200, 220)
	pdf.AddPage()
	return
}

func ExampleRegister() {
	pdf := createPdf()

	fileStr := example.Filename("contrib_barcode_Register")

	bcode, err := code128.Encode("gofpdf")

	if err == nil {
		key := barcode.Register(bcode)
		barcode.Barcode(pdf, key, 15, 15, 100, 10, false)
	}

	err = pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_Register.pdf
}

func ExampleRegisterCodabar() {
	pdf := createPdf()

	key := barcode.RegisterCode128(pdf, "codabar")
	barcode.Barcode(pdf, key, 15, 15, 100, 10, false)

	fileStr := example.Filename("contrib_barcode_RegisterCodabar")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterCodabar.pdf
}

func ExampleRegisterAztec() {
	pdf := createPdf()

	key := barcode.RegisterAztec(pdf, "aztec", 33, 0)
	barcode.Barcode(pdf, key, 15, 15, 100, 100, false)

	fileStr := example.Filename("contrib_barcode_RegisterAztec")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterAztec.pdf
}

func ExampleRegisterCode128() {
	pdf := createPdf()

	key := barcode.RegisterCode128(pdf, "code128")
	barcode.Barcode(pdf, key, 15, 15, 100, 10, false)

	fileStr := example.Filename("contrib_barcode_RegisterCode128")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterCode128.pdf
}

func ExampleRegisterCode39() {
	pdf := createPdf()

	key := barcode.RegisterCode39(pdf, "CODE39", false, true)
	barcode.Barcode(pdf, key, 15, 15, 100, 10, false)

	fileStr := example.Filename("contrib_barcode_RegisterCode39")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterCode39.pdf
}

func ExampleRegisterDataMatrix() {
	pdf := createPdf()

	key := barcode.RegisterDataMatrix(pdf, "datamatrix")
	barcode.Barcode(pdf, key, 15, 15, 20, 20, false)

	fileStr := example.Filename("contrib_barcode_RegisterDataMatrix")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterDataMatrix.pdf
}

func ExampleRegisterEAN() {
	pdf := createPdf()

	key := barcode.RegisterEAN(pdf, "96385074")
	barcode.Barcode(pdf, key, 15, 15, 100, 10, false)

	fileStr := example.Filename("contrib_barcode_RegisterEAN")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterEAN.pdf
}

func ExampleRegisterQR() {
	pdf := createPdf()

	key := barcode.RegisterQR(pdf, "qrcode", qr.H, qr.Unicode)
	barcode.Barcode(pdf, key, 15, 15, 20, 20, false)

	fileStr := example.Filename("contrib_barcode_RegisterQR")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterQR.pdf
}

func ExampleRegisterTwoOfFive() {
	pdf := createPdf()

	key := barcode.RegisterTwoOfFive(pdf, "1234567895", true)
	barcode.Barcode(pdf, key, 15, 15, 100, 20, false)

	fileStr := example.Filename("contrib_barcode_RegisterTwoOfFive")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterTwoOfFive.pdf
}

func ExampleRegisterPdf417() {
	pdf := createPdf()

	key := barcode.RegisterPdf417(pdf, "1234567895", 10, 5)
	barcode.Barcode(pdf, key, 15, 15, 100, 20, false)

	fileStr := example.Filename("contrib_barcode_RegisterPdf417")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterPdf417.pdf
}
