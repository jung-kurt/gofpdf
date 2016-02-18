package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jung-kurt/gofpdf"
)

func report(fileStr string, err error) {
	if err == nil {
		var info os.FileInfo
		info, err = os.Stat(fileStr)
		if err == nil {
			fmt.Printf("%s: OK, size %d\n", fileStr, info.Size())
		} else {
			fmt.Printf("%s: bad stat\n", fileStr)
		}
	} else {
		fmt.Printf("%s: %s\n", fileStr, err)
	}
}

func newPdf() (pdf *gofpdf.Fpdf) {
	pdf = gofpdf.New("P", "mm", "A4", "../../font")
	pdf.SetCompression(false)
	pdf.AddFont("Calligrapher", "", "calligra.json")
	pdf.AddPage()
	pdf.SetFont("Calligrapher", "", 35)
	pdf.Cell(0, 10, "Enjoy new fonts with FPDF!")
	return
}

func full() {
	const name = "full.pdf"
	report(name, newPdf().OutputFileAndClose(name))
}

func min() {
	const name = "min.pdf"
	cmd := exec.Command("gs", "-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS=/screen", "-dNOPAUSE", "-dQUIET",
		"-dBATCH", "-sOutputFile="+name, "-")
	inPipe, err := cmd.StdinPipe()
	if err == nil {
		errChan := make(chan error, 1)
		go func() {
			errChan <- cmd.Start()
		}()
		newPdf().Output(inPipe)
		report(name, <-errChan)
	} else {
		report(name, err)
	}
}

func main() {
	full()
	min()
}
