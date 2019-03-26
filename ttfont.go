package gofpdf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

// Define the value used in the "head" table of a created TTF file
// 0x74727565 "true" for Mac
// 0x00010000 for Windows
// Either seems to work for a font embedded in a PDF file
// when read by Adobe Reader on a Windows PC(!)
const ttfMacHeader = false

// TrueType Font Glyph operators
const gfWords = 1 << 0
const gfScale = 1 << 3
const gfMore = 1 << 5
const gfXYScale = 1 << 6
const gfTwoByTwo = 1 << 7

//TTFontFile contain .ttf file data
type TTFontFile struct {
	MaxUni             int
	pos                int
	numTables          int
	searchRange        int
	entrySelector      int
	rangeShift         int
	tables             map[string]*tableRecord
	oTables            map[string][]byte
	filename           string
	fh                 *ttfFile
	hMetrics           int
	glyphPos           []int
	charToGlyph        map[int]int
	Ascent             int
	Descent            int
	name               string
	familyName         string
	styleName          string
	FullName           string
	uniqueFontID       string
	unitsPerEm         int
	Bbox               []float64
	CapHeight          int
	StemV              int
	ItalicAngle        int
	Flags              int
	UnderlinePosition  float64
	UnderlineThickness float64
	CharWidths         []int
	DefaultWidth       float64
	maxStrLenRead      int
	glyphData          map[int]map[string][]int
	CodeToGlyph        map[int]int
}

type tableRecord struct {
	tag      string
	checksum []int
	offset   int
	length   int
}

type ttfFile struct {
	pos  int64
	arr  []byte
	file *os.File
}

func (ttf *ttfFile) Read(s int) []byte {
	if ttf.arr != nil {
		a := ttf.arr[ttf.pos : ttf.pos+int64(s)]
		ttf.pos += int64(s)
		return a
	}
	answer := make([]byte, s)
	_, _ = ttf.file.Read(answer)
	return answer

}

func (ttf *ttfFile) Seek(shift int64, flag int) (int64, error) {
	if ttf.arr != nil {
		if flag == 0 {
			ttf.pos = shift
		} else if flag == 1 {
			ttf.pos += shift
		} else if flag == 2 {
			ttf.pos = int64(len(ttf.arr)) - shift
		}
		return int64(ttf.pos), nil
	}
	return ttf.file.Seek(shift, flag)
}

func (ttf *ttfFile) Close() error {
	if ttf.file != nil {
		return ttf.file.Close()
	}
	ttf.arr = nil
	return nil
}

//NewTTFontFile make new empty TTFont
func NewTTFontFile() *TTFontFile {
	ttf := TTFontFile{
		maxStrLenRead: 20000000,
	}
	return &ttf
}

//GetBytesMetrics fill TTFile metrics from .ttf file byte array
func (ttf *TTFontFile) GetBytesMetrics(arr []byte) error {
	return ttf.getMetrics("", arr)
}

//GetFileMetrics fill TTFile metrics from .ttf file
func (ttf *TTFontFile) GetFileMetrics(file string) error {
	return ttf.getMetrics(file, nil)
}

func (ttf *TTFontFile) getMetrics(file string, arr []byte) error {
	ttf.filename = file
	var err error
	if arr != nil {
		a := ttfFile{pos: 0, arr: arr, file: nil}
		ttf.fh = &a
	} else {
		f, _ := os.Open(file)
		defer f.Close()
		a := ttfFile{pos: 0, arr: nil, file: f}
		ttf.fh = &a
	}
	if err != nil {
		return fmt.Errorf("Cant open file %s\n ", file)
	}
	ttf.pos = 0
	ttf.CharWidths = make([]int, 0)
	ttf.glyphPos = make([]int, 0)
	ttf.charToGlyph = make(map[int]int)
	ttf.tables = make(map[string]*tableRecord)
	ttf.oTables = make(map[string][]byte)
	ttf.Ascent = 0
	ttf.Descent = 0
	version := uint32(ttf.readUlong())
	if version == 0x4F54544F {
		return fmt.Errorf("Postscript outlines are not supported\n ")
	}
	if version == 0x74746366 {
		return fmt.Errorf("ERROR - TrueType Fonts Collections not supported\n ")
	}
	if version != 0x00010000 && version != 0x74727565 {
		return fmt.Errorf("Not a TrueType font: version=%v\n ", version)
	}
	ttf.readTableDirectory()
	ttf.extractInfo()
	return nil
}

func (ttf *TTFontFile) readTableDirectory() {

	ttf.numTables = ttf.readUshort()
	ttf.searchRange = ttf.readUshort()
	ttf.entrySelector = ttf.readUshort()
	ttf.rangeShift = ttf.readUshort()
	ttf.tables = make(map[string]*tableRecord)

	for i := 0; i < ttf.numTables; i++ {
		record := tableRecord{
			tag:      ttf.readTag(),
			checksum: []int{ttf.readUshort(), ttf.readUshort()},
			offset:   ttf.readUlong(),
			length:   ttf.readUlong(),
		}
		ttf.tables[record.tag] = &record
	}
}

func (ttf *TTFontFile) readTag() string {
	ttf.pos += 4
	answer := ttf.fh.Read(4)
	return string(answer)
}

func (ttf *TTFontFile) readUshort() int {
	ttf.pos += 2
	s := ttf.fh.Read(2)
	return (int(s[0]) << 8) + int(s[1])
}

func (ttf *TTFontFile) readUlong() int {
	ttf.pos += 4
	s := ttf.fh.Read(4)
	return (int(s[0]) * 16777216) + (int(s[1]) << 16) + (int(s[2]) << 8) + int(s[3]) // 	16777216  = 1<<24
}

func (ttf *TTFontFile) sub32(x, y []int) []int {
	xlo := x[1]
	xhi := x[0]
	ylo := y[1]
	yhi := y[0]
	if ylo > xlo {
		xlo += 1 << 16
		yhi++
	}

	reslo := xlo - ylo
	if yhi > xhi {
		xhi += 1 << 16
	}
	reshi := xhi - yhi
	reshi = reshi & 0xFFFF
	return []int{reshi, reslo}
}

func (ttf *TTFontFile) calcChecksum(data []byte) []int {
	if (len(data) % 4) != 0 {
		for i := 0; (len(data) % 4) != 0; i++ {
			data = append(data, 0)
		}
	}
	hi := 0x0000
	lo := 0x0000
	for i := 0; i < len(data); i += 4 {
		hi += (int(data[i]) << 8) + int(data[i+1])
		lo += (int(data[i+2]) << 8) + int(data[i+3])
		hi += lo >> 16
		lo = lo & 0xFFFF
		hi = hi & 0xFFFF
	}
	return []int{hi, lo}
}

func (ttf *TTFontFile) getTablePos(tag string) []int {
	offset := ttf.tables[tag].offset
	length := ttf.tables[tag].length
	return []int{offset, length}
}

func (ttf *TTFontFile) seek(pos int) {
	ttf.pos = pos
	_, _ = ttf.fh.Seek(int64(ttf.pos), 0)
}

func (ttf *TTFontFile) skip(delta int) {
	ttf.pos += delta
	_, _ = ttf.fh.Seek(int64(ttf.pos), 0)
}

//SeekTable position
func (ttf *TTFontFile) SeekTable(tag string) int {
	return ttf.seekTable(tag, 0)
}

func (ttf *TTFontFile) seekTable(tag string, offsetInTable int) int {
	tpos := ttf.getTablePos(tag)
	ttf.pos = tpos[0] + offsetInTable
	_, _ = ttf.fh.Seek(int64(ttf.pos), 0)
	return ttf.pos
}

func (ttf *TTFontFile) readShort() int16 {
	ttf.pos += 2
	s := ttf.fh.Read(2)
	a := (int16(s[0]) << 8) + int16(s[1])
	if (int(a) & (1 << 15)) == 0 {
		a = int16(int(a) - (1 << 16))
	}
	return a
}

func (ttf *TTFontFile) getUshort(pos int) int {
	_, _ = ttf.fh.Seek(int64(pos), 0)
	s := ttf.fh.Read(2)
	return (int(s[0]) << 8) + int(s[1])
}

func (ttf *TTFontFile) splice(stream []byte, offset int, value []byte) []byte {
	return append(append(stream[:offset], value...), stream[offset+len(value):]...)
}

func (ttf *TTFontFile) setUshort(stream []byte, offset int, value int) []byte {
	up := make([]byte, 2)
	binary.BigEndian.PutUint16(up, uint16(value))
	return ttf.splice(stream, offset, up)
}

func (ttf *TTFontFile) getChunk(pos, length int) []byte {
	ttf.fh.Seek(int64(pos), 0)
	if length < 1 {
		return make([]byte, 0)
	}
	s := ttf.fh.Read(length)
	return s
}

func (ttf *TTFontFile) getTable(tag string) []byte {
	a := ttf.getTablePos(tag)
	pos, length := a[0], a[1]
	if length == 0 {
		return nil //, fmt.Errorf("Truetype font (%s): error reading table: %s\n", ttf.filename, tag)
	}
	ttf.fh.Seek(int64(pos), 0)
	s := ttf.fh.Read(length)
	return s //, nil
}

func (ttf *TTFontFile) add(tag string, data []byte) {
	if tag == "head" {
		data = ttf.splice(data, 8, []byte{0, 0, 0, 0})
	}
	ttf.oTables[tag] = data
}

func arrayKeys(arr map[int]string) []int {
	answer := make([]int, len(arr))
	i := 0
	for key := range arr {
		answer[i] = key
		i++
	}
	return answer
}

func inArray(s int, arr []int) bool {
	for _, i := range arr {
		if s == i {
			return true
		}
	}
	return false
}

func (ttf *TTFontFile) extractInfo() {
	///////////////////////////////////
	// name - Naming table
	///////////////////////////////////
	nameOffset := ttf.SeekTable("name")
	format := ttf.readUshort()
	if format != 0 {
		fmt.Printf("Unknown name table format %d\n", format)
		return
	}
	numRecords := ttf.readUshort()
	stringDataOffset := nameOffset + ttf.readUshort()
	names := map[int]string{1: "", 2: "", 3: "", 4: "", 6: ""}
	K := arrayKeys(names)
	nameCount := len(names)
	for i := 0; i < numRecords; i++ {
		platformID := ttf.readUshort()
		encodingID := ttf.readUshort()
		languageID := ttf.readUshort()
		nameID := ttf.readUshort()
		length := ttf.readUshort()
		offset := ttf.readUshort()
		if !inArray(nameID, K) {
			continue
		}
		N := ""
		if platformID == 3 && encodingID == 1 && languageID == 0x409 { // Microsoft, Unicode, US English, PS Name
			opos := ttf.pos
			ttf.seek(stringDataOffset + offset)
			if length%2 != 0 {
				fmt.Printf("PostScript name is UTF-16BE string of odd length\n")
				return
			}
			length /= 2
			N = ""
			for length > 0 {
				char := ttf.readUshort()
				N += string(rune(char))
				length--
			}
			ttf.pos = opos
			ttf.seek(opos)
		} else if platformID == 1 && encodingID == 0 && languageID == 0 { // Macintosh, Roman, English, PS Name
			opos := ttf.pos
			N = string(ttf.getChunk(stringDataOffset+offset, length))
			ttf.pos = opos
			ttf.seek(opos)
		}
		if N != "" && names[nameID] == "" {
			names[nameID] = N
			nameCount--
			if nameCount == 0 {
				break
			}
		}
	}
	var psName string
	if names[6] != "" {
		psName = names[6]
	} else if names[4] != "" {
		psName = strings.Replace(names[4], " ", "-", -1)
	} else if names[1] != "" {
		psName = strings.Replace(names[1], " ", "-", -1)
	} else {
		psName = ""
	}
	if psName == "" {
		fmt.Printf("Could not find PostScript font name\n")
		return
	}
	ttf.name = psName
	if names[1] != "" {
		ttf.familyName = names[1]
	} else {
		ttf.familyName = psName
	}
	if names[2] != "" {
		ttf.styleName = names[2]
	} else {
		ttf.styleName = "Regular"
	}
	if names[4] != "" {
		ttf.FullName = names[4]
	} else {
		ttf.FullName = psName
	}
	if names[3] != "" {
		ttf.uniqueFontID = names[3]
	} else {
		ttf.uniqueFontID = psName
	}
	if names[6] != "" {
		ttf.FullName = names[6]
	}

	///////////////////////////////////
	// head - Font header table
	///////////////////////////////////
	ttf.SeekTable("head")
	ttf.skip(18)
	ttf.unitsPerEm = ttf.readUshort()
	scale := 1000.0 / float64(ttf.unitsPerEm)
	ttf.skip(16)
	xMin := ttf.readShort()
	yMin := ttf.readShort()
	xMax := ttf.readShort()
	yMax := ttf.readShort()
	ttf.Bbox = []float64{float64(xMin) * scale, float64(yMin) * scale, float64(xMax) * scale, float64(yMax) * scale}
	ttf.skip(3 * 2)
	_ = ttf.readUshort()
	glyphDataFormat := ttf.readUshort()
	if glyphDataFormat != 0 {
		fmt.Printf("Unknown glyph data format %d\n", glyphDataFormat)
		return
	}

	///////////////////////////////////
	// hhea metrics table
	///////////////////////////////////
	// ttf2t1 seems to use this value rather than the one in OS/2 - so put in for compatibility
	if _, OK := ttf.tables["hhea"]; OK {
		ttf.SeekTable("hhea")
		ttf.skip(4)
		hheaAscender := ttf.readShort()
		hheaDescender := ttf.readShort()
		ttf.Ascent = int(float64(hheaAscender) * scale)
		ttf.Descent = int(float64(hheaDescender) * scale)
	}

	///////////////////////////////////
	// OS/2 - OS/2 and Windows metrics table
	///////////////////////////////////
	usWeightClass := 0
	if _, OK := ttf.tables["OS/2"]; OK {
		ttf.SeekTable("OS/2")
		version := ttf.readUshort()
		ttf.skip(2)
		usWeightClass = ttf.readUshort()
		ttf.skip(2)
		fsType := ttf.readUshort()
		if fsType == 0x0002 || (fsType&0x0300) != 0 {
			fmt.Printf("ERROR - Font file %s cannot be embedded due to copyright restrictions.\n", ttf.filename)
			return
		}
		ttf.skip(20)
		_ = ttf.readShort()
		ttf.pos += 10 //PANOSE = 10 byte length
		_ = ttf.fh.Read(10)
		ttf.skip(26)
		sTypoAscender := ttf.readShort()
		sTypoDescender := ttf.readShort()
		if ttf.Ascent == 0 {
			ttf.Ascent = int(float64(sTypoAscender) * scale)
		}
		if ttf.Descent == 0 {
			ttf.Descent = int(float64(sTypoDescender) * scale)
		}
		if version > 1 {
			ttf.skip(16)
			sCapHeight := ttf.readShort()
			ttf.CapHeight = int(float64(sCapHeight) * scale)
		} else {
			ttf.CapHeight = ttf.Ascent
		}
	} else {
		usWeightClass = 500
		if ttf.Ascent == 0 {
			ttf.Ascent = int(float64(yMax) * scale)
		}
		if ttf.Descent == 0 {
			ttf.Descent = int(float64(yMin) * scale)
		}
		ttf.CapHeight = ttf.Ascent
	}
	ttf.StemV = 50 + int(math.Pow(float64(usWeightClass)/65.0, 2))

	///////////////////////////////////
	// post - PostScript table
	///////////////////////////////////
	ttf.SeekTable("post")
	ttf.skip(4)
	ttf.ItalicAngle = int(ttf.readShort()) + ttf.readUshort()/65536.0
	ttf.UnderlinePosition = float64(ttf.readShort()) * scale
	ttf.UnderlineThickness = float64(ttf.readShort()) * scale
	isFixedPitch := ttf.readUlong()

	ttf.Flags = 4

	if ttf.ItalicAngle != 0 {
		ttf.Flags = ttf.Flags | 64
	}
	if usWeightClass >= 600 {
		ttf.Flags = ttf.Flags | 262144
	}
	if isFixedPitch != 0 {
		ttf.Flags = ttf.Flags | 1
	}

	///////////////////////////////////
	// hhea - Horizontal header table
	///////////////////////////////////
	ttf.SeekTable("hhea")
	ttf.skip(32)
	metricDataFormat := ttf.readUshort()
	if metricDataFormat != 0 {
		fmt.Printf("Unknown horizontal metric data format %d\n", metricDataFormat)
		return
	}
	numberOfHMetrics := ttf.readUshort()
	if numberOfHMetrics == 0 {
		fmt.Printf("Number of horizontal metrics is 0\n")
		return
	}

	///////////////////////////////////
	// maxp - Maximum profile table
	///////////////////////////////////
	ttf.SeekTable("maxp")
	ttf.skip(4)
	numGlyphs := ttf.readUshort()

	///////////////////////////////////
	// cmap - Character to glyph index mapping table
	///////////////////////////////////
	cmapOffset := ttf.SeekTable("cmap")
	ttf.skip(2)
	cmapTableCount := ttf.readUshort()
	unicodeCmapOffset := 0
	for i := 0; i < cmapTableCount; i++ {
		platformID := ttf.readUshort()
		encodingID := ttf.readUshort()
		offset := ttf.readUlong()
		savePos := ttf.pos
		if (platformID == 3 && encodingID == 1) || platformID == 0 { // Microsoft, Unicode
			format = ttf.getUshort(cmapOffset + offset)
			if format == 4 {
				if unicodeCmapOffset == 0 {
					unicodeCmapOffset = cmapOffset + offset
				}
				break
			}
		}
		ttf.seek(savePos)
	}
	if unicodeCmapOffset == 0 {
		fmt.Printf("Font (%s) does not have cmap for Unicode (platform 3, encoding 1, format 4, or platform 0, any encoding, format 4)\n", ttf.filename)
		return
	}

	glyphToChar := make(map[int][]int)
	charToGlyph := make(map[int]int)
	ttf.getCMAP4(unicodeCmapOffset, glyphToChar, charToGlyph)

	///////////////////////////////////
	// hmtx - Horizontal metrics table
	///////////////////////////////////
	ttf.getHMTX(numberOfHMetrics, numGlyphs, glyphToChar, scale)

}

///////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////

//MakeArraySubset fill TTFontFile from .ttf file byte array, only with runes from subset
func (ttf *TTFontFile) MakeArraySubset(arr []byte, subset map[int]int) []byte {
	return ttf.makeSubset("", subset, arr)
}

//MakeSubset fill TTFontFile from .ttf file, only with runes from subset
func (ttf *TTFontFile) MakeSubset(file string, subset map[int]int) []byte {
	return ttf.makeSubset(file, subset, nil)
}

func (ttf *TTFontFile) makeSubset(file string, subset map[int]int, arr []byte) []byte {
	ttf.filename = file
	var err error
	if arr != nil {
		a := ttfFile{pos: 0, arr: arr, file: nil}
		ttf.fh = &a
	} else {
		f, _ := os.Open(file)
		defer f.Close()
		a := ttfFile{pos: 0, arr: nil, file: f}
		ttf.fh = &a
	}
	defer ttf.fh.Close()
	if err != nil {
		return nil
	}

	ttf.pos = 0
	ttf.CharWidths = make([]int, 0)
	ttf.glyphPos = make([]int, 0)
	ttf.charToGlyph = make(map[int]int)
	ttf.tables = make(map[string]*tableRecord)
	ttf.oTables = make(map[string][]byte)
	ttf.Ascent = 0
	ttf.Descent = 0
	ttf.skip(4)
	ttf.MaxUni = 0
	ttf.readTableDirectory()

	///////////////////////////////////
	// head - Font header table
	///////////////////////////////////
	ttf.SeekTable("head")
	ttf.skip(50)
	indexToLocFormat := ttf.readUshort()
	_ = ttf.readUshort()

	///////////////////////////////////
	// hhea - Horizontal header table
	///////////////////////////////////
	ttf.SeekTable("hhea")
	ttf.skip(32)
	_ = ttf.readUshort()
	numberOfHMetrics := ttf.readUshort()
	orignHmetrics := numberOfHMetrics
	///////////////////////////////////
	// maxp - Maximum profile table
	///////////////////////////////////
	ttf.SeekTable("maxp")
	ttf.skip(4)
	numGlyphs := ttf.readUshort()

	///////////////////////////////////
	// cmap - Character to glyph index mapping table
	///////////////////////////////////
	cmapOffset := ttf.SeekTable("cmap")
	ttf.skip(2)
	cmapTableCount := ttf.readUshort()
	unicodeCmapOffset := 0
	for i := 0; i < cmapTableCount; i++ {
		platformID := ttf.readUshort()
		encodingID := ttf.readUshort()
		offset := ttf.readUlong()
		savePos := ttf.pos
		if (platformID == 3 && encodingID == 1) || platformID == 0 { // Microsoft, Unicode
			format := ttf.getUshort(cmapOffset + offset)
			if format == 4 {
				unicodeCmapOffset = cmapOffset + offset
				break
			}
		}
		ttf.seek(savePos)
	}

	if unicodeCmapOffset == 0 {
		fmt.Printf("Font (%s) does not have cmap for Unicode (platform 3, encoding 1, format 4, or platform 0, any encoding, format 4)\n", ttf.filename)
		return nil
	}

	glyphToChar := make(map[int][]int)
	charToGlyph := make(map[int]int)
	ttf.getCMAP4(unicodeCmapOffset, glyphToChar, charToGlyph)

	ttf.charToGlyph = charToGlyph

	///////////////////////////////////
	// hmtx - Horizontal metrics table
	///////////////////////////////////
	scale := 1.0 // not used
	ttf.getHMTX(numberOfHMetrics, numGlyphs, glyphToChar, scale)

	///////////////////////////////////
	// loca - Index to location
	///////////////////////////////////
	ttf.getLOCA(indexToLocFormat, numGlyphs)

	subsetglyphs := map[int]int{0: 0}
	subsetCharToGlyph := make(map[int]int)
	for _, code := range subset {
		if _, OK := ttf.charToGlyph[code]; OK {
			subsetglyphs[ttf.charToGlyph[code]] = code      // Old Glyph ID => Unicode
			subsetCharToGlyph[code] = ttf.charToGlyph[code] // Unicode to old GlyphID

		}
		ttf.MaxUni = max(ttf.MaxUni, code)
	}

	l := ttf.getTablePos("glyf")
	start := l[0]

	glyphSet := make(map[int]int)
	subsetglyphsKeys := keySortInt(subsetglyphs)

	n := 0
	fsLastCharIndex := 0 // maximum Unicode index (character code) in this font, according to the cmap subtable for platform ID 3 and platform- specific encoding ID 0 or 1.
	for _, originalGlyphIdx := range subsetglyphsKeys {
		fsLastCharIndex = max(fsLastCharIndex, subsetglyphs[originalGlyphIdx])
		glyphSet[originalGlyphIdx] = n // old glyphID to new glyphID
		n++
	}
	subsetCharToGlyphKeys := keySortInt(subsetCharToGlyph)
	codeToGlyph := make(map[int]int)
	for _, uni := range subsetCharToGlyphKeys {
		codeToGlyph[uni] = glyphSet[subsetCharToGlyph[uni]]
	}
	ttf.CodeToGlyph = codeToGlyph

	subsetglyphsKeys = keySortInt(subsetglyphs)
	for _, originalGlyphIdx := range subsetglyphsKeys {
		_, glyphSet, subsetglyphs, subsetglyphsKeys = ttf.getGlyphs(originalGlyphIdx, &start, glyphSet, subsetglyphs, subsetglyphsKeys)
	}
	numberOfHMetrics = len(subsetglyphs)
	numGlyphs = numberOfHMetrics

	//tables copied from the original
	tags := []string{"name"}
	for _, tag := range tags {
		ttf.add(tag, ttf.getTable(tag))
	}
	tags = []string{"cvt ", "fpgm", "prep", "gasp"}
	for _, tag := range tags {
		if _, OK := ttf.tables[tag]; OK {
			ttf.add(tag, ttf.getTable(tag))
		}
	}
	// post - PostScript
	opost := ttf.getTable("post")
	post := append(append([]byte{0x00, 0x03, 0x00, 0x00}, opost[4:16]...), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
	ttf.add("post", post)

	// Sort CID2GID map into segments of contiguous codes
	delete(codeToGlyph, 0)
	codeToGlyphKeys := keySortInt(codeToGlyph)

	//unset(codeToGlyph[65535]);
	rangeID := 0
	arrayRange := make(map[int][]int)
	prevCid := -2
	prevGlyph := -1
	// for each character
	for _, cid := range codeToGlyphKeys {
		//for cid, glidx := range codeToGlyph {
		if cid == (prevCid+1) && codeToGlyph[cid] == (prevGlyph+1) {
			if n, OK := arrayRange[rangeID]; !OK || n == nil {
				arrayRange[rangeID] = make([]int, 0)
			}
			arrayRange[rangeID] = append(arrayRange[rangeID], codeToGlyph[cid])
		} else {
			// new arrayRange
			rangeID = cid
			arrayRange[rangeID] = make([]int, 0)
			arrayRange[rangeID] = append(arrayRange[rangeID], codeToGlyph[cid])
		}
		prevCid = cid
		prevGlyph = codeToGlyph[cid]
	}

	arrayRangeKeys := keySortArrayRangeMap(arrayRange)
	// cmap - Character to glyph mapping - Format 4 (MS / )
	segCount := len(arrayRange) + 1 // + 1 Last segment has missing character 0xFFFF
	searchRange := 1
	entrySelector := 0
	for searchRange*2 <= segCount {
		searchRange = searchRange * 2
		entrySelector = entrySelector + 1
	}
	searchRange = searchRange * 2
	rangeShift := segCount*2 - searchRange
	length := 16 + (8 * segCount) + (numGlyphs + 1)
	cmap := []int{0, 1, // Index : version, number of encoding subtables
		3, 1, // Encoding Subtable : platform (MS=3), encoding (Unicode)
		0, 12, // Encoding Subtable : offset (hi,lo)
		4, length, 0, // Format 4 Mapping subtable: format, length, language
		segCount * 2,
		searchRange,
		entrySelector,
		rangeShift}

	// endCode(s)
	for _, start := range arrayRangeKeys {
		endCode := start + (len(arrayRange[start]) - 1)
		cmap = append(cmap, endCode) // endCode(s)
	}
	cmap = append(cmap, 0xFFFF) // endCode of last Segment
	cmap = append(cmap, 0)      // reservedPad

	// startCode(s)
	for _, start := range arrayRangeKeys {
		cmap = append(cmap, start) // startCode(s)
	}
	cmap = append(cmap, 0xFFFF) // startCode of last Segment
	// idDelta(s)
	for _, start := range arrayRangeKeys {
		idDelta := -(start - arrayRange[start][0])
		n += len(arrayRange[start])
		cmap = append(cmap, idDelta) // idDelta(s)
	}
	cmap = append(cmap, 1) // idDelta of last Segment
	// idRangeOffset(s)
	for range arrayRange {
		cmap = append(cmap, 0) // idRangeOffset[segCount]  	Offset in bytes to glyph indexArray, or 0

	}
	cmap = append(cmap, 0) // idRangeOffset of last Segment
	for _, start := range arrayRangeKeys {
		for _, glidx := range arrayRange[start] {
			cmap = append(cmap, glidx)
		}
	}
	cmap = append(cmap, 0) // Mapping for last character
	cmapstr := make([]byte, 0)
	for _, cm := range cmap {
		cmapstr = append(cmapstr, packN16(cm)...)
	}
	ttf.add("cmap", cmapstr)

	// glyf - Glyph data
	l = ttf.getTablePos("glyf")
	glyfOffset, glyfLength := l[0], l[1]
	glyphData := make([]byte, 0)
	if glyfLength < ttf.maxStrLenRead {
		glyphData = ttf.getTable("glyf")
	}

	offsets := make([]int, 0)
	glyfData := make([]byte, 0)
	pos := 0
	hmtxData := make([]byte, 0)
	ttf.glyphData = make(map[int]map[string][]int, 0)

	for _, originalGlyphIdx := range subsetglyphsKeys {
		// hmtx - Horizontal Metrics
		hm := ttf.getHMetric(orignHmetrics, originalGlyphIdx)
		hmtxData = append(hmtxData, hm...)

		offsets = append(offsets, pos)
		glyphPos := ttf.glyphPos[originalGlyphIdx]
		glyphLen := ttf.glyphPos[originalGlyphIdx+1] - glyphPos
		data := make([]byte, 0)
		if glyfLength < ttf.maxStrLenRead {
			data = glyphData[glyphPos : glyphPos+glyphLen]
		} else {
			if glyphLen > 0 {
				data = ttf.getChunk(glyfOffset+glyphPos, glyphLen)
			}
		}

		var up int
		if glyphLen > 0 {
			up = unpackN16(data[0:2])
		}

		if glyphLen > 2 && (up&(1<<15)) != 0 { // If number of contours <= -1 i.e. composite glyph
			posInGlyph := 10
			flags := gfMore
			nComponentElements := 0
			for (flags & gfMore) != 0 {
				nComponentElements++ // number of glyphs referenced at top level
				up = unpackN16(data[posInGlyph : posInGlyph+2])
				flags = up
				up = unpackN16(data[posInGlyph+2 : posInGlyph+4])
				glyphIdx := up
				if _, OK := ttf.glyphData[originalGlyphIdx]; !OK {
					ttf.glyphData[originalGlyphIdx] = make(map[string][]int)
				}
				if _, OK := ttf.glyphData[originalGlyphIdx]["compGlyphs"]; !OK {
					ttf.glyphData[originalGlyphIdx]["compGlyphs"] = make([]int, 0)
				}
				ttf.glyphData[originalGlyphIdx]["compGlyphs"] = append(ttf.glyphData[originalGlyphIdx]["compGlyphs"], glyphIdx)
				data = ttf.setUshort(data, posInGlyph+2, glyphSet[glyphIdx])
				posInGlyph += 4
				if (flags & gfWords) != 0 {
					posInGlyph += 4
				} else {
					posInGlyph += 2
				}
				if (flags & gfScale) != 0 {
					posInGlyph += 2
				} else if (flags & gfXYScale) != 0 {
					posInGlyph += 4
				} else if (flags & gfTwoByTwo) != 0 {
					posInGlyph += 8
				}
			}
			//maxComponentElements = max(maxComponentElements, nComponentElements);
		}

		glyfData = append(glyfData, data...)
		pos += glyphLen
		if pos%4 != 0 {
			padding := 4 - (pos % 4)
			glyfData = append(glyfData, make([]byte, padding)...)
			pos += padding
		}
	}

	offsets = append(offsets, pos)
	ttf.add("glyf", glyfData)

	// hmtx - Horizontal Metrics
	ttf.add("hmtx", hmtxData)

	// loca - Index to location
	locaData := make([]byte, 0)
	if ((pos + 1) >> 1) > 0xFFFF {
		indexToLocFormat = 1 // long format
		for _, offset := range offsets {
			locaData = append(locaData, packN32(offset)...)
		}
	} else {
		indexToLocFormat = 0 // short format
		for _, offset := range offsets {
			locaData = append(locaData, packN16(offset/2)...)
		}
	}
	ttf.add("loca", locaData)

	// head - Font header
	headData := ttf.getTable("head")
	headData = ttf.setUshort(headData, 50, indexToLocFormat)
	ttf.add("head", headData)

	// hhea - Horizontal Header
	hheaData := ttf.getTable("hhea")
	hheaData = ttf.setUshort(hheaData, 34, numberOfHMetrics)
	ttf.add("hhea", hheaData)

	// maxp - Maximum Profile
	maxp := ttf.getTable("maxp")
	maxp = ttf.setUshort(maxp, 4, numGlyphs)
	ttf.add("maxp", maxp)

	// OS/2 - OS/2
	os2Data := ttf.getTable("OS/2")
	ttf.add("OS/2", os2Data)
	//fclose(ttf.fh);
	defer ttf.fh.Close()

	// Put the TTF file together
	stm := make([]byte, 0)
	stm = ttf.endTTFile(stm)
	return stm
}

//////////////////////////////////////////////////////////////////////////////////
// Recursively get composite glyphs
//////////////////////////////////////////////////////////////////////////////////
func (ttf *TTFontFile) getGlyphs(originalGlyphIdx int, start *int, glyphSet map[int]int, subsetGlyphs map[int]int, subsetGlyphsKeys []int) (*int, map[int]int, map[int]int, []int) {
	glyphPos := ttf.glyphPos[originalGlyphIdx]
	glyphLen := ttf.glyphPos[originalGlyphIdx+1] - glyphPos
	if glyphLen == 0 {
		return start, glyphSet, subsetGlyphs, subsetGlyphsKeys
	}
	ttf.seek(*start + glyphPos)

	numberOfContours := ttf.readShort()

	if numberOfContours < 0 {
		ttf.skip(8)
		flags := gfMore
		for flags&gfMore != 0 {
			flags = ttf.readUshort()
			glyphIdx := ttf.readUshort()
			if _, OK := glyphSet[glyphIdx]; !OK {
				glyphSet[glyphIdx] = len(subsetGlyphs) // old glyphID to new glyphID
				subsetGlyphs[glyphIdx] = 1
				subsetGlyphsKeys = append(subsetGlyphsKeys, glyphIdx)
			}
			savePos, _ := ttf.fh.Seek(0, 1)
			_, _, _, subsetGlyphsKeys = ttf.getGlyphs(glyphIdx, start, glyphSet, subsetGlyphs, subsetGlyphsKeys)
			ttf.seek(int(savePos))
			if flags&gfWords != 0 {
				ttf.skip(4)
			} else {
				ttf.skip(2)
			}
			if flags&gfScale != 0 {
				ttf.skip(2)
			} else if flags&gfXYScale != 0 {
				ttf.skip(4)
			} else if flags&gfTwoByTwo != 0 {
				ttf.skip(8)
			}
		}
	}
	return start, glyphSet, subsetGlyphs, subsetGlyphsKeys
}

//////////////////////////////////////////////////////////////////////////////////

func (ttf *TTFontFile) getHMTX(numberOfHMetrics, numGlyphs int, glyphToChar map[int][]int, scale float64) {
	start := ttf.SeekTable("hmtx")
	aw := 0
	w := 0
	var arr []int
	ttf.CharWidths = make([]int, 256*256)
	nCharWidths := 0
	if (numberOfHMetrics * 4) < ttf.maxStrLenRead {
		data := ttf.getChunk(start, numberOfHMetrics*4)
		arr = unpackN16Array(data)
	} else {
		ttf.seek(start)
	}
	for glyph := 0; glyph < numberOfHMetrics; glyph++ {

		if (numberOfHMetrics * 4) < ttf.maxStrLenRead {
			aw = arr[(glyph*2)+1]
		} else {
			aw = ttf.readUshort()
			_ = ttf.readUshort()
		}
		if _, OK := glyphToChar[glyph]; OK || glyph == 0 {

			if aw >= (1 << 15) {
				aw = 0
			} // 1.03 Some (arabic) fonts have -ve values for width
			// although should be unsigned value - comes out as e.g. 65108 (intended -50)
			if glyph == 0 {
				ttf.DefaultWidth = scale * float64(aw)
				continue
			}
			for _, char := range glyphToChar[glyph] {
				if char != 0 && char != 65535 {
					w = int(math.Round(scale * float64(aw)))
					if w == 0 {
						w = 65535
					}
					if char < 196608 {
						ttf.CharWidths[char] = w
						nCharWidths++
					}
				}
			}
		}
	}
	data := ttf.getChunk(start+numberOfHMetrics*4, numGlyphs*2)
	arr = unpackN16Array(data)
	diff := numGlyphs - numberOfHMetrics
	for pos := 0; pos < diff; pos++ {
		glyph := pos + numberOfHMetrics
		if _, OK := glyphToChar[glyph]; OK {
			for _, char := range glyphToChar[glyph] {
				if char != 0 && char != 65535 {
					w = int(math.Round(scale * float64(aw)))
					if w == 0 {
						w = 65535
					}
					if char < 196608 {
						ttf.CharWidths[char] = w
						nCharWidths++
					}
				}
			}
		}
	}
	// NB 65535 is a set width of 0
	// First bytes define number of chars in font
	ttf.CharWidths[0] = nCharWidths
}

func (ttf *TTFontFile) getHMetric(numberOfHMetrics, gid int) []byte {
	start := ttf.SeekTable("hmtx")
	var hm []byte
	if gid < numberOfHMetrics {
		ttf.seek(start + (gid * 4))
		hm = fRead(ttf.fh, 4)
	} else {
		ttf.seek(start + ((numberOfHMetrics - 1) * 4))
		hm = fRead(ttf.fh, 2)
		ttf.seek(start + (numberOfHMetrics * 2) + (gid * 2))
		hm = append(hm, fRead(ttf.fh, 2)...)
	}
	return hm
}

func (ttf *TTFontFile) getLOCA(indexToLocFormat, numGlyphs int) {
	start := ttf.SeekTable("loca")
	ttf.glyphPos = make([]int, 0)
	if indexToLocFormat == 0 {
		data := ttf.getChunk(start, (numGlyphs*2)+2)
		arr := unpackN16Array(data)
		for n := 0; n <= numGlyphs; n++ {
			ttf.glyphPos = append(ttf.glyphPos, arr[n+1]*2)
		}
	} else if indexToLocFormat == 1 {
		data := ttf.getChunk(start, (numGlyphs*4)+4)
		arr := unpackN32Array(data)
		for n := 0; n <= numGlyphs; n++ {
			ttf.glyphPos = append(ttf.glyphPos, arr[n+1])
		}
	} else {
		fmt.Printf("Unknown location table format %d\n", indexToLocFormat)
		return
	}
}

// CMAP Format 4
func (ttf *TTFontFile) getCMAP4(unicodeCmapOffset int, glyphToChar map[int][]int, charToGlyph map[int]int) {
	maxUniChar := 0
	ttf.seek(unicodeCmapOffset + 2)
	length := ttf.readUshort()
	limit := unicodeCmapOffset + length
	ttf.skip(2)

	segCount := ttf.readUshort() / 2
	ttf.skip(6)
	endCount := make([]int, 0)
	for i := 0; i < segCount; i++ {
		endCount = append(endCount, ttf.readUshort())
	}
	ttf.skip(2)
	startCount := make([]int, 0)
	for i := 0; i < segCount; i++ {
		startCount = append(startCount, ttf.readUshort())
	}
	idDelta := make([]int, 0)
	for i := 0; i < segCount; i++ {
		idDelta = append(idDelta, int(ttf.readShort()))
	} // ???? was unsigned short
	idRangeOffsetStart := ttf.pos
	idRangeOffset := make([]int, 0)
	for i := 0; i < segCount; i++ {
		idRangeOffset = append(idRangeOffset, ttf.readUshort())
	}
	glyph := 0
	for n := 0; n < segCount; n++ {
		endpoint := endCount[n] + 1
		for unichar := startCount[n]; unichar < endpoint; unichar++ {
			if idRangeOffset[n] == 0 {
				glyph = (unichar + idDelta[n]) & 0xFFFF
			} else {
				offset := (unichar-startCount[n])*2 + idRangeOffset[n]
				offset = idRangeOffsetStart + 2*n + offset
				if offset >= limit {
					glyph = 0
				} else {
					glyph = ttf.getUshort(offset)
					if glyph != 0 {
						glyph = (glyph + idDelta[n]) & 0xFFFF
					}
				}
			}
			charToGlyph[unichar] = glyph
			if unichar < 196608 {
				maxUniChar = max(unichar, maxUniChar)
			}
			glyphToChar[glyph] = append(glyphToChar[glyph], unichar)
		}
	}
}

func max(i, n int) int {
	if n > i {
		return n
	}
	return i
}

// Put the TTF file together
func (ttf *TTFontFile) endTTFile(stm []byte) []byte {
	stm = make([]byte, 0)
	numTables := len(ttf.oTables)
	searchRange := 1
	entrySelector := 0
	for searchRange*2 <= numTables {
		searchRange = searchRange * 2
		entrySelector = entrySelector + 1
	}
	searchRange = searchRange * 16
	rangeShift := numTables*16 - searchRange

	// Header
	if ttfMacHeader {
		stm = append(stm, packHeader(0x74727565, numTables, searchRange, entrySelector, rangeShift)...) // Mac
	} else {
		stm = append(stm, packHeader(0x00010000, numTables, searchRange, entrySelector, rangeShift)...) // Windows
	}

	// Table directory
	tables := ttf.oTables
	tablesKeys := keySortStrings(tables)

	offset := 12 + numTables*16
	headStart := 0

	for _, tag := range tablesKeys {
		if tag == "head" {
			headStart = offset
		}
		stm = append(stm, []byte(tag)...)
		checksum := ttf.calcChecksum(tables[tag])
		stm = append(stm, pack2N16(checksum[0], checksum[1])...)
		stm = append(stm, pack2N32(offset, len(tables[tag]))...)
		paddedLength := (len(tables[tag]) + 3) &^ 3
		offset = offset + paddedLength
	}

	// Table data
	for _, key := range tablesKeys {
		data := tables[key]
		data = append(data, []byte{0, 0, 0}...)
		stm = append(stm, data[:(len(data)&^3)]...)
	}

	checksum := ttf.calcChecksum([]byte(stm))
	checksum = ttf.sub32([]int{0xB1B0, 0xAFBA}, checksum)
	chk := pack2N16(checksum[0], checksum[1])
	stm = ttf.splice(stm, (headStart + 8), chk)
	return stm
}

func fRead(r *ttfFile, count int) []byte {
	answer := r.Read(count)
	return answer
}

func unpackN16Array(data []byte) []int {
	answer := make([]int, 1)
	r := bytes.NewReader(data)
	bs := make([]byte, 2)
	var e error
	var c int
	c, e = r.Read(bs)
	for e == nil && c > 0 {
		answer = append(answer, int(binary.BigEndian.Uint16(bs)))
		c, e = r.Read(bs)
	}
	return answer
}

func unpackN32Array(data []byte) []int {
	answer := make([]int, 1)
	r := bytes.NewReader(data)
	bs := make([]byte, 4)
	var e error
	var c int
	c, e = r.Read(bs)
	for e == nil && c > 0 {
		answer = append(answer, int(binary.BigEndian.Uint32(bs)))
		c, e = r.Read(bs)
	}
	return answer
}

func unpackN16(data []byte) int {
	return int(binary.BigEndian.Uint16(data))
}

func packHeader(N uint32, n1, n2, n3, n4 int) []byte {
	answer := make([]byte, 0)
	bs4 := make([]byte, 4)
	binary.BigEndian.PutUint32(bs4, N)
	answer = append(answer, bs4...)
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(n1))
	answer = append(answer, bs...)
	binary.BigEndian.PutUint16(bs, uint16(n2))
	answer = append(answer, bs...)
	binary.BigEndian.PutUint16(bs, uint16(n3))
	answer = append(answer, bs...)
	binary.BigEndian.PutUint16(bs, uint16(n4))
	answer = append(answer, bs...)
	return answer
}

func pack2N16(n1, n2 int) []byte {
	answer := make([]byte, 0)
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(n1))
	answer = append(answer, bs...)
	binary.BigEndian.PutUint16(bs, uint16(n2))
	answer = append(answer, bs...)
	return answer
}

func pack2N32(n1, n2 int) []byte {
	answer := make([]byte, 0)
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(n1))
	answer = append(answer, bs...)
	binary.BigEndian.PutUint32(bs, uint32(n2))
	answer = append(answer, bs...)
	return answer
}

func packN32(n1 int) []byte {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(n1))
	return bs
}

func packN16(n1 int) []byte {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(n1))
	return bs
}

func keySortStrings(s map[string][]byte) []string {
	keys := make([]string, len(s))
	i := 0
	for key := range s {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

func keySortInt(s map[int]int) []int {
	keys := make([]int, len(s))
	i := 0
	for key := range s {
		keys[i] = key
		i++
	}
	sort.Ints(keys)
	return keys
}

func keySortArrayRangeMap(s map[int][]int) []int {
	keys := make([]int, len(s))
	i := 0
	for key := range s {
		keys[i] = key
		i++
	}
	sort.Ints(keys)
	return keys
}
