package gofpdi

import (
	realgofpdi "github.com/phpdave11/gofpdi"
)

// Create new gofpdi instance
var fpdi = realgofpdi.NewImporter()

// gofpdiPdf is a partial interface that only implements the functions we need
// from the PDF generator to put the imported PDF templates on the PDF.
type gofpdiPdf interface {
	ImportObjects(objs map[string][]byte)
	ImportObjPos(objs map[string]map[int]string)
	ImportTemplates(tpls map[string]string)
	UseImportedTemplate(tplName string, x float64, y float64, w float64, h float64)
	SetError(err error)
}

// ImportPage imports a page of a PDF file with the specified box (/MediaBox,
// /TrimBox, /ArtBox, /CropBox, or /BleedBox). Returns a template id that can
// be used with UseImportedTemplate to draw the template onto the page.
func ImportPage(f gofpdiPdf, sourceFile string, pageno int, box string) int {
	// Set source file for fpdi
	fpdi.SetSourceFile(sourceFile)

	// Import page
	tpl := fpdi.ImportPage(pageno, box)

	// Import objects into current pdf document
	// Unordered means that the objects will be returned with a sha1 hash instead of an integer
	// The objects themselves may have references to other hashes which will be replaced in ImportObjects()
	tplObjIDs := fpdi.PutFormXobjectsUnordered()

	// Set template names and ids (hashes) in gofpdf
	f.ImportTemplates(tplObjIDs)

	// Get a map[string]string of the imported objects.
	// The map keys will be the ID of each object.
	imported := fpdi.GetImportedObjectsUnordered()

	// Import gofpdi objects into gofpdf
	f.ImportObjects(imported)

	// Get a map[string]map[int]string of the object hashes and their positions within each object,
	// to be replaced with object ids (integers).
	importedObjPos := fpdi.GetImportedObjHashPos()

	// Import gofpdi object hashes and their positions into gopdf
	f.ImportObjPos(importedObjPos)

	return tpl
}

// UseImportedTemplate draws the template onto the page at x,y. If w is 0, the
// template will be scaled to fit based on h. If h is 0, the template will be
// scaled to fit based on w.
func UseImportedTemplate(f gofpdiPdf, tplid int, x float64, y float64, w float64, h float64) {
	// Get values from fpdi
	tplName, scaleX, scaleY, tX, tY := fpdi.UseTemplate(tplid, x, y, w, h)

	f.UseImportedTemplate(tplName, scaleX, scaleY, tX, tY)
}
