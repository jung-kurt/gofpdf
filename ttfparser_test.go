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
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	// "testing"
)

func ExampleTtfParse() {
	ttf, err := gofpdf.TtfParse(cnFontDir + "/calligra.ttf")
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

func hexStr(s string) string {
	var b bytes.Buffer
	b.WriteString("\"")
	for _, ch := range []byte(s) {
		b.WriteString(fmt.Sprintf("\\x%02x", ch))
	}
	b.WriteString("\":")
	return b.String()
}

func ExampleFpdf_GetStringWidth() {
	pdf := gofpdf.New("", "", "", cnFontDir)
	pdf.SetFont("Helvetica", "", 12)
	pdf.AddPage()
	for _, s := range []string{"Hello", "世界", "\xe7a va?"} {
		fmt.Printf("%-32s width %5.2f, bytes %2d, runes %2d\n",
			hexStr(s), pdf.GetStringWidth(s), len(s), len([]rune(s)))
		if pdf.Err() {
			fmt.Println(pdf.Error())
		}
	}
	pdf.Close()
	// Output:
	// "\x48\x65\x6c\x6c\x6f":          width  9.64, bytes  5, runes  5
	// "\xe4\xb8\x96\xe7\x95\x8c":      width 13.95, bytes  6, runes  2
	// "\xe7\x61\x20\x76\x61\x3f":      width 12.47, bytes  6, runes  6
}
