/*
 * Copyright (c) 2013-2014 Kurt Jung (Gmail: kurt.w.jung)
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package gofpdf_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// Absolute path needed for gocov tool; relative OK for test
const (
	cnGofpdfDir  = "."
	cnFontDir    = cnGofpdfDir + "/font"
	cnImgDir     = cnGofpdfDir + "/image"
	cnTextDir    = cnGofpdfDir + "/text"
	cnExampleDir = cnGofpdfDir + "/pdf"
)

func init() {
	cleanup()
}

func cleanup() {
	filepath.Walk(cnExampleDir,
		func(path string, info os.FileInfo, err error) (reterr error) {
			if path[len(path)-4:] == ".pdf" {
				os.Remove(path)
			}
			return
		})
}

func imageFile(fileStr string) string {
	return filepath.Join(cnImgDir, fileStr)
}

func fontFile(fileStr string) string {
	return filepath.Join(cnFontDir, fileStr)
}

func textFile(fileStr string) string {
	return filepath.Join(cnTextDir, fileStr)
}

func exampleFile(fileStr string) string {
	return filepath.Join(cnExampleDir, fileStr)
}

type fontResourceType struct {
}

func (f fontResourceType) Open(name string) (rdr io.Reader, err error) {
	var buf []byte
	buf, err = ioutil.ReadFile(fontFile(name))
	if err == nil {
		rdr = bytes.NewReader(buf)
		fmt.Printf("Generalized font loader reading %s\n", name)
	}
	return
}

// Convert 'ABCDEFG' to, for example, 'A,BCD,EFG'
func strDelimit(str string, sepstr string, sepcount int) string {
	pos := len(str) - sepcount
	for pos > 0 {
		str = str[:pos] + sepstr + str[pos:]
		pos = pos - sepcount
	}
	return str
}

func lorem() string {
	return "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod " +
		"tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis " +
		"nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis " +
		"aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat " +
		"nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui " +
		"officia deserunt mollit anim id est laborum."
}

func exampleFilename(baseStr string) string {
	return filepath.Join(cnExampleDir, baseStr+".pdf")
}

func summary(err error, fileStr string) {
	if err == nil {
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		fmt.Println(err)
	}
}

// This example demonstrates the generation of a simple PDF document. Note that
// since only core fonts are used (in this case Arial, a synonym for
// Helvetica), an empty string can be specified for the font directory in the
// call to New(). Note also that the exampleFilename and summary functions are
// local to the test file and are not part of the gofpdf library. If an error
// occurs at some point during the construction of the document, subsequent
// method calls exit immediately and the error is finally retreived with the
// output call where it can be handled by the application.
func Example() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello World!")
	fileStr := exampleFilename("basic")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/basic.pdf
}

// This example demonsrates the generation of headers, footers and page breaks.
func ExampleFpdf_AddPage() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetHeaderFunc(func() {
		pdf.Image(imageFile("logo.png"), 10, 6, 30, 0, false, "", 0, "")
		pdf.SetY(5)
		pdf.SetFont("Arial", "B", 15)
		pdf.Cell(80, 0, "")
		pdf.CellFormat(30, 10, "Title", "1", 0, "C", false, 0, "")
		pdf.Ln(20)
	})
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.SetFont("Times", "", 12)
	for j := 1; j <= 40; j++ {
		pdf.CellFormat(0, 10, fmt.Sprintf("Printing line number %d", j),
			"", 1, "", false, 0, "")
	}
	fileStr := exampleFilename("Fpdf_AddPage")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_AddPage.pdf
}

// This example demonstrates word-wrapping, line justification and
// page-breaking.
func ExampleFpdf_MultiCell() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	titleStr := "20000 Leagues Under the Seas"
	pdf.SetTitle(titleStr, false)
	pdf.SetAuthor("Jules Verne", false)
	pdf.SetHeaderFunc(func() {
		// Arial bold 15
		pdf.SetFont("Arial", "B", 15)
		// Calculate width of title and position
		wd := pdf.GetStringWidth(titleStr) + 6
		pdf.SetX((210 - wd) / 2)
		// Colors of frame, background and text
		pdf.SetDrawColor(0, 80, 180)
		pdf.SetFillColor(230, 230, 0)
		pdf.SetTextColor(220, 50, 50)
		// Thickness of frame (1 mm)
		pdf.SetLineWidth(1)
		// Title
		pdf.CellFormat(wd, 9, titleStr, "1", 1, "C", true, 0, "")
		// Line break
		pdf.Ln(10)
	})
	pdf.SetFooterFunc(func() {
		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 8)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
	chapterTitle := func(chapNum int, titleStr string) {
		// 	// Arial 12
		pdf.SetFont("Arial", "", 12)
		// Background color
		pdf.SetFillColor(200, 220, 255)
		// Title
		pdf.CellFormat(0, 6, fmt.Sprintf("Chapter %d : %s", chapNum, titleStr),
			"", 1, "L", true, 0, "")
		// Line break
		pdf.Ln(4)
	}
	chapterBody := func(fileStr string) {
		// Read text file
		txtStr, err := ioutil.ReadFile(fileStr)
		if err != nil {
			pdf.SetError(err)
		}
		// Times 12
		pdf.SetFont("Times", "", 12)
		// Output justified text
		pdf.MultiCell(0, 5, string(txtStr), "", "", false)
		// Line break
		pdf.Ln(-1)
		// Mention in italics
		pdf.SetFont("", "I", 0)
		pdf.Cell(0, 5, "(end of excerpt)")
	}
	printChapter := func(chapNum int, titleStr, fileStr string) {
		pdf.AddPage()
		chapterTitle(chapNum, titleStr)
		chapterBody(fileStr)
	}
	printChapter(1, "A RUNAWAY REEF", textFile("20k_c1.txt"))
	printChapter(2, "THE PROS AND CONS", textFile("20k_c2.txt"))
	fileStr := exampleFilename("Fpdf_MultiCell")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_MultiCell.pdf
}

// This example demonstrates the generation of a PDF document that has multiple
// columns. This is accomplished with the SetLeftMargin() and Cell() methods.
func ExampleFpdf_SetLeftMargin() {
	var y0 float64
	var crrntCol int
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetDisplayMode("fullpage", "TwoColumnLeft")
	titleStr := "20000 Leagues Under the Seas"
	pdf.SetTitle(titleStr, false)
	pdf.SetAuthor("Jules Verne", false)
	setCol := func(col int) {
		// Set position at a given column
		crrntCol = col
		x := 10.0 + float64(col)*65.0
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
	}
	chapterTitle := func(chapNum int, titleStr string) {
		// Arial 12
		pdf.SetFont("Arial", "", 12)
		// Background color
		pdf.SetFillColor(200, 220, 255)
		// Title
		pdf.CellFormat(0, 6, fmt.Sprintf("Chapter %d : %s", chapNum, titleStr),
			"", 1, "L", true, 0, "")
		// Line break
		pdf.Ln(4)
		y0 = pdf.GetY()
	}
	chapterBody := func(fileStr string) {
		// Read text file
		txtStr, err := ioutil.ReadFile(fileStr)
		if err != nil {
			pdf.SetError(err)
		}
		// Font
		pdf.SetFont("Times", "", 12)
		// Output text in a 6 cm width column
		pdf.MultiCell(60, 5, string(txtStr), "", "", false)
		pdf.Ln(-1)
		// Mention
		pdf.SetFont("", "I", 0)
		pdf.Cell(0, 5, "(end of excerpt)")
		// Go back to first column
		setCol(0)
	}
	printChapter := func(num int, titleStr, fileStr string) {
		// Add chapter
		pdf.AddPage()
		chapterTitle(num, titleStr)
		chapterBody(fileStr)
	}
	pdf.SetAcceptPageBreakFunc(func() bool {
		// Method accepting or not automatic page break
		if crrntCol < 2 {
			// Go to next column
			setCol(crrntCol + 1)
			// Set ordinate to top
			pdf.SetY(y0)
			// Keep on page
			return false
		}
		// Go back to first column
		setCol(0)
		// Page break
		return true
	})
	pdf.SetHeaderFunc(func() {
		// Arial bold 15
		pdf.SetFont("Arial", "B", 15)
		// Calculate width of title and position
		wd := pdf.GetStringWidth(titleStr) + 6
		pdf.SetX((210 - wd) / 2)
		// Colors of frame, background and text
		pdf.SetDrawColor(0, 80, 180)
		pdf.SetFillColor(230, 230, 0)
		pdf.SetTextColor(220, 50, 50)
		// Thickness of frame (1 mm)
		pdf.SetLineWidth(1)
		// Title
		pdf.CellFormat(wd, 9, titleStr, "1", 1, "C", true, 0, "")
		// Line break
		pdf.Ln(10)
		// Save ordinate
		y0 = pdf.GetY()
	})
	pdf.SetFooterFunc(func() {
		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 8)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
	printChapter(1, "A RUNAWAY REEF", textFile("20k_c1.txt"))
	printChapter(2, "THE PROS AND CONS", textFile("20k_c2.txt"))
	fileStr := exampleFilename("Fpdf_SetLeftMargin_multicolumn")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_SetLeftMargin_multicolumn.pdf
}

// This example demonstrates various table styles.
func ExampleFpdf_CellFormat_1() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	type countryType struct {
		nameStr, capitalStr, areaStr, popStr string
	}
	countryList := make([]countryType, 0, 8)
	header := []string{"Country", "Capital", "Area (sq km)", "Pop. (thousands)"}
	loadData := func(fileStr string) {
		fl, err := os.Open(fileStr)
		if err == nil {
			scanner := bufio.NewScanner(fl)
			var c countryType
			for scanner.Scan() {
				// Austria;Vienna;83859;8075
				lineStr := scanner.Text()
				list := strings.Split(lineStr, ";")
				if len(list) == 4 {
					c.nameStr = list[0]
					c.capitalStr = list[1]
					c.areaStr = list[2]
					c.popStr = list[3]
					countryList = append(countryList, c)
				} else {
					err = fmt.Errorf("error tokenizing %s", lineStr)
				}
			}
			fl.Close()
			if len(countryList) == 0 {
				err = fmt.Errorf("error loading data from %s", fileStr)
			}
		}
		if err != nil {
			pdf.SetError(err)
		}
	}
	// Simple table
	basicTable := func() {
		for _, str := range header {
			pdf.CellFormat(40, 7, str, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
		for _, c := range countryList {
			pdf.CellFormat(40, 6, c.nameStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 6, c.capitalStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 6, c.areaStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 6, c.popStr, "1", 0, "", false, 0, "")
			pdf.Ln(-1)
		}
	}
	// Better table
	improvedTable := func() {
		// Column widths
		w := []float64{40.0, 35.0, 40.0, 45.0}
		wSum := 0.0
		for _, v := range w {
			wSum += v
		}
		// 	Header
		for j, str := range header {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
		// Data
		for _, c := range countryList {
			pdf.CellFormat(w[0], 6, c.nameStr, "LR", 0, "", false, 0, "")
			pdf.CellFormat(w[1], 6, c.capitalStr, "LR", 0, "", false, 0, "")
			pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3),
				"LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
				"LR", 0, "R", false, 0, "")
			pdf.Ln(-1)
		}
		pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	}
	// Colored table
	fancyTable := func() {
		// Colors, line width and bold font
		pdf.SetFillColor(255, 0, 0)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetDrawColor(128, 0, 0)
		pdf.SetLineWidth(.3)
		pdf.SetFont("", "B", 0)
		// 	Header
		w := []float64{40, 35, 40, 45}
		wSum := 0.0
		for _, v := range w {
			wSum += v
		}
		for j, str := range header {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)
		// Color and font restoration
		pdf.SetFillColor(224, 235, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("", "", 0)
		// 	Data
		fill := false
		for _, c := range countryList {
			pdf.CellFormat(w[0], 6, c.nameStr, "LR", 0, "", fill, 0, "")
			pdf.CellFormat(w[1], 6, c.capitalStr, "LR", 0, "", fill, 0, "")
			pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3),
				"LR", 0, "R", fill, 0, "")
			pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
				"LR", 0, "R", fill, 0, "")
			pdf.Ln(-1)
			fill = !fill
		}
		pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	}
	loadData(textFile("countries.txt"))
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()
	basicTable()
	pdf.AddPage()
	improvedTable()
	pdf.AddPage()
	fancyTable()
	fileStr := exampleFilename("Fpdf_CellFormat_1_tables")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_CellFormat_1_tables.pdf
}

// This example demonstrates internal and external links with and without basic
// HTML.
func ExampleFpdf_HTMLBasicNew() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	// First page: manual local link
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 20)
	_, lineHt := pdf.GetFontSize()
	pdf.Write(lineHt, "To find out what's new in this tutorial, click ")
	pdf.SetFont("", "U", 0)
	link := pdf.AddLink()
	pdf.WriteLinkID(lineHt, "here", link)
	pdf.SetFont("", "", 0)
	// Second page: image link and basic HTML with link
	pdf.AddPage()
	pdf.SetLink(link, 0, -1)
	pdf.Image(imageFile("logo.png"), 10, 12, 30, 0, false, "", 0, "http://www.fpdf.org")
	pdf.SetLeftMargin(45)
	pdf.SetFontSize(14)
	_, lineHt = pdf.GetFontSize()
	htmlStr := `You can now easily print text mixing different styles: <b>bold</b>, ` +
		`<i>italic</i>, <u>underlined</u>, or <b><i><u>all at once</u></i></b>!<br><br>` +
		`You can also insert links on text, such as ` +
		`<a href="http://www.fpdf.org">www.fpdf.org</a>, or on an image: click on the logo.`
	html := pdf.HTMLBasicNew()
	html.Write(lineHt, htmlStr)
	fileStr := exampleFilename("Fpdf_HTMLBasicNew")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_HTMLBasicNew.pdf
}

// This example demonstrates the use of a non-standard font.
func ExampleFpdf_AddFont() {
	pdf := gofpdf.New("P", "mm", "A4", cnFontDir)
	pdf.AddFont("Calligrapher", "", "calligra.json")
	pdf.AddPage()
	pdf.SetFont("Calligrapher", "", 35)
	pdf.Cell(0, 10, "Enjoy new fonts with FPDF!")
	fileStr := exampleFilename("Fpdf_AddFont")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_AddFont.pdf
}

// This example demonstrates how images are included in documents.
func ExampleFpdf_Image() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 11)
	pdf.Image(imageFile("logo.png"), 10, 10, 30, 0, false, "", 0, "")
	pdf.Text(50, 20, "logo.png")
	pdf.Image(imageFile("logo.gif"), 10, 40, 30, 0, false, "", 0, "")
	pdf.Text(50, 50, "logo.gif")
	pdf.Image(imageFile("logo-gray.png"), 10, 70, 30, 0, false, "", 0, "")
	pdf.Text(50, 80, "logo-gray.png")
	pdf.Image(imageFile("logo-rgb.png"), 10, 100, 30, 0, false, "", 0, "")
	pdf.Text(50, 110, "logo-rgb.png")
	pdf.Image(imageFile("logo.jpg"), 10, 130, 30, 0, false, "", 0, "")
	pdf.Text(50, 140, "logo.jpg")
	fileStr := exampleFilename("Fpdf_Image")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_Image.pdf
}

// This examples demonstrates Landscape mode with images.
func ExampleFpdf_SetAcceptPageBreakFunc() {
	var y0 float64
	var crrntCol int
	loremStr := lorem()
	pdf := gofpdf.New("L", "mm", "A4", "")
	const (
		pageWd = 297.0 // A4 210.0 x 297.0
		margin = 10.0
		gutter = 4
		colNum = 3
		colWd  = (pageWd - 2*margin - (colNum-1)*gutter) / colNum
	)
	setCol := func(col int) {
		crrntCol = col
		x := margin + float64(col)*(colWd+gutter)
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
	}
	pdf.SetHeaderFunc(func() {
		titleStr := "gofpdf"
		pdf.SetFont("Helvetica", "B", 48)
		wd := pdf.GetStringWidth(titleStr) + 6
		pdf.SetX((pageWd - wd) / 2)
		pdf.SetTextColor(128, 128, 160)
		pdf.Write(12, titleStr[:2])
		pdf.SetTextColor(128, 128, 128)
		pdf.Write(12, titleStr[2:])
		pdf.Ln(20)
		y0 = pdf.GetY()
	})
	pdf.SetAcceptPageBreakFunc(func() bool {
		if crrntCol < colNum-1 {
			setCol(crrntCol + 1)
			pdf.SetY(y0)
			// Start new column, not new page
			return false
		}
		setCol(0)
		return true
	})
	pdf.AddPage()
	pdf.SetFont("Times", "", 12)
	for j := 0; j < 20; j++ {
		if j == 1 {
			pdf.Image(imageFile("fpdf.png"), -1, 0, colWd, 0, true, "", 0, "")
		} else if j == 5 {
			pdf.Image(imageFile("golang-gopher.png"),
				-1, 0, colWd, 0, true, "", 0, "")
		}
		pdf.MultiCell(colWd, 5, loremStr, "", "", false)
		pdf.Ln(-1)
	}
	fileStr := exampleFilename("Fpdf_SetAcceptPageBreakFunc_landscape")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_SetAcceptPageBreakFunc_landscape.pdf
}

// This examples tests corner cases as reported by the gocov tool.
func ExampleFpdf_SetKeywords() {
	var err error
	fileStr := exampleFilename("Fpdf_SetKeywords")
	err = gofpdf.MakeFont(fontFile("CalligrapherRegular.pfb"),
		fontFile("cp1252.map"), cnFontDir, nil, true)
	if err == nil {
		err = gofpdf.MakeFont(fontFile("calligra.ttf"),
			fontFile("cp1252.map"), cnFontDir, nil, true)
		if err == nil {
			pdf := gofpdf.New("", "", "", "")
			pdf.SetFontLocation(cnFontDir)
			pdf.SetTitle("世界", true)
			pdf.SetAuthor("世界", true)
			pdf.SetSubject("世界", true)
			pdf.SetCreator("世界", true)
			pdf.SetKeywords("世界", true)
			pdf.AddFont("Calligrapher", "", "CalligrapherRegular.json")
			pdf.AddPage()
			pdf.SetFont("Calligrapher", "", 16)
			pdf.Writef(5, "\x95 %s \x95", pdf)
			err = pdf.OutputFileAndClose(fileStr)
		}
	}
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_SetKeywords.pdf
}

// This example demonstrates the construction of various geometric figures,
func ExampleFpdf_Circle() {
	const (
		thin  = 0.2
		thick = 3.0
	)
	pdf := gofpdf.New("", "", "", "")
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetFillColor(200, 200, 220)
	pdf.AddPage()

	y := 15.0
	pdf.Text(10, y, "Circles")
	pdf.SetFillColor(200, 200, 220)
	pdf.SetLineWidth(thin)
	pdf.Circle(20, y+15, 10, "D")
	pdf.Circle(45, y+15, 10, "F")
	pdf.Circle(70, y+15, 10, "FD")
	pdf.SetLineWidth(thick)
	pdf.Circle(95, y+15, 10, "FD")
	pdf.SetLineWidth(thin)

	y += 40.0
	pdf.Text(10, y, "Ellipses")
	pdf.SetFillColor(220, 200, 200)
	pdf.Ellipse(30, y+15, 20, 10, 0, "D")
	pdf.Ellipse(75, y+15, 20, 10, 0, "F")
	pdf.Ellipse(120, y+15, 20, 10, 0, "FD")
	pdf.SetLineWidth(thick)
	pdf.Ellipse(165, y+15, 20, 10, 0, "FD")
	pdf.SetLineWidth(thin)

	y += 40.0
	pdf.Text(10, y, "Curves (quadratic)")
	pdf.SetFillColor(220, 220, 200)
	pdf.Curve(10, y+30, 15, y-20, 40, y+30, "D")
	pdf.Curve(45, y+30, 50, y-20, 75, y+30, "F")
	pdf.Curve(80, y+30, 85, y-20, 110, y+30, "FD")
	pdf.SetLineWidth(thick)
	pdf.Curve(115, y+30, 120, y-20, 145, y+30, "FD")
	pdf.SetLineCapStyle("round")
	pdf.Curve(150, y+30, 155, y-20, 180, y+30, "FD")
	pdf.SetLineWidth(thin)
	pdf.SetLineCapStyle("butt")

	y += 40.0
	pdf.Text(10, y, "Curves (cubic)")
	pdf.SetFillColor(220, 200, 220)
	pdf.CurveBezierCubic(10, y+30, 15, y-20, 10, y+30, 40, y+30, "D")
	pdf.CurveBezierCubic(45, y+30, 50, y-20, 45, y+30, 75, y+30, "F")
	pdf.CurveBezierCubic(80, y+30, 85, y-20, 80, y+30, 110, y+30, "FD")
	pdf.SetLineWidth(thick)
	pdf.CurveBezierCubic(115, y+30, 120, y-20, 115, y+30, 145, y+30, "FD")
	pdf.SetLineCapStyle("round")
	pdf.CurveBezierCubic(150, y+30, 155, y-20, 150, y+30, 180, y+30, "FD")
	pdf.SetLineWidth(thin)
	pdf.SetLineCapStyle("butt")

	y += 40.0
	pdf.Text(10, y, "Arcs")
	pdf.SetFillColor(200, 220, 220)
	pdf.SetLineWidth(thick)
	pdf.Arc(45, y+35, 20, 10, 0, 0, 180, "FD")
	pdf.SetLineWidth(thin)
	pdf.Arc(45, y+35, 25, 15, 0, 90, 270, "D")
	pdf.SetLineWidth(thick)
	pdf.Arc(45, y+35, 30, 20, 0, 0, 360, "D")
	pdf.SetLineCapStyle("round")
	pdf.Arc(135, y+35, 20, 10, 135, 0, 180, "FD")
	pdf.SetLineWidth(thin)
	pdf.Arc(135, y+35, 25, 15, 135, 90, 270, "D")
	pdf.SetLineWidth(thick)
	pdf.Arc(135, y+35, 30, 20, 135, 0, 360, "D")
	pdf.SetLineWidth(thin)
	pdf.SetLineCapStyle("butt")

	fileStr := exampleFilename("Fpdf_Circle_figures")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_Circle_figures.pdf
}

// This example demonstrates alpha transparency.
func ExampleFpdf_SetAlpha() {
	const (
		gapX  = 10.0
		gapY  = 9.0
		rectW = 40.0
		rectH = 58.0
		pageW = 210
		pageH = 297
	)
	modeList := []string{"Normal", "Multiply", "Screen", "Overlay",
		"Darken", "Lighten", "ColorDodge", "ColorBurn", "HardLight", "SoftLight",
		"Difference", "Exclusion", "Hue", "Saturation", "Color", "Luminosity"}
	pdf := gofpdf.New("", "", "", "")
	pdf.SetLineWidth(2)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 18)
	pdf.SetXY(0, gapY)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(pageW, gapY, "Alpha Blending Modes", "", 0, "C", false, 0, "")
	j := 0
	y := 3 * gapY
	for col := 0; col < 4; col++ {
		x := gapX
		for row := 0; row < 4; row++ {
			pdf.Rect(x, y, rectW, rectH, "D")
			pdf.SetFont("Helvetica", "B", 12)
			pdf.SetFillColor(0, 0, 0)
			pdf.SetTextColor(250, 250, 230)
			pdf.SetXY(x, y+rectH-4)
			pdf.CellFormat(rectW, 5, modeList[j], "", 0, "C", true, 0, "")
			pdf.SetFont("Helvetica", "I", 150)
			pdf.SetTextColor(80, 80, 120)
			pdf.SetXY(x, y+2)
			pdf.CellFormat(rectW, rectH, "A", "", 0, "C", false, 0, "")
			pdf.SetAlpha(0.5, modeList[j])
			pdf.Image(imageFile("golang-gopher.png"),
				x-gapX, y, rectW+2*gapX, 0, false, "", 0, "")
			pdf.SetAlpha(1.0, "Normal")
			x += rectW + gapX
			j++
		}
		y += rectH + gapY
	}
	fileStr := exampleFilename("Fpdf_SetAlpha_transparency")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_SetAlpha_transparency.pdf
}

// This example deomstrates various gradients.
func ExampleFpdf_LinearGradient() {
	pdf := gofpdf.New("", "", "", "")
	pdf.SetFont("Helvetica", "", 12)
	pdf.AddPage()
	pdf.LinearGradient(0, 0, 210, 100, 250, 250, 255, 220, 220, 225, 0, 0, 0, .5)
	pdf.LinearGradient(20, 25, 75, 75, 220, 220, 250, 80, 80, 220, 0, .2, 0, .8)
	pdf.Rect(20, 25, 75, 75, "D")
	pdf.LinearGradient(115, 25, 75, 75, 220, 220, 250, 80, 80, 220, 0, 0, 1, 1)
	pdf.Rect(115, 25, 75, 75, "D")
	pdf.RadialGradient(20, 120, 75, 75, 220, 220, 250, 80, 80, 220,
		0.25, 0.75, 0.25, 0.75, 1)
	pdf.Rect(20, 120, 75, 75, "D")
	pdf.RadialGradient(115, 120, 75, 75, 220, 220, 250, 80, 80, 220,
		0.25, 0.75, 0.75, 0.75, 0.75)
	pdf.Rect(115, 120, 75, 75, "D")
	fileStr := exampleFilename("Fpdf_LinearGradient_gradient")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_LinearGradient_gradient.pdf
}

// This example demonstrates clipping.
func ExampleFpdf_ClipText() {
	pdf := gofpdf.New("", "", "", "")
	y := 10.0
	pdf.AddPage()

	pdf.SetFont("Helvetica", "", 24)
	pdf.SetXY(0, y)
	pdf.ClipText(10, y+12, "Clipping examples", false)
	pdf.RadialGradient(10, y, 100, 20, 128, 128, 160, 32, 32, 48,
		0.25, 0.5, 0.25, 0.5, 0.2)
	pdf.ClipEnd()

	y += 12
	pdf.SetFont("Helvetica", "B", 120)
	pdf.SetDrawColor(64, 80, 80)
	pdf.SetLineWidth(.5)
	pdf.ClipText(10, y+40, pdf.String(), true)
	pdf.RadialGradient(10, y, 200, 50, 220, 220, 250, 80, 80, 220,
		0.25, 0.5, 0.25, 0.5, 1)
	pdf.ClipEnd()

	y += 55
	pdf.ClipRect(10, y, 105, 20, true)
	pdf.SetFillColor(255, 255, 255)
	pdf.Rect(10, y, 105, 20, "F")
	pdf.ClipCircle(40, y+10, 15, false)
	pdf.RadialGradient(25, y, 30, 30, 220, 250, 220, 40, 60, 40, 0.3,
		0.85, 0.3, 0.85, 0.5)
	pdf.ClipEnd()
	pdf.ClipEllipse(80, y+10, 20, 15, false)
	pdf.RadialGradient(60, y, 40, 30, 250, 220, 220, 60, 40, 40, 0.3,
		0.85, 0.3, 0.85, 0.5)
	pdf.ClipEnd()
	pdf.ClipEnd()

	y += 28
	pdf.ClipEllipse(26, y+10, 16, 10, true)
	pdf.Image(imageFile("logo.jpg"), 10, y, 32, 0, false, "JPG", 0, "")
	pdf.ClipEnd()

	pdf.ClipCircle(60, y+10, 10, true)
	pdf.RadialGradient(50, y, 20, 20, 220, 220, 250, 40, 40, 60, 0.3,
		0.7, 0.3, 0.7, 0.5)
	pdf.ClipEnd()

	pdf.ClipPolygon([]gofpdf.PointType{{80, y + 20}, {90, y},
		{100, y + 20}}, true)
	pdf.LinearGradient(80, y, 20, 20, 250, 220, 250, 60, 40, 60, 0.5,
		1, 0.5, 0.5)
	pdf.ClipEnd()

	y += 30
	pdf.SetLineWidth(.1)
	pdf.SetDrawColor(180, 180, 180)
	pdf.ClipRoundedRect(10, y, 120, 20, 5, true)
	pdf.RadialGradient(10, y, 120, 20, 255, 255, 255, 240, 240, 220,
		0.25, 0.75, 0.25, 0.75, 0.5)
	pdf.SetXY(5, y-5)
	pdf.SetFont("Times", "", 12)
	pdf.MultiCell(130, 5, lorem(), "", "", false)
	pdf.ClipEnd()

	fileStr := exampleFilename("Fpdf_ClipText")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_ClipText.pdf
}

// This example generates a PDF document with various page sizes.
func ExampleFpdf_PageSize() {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:    "in",
		Size:       gofpdf.SizeType{Wd: 6, Ht: 6},
		FontDirStr: cnFontDir,
	})
	pdf.SetMargins(0.5, 1, 0.5)
	pdf.SetFont("Times", "", 14)
	pdf.AddPageFormat("L", gofpdf.SizeType{Wd: 3, Ht: 12})
	pdf.SetXY(0.5, 1.5)
	pdf.CellFormat(11, 0.2, "12 in x 3 in", "", 0, "C", false, 0, "")
	pdf.AddPage() // Default size established in NewCustom()
	pdf.SetXY(0.5, 3)
	pdf.CellFormat(5, 0.2, "6 in x 6 in", "", 0, "C", false, 0, "")
	pdf.AddPageFormat("P", gofpdf.SizeType{Wd: 3, Ht: 12})
	pdf.SetXY(0.5, 6)
	pdf.CellFormat(2, 0.2, "3 in x 12 in", "", 0, "C", false, 0, "")
	for j := 0; j <= 3; j++ {
		wd, ht, u := pdf.PageSize(j)
		fmt.Printf("%d: %6.2f %s, %6.2f %s\n", j, wd, u, ht, u)
	}
	fileStr := exampleFilename("Fpdf_PageSize")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// 0:   6.00 in,   6.00 in
	// 1:  12.00 in,   3.00 in
	// 2:   6.00 in,   6.00 in
	// 3:   3.00 in,  12.00 in
	// Successfully generated pdf/Fpdf_PageSize.pdf
}

// This example demonstrates the Bookmark method.
func ExampleFpdf_Bookmark() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 15)
	pdf.Bookmark("Page 1", 0, 0)
	pdf.Bookmark("Paragraph 1", 1, -1)
	pdf.Cell(0, 6, "Paragraph 1")
	pdf.Ln(50)
	pdf.Bookmark("Paragraph 2", 1, -1)
	pdf.Cell(0, 6, "Paragraph 2")
	pdf.AddPage()
	pdf.Bookmark("Page 2", 0, 0)
	pdf.Bookmark("Paragraph 3", 1, -1)
	pdf.Cell(0, 6, "Paragraph 3")
	fileStr := exampleFilename("Fpdf_Bookmark")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_Bookmark.pdf
}

// This example demonstrates various transformations. It is adapted from an
// example script by Moritz Wagner and Andreas Würmser.
func ExampleFpdf_TransformBegin() {
	const (
		light = 200
		dark  = 0
	)
	var refX, refY float64
	var refStr string
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	color := func(val int) {
		pdf.SetDrawColor(val, val, val)
		pdf.SetTextColor(val, val, val)
	}
	reference := func(str string, x, y float64, val int) {
		color(val)
		pdf.Rect(x, y, 40, 10, "D")
		pdf.Text(x, y-1, str)
	}
	refDraw := func(str string, x, y float64) {
		refStr = str
		refX = x
		refY = y
		reference(str, x, y, light)
	}
	refDupe := func() {
		reference(refStr, refX, refY, dark)
	}

	titleStr := "Transformations"
	titlePt := 36.0
	titleHt := pdf.PointConvert(titlePt)
	pdf.SetFont("Helvetica", "", titlePt)
	titleWd := pdf.GetStringWidth(titleStr)
	titleX := (210 - titleWd) / 2
	pdf.Text(titleX, 10+titleHt, titleStr)
	pdf.TransformBegin()
	pdf.TransformMirrorVertical(10 + titleHt + 0.5)
	pdf.ClipText(titleX, 10+titleHt, titleStr, false)
	// Remember that the transform will mirror the gradient box too
	pdf.LinearGradient(titleX, 10, titleWd, titleHt+4, 120, 120, 120,
		255, 255, 255, 0, 0, 0, 0.6)
	pdf.ClipEnd()
	pdf.TransformEnd()

	pdf.SetFont("Helvetica", "", 12)

	// Scale by 150% centered by lower left corner of the rectangle
	refDraw("Scale", 50, 60)
	pdf.TransformBegin()
	pdf.TransformScaleXY(150, 50, 70)
	refDupe()
	pdf.TransformEnd()

	// Translate 7 to the right, 5 to the bottom
	refDraw("Translate", 125, 60)
	pdf.TransformBegin()
	pdf.TransformTranslate(7, 5)
	refDupe()
	pdf.TransformEnd()

	// Rotate 20 degrees counter-clockwise centered by the lower left corner of
	// the rectangle
	refDraw("Rotate", 50, 110)
	pdf.TransformBegin()
	pdf.TransformRotate(20, 50, 120)
	refDupe()
	pdf.TransformEnd()

	// Skew 30 degrees along the x-axis centered by the lower left corner of the
	// rectangle
	refDraw("Skew", 125, 110)
	pdf.TransformBegin()
	pdf.TransformSkewX(30, 125, 110)
	refDupe()
	pdf.TransformEnd()

	// Mirror horizontally with axis of reflection at left side of the rectangle
	refDraw("Mirror horizontal", 50, 160)
	pdf.TransformBegin()
	pdf.TransformMirrorHorizontal(50)
	refDupe()
	pdf.TransformEnd()

	// Mirror vertically with axis of reflection at bottom side of the rectangle
	refDraw("Mirror vertical", 125, 160)
	pdf.TransformBegin()
	pdf.TransformMirrorVertical(170)
	refDupe()
	pdf.TransformEnd()

	// Reflect against a point at the lower left point of rectangle
	refDraw("Mirror point", 50, 210)
	pdf.TransformBegin()
	pdf.TransformMirrorPoint(50, 220)
	refDupe()
	pdf.TransformEnd()

	// Mirror against a straight line described by a point and an angle
	angle := -20.0
	px := 120.0
	py := 220.0
	refDraw("Mirror line", 125, 210)
	pdf.TransformBegin()
	pdf.TransformRotate(angle, px, py)
	pdf.Line(px-1, py-1, px+1, py+1)
	pdf.Line(px-1, py+1, px+1, py-1)
	pdf.Line(px-5, py, px+60, py)
	pdf.TransformEnd()
	pdf.TransformBegin()
	pdf.TransformMirrorLine(angle, px, py)
	refDupe()
	pdf.TransformEnd()

	fileStr := exampleFilename("Fpdf_TransformBegin")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_TransformBegin.pdf
}

// This example demonstrates Lawrence Kesteloot's image registration code.
func ExampleFpdf_RegisterImage() {
	const (
		margin = 10
		wd     = 210
		ht     = 297
	)
	fileList := []string{
		"logo-gray.png",
		"logo.jpg",
		"logo.png",
		"logo-rgb.png",
		"logo-progressive.jpg",
	}
	var infoPtr *gofpdf.ImageInfoType
	var imageFileStr string
	var imgWd, imgHt, lf, tp float64
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(10, 10, 10)
	pdf.SetFont("Helvetica", "", 15)
	for j, str := range fileList {
		imageFileStr = imageFile(str)
		infoPtr = pdf.RegisterImage(imageFileStr, "")
		imgWd, imgHt = infoPtr.Extent()
		switch j {
		case 0:
			lf = margin
			tp = margin
		case 1:
			lf = wd - margin - imgWd
			tp = margin
		case 2:
			lf = (wd - imgWd) / 2.0
			tp = (ht - imgHt) / 2.0
		case 3:
			lf = margin
			tp = ht - imgHt - margin
		case 4:
			lf = wd - imgWd - margin
			tp = ht - imgHt - margin
		}
		pdf.Image(imageFileStr, lf, tp, imgWd, imgHt, false, "", 0, "")
	}
	fileStr := exampleFilename("Fpdf_RegisterImage")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_RegisterImage.pdf
}

// This example demonstrates Bruno Michel's line splitting function.
func ExampleFpdf_SplitLines() {
	const (
		fontPtSize = 18.0
		wd         = 100.0
	)
	pdf := gofpdf.New("P", "mm", "A4", "") // A4 210.0 x 297.0
	pdf.SetFont("Times", "", fontPtSize)
	_, lineHt := pdf.GetFontSize()
	pdf.AddPage()
	pdf.SetMargins(10, 10, 10)
	lines := pdf.SplitLines([]byte(lorem()), wd)
	ht := float64(len(lines)) * lineHt
	y := (297.0 - ht) / 2.0
	pdf.SetDrawColor(128, 128, 128)
	pdf.SetFillColor(255, 255, 210)
	x := (210.0 - (wd + 40.0)) / 2.0
	pdf.Rect(x, y-20.0, wd+40.0, ht+40.0, "FD")
	pdf.SetY(y)
	for _, line := range lines {
		pdf.CellFormat(190.0, lineHt, string(line), "", 1, "C", false, 0, "")
	}
	fileStr := exampleFilename("Fpdf_Splitlines")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_Splitlines.pdf
}

// This example demonstrates how to render a simple path-only SVG image of the
// type generated by the jSignature web control.
func ExampleFpdf_SVGBasicWrite() {
	const (
		fontPtSize = 16.0
		wd         = 100.0
		sigFileStr = "signature.svg"
	)
	var (
		sig gofpdf.SVGBasicType
		err error
	)
	pdf := gofpdf.New("P", "mm", "A4", "") // A4 210.0 x 297.0
	pdf.SetFont("Times", "", fontPtSize)
	lineHt := pdf.PointConvert(fontPtSize)
	pdf.AddPage()
	pdf.SetMargins(10, 10, 10)
	htmlStr := `This example renders a simple ` +
		`<a href="http://www.w3.org/TR/SVG/">SVG</a> (scalable vector graphics) ` +
		`image that contains only basic path commands without any styling, ` +
		`color fill, reflection or endpoint closures. In particular, the ` +
		`type of vector graphic returned from a ` +
		`<a href="http://willowsystems.github.io/jSignature/#/demo/">jSignature</a> ` +
		`web control is supported and is used in this example.`
	html := pdf.HTMLBasicNew()
	html.Write(lineHt, htmlStr)
	sig, err = gofpdf.SVGBasicFileParse(imageFile(sigFileStr))
	if err == nil {
		scale := 100 / sig.Wd
		scaleY := 30 / sig.Ht
		if scale > scaleY {
			scale = scaleY
		}
		pdf.SetLineCapStyle("round")
		pdf.SetLineWidth(0.25)
		pdf.SetDrawColor(0, 0, 128)
		pdf.SetXY((210.0-scale*sig.Wd)/2.0, pdf.GetY()+10)
		pdf.SVGBasicWrite(&sig, scale)
	} else {
		pdf.SetError(err)
	}
	fileStr := exampleFilename("Fpdf_SVGBasicWrite")
	err = pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_SVGBasicWrite.pdf
}

// This example demonstrates Stefan Schroeder's code to control vertical
// alignment.
func ExampleFpdf_CellFormat_2() {
	type recType struct {
		align, txt string
	}
	recList := []recType{
		recType{"TL", "top left"},
		recType{"TC", "top center"},
		recType{"TR", "top right"},
		recType{"LM", "middle left"},
		recType{"CM", "middle center"},
		recType{"RM", "middle right"},
		recType{"BL", "bottom left"},
		recType{"BC", "bottom center"},
		recType{"BR", "bottom right"},
	}
	recListBaseline := []recType{
		recType{"AL", "baseline left"},
		recType{"AC", "baseline center"},
		recType{"AR", "baseline right"},
	}
	var formatRect = func(pdf *gofpdf.Fpdf, recList []recType) {
		linkStr := ""
		for pageJ := 0; pageJ < 2; pageJ++ {
			pdf.AddPage()
			pdf.SetMargins(10, 10, 10)
			pdf.SetAutoPageBreak(false, 0)
			borderStr := "1"
			for _, rec := range recList {
				pdf.SetXY(20, 20)
				pdf.CellFormat(170, 257, rec.txt, borderStr, 0, rec.align, false, 0, linkStr)
				borderStr = ""
			}
			linkStr = "https://github.com/jung-kurt/gofpdf"
		}
	}
	pdf := gofpdf.New("P", "mm", "A4", "") // A4 210.0 x 297.0
	pdf.SetFont("Helvetica", "", 16)
	formatRect(pdf, recList)
	formatRect(pdf, recListBaseline)
	var fr fontResourceType
	pdf.SetFontLoader(fr)
	pdf.AddFont("Calligrapher", "", "calligra.json")
	pdf.SetFont("Calligrapher", "", 16)
	formatRect(pdf, recListBaseline)
	fileStr := exampleFilename("Fpdf_CellFormat_2_align")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Generalized font loader reading calligra.json
	// Generalized font loader reading calligra.z
	// Successfully generated pdf/Fpdf_CellFormat_2_align.pdf
}

// This example demonstrates the use of characters in the high range of the
// Windows-1252 code page (gofdpf default). See the example for CellFormat (4)
// for a way to do this automatically.
func ExampleFpdf_CellFormat_3() {
	pdf := gofpdf.New("P", "mm", "A4", "") // A4 210.0 x 297.0
	fontSize := 16.0
	pdf.SetFont("Helvetica", "", fontSize)
	ht := pdf.PointConvert(fontSize)
	write := func(str string) {
		pdf.CellFormat(190, ht, str, "", 1, "C", false, 0, "")
		pdf.Ln(ht)
	}
	pdf.AddPage()
	htmlStr := `Until gofpdf supports UTF-8 encoded source text, source text needs ` +
		`to be specified with all special characters escaped to match the code page ` +
		`layout of the currently selected font. By default, gofdpf uses code page 1252.` +
		` See <a href="http://en.wikipedia.org/wiki/Windows-1252">Wikipedia</a> for ` +
		`a table of this layout.`
	html := pdf.HTMLBasicNew()
	html.Write(ht, htmlStr)
	pdf.Ln(2 * ht)
	write("Voix ambigu\xeb d'un c\x9cur qui au z\xe9phyr pr\xe9f\xe8re les jattes de kiwi.")
	write("Falsches \xdcben von Xylophonmusik qu\xe4lt jeden gr\xf6\xdferen Zwerg.")
	write("Heiz\xf6lr\xfccksto\xdfabd\xe4mpfung")
	write("For\xe5rsj\xe6vnd\xf8gn / Efter\xe5rsj\xe6vnd\xf8gn")
	fileStr := exampleFilename("Fpdf_CellFormat_3_codepageescape")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_CellFormat_3_codepageescape.pdf
}

// This example demonstrates the automatic conversion of UTF-8 strings to an
// 8-bit font encoding.
func ExampleFpdf_CellFormat_4() {
	pdf := gofpdf.New("P", "mm", "A4", cnFontDir) // A4 210.0 x 297.0
	// See documentation for details on how to generate fonts
	pdf.AddFont("Helvetica-1251", "", "helvetica_1251.json")
	pdf.AddFont("Helvetica-1253", "", "helvetica_1253.json")
	fontSize := 16.0
	pdf.SetFont("Helvetica", "", fontSize)
	ht := pdf.PointConvert(fontSize)
	tr := pdf.UnicodeTranslatorFromDescriptor("") // "" defaults to "cp1252"
	write := func(str string) {
		pdf.CellFormat(190, ht, tr(str), "", 1, "C", false, 0, "")
		pdf.Ln(ht)
	}
	pdf.AddPage()
	str := `Gofpdf provides a translator that will convert any UTF-8 code point ` +
		`that is present in the specified code page.`
	pdf.MultiCell(190, ht, str, "", "L", false)
	pdf.Ln(2 * ht)
	write("Voix ambiguë d'un cœur qui au zéphyr préfère les jattes de kiwi.")
	write("Falsches Üben von Xylophonmusik quält jeden größeren Zwerg.")
	write("Heizölrückstoßabdämpfung")
	write("Forårsjævndøgn / Efterårsjævndøgn")

	pdf.SetFont("Helvetica-1251", "", fontSize) // Name matches one specified in AddFont()
	tr = pdf.UnicodeTranslatorFromDescriptor("cp1251")
	write("Съешь же ещё этих мягких французских булок, да выпей чаю.")

	pdf.SetFont("Helvetica-1253", "", fontSize)
	tr = pdf.UnicodeTranslatorFromDescriptor("cp1253")
	write("Θέλει αρετή και τόλμη η ελευθερία. (Ανδρέας Κάλβος)")

	fileStr := exampleFilename("Fpdf_CellFormat_4_codepage")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_CellFormat_4_codepage.pdf
}

// This example demonstrates password protection for documents.
func ExampleFpdf_SetProtection() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetProtection(gofpdf.CnProtectPrint, "123", "abc")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	pdf.Write(10, "Password-protected.")
	fileStr := exampleFilename("Fpdf_SetProtection")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_SetProtection.pdf
}

// This example displays equilateral polygons in a demonstration of the Polygon
// function.
func ExampleFpdf_Polygon() {
	const rowCount = 5
	const colCount = 4
	const ptSize = 36
	var x, y, radius, gap, advance float64
	var rgVal int
	var pts []gofpdf.PointType
	vertices := func(count int) (res []gofpdf.PointType) {
		var pt gofpdf.PointType
		res = make([]gofpdf.PointType, 0, count)
		mlt := 2.0 * math.Pi / float64(count)
		for j := 0; j < count; j++ {
			pt.Y, pt.X = math.Sincos(float64(j) * mlt)
			res = append(res, gofpdf.PointType{
				X: x + radius*pt.X,
				Y: y + radius*pt.Y})
		}
		return
	}
	pdf := gofpdf.New("P", "mm", "A4", "") // A4 210.0 x 297.0
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", ptSize)
	pdf.SetDrawColor(0, 80, 180)
	gap = 12.0
	pdf.SetY(gap)
	pdf.CellFormat(190.0, gap, "Equilateral polygons", "", 1, "C", false, 0, "")
	radius = (210.0 - float64(colCount+1)*gap) / (2.0 * float64(colCount))
	advance = gap + 2.0*radius
	y = 2*gap + pdf.PointConvert(ptSize) + radius
	rgVal = 230
	for row := 0; row < rowCount; row++ {
		pdf.SetFillColor(rgVal, rgVal, 0)
		rgVal -= 12
		x = gap + radius
		for col := 0; col < colCount; col++ {
			pts = vertices(row*colCount + col + 3)
			pdf.Polygon(pts, "FD")
			x += advance
		}
		y += advance
	}
	fileStr := exampleFilename("Fpdf_Polygon")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_Polygon.pdf
}

// This example demonstrates document layers. The initial visibility of a layer
// is specified with the second parameter to AddLayer(). The layer list
// displayed by the document reader allows layer visibility to be controlled
// interactively.
func ExampleFpdf_AddLayer() {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 15)
	pdf.Write(8, "This line doesn't belong to any layer.\n")

	// Define layers
	l1 := pdf.AddLayer("Layer 1", true)
	l2 := pdf.AddLayer("Layer 2", true)

	// Open layer pane in PDF viewer
	pdf.OpenLayerPane()

	// First layer
	pdf.BeginLayer(l1)
	pdf.Write(8, "This line belongs to layer 1.\n")
	pdf.EndLayer()

	// Second layer
	pdf.BeginLayer(l2)
	pdf.Write(8, "This line belongs to layer 2.\n")
	pdf.EndLayer()

	// First layer again
	pdf.BeginLayer(l1)
	pdf.Write(8, "This line belongs to layer 1 again.\n")
	pdf.EndLayer()

	fileStr := exampleFilename("Fpdf_AddLayer")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_AddLayer.pdf
}

// This example demonstrates the use of an image that is retrieved from a web
// server.
func ExampleFpdf_RegisterImageReader() {

	const (
		margin   = 10
		wd       = 210
		ht       = 297
		fontSize = 15
		urlStr   = "https://github.com/jung-kurt/gofpdf/blob/master/image/gofpdf.png?raw=true"
		msgStr   = `Images from the web can be easily embedded when a PDF document is generated.`
	)

	var (
		rsp *http.Response
		err error
		tp  string
	)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", fontSize)
	ln := pdf.PointConvert(fontSize)
	pdf.MultiCell(wd-margin-margin, ln, msgStr, "", "L", false)
	rsp, err = http.Get(urlStr)
	if err == nil {
		tp = pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
		infoPtr := pdf.RegisterImageReader(urlStr, tp, rsp.Body)
		if pdf.Ok() {
			imgWd, imgHt := infoPtr.Extent()
			pdf.Image(urlStr, (wd-imgWd)/2.0, pdf.GetY()+ln,
				imgWd, imgHt, false, tp, 0, "")
		}
	} else {
		pdf.SetError(err)
	}
	fileStr := exampleFilename("Fpdf_RegisterImageReader_url")
	err = pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_RegisterImageReader_url.pdf

}

// This example demonstrates the Beziergon function.
func ExampleFpdf_Beziergon() {

	const (
		margin      = 10
		wd          = 210
		unit        = (wd - 2*margin) / 6
		ht          = 297
		fontSize    = 15
		msgStr      = `Demonstration of Beziergon function`
		coefficient = 0.6
		delta       = coefficient * unit
		ln          = fontSize * 25.4 / 72
		offsetX     = (wd - 4*unit) / 2.0
		offsetY     = offsetX + 2*ln
	)

	srcList := []gofpdf.PointType{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 2, Y: 2},
		{X: 3, Y: 2},
		{X: 3, Y: 3},
		{X: 4, Y: 3},
		{X: 4, Y: 4},
		{X: 1, Y: 4},
		{X: 1, Y: 3},
		{X: 0, Y: 3},
	}

	ctrlList := []gofpdf.PointType{
		{X: 1, Y: -1},
		{X: 1, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: 1},
		{X: -1, Y: 1},
		{X: -1, Y: -1},
		{X: -1, Y: -1},
		{X: -1, Y: -1},
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", fontSize)
	for j, src := range srcList {
		srcList[j].X = offsetX + src.X*unit
		srcList[j].Y = offsetY + src.Y*unit
	}
	for j, ctrl := range ctrlList {
		ctrlList[j].X = ctrl.X * delta
		ctrlList[j].Y = ctrl.Y * delta
	}
	jPrev := len(srcList) - 1
	srcPrev := srcList[jPrev]
	curveList := []gofpdf.PointType{srcPrev} // point [, control 0, control 1, point]*
	control := func(x, y float64) {
		curveList = append(curveList, gofpdf.PointType{X: x, Y: y})
	}
	for j, src := range srcList {
		ctrl := ctrlList[jPrev]
		control(srcPrev.X+ctrl.X, srcPrev.Y+ctrl.Y) // Control 0
		ctrl = ctrlList[j]
		control(src.X-ctrl.X, src.Y-ctrl.Y) // Control 1
		curveList = append(curveList, src)  // Destination
		jPrev = j
		srcPrev = src
	}
	pdf.MultiCell(wd-margin-margin, ln, msgStr, "", "C", false)
	pdf.SetDashPattern([]float64{0.8, 0.8}, 0)
	pdf.SetDrawColor(160, 160, 160)
	pdf.Polygon(srcList, "D")
	pdf.SetDashPattern([]float64{}, 0)
	pdf.SetDrawColor(64, 64, 128)
	pdf.SetLineWidth(pdf.GetLineWidth() * 3)
	pdf.Beziergon(curveList, "D")
	fileStr := exampleFilename("Fpdf_Beziergon")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_Beziergon.pdf

}

// This example demonstrates loading a non-standard font using a generalized
// font loader. fontResourceType implements the FontLoader interface and is
// defined locally in the test source code.
func ExampleFpdf_SetFontLoader() {
	var fr fontResourceType
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFontLoader(fr)
	pdf.AddFont("Calligrapher", "", "calligra.json")
	pdf.AddPage()
	pdf.SetFont("Calligrapher", "", 35)
	pdf.Cell(0, 10, "Load fonts from any source")
	fileStr := exampleFilename("Fpdf_SetFontLoader")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Generalized font loader reading calligra.json
	// Generalized font loader reading calligra.z
	// Successfully generated pdf/Fpdf_SetFontLoader.pdf
}

// This example demonstrates the Path Drawing functions, such as: MoveTo,
// LineTo, CurveTo, ..., ClosePath and DrawPath.
func ExampleFpdf_MoveTo() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.MoveTo(20, 20)
	pdf.LineTo(170, 20)
	pdf.ArcTo(170, 40, 20, 20, 0, 90, 0)
	pdf.CurveTo(190, 100, 105, 100)
	pdf.CurveBezierCubicTo(20, 100, 105, 200, 20, 200)
	pdf.ClosePath()
	pdf.SetFillColor(200, 200, 200)
	pdf.SetLineWidth(3)
	pdf.DrawPath("DF")
	fileStr := exampleFilename("Fpdf_MoveTo_path")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_MoveTo_path.pdf
}

// This example demonstrates various line cap and line join styles.
func ExampleFpdf_SetLineJoinStyle() {
	const offset = 75.0
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	var draw = func(cap, join string, x0, y0, x1, y1 float64) {
		// transform begin & end needed to isolate caps and joins
		pdf.SetLineCapStyle(cap)
		pdf.SetLineJoinStyle(join)

		// Draw thick line
		pdf.SetDrawColor(0x33, 0x33, 0x33)
		pdf.SetLineWidth(30.0)
		pdf.MoveTo(x0, y0)
		pdf.LineTo((x0+x1)/2+offset, (y0+y1)/2)
		pdf.LineTo(x1, y1)
		pdf.DrawPath("D")

		// Draw thin helping line
		pdf.SetDrawColor(0xFF, 0x33, 0x33)
		pdf.SetLineWidth(2.56)
		pdf.MoveTo(x0, y0)
		pdf.LineTo((x0+x1)/2+offset, (y0+y1)/2)
		pdf.LineTo(x1, y1)
		pdf.DrawPath("D")

	}
	x := 35.0
	caps := []string{"butt", "square", "round"}
	joins := []string{"bevel", "miter", "round"}
	for i := range caps {
		draw(caps[i], joins[i], x, 50, x, 160)
		x += offset
	}
	fileStr := exampleFilename("Fpdf_SetLineJoinStyle_caps")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_SetLineJoinStyle_caps.pdf
}

// This example demonstrates various fill modes.
func ExampleFpdf_DrawPath() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetDrawColor(0xff, 0x00, 0x00)
	pdf.SetFillColor(0x99, 0x99, 0x99)
	pdf.SetFont("Helvetica", "", 15)
	pdf.AddPage()
	pdf.SetAlpha(1, "Multiply")
	var (
		polygon = func(cx, cy, r, n, dir float64) {
			da := 2 * math.Pi / n
			pdf.MoveTo(cx+r, cy)
			pdf.Text(cx+r, cy, "0")
			i := 1
			for a := da; a < 2*math.Pi; a += da {
				x, y := cx+r*math.Cos(dir*a), cy+r*math.Sin(dir*a)
				pdf.LineTo(x, y)
				pdf.Text(x, y, strconv.Itoa(i))
				i++
			}
			pdf.ClosePath()
		}
		polygons = func(cx, cy, r, n, dir float64) {
			d := 1.0
			for rf := r; rf > 0; rf -= 10 {
				polygon(cx, cy, rf, n, d)
				d *= dir
			}
		}
		star = func(cx, cy, r, n float64) {
			da := 4 * math.Pi / n
			pdf.MoveTo(cx+r, cy)
			for a := da; a < 4*math.Pi+da; a += da {
				x, y := cx+r*math.Cos(a), cy+r*math.Sin(a)
				pdf.LineTo(x, y)
			}
			pdf.ClosePath()
		}
	)
	// triangle
	polygons(55, 45, 40, 3, 1)
	pdf.DrawPath("B")
	pdf.Text(15, 95, "B (same direction, non zero winding)")

	// square
	polygons(155, 45, 40, 4, 1)
	pdf.DrawPath("B*")
	pdf.Text(115, 95, "B* (same direction, even odd)")

	// pentagon
	polygons(55, 145, 40, 5, -1)
	pdf.DrawPath("B")
	pdf.Text(15, 195, "B (different direction, non zero winding)")

	// hexagon
	polygons(155, 145, 40, 6, -1)
	pdf.DrawPath("B*")
	pdf.Text(115, 195, "B* (different direction, even odd)")

	// star
	star(55, 245, 40, 5)
	pdf.DrawPath("B")
	pdf.Text(15, 290, "B (non zero winding)")

	// star
	star(155, 245, 40, 5)
	pdf.DrawPath("B*")
	pdf.Text(115, 290, "B* (even odd)")

	fileStr := exampleFilename("Fpdf_DrawPath_fill")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_DrawPath_fill.pdf
}
