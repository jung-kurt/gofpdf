/*
 * Copyright (c) 2013 Kurt Jung (Gmail: kurt.w.jung)
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

package gofpdf

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Absolute path needed for gocov tool; relative OK for test
const (
	GOFPDF_DIR = "."
	FONT_DIR   = GOFPDF_DIR + "/font"
	IMG_DIR    = GOFPDF_DIR + "/image"
	TEXT_DIR   = GOFPDF_DIR + "/text"
)

type nullWriter struct {
}

func (nw *nullWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	return
}

func (nw *nullWriter) Close() (err error) {
	return
}

type pdfWriter struct {
	pdf *Fpdf
	fl  *os.File
	idx int
}

func (pw *pdfWriter) Write(p []byte) (n int, err error) {
	if pw.pdf.Ok() {
		return pw.fl.Write(p)
	}
	return
}

func (pw *pdfWriter) Close() (err error) {
	if pw.fl != nil {
		pw.fl.Close()
		pw.fl = nil
	}
	if pw.pdf.Ok() {
		fmt.Printf("Successfully generated pdf/tutorial%02d.pdf\n", pw.idx)
	} else {
		fmt.Printf("%s\n", pw.pdf.Error())
	}
	return
}

func docWriter(pdf *Fpdf, idx int) *pdfWriter {
	pw := new(pdfWriter)
	pw.pdf = pdf
	pw.idx = idx
	if pdf.Ok() {
		var err error
		fileStr := fmt.Sprintf("%s/pdf/tutorial%02d.pdf", GOFPDF_DIR, idx)
		pw.fl, err = os.Create(fileStr)
		if err != nil {
			pdf.SetErrorf("Error opening output file %s", fileStr)
		}
	}
	return pw
}

// Hello, world
func ExampleFpdf_tutorial01() {
	pdf := New("P", "mm", "A4", FONT_DIR)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello World!")
	pdf.OutputAndClose(docWriter(pdf, 1))
	// Output:
	// Successfully generated pdf/tutorial01.pdf
}

// Header, footer and page-breaking
func ExampleFpdf_tutorial02() {
	pdf := New("P", "mm", "A4", FONT_DIR)
	pdf.SetHeaderFunc(func() {
		pdf.Image(IMG_DIR+"/logo.png", 10, 6, 30, 0, false, "", 0, "")
		pdf.SetY(5)
		pdf.SetFont("Arial", "B", 15)
		pdf.Cell(80, 0, "")
		pdf.CellFormat(30, 10, "Title", "1", 0, "C", false, 0, "")
		pdf.Ln(20)
	})
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.SetFont("Times", "", 12)
	for j := 1; j <= 40; j++ {
		pdf.CellFormat(0, 10, fmt.Sprintf("Printing line number %d", j), "", 1, "", false, 0, "")
	}
	pdf.OutputAndClose(docWriter(pdf, 2))
	// Output:
	// Successfully generated pdf/tutorial02.pdf
}

// Word-wrapping, line justification and page-breaking
func ExampleFpdf_tutorial03() {
	pdf := New("P", "mm", "A4", FONT_DIR)
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
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	chapterTitle := func(chapNum int, titleStr string) {
		// 	// Arial 12
		pdf.SetFont("Arial", "", 12)
		// Background color
		pdf.SetFillColor(200, 220, 255)
		// Title
		pdf.CellFormat(0, 6, fmt.Sprintf("Chapter %d : %s", chapNum, titleStr), "", 1, "L", true, 0, "")
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
	printChapter(1, "A RUNAWAY REEF", TEXT_DIR+"/20k_c1.txt")
	printChapter(2, "THE PROS AND CONS", TEXT_DIR+"/20k_c2.txt")
	pdf.OutputAndClose(docWriter(pdf, 3))
	// Output:
	// Successfully generated pdf/tutorial03.pdf
}

// Multiple column layout
func ExampleFpdf_tutorial04() {
	var y0 float64
	var crrntCol int
	pdf := New("P", "mm", "A4", FONT_DIR)
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
		pdf.CellFormat(0, 6, fmt.Sprintf("Chapter %d : %s", chapNum, titleStr), "", 1, "L", true, 0, "")
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
		} else {
			// Go back to first column
			setCol(0)
			// Page break
			return true
		}
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
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	printChapter(1, "A RUNAWAY REEF", TEXT_DIR+"/20k_c1.txt")
	printChapter(2, "THE PROS AND CONS", TEXT_DIR+"/20k_c2.txt")
	pdf.OutputAndClose(docWriter(pdf, 4))
	// Output:
	// Successfully generated pdf/tutorial04.pdf
}

// Various table styles
func ExampleFpdf_tutorial05() {
	pdf := New("P", "mm", "A4", FONT_DIR)
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
					err = fmt.Errorf("Error tokenizing %s", lineStr)
				}
			}
			fl.Close()
			if len(countryList) == 0 {
				err = fmt.Errorf("Error loading data from %s", fileStr)
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
			pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3), "LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3), "LR", 0, "R", false, 0, "")
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
			pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3), "LR", 0, "R", fill, 0, "")
			pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3), "LR", 0, "R", fill, 0, "")
			pdf.Ln(-1)
			fill = !fill
		}
		pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	}
	loadData(TEXT_DIR + "/countries.txt")
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()
	basicTable()
	pdf.AddPage()
	improvedTable()
	pdf.AddPage()
	fancyTable()
	pdf.OutputAndClose(docWriter(pdf, 5))
	// Output:
	// Successfully generated pdf/tutorial05.pdf
}

// Internal and external links
func ExampleFpdf_tutorial06() {
	var boldLvl, italicLvl, underscoreLvl int
	var hrefStr string
	pdf := New("P", "mm", "A4", FONT_DIR)
	setStyle := func(boldAdj, italicAdj, underscoreAdj int) {
		styleStr := ""
		boldLvl += boldAdj
		if boldLvl > 0 {
			styleStr += "B"
		}
		italicLvl += italicAdj
		if italicLvl > 0 {
			styleStr += "I"
		}
		underscoreLvl += underscoreAdj
		if underscoreLvl > 0 {
			styleStr += "U"
		}
		pdf.SetFont("", styleStr, 0)
	}
	putLink := func(urlStr, txtStr string) {
		// Put a hyperlink
		pdf.SetTextColor(0, 0, 255)
		setStyle(0, 0, 1)
		pdf.WriteLinkString(5, txtStr, urlStr)
		setStyle(0, 0, -1)
		pdf.SetTextColor(0, 0, 0)
	}

	writeHtml := func(htmlStr string) {
		list := htmlTokenize(htmlStr)
		var ok bool
		for _, el := range list {
			switch el.cat {
			case 'T':
				if len(hrefStr) > 0 {
					putLink(hrefStr, el.str)
					hrefStr = ""
				} else {
					pdf.Write(5, el.str)
				}
			case 'O':
				switch el.str {
				case "b":
					setStyle(1, 0, 0)
				case "i":
					setStyle(0, 1, 0)
				case "u":
					setStyle(0, 0, 1)
				case "br":
					pdf.Ln(5)
				case "a":
					hrefStr, ok = el.attr["href"]
					if !ok {
						hrefStr = ""
					}
				}
			case 'C':
				switch el.str {
				case "b":
					setStyle(-1, 0, 0)
				case "i":
					setStyle(0, -1, 0)
				case "u":
					setStyle(0, 0, -1)

				}
			}
		}
	}
	// First page
	pdf.AddPage()
	pdf.SetFont("Arial", "", 20)
	pdf.Write(5, "To find out what's new in this tutorial, click ")
	pdf.SetFont("", "U", 0)
	link := pdf.AddLink()
	pdf.WriteLinkId(5, "here", link)
	pdf.SetFont("", "", 0)
	// Second page
	pdf.AddPage()
	pdf.SetLink(link, 0, -1)
	pdf.Image(IMG_DIR+"/logo.png", 10, 12, 30, 0, false, "", 0, "http://www.fpdf.org")
	pdf.SetLeftMargin(45)
	pdf.SetFontSize(14)
	htmlStr := `You can now easily print text mixing different styles: <b>bold</b>, ` +
		`<i>italic</i>, <u>underlined</u>, or <b><i><u>all at once</u></i></b>!<br><br>` +
		`You can also insert links on text, such as ` +
		`<a href="http://www.fpdf.org">www.fpdf.org</a>, or on an image: click on the logo.`
	writeHtml(htmlStr)
	pdf.OutputAndClose(docWriter(pdf, 6))
	// Output:
	// Successfully generated pdf/tutorial06.pdf
}

// Non-standard font
func ExampleFpdf_tutorial07() {
	pdf := New("P", "mm", "A4", FONT_DIR)
	pdf.AddFont("Calligrapher", "", "calligra.json")
	pdf.AddPage()
	pdf.SetFont("Calligrapher", "", 35)
	pdf.Cell(0, 10, "Enjoy new fonts with FPDF!")
	pdf.OutputAndClose(docWriter(pdf, 7))
	// Output:
	// Successfully generated pdf/tutorial07.pdf
}

// Various image types
func ExampleFpdf_tutorial08() {
	pdf := New("P", "mm", "A4", FONT_DIR)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 11)
	pdf.Image(IMG_DIR+"/logo.png", 10, 10, 30, 0, false, "", 0, "")
	pdf.Text(50, 20, "logo.png")
	pdf.Image(IMG_DIR+"/logo.gif", 10, 40, 30, 0, false, "", 0, "")
	pdf.Text(50, 50, "logo.gif")
	pdf.Image(IMG_DIR+"/logo-gray.png", 10, 70, 30, 0, false, "", 0, "")
	pdf.Text(50, 80, "logo-gray.png")
	pdf.Image(IMG_DIR+"/logo-rgb.png", 10, 100, 30, 0, false, "", 0, "")
	pdf.Text(50, 110, "logo-rgb.png")
	pdf.Image(IMG_DIR+"/logo.jpg", 10, 130, 30, 0, false, "", 0, "")
	pdf.Text(50, 140, "logo.jpg")
	pdf.OutputAndClose(docWriter(pdf, 8))
	// Output:
	// Successfully generated pdf/tutorial08.pdf
}

// Landscape mode with logos
func ExampleFpdf_tutorial09() {
	var y0 float64
	var crrntCol int
	loremStr := "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod " +
		"tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis " +
		"nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis " +
		"aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat " +
		"nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui " +
		"officia deserunt mollit anim id est laborum."
	pdf := New("L", "mm", "A4", FONT_DIR)
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
		} else {
			setCol(0)
			return true
		}
	})
	pdf.AddPage()
	pdf.SetFont("Times", "", 12)
	for j := 0; j < 20; j++ {
		if j == 1 {
			pdf.Image(IMG_DIR+"/fpdf.png", -1, 0, colWd, 0, true, "", 0, "")
		} else if j == 5 {
			pdf.Image(IMG_DIR+"/golang-gopher.png", -1, 0, colWd, 0, true, "", 0, "")
		}
		pdf.MultiCell(colWd, 5, loremStr, "", "", false)
		pdf.Ln(-1)
	}
	pdf.OutputAndClose(docWriter(pdf, 9))
	// Output:
	// Successfully generated pdf/tutorial09.pdf
}

// Test the corner cases as reported by the gocov tool
func ExampleFpdf_tutorial10() {
	MakeFont(FONT_DIR+"/calligra.ttf", FONT_DIR+"/cp1252.map", FONT_DIR, nil, true)
	pdf := New("", "", "", "")
	pdf.SetFontLocation(FONT_DIR)
	pdf.SetTitle("世界", true)
	pdf.SetAuthor("世界", true)
	pdf.SetSubject("世界", true)
	pdf.SetCreator("世界", true)
	pdf.SetKeywords("世界", true)
	pdf.AddFont("Calligrapher", "", "calligra.json")
	pdf.AddPage()
	pdf.SetFont("Calligrapher", "", 16)
	pdf.Writef(5, "\x95 %s \x95", pdf)
	pdf.OutputAndClose(docWriter(pdf, 10))
	// Output:
	// Successfully generated pdf/tutorial10.pdf
}

// Geometric figures
func ExampleFpdf_tutorial11() {
	const (
		thin  = 0.2
		thick = 3.0
	)
	pdf := New("", "", "", FONT_DIR)
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
	pdf.CurveCubic(10, y+30, 15, y-20, 40, y+30, 10, y+30, "D")
	pdf.CurveCubic(45, y+30, 50, y-20, 75, y+30, 45, y+30, "F")
	pdf.CurveCubic(80, y+30, 85, y-20, 110, y+30, 80, y+30, "FD")
	pdf.SetLineWidth(thick)
	pdf.CurveCubic(115, y+30, 120, y-20, 145, y+30, 115, y+30, "FD")
	pdf.SetLineCapStyle("round")
	pdf.CurveCubic(150, y+30, 155, y-20, 180, y+30, 150, y+30, "FD")
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

	pdf.OutputAndClose(docWriter(pdf, 11))
	// Output:
	// Successfully generated pdf/tutorial11.pdf
}

// Transparency
func ExampleFpdf_tutorial12() {
	pdf := New("", "", "", FONT_DIR)
	pdf.SetFont("Helvetica", "", 48)
	pdf.SetLineWidth(4)
	pdf.SetFillColor(180, 180, 180)
	pdf.AddPage()
	pdf.SetXY(55, 60)
	pdf.CellFormat(100, 40, "Go", "1", 0, "C", true, 0, "")
	pdf.SetAlpha(0.5, "Normal")
	pdf.Image(IMG_DIR+"/golang-gopher.png", 30, 10, 150, 0, false, "", 0, "")
	pdf.SetAlpha(1.0, "Normal")
	pdf.OutputAndClose(docWriter(pdf, 12))
	// Output:
	// Successfully generated pdf/tutorial12.pdf
}

// Gradients
func ExampleFpdf_tutorial13() {
	pdf := New("", "", "", FONT_DIR)
	pdf.SetFont("Helvetica", "", 12)
	pdf.AddPage()
	pdf.LinearGradient(0, 0, 210, 100, 250, 250, 255, 220, 220, 225, 0, 0, 0, .5)
	pdf.LinearGradient(20, 25, 75, 75, 220, 220, 250, 80, 80, 220, 0, .2, 0, .8)
	pdf.Rect(20, 25, 75, 75, "D")
	pdf.LinearGradient(115, 25, 75, 75, 220, 220, 250, 80, 80, 220, 0, 0, 1, 1)
	pdf.Rect(115, 25, 75, 75, "D")
	pdf.RadialGradient(20, 120, 75, 75, 220, 220, 250, 80, 80, 220, 0.25, 0.75, 0.25, 0.75, 1)
	pdf.Rect(20, 120, 75, 75, "D")
	pdf.RadialGradient(115, 120, 75, 75, 220, 220, 250, 80, 80, 220, 0.25, 0.75, 0.75, 0.75, 0.75)
	pdf.Rect(115, 120, 75, 75, "D")
	pdf.OutputAndClose(docWriter(pdf, 13))
	// Output:
	// Successfully generated pdf/tutorial13.pdf
}
