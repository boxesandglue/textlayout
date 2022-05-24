package truetype

import (
	"encoding/binary"
	"errors"
)

var errInvalidMaxpTable = errors.New("invalid maxp table")

// TableHVmtx contains advance width and side bearings for each glyph.
type TableHVmtx []Metric // with length numGlyphs

// return the base side bearing, handling invalid glyph index
func (t TableHVmtx) getSideBearing(glyph GID) int16 {
	if int(glyph) >= len(t) {
		return 0
	}
	return t[glyph].SideBearing
}

// Metric contains the advance and side bearing of each glyph.
type Metric struct {
	Advance, SideBearing int16
}

// pad the width if numberOfHMetrics < numGlyphs
func parseHVmtxTable(input []byte, numberOfHMetrics, numGlyphs uint16) (TableHVmtx, error) {
	if numberOfHMetrics == 0 {
		return nil, errors.New("number of glyph metrics is 0")
	}

	if len(input) < 4*int(numberOfHMetrics) {
		return nil, errors.New("invalid h/vmtx table (EOF)")
	}
	widths := make(TableHVmtx, numberOfHMetrics)
	for i := range widths {
		widths[i].Advance = int16(binary.BigEndian.Uint16(input[4*i:]))
		widths[i].SideBearing = int16(binary.BigEndian.Uint16(input[4*i+2:]))
	}
	if left := int(numGlyphs) - int(numberOfHMetrics); left > 0 {
		// advances are padded with the last value
		// side bearings are given
		widths = append(widths, make(TableHVmtx, numGlyphs-numberOfHMetrics)...)
		lastWidth := widths[numberOfHMetrics-1]
		for i := numberOfHMetrics; i < numGlyphs; i++ {
			widths[i] = lastWidth
		}
		input = input[4*int(numberOfHMetrics):]
		if len(input) < 2*left {
			return nil, errors.New("invalid h/vmtx table (EOF)")
		}
		subslice := widths[numberOfHMetrics:]
		for i := range subslice {
			subslice[i].SideBearing = int16(binary.BigEndian.Uint16(input[2*i:]))
		}
	}
	return widths, nil
}
