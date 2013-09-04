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

package gofpdf_test

import (
	"code.google.com/p/gofpdf"
	"fmt"
	// "testing"
)

func ExampleTtfParse() {
	ttf, err := gofpdf.TtfParse(FONT_DIR + "/calligra.ttf")
	if err == nil {
		fmt.Printf("Postscript name:  %s\n", ttf.PostScriptName)
		fmt.Printf("unitsPerEm:       %8d\n", ttf.UnitsPerEm)
		fmt.Printf("Xmin:             %8d\n", ttf.Xmin)
		fmt.Printf("Ymin:             %8d\n", ttf.Ymin)
		fmt.Printf("Xmax:             %8d\n", ttf.Xmax)
		fmt.Printf("Ymax:             %8d\n", ttf.Ymax)
	} else {
		fmt.Printf("%s\n", err)
	}
	// Output:
	// Postscript name:  CalligrapherRegular
	// unitsPerEm:           1000
	// Xmin:                 -173
	// Ymin:                 -234
	// Xmax:                 1328
	// Ymax:                  899
}

// func TestLoadMap(t *testing.T) {
// 	expectList := []string{
// 		"164: 0x0E04 khokhwaithai",
// 		"165: 0x0E05 khokhonthai",
// 		"166: 0x0E06 khorakhangthai",
// 		"167: 0x0E07 ngonguthai",
// 		"168: 0x0E08 chochanthai",
// 		"169: 0x0E09 chochingthai",
// 	}
// 	list, err := loadMap(FONT_DIR + "/iso-8859-11.map")
// 	if err == nil {
// 		pos := 0
// 		for j := 164; j < 170; j++ {
// 			enc := list[j]
// 			str := fmt.Sprintf("%3d: 0x%04X %s", j, enc.uv, enc.name)
// 			// fmt.Printf("Expect [%s], Got [%s]\n", expectList[pos], str)
// 			if expectList[pos] != str {
// 				t.Fatalf("Unexpected output from loadMap")
// 			}
// 			pos++
// 		}
// 	}
// }

func ExampleFpdf_GetStringWidth() {
	pdf := gofpdf.New("", "", "", FONT_DIR)
	pdf.SetFont("Helvetica", "", 12)
	pdf.AddPage()
	for _, s := range []string{"Hello", "世界"} {
		fmt.Printf("Width of \"%s\" is %.2f\n", s, pdf.GetStringWidth(s))
		if pdf.Err() {
			fmt.Println(pdf.Error())
		}
	}
	pdf.Close()
	// Output:
	// Width of "Hello" is 9.64
	// Width of "世界" is 0.00
	// Unicode strings not supported
}
