// Copyright (c) 2015 Jelmer Snoeck (Gmail: jelmer.snoeck)
//
// Permission to use, copy, modify, and distribute this software for any purpose
// with or without fee is hereby granted, provided that the above copyright notice
// and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
// REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
// FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
// INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
// LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
// OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THIS SOFTWARE.

// Package barcode provides helper methods for adding barcodes of different
// types to your pdf document. It relies on the github.com/boombuler/barcode
// package for the barcode creation.
package barcode

import (
	"bytes"
	"errors"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/codabar"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode/twooffive"
	"github.com/jung-kurt/gofpdf"
	"image/jpeg"
	"strconv"
)

// barcodes represents the barcodes that have been registered through this
// package. They will later be used to be scaled and put on the page.
var barcodes map[string]barcode.Barcode

// Barcode puts a registered barcode in the current page.
//
// The size should be specified in the units used to create the PDF document.
//
// Positioning with x, y and flow is inherited from Fpdf.Image().
func Barcode(pdf *gofpdf.Fpdf, code string, x, y, w, h float64, flow bool) {
	unscaled, ok := barcodes[code]

	if !ok {
		err := errors.New("Barcode not found")
		pdf.SetError(err)
		return
	}

	bname := uniqueBarcodeName(code, x, y)
	info := pdf.GetImageInfo(bname)

	if info == nil {
		bcode, err := barcode.Scale(
			unscaled,
			int(convertTo96Dpi(pdf, w)),
			int(convertTo96Dpi(pdf, h)),
		)

		if err != nil {
			pdf.SetError(err)
			return
		}

		err = registerScaledBarcode(pdf, bname, bcode)
		if err != nil {
			pdf.SetError(err)
			return
		}
	}

	pdf.Image(bname, x, y, 0, 0, flow, "jpg", 0, "")
}

// Register registers a barcode but does not put it on the page. Use Barcode()
// with the same code to put the barcode on the PDF page.
func Register(bcode barcode.Barcode) string {
	if len(barcodes) == 0 {
		barcodes = make(map[string]barcode.Barcode)
	}

	key := barcodeKey(bcode)
	barcodes[key] = bcode
	return key
}

// RegisterCodabar registers a barcode of type Codabar to the PDF, but not to
// the page. Use Barcode() with the return value to put the barcode on the page.
func RegisterCodabar(pdf *gofpdf.Fpdf, code string) string {
	bcode, err := codabar.Encode(code)
	return registerBarcode(pdf, bcode, err)
}

// RegisterCode128 registers a barcode of type Code128 to the PDF, but not to
// the page. Use Barcode() with the return value to put the barcode on the page.
func RegisterCode128(pdf *gofpdf.Fpdf, code string) string {
	bcode, err := code128.Encode(code)
	return registerBarcode(pdf, bcode, err)
}

// RegisterCode39 registers a barcode of type Code39 to the PDF, but not to
// the page. Use Barcode() with the return value to put the barcode on the page.
//
// includeChecksum and fullASCIIMode are inherited from code39.Encode().
func RegisterCode39(pdf *gofpdf.Fpdf, code string, includeChecksum, fullASCIIMode bool) string {
	bcode, err := code39.Encode(code, includeChecksum, fullASCIIMode)
	return registerBarcode(pdf, bcode, err)
}

// RegisterDataMatrix registers a barcode of type DataMatrix to the PDF, but not
// to the page. Use Barcode() with the return value to put the barcode on the
// page.
func RegisterDataMatrix(pdf *gofpdf.Fpdf, code string) string {
	bcode, err := datamatrix.Encode(code)
	return registerBarcode(pdf, bcode, err)
}

// RegisterEAN registers a barcode of type EAN to the PDF, but not to the page.
// It will automatically detect if the barcode is EAN8 or EAN13. Use Barcode()
// with the return value to put the barcode on the page.
func RegisterEAN(pdf *gofpdf.Fpdf, code string) string {
	bcode, err := ean.Encode(code)
	return registerBarcode(pdf, bcode, err)
}

// RegisterQR registers a barcode of type QR to the PDF, but not to the page.
// Use Barcode() with the return value to put the barcode on the page.
//
// The ErrorCorrectionLevel and Encoding mode are inherited from qr.Encode().
func RegisterQR(pdf *gofpdf.Fpdf, code string, ecl qr.ErrorCorrectionLevel, mode qr.Encoding) string {
	bcode, err := qr.Encode(code, ecl, mode)
	return registerBarcode(pdf, bcode, err)
}

// RegisterTwoOfFive registers a barcode of type TwoOfFive to the PDF, but not
// to the page. Use Barcode() with the return value to put the barcode on the
// page.
//
// The interleaved bool is inherited from twooffive.Encode().
func RegisterTwoOfFive(pdf *gofpdf.Fpdf, code string, interleaved bool) string {
	bcode, err := twooffive.Encode(code, interleaved)
	return registerBarcode(pdf, bcode, err)
}

// registerBarcode registers a barcode internally using the Register() function.
// In case of an error generating the barcode it will not be registered and will
// set an error on the PDF. It will return a unique key for the barcode type and
// content that can be used to put the barcode on the page.
func registerBarcode(pdf *gofpdf.Fpdf, bcode barcode.Barcode, err error) string {
	if err != nil {
		pdf.SetError(err)
	}

	return Register(bcode)
}

// uniqueBarcodeName makes sure every barcode has a unique name for its
// dimensions. Scaling a barcode image results in quality loss, which could be
// a problem for barcode readers.
func uniqueBarcodeName(code string, x, y float64) string {
	xStr := strconv.FormatFloat(x, 'E', -1, 64)
	yStr := strconv.FormatFloat(y, 'E', -1, 64)

	return "barcode-" + code + "-" + xStr + yStr
}

// barcodeKey combines the code type and code value into a unique identifier for
// a barcode type. This is so that we can store several barcodes with the same
// code but different type in the barcodes map.
func barcodeKey(bcode barcode.Barcode) string {
	return bcode.Metadata().CodeKind + bcode.Content()
}

// registerScaledBarcode registers a barcode with its exact dimensions to the
// PDF but does not put it on the page. Use Fpdf.Image() with the same code to
// add the barcode to the page.
func registerScaledBarcode(pdf *gofpdf.Fpdf, code string, bcode barcode.Barcode) error {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, bcode, nil)

	if err != nil {
		return err
	}

	reader := bytes.NewReader(buf.Bytes())
	pdf.RegisterImageReader(code, "jpg", reader)

	return nil
}

// convertTo96DPI converts the given value, which is based on a 72 DPI value
// like the rest of the PDF document, to a 96 DPI value that is required for
// an Image.
//
// Doing this through the Fpdf.Image() function would mean that it uses a 72 DPI
// value and stretches it to a 96 DPI value. This results in quality loss which
// could be problematic for barcode scanners.
func convertTo96Dpi(pdf *gofpdf.Fpdf, value float64) float64 {
	return value * pdf.GetConversionRatio() / 72 * 96
}
