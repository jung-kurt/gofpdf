package barcode_test

import (
	"testing"

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
		var width float64 = 100
		var height float64 = 10.0
		barcode.BarcodeUnscalable(pdf, key, 15, 15, &width, &height, false)
	}

	err = pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_Register.pdf
}

func ExampleRegisterCodabar() {
	pdf := createPdf()

	key := barcode.RegisterCode128(pdf, "codabar")
	var width float64 = 100
	var height float64 = 10
	barcode.BarcodeUnscalable(pdf, key, 15, 15, &width, &height, false)

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
	barcode.Barcode(pdf, key, 15, 15, 100, 10, false)

	fileStr := example.Filename("contrib_barcode_RegisterQR")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterQR.pdf
}

func ExampleRegisterTwoOfFive() {
	pdf := createPdf()

	key := barcode.RegisterTwoOfFive(pdf, "1234567895", true)
	barcode.Barcode(pdf, key, 15, 15, 100, 10, false)

	fileStr := example.Filename("contrib_barcode_RegisterTwoOfFive")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterTwoOfFive.pdf
}

func ExampleRegisterPdf417() {
	pdf := createPdf()

	key := barcode.RegisterPdf417(pdf, "1234567895", 10, 5)
	barcode.Barcode(pdf, key, 15, 15, 100, 10, false)

	fileStr := example.Filename("contrib_barcode_RegisterPdf417")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterPdf417.pdf
}

// TestRegisterCode128 ensures that no panic arises when an invalid barcode is registered.
func TestRegisterCode128(t *testing.T) {
	pdf := createPdf()
	barcode.RegisterCode128(pdf, "Invalid character: é")
}

// TestBarcodeUnscalable shows that the barcode may be scaled or not by providing optional heights and widths.
func TestBarcodeUnscalable(t *testing.T) {
	pdf := createPdf()

	key := barcode.RegisterCode128(pdf, "code128")
	var width float64 = 100
	var height float64 = 10
	barcode.BarcodeUnscalable(pdf, key, 15, 15, &width, &height, false)
	barcode.BarcodeUnscalable(pdf, key, 15, 35, nil, &height, false)
	barcode.BarcodeUnscalable(pdf, key, 15, 55, &width, nil, false)
	barcode.BarcodeUnscalable(pdf, key, 15, 75, nil, nil, false)

	fileStr := example.Filename("contrib_barcode_Barcode")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_Barcode.pdf
}

// TestGetUnscaledBarcodeDimensions shows that the width and height returned by the function match that of the barcode
func TestGetUnscaledBarcodeDimensions(t *testing.T) {
	pdf := createPdf()

	key := barcode.RegisterQR(pdf, "qrcode", qr.H, qr.Unicode)
	barcode.BarcodeUnscalable(pdf, key, 15, 15, nil, nil, false)
	w, h := barcode.GetUnscaledBarcodeDimensions(pdf, key)

	pdf.SetDrawColor(255, 0, 0)
	pdf.Line(15, 15, 15+w, 15+h)

	fileStr := example.Filename("contrib_barcode_GetBarcodeDimensions")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_GetBarcodeDimensions.pdf
}

// TestBarcodeNonIntegerScalingFactors shows that the barcode may be scaled to non-integer sizes
func TestBarcodeNonIntegerScalingFactors(t *testing.T) {
	pdf := gofpdf.New("L", "in", "A4", "")
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetFillColor(200, 200, 220)
	pdf.AddPage()

	key := barcode.RegisterQR(pdf, "qrcode", qr.H, qr.Unicode)
	var scale float64 = 1.5
	barcode.BarcodeUnscalable(pdf, key, 0.5, 0.5, &scale, &scale, false)

	pdf.SetDrawColor(255, 0, 0)
	pdf.Line(0.5, 0.5, 0.5+scale, 0.5+scale)

	fileStr := example.Filename("contrib_barcode_BarcodeScaling")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/contrib_barcode_BarcodeScaling.pdf
}
