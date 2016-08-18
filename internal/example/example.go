// Copyright (c) 2015 Kurt Jung (Gmail: kurt.w.jung)
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

// Package example provides some helper routines for the test packages of
// gofpdf and its various contributed packages located beneath the contrib
// directory.
package example

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

var gofpdfDir string

func init() {
	setRoot()
	gofpdf.SetDefaultCompression(false)
	gofpdf.SetDefaultCatalogSort(true)
	gofpdf.SetDefaultCreationDate(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
}

// Assign the relative path to the gofpdfDir directory based on current working
// directory
func setRoot() {
	wdStr, err := os.Getwd()
	if err == nil {
		gofpdfDir = ""
		list := strings.Split(filepath.ToSlash(wdStr), "/")
		for j := len(list) - 1; j >= 0 && list[j] != "gofpdf"; j-- {
			gofpdfDir = filepath.Join(gofpdfDir, "..")
		}
	} else {
		panic(err)
	}
}

// ImageFile returns a qualified filename in which the path to the image
// directory is prepended to the specified filename.
func ImageFile(fileStr string) string {
	return filepath.Join(gofpdfDir, "image", fileStr)
}

// FontDir returns the path to the font directory.
func FontDir() string {
	return filepath.Join(gofpdfDir, "font")
}

// FontFile returns a qualified filename in which the path to the font
// directory is prepended to the specified filename.
func FontFile(fileStr string) string {
	return filepath.Join(FontDir(), fileStr)
}

// TextFile returns a qualified filename in which the path to the text
// directory is prepended to the specified filename.
func TextFile(fileStr string) string {
	return filepath.Join(gofpdfDir, "text", fileStr)
}

// PdfDir returns the path to the PDF output directory.
func PdfDir() string {
	return filepath.Join(gofpdfDir, "pdf")
}

// PdfFile returns a qualified filename in which the path to the PDF output
// directory is prepended to the specified filename.
func PdfFile(fileStr string) string {
	return filepath.Join(PdfDir(), fileStr)
}

// Filename returns a qualified filename in which the example PDF directory
// path is prepended and the suffix ".pdf" is appended to the specified
// filename.
func Filename(baseStr string) string {
	return PdfFile(baseStr + ".pdf")
}

// referenceCompare compares the specified file with the file's reference copy
// located in the 'reference' subdirectory. All bytes of the two files are
// compared except for the value of the /CreationDate field in the PDF. This
// function succeeds if both files are equivalent except for their
// /CreationDate values or if the reference file does not exist.
func referenceCompare(fileStr string) (err error) {
	var refFileStr, refDirStr, dirStr, baseFileStr string
	dirStr, baseFileStr = filepath.Split(fileStr)
	refDirStr = filepath.Join(dirStr, "reference")
	err = os.MkdirAll(refDirStr, 0755)
	if err == nil {
		refFileStr = filepath.Join(refDirStr, baseFileStr)
		err = gofpdf.ComparePDFFiles(fileStr, refFileStr)
	}
	return
}

// Summary generates a predictable report for use by test examples. If the
// specified error is nil, the filename delimiters are normalized and the
// filename printed to standard output with a success message. If the specified
// error is not nil, its String() value is printed to standard output.
func Summary(err error, fileStr string) {
	if err == nil {
		err = referenceCompare(fileStr)
	}
	if err == nil {
		fileStr = filepath.ToSlash(fileStr)
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		fmt.Println(err)
	}
}
