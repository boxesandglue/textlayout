package truetype

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/boxesandglue/textlayout/fonts"
)

func binarywrite(w io.Writer, data interface{}) error {
	err := binary.Write(w, binary.BigEndian, data)
	if err != nil {
		return err
	}
	return nil
}

func (fnt *Font) getAdditionalCodepoints(codepoint GID) []GID {
	var additionalCodepoints []GID
	cp := fnt.Glyf[codepoint]
	switch t := cp.data.(type) {
	case simpleGlyphData:
		additionalCodepoints = append(additionalCodepoints, codepoint)
	case *compositeGlyphData:
		for _, g := range t.glyphs {
			additionalCodepoints = append(additionalCodepoints, fnt.getAdditionalCodepoints(g.glyphIndex)...)
		}
	}
	return additionalCodepoints
}

// subsetTrueType removes all data from the font file that is not necessary to
// render the given code points.
func (fnt *Font) subsetTrueType(codepoints []GID) error {
	var additionalCodepoints []GID
	for _, gid := range codepoints {
		cp := fnt.Glyf[gid]
		switch t := cp.data.(type) {
		case compositeGlyphData:
			for _, g := range t.glyphs {
				additionalCodepoints = append(additionalCodepoints, fnt.getAdditionalCodepoints(g.glyphIndex)...)
			}
		}
	}
	codepoints = append(codepoints, additionalCodepoints...)
	fonts.RemoveDuplicates(codepoints)

	maxCP := codepoints[len(codepoints)-1] + 1
	// the codepoints not used in the subset (or used from one of these glyphs)
	// are replaced by an empty glyph.
	glyphs := make([]GlyphData, maxCP)

	emptyGlyph := GlyphData{}
	for i, c := GID(0), 0; i < maxCP; i++ {
		if i == codepoints[c] {
			glyphs[i] = fnt.Glyf[i]
			c++
		} else {
			fnt.Hmtx[i].Advance = 0
			fnt.Hmtx[i].SideBearing = 0
			glyphs[i] = emptyGlyph
		}
	}
	fnt.Glyf = glyphs
	fnt.Hmtx = fnt.Hmtx[:maxCP]
	fnt.NumGlyphs = int(maxCP)
	fnt.Head.indexToLocFormat = 1
	fnt.hhea.NumberOfHMetrics = uint16(maxCP)
	fnt.subsetCodepoints = codepoints
	return nil
}

func (fnt *Font) subsetCFF(codepoints []GID) error {
	fnt.subsetCodepoints = codepoints
	fnt.cff.Subset(codepoints)
	return nil
}

// WidthsPDF returns a width entry suitable for embedding in a PDF file.
func (fnt *Font) WidthsPDF() string {
	r := float64(fnt.upem) / 1000

	getWd := func(cp GID) string {
		rounded := math.Round(10*float64(fnt.Hmtx[cp].Advance)/r) / 10
		return strconv.FormatFloat(rounded, 'f', -1, 64)
	}

	var b strings.Builder
	b.WriteString("[")
	c := 0
	for {
		if c >= len(fnt.subsetCodepoints) {
			break
		}
		cp := fnt.subsetCodepoints[c]
		fmt.Fprintf(&b, "%d[%s", cp, getWd(cp))
		c++
		for c < len(fnt.subsetCodepoints) && fnt.subsetCodepoints[c] == cp+1 {
			cp++
			fmt.Fprintf(&b, " %s", getWd(cp))
			c++
		}
		fmt.Fprintf(&b, "]")
	}
	b.WriteString("]")
	return b.String()
}

// CMapPDF returns a CMap string to be used in a PDF file
func (fnt *Font) CMapPDF() string {
	var numGlyphs int

	if fnt.cff != nil {
		// numGlyphs = len(fnt..Font[0].CharStrings)
	} else {
		numGlyphs = fnt.NumGlyphs
	}

	cm, _ := fnt.Cmap()
	iter := cm.Iter()
	toUni := make(map[GID]rune)
	for iter.Next() {
		a, b := iter.Char()
		toUni[b] = a
	}

	var b strings.Builder
	b.WriteString(`/CIDInit /ProcSet findresource begin
12 dict begin
begincmap
/CIDSystemInfo << /Registry (Adobe)/Ordering (UCS)/Supplement 0>> def
/CMapName /Adobe-Identity-UCS def /CMapType 2 def
1 begincodespacerange
`)
	fmt.Fprintf(&b, "<0001><%04X>\n", numGlyphs)
	b.WriteString("endcodespacerange\n")
	fmt.Fprintf(&b, "%d beginbfchar\n", len(fnt.subsetCodepoints))
	for _, cp := range fnt.subsetCodepoints {

		fmt.Fprintf(&b, "<%04X><%04X>\n", cp, toUni[cp])
	}
	b.WriteString(`endbfchar
endcmap CMapName currentdict /CMap defineresource pop end end`)
	return b.String()
}

// NamePDF returns the PDF name of the font file
func (fnt *Font) NamePDF() string {
	return fmt.Sprintf("/%s-%s", fnt.SubsetID, fnt.PostscriptName())
}

// AscenderPDF returns the /Ascent value for the PDF file
func (fnt *Font) AscenderPDF() int {
	return int(fnt.hhea.Ascent)
}

// DescenderPDF returns the /Descent value for the PDF file
func (fnt *Font) DescenderPDF() int {
	return int(fnt.hhea.Descent)
}

// CapHeightPDF returns the /CapHeight value for the PDF file
func (fnt *Font) CapHeightPDF() int {
	ch := int(fnt.OS2.SCapHeight)
	return ch
}

// BoundingBoxPDF returns the /FontBBox value for the PDF file
func (fnt *Font) BoundingBoxPDF() string {
	return fmt.Sprintf("[%d %d %d %d]", 0, fnt.hhea.Descent, 1000, fnt.hhea.Ascent)
}

// FlagsPDF returns the /Flags value for the PDF file
func (fnt *Font) FlagsPDF() int {
	return 4
}

// ItalicAnglePDF returns the /ItalicAngle value for the PDF file
func (fnt *Font) ItalicAnglePDF() int {
	return int(fnt.post.ItalicAngle)
}

// StemVPDF returns the /StemV value for the PDF file
func (fnt *Font) StemVPDF() int {
	return 0
}

// XHeightPDF returns the /XHeight value for the PDF file
func (fnt *Font) XHeightPDF() int {
	xh := int(fnt.OS2.SxHeigh)
	return xh
}

// Subset removes all data from the font except the one needed for the given
// code points.
func (fnt *Font) Subset(codepoints []GID) error {
	fnt.SubsetID = getCharTag(codepoints)
	if fnt.cff == nil {
		err := fnt.subsetTrueType(codepoints)
		return err
	}
	return fnt.subsetCFF(codepoints)
}

type tableOffsetLength struct {
	offset    uint32
	length    uint32
	tag       Tag
	checksum  uint32
	tabledata []byte
}

func (fnt *Font) writeHead(w io.Writer) error {
	type head struct {
		majorVersion       uint16 // 1
		minorVersion       uint16 // 0
		fontRevision       uint32 // fixed
		checksumAdjustment uint32 // to be calculated
		magicNumber        uint32 // 0x5F0F3CF5
		flags              uint16
		unitsPerEm         uint16
		created            uint64
		modified           uint64
		xMin               uint16
		yMin               uint16
		xMax               uint16
		yMax               uint16
		macStyle           uint16
		lowestRecPPEM      uint16
		fontDirectionHint  int16
		indexToLocFormat   int16
		glyphDataFormat    int16
	}
	h := head{
		majorVersion:       1,
		minorVersion:       0,
		fontRevision:       fnt.Head.FontRevision,
		checksumAdjustment: 0,
		magicNumber:        0x5F0F3CF5,
		flags:              fnt.Head.Flags,
		unitsPerEm:         fnt.Head.UnitsPerEm,
		created:            fnt.Head.Created.SecondsSince1904,
		modified:           fnt.Head.Updated.SecondsSince1904,
		xMin:               uint16(fnt.Head.XMin),
		yMin:               uint16(fnt.Head.YMin),
		xMax:               uint16(fnt.Head.XMax),
		yMax:               uint16(fnt.Head.YMax),
		macStyle:           fnt.Head.MacStyle,
		lowestRecPPEM:      fnt.Head.LowestRecPPEM,
		fontDirectionHint:  fnt.Head.FontDirection,
		indexToLocFormat:   fnt.Head.indexToLocFormat,
		glyphDataFormat:    int16(fnt.Head.glyphDataFormat),
	}
	err := binary.Write(w, binary.BigEndian, h)
	return err
}

func (fnt *Font) writeGlyf(w io.Writer) error {
	glyphOffsets := []uint32{}
	c := uint32(0)
	for i := 0; i < fnt.NumGlyphs; i++ {
		g := fnt.Glyf[i]
		glyphOffsets = append(glyphOffsets, c)
		w.Write(g.rawdata)
		c += uint32(len(g.rawdata))
	}
	glyphOffsets = append(glyphOffsets, c)
	fnt.glyphOffsets = glyphOffsets
	return nil
}

func (fnt *Font) writeHHea(w io.Writer) error {
	tbl := fnt.hhea
	binarywrite(w, uint16(1))
	binarywrite(w, uint16(0))

	var reserved int16

	binarywrite(w, tbl.Ascent)
	binarywrite(w, tbl.Descent)
	binarywrite(w, tbl.LineGap)

	binarywrite(w, tbl.AdvanceMax)
	binarywrite(w, tbl.MinFirstSideBearing)
	binarywrite(w, tbl.MinSecondSideBearing)
	binarywrite(w, tbl.MaxExtent)
	binarywrite(w, tbl.CaretSlopeRise)
	binarywrite(w, tbl.CaretSlopeRun)
	binarywrite(w, tbl.CaretOffset)

	binarywrite(w, reserved)
	binarywrite(w, reserved)
	binarywrite(w, reserved)
	binarywrite(w, reserved)

	binarywrite(w, tbl.MetricDataFormat)
	binarywrite(w, tbl.NumberOfHMetrics)
	return nil
}

func (fnt *Font) writeCvt(w io.Writer) error {
	_, err := w.Write(fnt.cvt)
	return err
}

func (fnt *Font) writePrep(w io.Writer) error {
	_, err := w.Write(fnt.prep)
	return err
}

func (fnt *Font) writeMaxp(w io.Writer) error {
	tbl := fnt.Maxp
	tbl.NumGlyphs = uint16(fnt.NumGlyphs)

	binarywrite(w, tbl.Version)
	binarywrite(w, tbl.NumGlyphs)

	switch tbl.Version {
	case 0x10000:
		binarywrite(w, tbl.MaxPoints)
		binarywrite(w, tbl.MaxContours)
		binarywrite(w, tbl.MaxCompositePoints)
		binarywrite(w, tbl.MaxCompositeContours)
		binarywrite(w, tbl.MaxZones)
		binarywrite(w, tbl.MaxTwilightPoints)
		binarywrite(w, tbl.MaxStorage)
		binarywrite(w, tbl.MaxFunctionDefs)
		binarywrite(w, tbl.MaxInstructionDefs)
		binarywrite(w, tbl.MaxStackElements)
		binarywrite(w, tbl.MaxSizeOfInstructions)
		binarywrite(w, tbl.MaxComponentElements)
		binarywrite(w, tbl.MaxComponentDepth)
	default:
		// version 0.5 only has NumGlyphs
	}
	return nil
}

func (fnt *Font) writeLoca(w io.Writer) error {
	version := fnt.Head.indexToLocFormat
	switch version {
	case 0:
		var offset uint16
		for _, off := range fnt.glyphOffsets {
			offset = uint16(off / 2)
			binarywrite(w, offset)
		}
	case 1:
		for _, off := range fnt.glyphOffsets {
			binarywrite(w, off)
		}
	}

	return nil
}

func (fnt *Font) writeHmtx(w io.Writer) error {
	var err error
	tbl := fnt.Hmtx
	l := GID(fnt.NumGlyphs)
	for i := GID(0); i < l; i++ {
		if err = binarywrite(w, uint16(tbl[i].Advance)); err != nil {
			return err
		}
		if err = binarywrite(w, tbl[i].SideBearing); err != nil {
			return err
		}
	}

	return nil
}

// WriteTable writes the table to w.
func (fnt *Font) writeTable(w io.Writer, t Tag) error {
	var err error
	switch t {
	// case "CFF ":
	// 	err = tt.CFF.WriteCFFData(w)
	case tagLoca:
		err = fnt.writeLoca(w)
	case tagHhea:
		err = fnt.writeHHea(w)
	case tagHead:
		err = fnt.writeHead(w)
	case tagMaxp:
		err = fnt.writeMaxp(w)
	case tagHmtx:
		err = fnt.writeHmtx(w)
	// case tagfpg:
	// 	err = fnt.writeFpgm(w)
	case tagCvt:
		err = fnt.writeCvt(w)
	case tagPrep:
		err = fnt.writePrep(w)
	case tagGlyf:
		err = fnt.writeGlyf(w)
	// case tagPost:
	// 	err = fnt.writePost(w)
	// case tagOS2:
	// 	err = fnt.writeOs2(w)
	default:
		// fmt.Printf("    skip write table %s\n", tbl)
	}
	if err != nil {
		return err
	}
	return nil
}

func calcChecksum(data []byte) uint32 {
	sum := uint32(0)
	c := 0
	for c < len(data) {
		sum += uint32(data[c])<<3 + uint32(data[c+1])<<2 + uint32(data[c+2])<<1 + uint32(data[c+3])
		c += 4
	}
	return sum
}

// WriteSubset writes a valid font to w that is suitable for including in PDF
func (fnt *Font) WriteSubset(w io.Writer) error {
	if fnt.cff != nil {
		return fnt.cff.WriteSubset(w)
	}
	var err error
	var fontfile bytes.Buffer
	fnt.Head.checkSumAdjustment = 0

	tablesForPDF := []tableOffsetLength{}

	// put only those tables in PDF which are present in the font file
	for _, tblname := range []Tag{tagCvt, tagGlyf, tagHead, tagHhea, tagHmtx, tagLoca, tagMaxp, tagPrep} {
		if _, ok := fnt.knowTables[tblname]; ok {
			tbl := tableOffsetLength{}
			tbl.tag = tblname
			tablesForPDF = append(tablesForPDF, tbl)
		}
	}

	// tables start at 12 (header) + table toc
	tableOffset := uint32(12 + 16*len(tablesForPDF))
	var newTables []tableOffsetLength
	for _, tbl := range tablesForPDF {
		var tableData bytes.Buffer
		err = fnt.writeTable(&tableData, tbl.tag)
		if err != nil {
			return err
		}
		l := tableData.Len()
		nt := tableOffsetLength{
			length: uint32(l),
			tag:    tbl.tag,
			offset: tableOffset,
		}

		switch l & 3 {
		case 0:
			// ok, no alignment
		case 1:
			binarywrite(&tableData, uint16(0))
			binarywrite(&tableData, uint8(0))
			l += 3
		case 2:
			binarywrite(&tableData, uint16(0))
			l += 2
		case 3:
			binarywrite(&tableData, uint8(0))
			l++
		}
		nt.tabledata = tableData.Bytes()
		tableOffset += uint32(len(nt.tabledata))
		nt.checksum = calcChecksum(nt.tabledata)
		newTables = append(newTables, nt)
	}

	binarywrite(&fontfile, fnt.Type)
	cTablesRead := float64(len(newTables))
	searchRange := (math.Pow(2, math.Floor(math.Log2(cTablesRead))) * 16)
	entrySelector := math.Floor(math.Log2(cTablesRead))
	rangeShift := (cTablesRead * 16.0) - searchRange

	binarywrite(&fontfile, uint16(cTablesRead))
	binarywrite(&fontfile, uint16(searchRange))
	binarywrite(&fontfile, uint16(entrySelector))
	binarywrite(&fontfile, uint16(rangeShift))

	checksumAdjustmentOffset := 0
	for _, tbl := range newTables {
		binarywrite(&fontfile, []byte(tbl.tag.String()))
		binarywrite(&fontfile, tbl.checksum)
		binarywrite(&fontfile, tbl.offset)
		binarywrite(&fontfile, tbl.length)
		if tbl.tag == tagHead {
			checksumAdjustmentOffset = int(tbl.offset) + 8
		}
	}

	for _, tbl := range newTables {
		binarywrite(&fontfile, tbl.tabledata)
	}

	b := fontfile.Bytes()
	checksumFontFile := calcChecksum(b)
	if checksumAdjustmentOffset > 0 {
		// only if we write the head table
		binary.BigEndian.PutUint32(b[checksumAdjustmentOffset:], checksumFontFile-0xB1B0AFBA)
	}
	w.Write(b)

	return nil
}

// getCharTag returns a string of length 6 based on the characters in code point
// list. All returned characters are in the range A-Z.
func getCharTag(codepoints []GID) string {
	// sort the code points so we can create reproducible PDFs
	sort.Sort(fonts.SortByGID(codepoints))
	data := make([]byte, len(codepoints)*2)
	for i, r := range codepoints {
		data[i*2] = byte((r >> 8) & 0xff)
		data[i*2+1] = byte(r & 0xff)
	}

	sum := md5.Sum(data)
	ret := make([]rune, 6)
	for i := 0; i < 6; i++ {
		ret[i] = rune(sum[2*i]+sum[2*i+1])/26 + 'A'
	}
	return string(ret)
}
