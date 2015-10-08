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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var gofpdfDir string

func init() {
	setRoot()
}

// Assign the relative path to the gofpdfDir directory based on current working
// directory
func setRoot() {
	wdStr, err := os.Getwd()
	if err == nil {
		gofpdfDir = ""
		sepStr := string(os.PathSeparator)
		list := strings.Split(wdStr, sepStr)
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

var (
	// 00000230  44 46 20 31 2e 37 29 0a  2f 43 72 65 61 74 69 6f  |DF 1.7)./Creatio|
	// 00000240  6e 44 61 74 65 20 28 44  3a 32 30 31 35 31 30 30  |nDate (D:2015100|
	// 00000250  38 31 32 33 30 34 35 29  0a 3e 3e 0a 65 6e 64 6f  |8123045).>>.endo|
	creationDateRe *regexp.Regexp = regexp.MustCompile("/CreationDate \\(D:\\d{14}\\)")
	fixDate        []byte         = []byte("/CreationDate (D:20000101000000)")
)

// referenceCompare compares the specified file with the file's reference copy
// located in the 'reference' subdirectory. All bytes of the two files are
// compared except for the value of the /CreationDate field in the PDF. An
// error is returned if the two files do not match. If the reference file does
// not exist, a copy of the specified file is made and a non-nil error is
// returned only if this copy fails.
func referenceCompare(fileStr string) (err error) {
	var fileBuf, refFileBuf []byte
	var refFileStr, refDirStr, dirStr, baseFileStr string
	dirStr, baseFileStr = filepath.Split(fileStr)
	refDirStr = filepath.Join(dirStr, "reference")
	err = os.MkdirAll(refDirStr, 0755)
	if err == nil {
		refFileStr = filepath.Join(refDirStr, baseFileStr)
		fileBuf, err = ioutil.ReadFile(fileStr)
		if err == nil {
			// Replace the creation timestamp of this PDF with a fixed value
			fileBuf = creationDateRe.ReplaceAll(fileBuf, fixDate)
			refFileBuf, err = ioutil.ReadFile(refFileStr)
			if err == nil {
				if len(fileBuf) == len(refFileBuf) {
					if bytes.Equal(fileBuf, refFileBuf) {
						// Files match
					} else {
						err = fmt.Errorf("%s differs from %s", fileStr, refFileStr)
					}
				} else {
					err = fmt.Errorf("size of %s (%d) does not match size of %s (%d)",
						fileStr, len(fileBuf), refFileStr, len(refFileBuf))
				}
			} else {
				// Reference file is missing. Create it with a copy of the newly produced
				// file in which the creation date has been fixed. Overwrite error with copy
				// error.
				err = ioutil.WriteFile(refFileStr, fileBuf, 0644)
			}
		}
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
