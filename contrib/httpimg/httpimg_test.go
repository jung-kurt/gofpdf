package httpimg_test

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/httpimg"
	"os"
	"path/filepath"
)

const (
	cnGofpdfDir  = "./"
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

func exampleFilename(baseStr string) string {
	return filepath.Join(cnExampleDir, baseStr+".pdf")
}

func summary(err error, fileStr string) {
	if err == nil {
		fileStr = filepath.ToSlash(fileStr)
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		fmt.Println(err)
	}
}

func ExampleRegister() {
	pdf := gofpdf.New("", "", "", "")
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetFillColor(200, 200, 220)
	pdf.AddPage()

	url := "https://github.com/jung-kurt/gofpdf/raw/master/image/logo_gofpdf.jpg?raw=true"
	httpimg.Register(pdf, url, "")
	pdf.Image(url, 100, 100, 20, 20, false, "", 0, "")

	fileStr := exampleFilename("contrib_httpimg_Register")
	err := pdf.OutputFileAndClose(fileStr)
	summary(err, fileStr)
	// Output:
	// Successfully generated pdf/contrib_httpimg_Register.pdf
}
