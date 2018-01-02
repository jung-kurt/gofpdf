package gofpdf

// Page orientation
const (
	// PageOrientationPortrait ...
	PageOrientationPortrait = "P"
	// PageOrientationLandscape ...
	PageOrientationLandscape = "L"
)

// Size Unit
const (
	// UnitPoint ...
	UnitPoint = "pt"
	// UnitMillimeter ...
	UnitMillimeter = "mm"
	// UnitCentimeter
	UnitCentimeter = "cm"
	// UnitInch
	UnitInch = "inch"
)

// Page size
const (
	// PageSizeA3 ...
	PageSizeA3 = "A3"
	// PageSizeA4 ...
	PageSizeA4 = "A4"
	// PageSizeA5 ...
	PageSizeA5 = "A5"
	// PageSizeLetter ...
	PageSizeLetter = "Letter"
	// PageSizeLegal ...
	PageSizeLegal = "Legal"
)

// Border
const (
	// BorderNone ...
	BorderNone = ""
	// BorderFull ...
	BorderFull = "1"
	// BorderLeft ...
	BorderLeft = "L"
	// BorderTop ...
	BorderTop = "T"
	// BorderRight ...
	BorderRight = "R"
	// BorderBottom ...
	BorderBottom = "B"
)

// LineBreak
const (
	LineBreakNone   = 0
	LineBreakNormal = 1
	LineBreakBelow  = 2
)

// Alignment
const (
	// AlignLeft ...
	AlignLeft = "L"
	// AlignRight ...
	AlignRight = "R"
	// AlignCenter ...
	AlignCenter = "C"
	// AlignTop ...
	AlignTop = "T"
	// AlignBottom ...
	AlignBottom = "B"
	// AlignMiddle ...
	AlignMiddle = "M"
	// AlignBaseline ...
	AlignBaseline = "B"
)
