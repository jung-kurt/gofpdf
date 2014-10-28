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

package gofpdf

import (
	"bytes"
)

// Version of FPDF from which this package is derived
const (
	cnFpdfVersion = "1.7"
)

type blendModeType struct {
	strokeStr, fillStr, modeStr string
	objNum                      int
}

type gradientType struct {
	tp                int // 2: linear, 3: radial
	clr1Str, clr2Str  string
	x1, y1, x2, y2, r float64
	objNum            int
}

// SizeType fields Wd and Ht specify the horizontal and vertical extents of a
// document element such as a page.
type SizeType struct {
	Wd, Ht float64
}

// PointType fields X and Y specify the horizontal and vertical coordinates of
// a point, typically used in drawing.
type PointType struct {
	X, Y float64
}

// XY returns the X and Y components of the receiver point.
func (p PointType) XY() (float64, float64) {
	return p.X, p.Y
}

// ImageInfoType contains size, color and other information about an image
type ImageInfoType struct {
	data  []byte
	smask []byte
	i     int
	n     int
	w     float64
	h     float64
	cs    string
	pal   []byte
	bpc   int
	f     string
	dp    string
	trns  []int
	scale float64 // document scaling factor
}

// PointConvert returns the value of pt, expressed in points (1/72 inch), as a
// value expressed in the unit of measure specified in New(). Since font
// management in Fpdf uses points, this method can help with line height
// calculations and other methods that require user units.
func (f *Fpdf) PointConvert(pt float64) float64 {
	return pt / f.k
}

// Extent returns the width and height of the image in the units of the Fpdf
// object.
func (info *ImageInfoType) Extent() (wd, ht float64) {
	return info.w / info.scale, info.h / info.scale
}

// Width returns the width of the image in the units of the Fpdf object.
func (info *ImageInfoType) Width() float64 {
	return info.w / info.scale
}

// Height returns the height of the image in the units of the Fpdf object.
func (info *ImageInfoType) Height() float64 {
	return info.h / info.scale
}

type fontFileType struct {
	length1, length2 int64
	n                int
}

type linkType struct {
	x, y, wd, ht float64
	link         int    // Auto-generated internal link ID or...
	linkStr      string // ...application-provided external link string
}

type intLinkType struct {
	page int
	y    float64
}

// outlineType is used for a sidebar outline of bookmarks
type outlineType struct {
	text                                   string
	level, parent, first, last, next, prev int
	y                                      float64
	p                                      int
}

// InitType is used with NewCustom() to customize an Fpdf instance.
// OrientationStr, UnitStr, SizeStr and FontDirStr correspond to the arguments
// accepted by New(). If the Wd and Ht fields of Size are each greater than
// zero, Size will be used to set the default page size rather than SizeStr. Wd
// and Ht are specified in the units of measure indicated by UnitStr.
type InitType struct {
	OrientationStr string
	UnitStr        string
	SizeStr        string
	Size           SizeType
	FontDirStr     string
}

// Fpdf is the principal structure for creating a single PDF document
type Fpdf struct {
	page             int                       // current page number
	n                int                       // current object number
	offsets          []int                     // array of object offsets
	buffer           fmtBuffer                 // buffer holding in-memory PDF
	pages            []*bytes.Buffer           // slice[page] of page content; 1-based
	state            int                       // current document state
	compress         bool                      // compression flag
	k                float64                   // scale factor (number of points in user unit)
	defOrientation   string                    // default orientation
	curOrientation   string                    // current orientation
	stdPageSizes     map[string]SizeType       // standard page sizes
	defPageSize      SizeType                  // default page size
	curPageSize      SizeType                  // current page size
	pageSizes        map[int]SizeType          // used for pages with non default sizes or orientations
	unitStr          string                    // unit of measure for all rendered objects except fonts
	wPt, hPt         float64                   // dimensions of current page in points
	w, h             float64                   // dimensions of current page in user unit
	lMargin          float64                   // left margin
	tMargin          float64                   // top margin
	rMargin          float64                   // right margin
	bMargin          float64                   // page break margin
	cMargin          float64                   // cell margin
	x, y             float64                   // current position in user unit
	lasth            float64                   // height of last printed cell
	lineWidth        float64                   // line width in user unit
	fontpath         string                    // path containing fonts
	coreFonts        map[string]bool           // array of core font names
	fonts            map[string]fontDefType    // array of used fonts
	fontFiles        map[string]fontFileType   // array of font files
	diffs            []string                  // array of encoding differences
	fontFamily       string                    // current font family
	fontStyle        string                    // current font style
	underline        bool                      // underlining flag
	currentFont      fontDefType               // current font info
	fontSizePt       float64                   // current font size in points
	fontSize         float64                   // current font size in user unit
	ws               float64                   // word spacing
	images           map[string]*ImageInfoType // array of used images
	pageLinks        [][]linkType              // pageLinks[page][link], both 1-based
	links            []intLinkType             // array of internal links
	outlines         []outlineType             // array of outlines
	outlineRoot      int                       // root of outlines
	autoPageBreak    bool                      // automatic page breaking
	acceptPageBreak  func() bool               // returns true to accept page break
	pageBreakTrigger float64                   // threshold used to trigger page breaks
	inHeader         bool                      // flag set when processing header
	headerFnc        func()                    // function provided by app and called to write header
	inFooter         bool                      // flag set when processing footer
	footerFnc        func()                    // function provided by app and called to write footer
	zoomMode         string                    // zoom display mode
	layoutMode       string                    // layout display mode
	title            string                    // title
	subject          string                    // subject
	author           string                    // author
	keywords         string                    // keywords
	creator          string                    // creator
	aliasNbPagesStr  string                    // alias for total number of pages
	pdfVersion       string                    // PDF version number
	fontDirStr       string                    // location of font definition files
	capStyle         int                       // line cap style: butt 0, round 1, square 2
	joinStyle        int                       // line segment join style: miter 0, round 1, bevel 2
	blendList        []blendModeType           // slice[idx] of alpha transparency modes, 1-based
	blendMap         map[string]int            // map into blendList
	gradientList     []gradientType            // slice[idx] of gradient records
	clipNest         int                       // Number of active clipping contexts
	transformNest    int                       // Number of active transformation contexts
	err              error                     // Set if error occurs during life cycle of instance
	protect          protectType               // document protection structure
	layer            layerRecType              // manages optional layers in document
	colorFlag        bool                      // indicates whether fill and text colors are different
	color            struct {                  // Composite values of colors
		draw, fill, text clrType
	}
}

type encType struct {
	uv   int
	name string
}

type encListType [256]encType

type fontBoxType struct {
	Xmin, Ymin, Xmax, Ymax int
}

type fontDescType struct {
	Ascent       int
	Descent      int
	CapHeight    int
	Flags        int
	FontBBox     fontBoxType
	ItalicAngle  int
	StemV        int
	MissingWidth int
}

type fontDefType struct {
	Tp           string       // "Core", "TrueType", ...
	Name         string       // "Courier-Bold", ...
	Desc         fontDescType // Font descriptor
	Up           int          // Underline position
	Ut           int          // Underline thickness
	Cw           [256]int     // Character width by ordinal
	Enc          string       // "cp1252", ...
	Diff         string       // Differences from reference encoding
	File         string       // "Redressed.z"
	Size1, Size2 int          // Type1 values
	OriginalSize int          // Size of uncompressed font file
	I            int          // 1-based position in font list, set by font loader, not this program
	N            int          // Set by font loader
	DiffN        int          // Position of diff in app array, set by font loader
}

type fontInfoType struct {
	Data               []byte
	File               string
	OriginalSize       int
	FontName           string
	Bold               bool
	IsFixedPitch       bool
	UnderlineThickness int
	UnderlinePosition  int
	Widths             [256]int
	Size1, Size2       uint32
	Desc               fontDescType
}
