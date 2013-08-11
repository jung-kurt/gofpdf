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
	"bytes"
)

const (
	FPDF_VERSION = "1.7"
)

type sizeType struct {
	wd, ht float64
}

type imageInfoType struct {
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
}

type fontFileType struct {
	length1, length2 int64
	n                int
}

type linkType struct {
	x, y, wd, ht float64
	link         int    // Auto-generated link ID or...
	linkStr      string // ...application-provided link string
}

type intLinkType struct {
	page int
	y    float64
}

type Fpdf struct {
	page             int                      // current page number
	n                int                      // current object number
	offsets          []int                    // array of object offsets
	buffer           fmtBuffer                // buffer holding in-memory PDF
	pages            []*bytes.Buffer          // slice[page] of page content; 1-based
	state            int                      // current document state
	compress         bool                     // compression flag
	k                float64                  // scale factor (number of points in user unit)
	defOrientation   string                   // default orientation
	curOrientation   string                   // current orientation
	stdpageSizes     map[string]sizeType      // standard page sizes
	defPageSize      sizeType                 // default page size
	curPageSize      sizeType                 // current page size
	pageSizes        map[int]sizeType         // used for pages with non default sizes or orientations
	wPt, hPt         float64                  // dimensions of current page in points
	w, h             float64                  // dimensions of current page in user unit
	lMargin          float64                  // left margin
	tMargin          float64                  // top margin
	rMargin          float64                  // right margin
	bMargin          float64                  // page break margin
	cMargin          float64                  // cell margin
	x, y             float64                  // current position in user unit
	lasth            float64                  // height of last printed cell
	lineWidth        float64                  // line width in user unit
	fontpath         string                   // path containing fonts
	coreFonts        map[string]bool          // array of core font names
	fonts            map[string]fontDefType   // array of used fonts
	fontFiles        map[string]fontFileType  // array of font files
	diffs            []string                 // array of encoding differences
	fontFamily       string                   // current font family
	fontStyle        string                   // current font style
	underline        bool                     // underlining flag
	currentFont      fontDefType              // current font info
	fontSizePt       float64                  // current font size in points
	fontSize         float64                  // current font size in user unit
	drawColor        string                   // commands for drawing color
	fillColor        string                   // commands for filling color
	textColor        string                   // commands for text color
	colorFlag        bool                     // indicates whether fill and text colors are different
	ws               float64                  // word spacing
	images           map[string]imageInfoType // array of used images
	pageLinks        [][]linkType             // pageLinks[page][link], both 1-based
	links            []intLinkType            // array of internal links
	autoPageBreak    bool                     // automatic page breaking
	acceptPageBreak  func() bool              // returns true to accept page break
	pageBreakTrigger float64                  // threshold used to trigger page breaks
	inHeader         bool                     // flag set when processing header
	headerFnc        func()                   // function provided by app and called to write header
	inFooter         bool                     // flag set when processing footer
	footerFnc        func()                   // function provided by app and called to write footer
	zoomMode         string                   // zoom display mode
	layoutMode       string                   // layout display mode
	title            string                   // title
	subject          string                   // subject
	author           string                   // author
	keywords         string                   // keywords
	creator          string                   // creator
	aliasNbPagesStr  string                   // alias for total number of pages
	pdfVersion       string                   // PDF version number
	fontDirStr       string                   // location of font definition files
	capStyle         int                      // line cap style: butt 0, round 1, square 2
	joinStyle        int                      // line segment join style: miter 0, round 1, bevel 2
	err              error                    // Set if error occurs during life cycle of instance
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
