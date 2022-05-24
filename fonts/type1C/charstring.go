package type1c

import (
	"fmt"

	"github.com/speedata/textlayout/fonts"
	ps "github.com/speedata/textlayout/fonts/psinterpreter"
)

// LoadGlyph parses the glyph charstring to compute segments and path bounds.
// It returns an error if the glyph is invalid or if decoding the charstring fails.
func (f *Font) LoadGlyph(glyph fonts.GID) ([]fonts.Segment, ps.PathBounds, error) {
	var (
		psi    ps.Machine
		loader type2CharstringHandler
		err    error
	)
	// if f.fdSelect != nil {
	// 	index, err = f.fdSelect.fontDictIndex(glyph)
	// 	if err != nil {
	// 		return nil, ps.PathBounds{}, err
	// 	}
	// }
	if int(glyph) >= len(f.CharStrings) {
		return nil, ps.PathBounds{}, fmt.Errorf("invalid glyph index %d", glyph)
	}

	subrs := f.subrsIndex
	err = psi.Run(f.CharStrings[glyph], subrs, f.global.globalSubrIndex, &loader)
	return loader.cs.Segments, loader.cs.Bounds, err
}

// type2CharstringHandler implements operators needed to fetch Type2 charstring metrics
type type2CharstringHandler struct {
	cs ps.CharstringReader

	// found in private DICT, needed since we can't differenciate
	// no width set from 0 width
	// `width` must be initialized to default width
	nominalWidthX int32
	width         int32
}

func (type2CharstringHandler) Context() ps.PsContext { return ps.Type2Charstring }

func (met *type2CharstringHandler) Apply(op ps.PsOperator, state *ps.Machine) error {
	var err error
	if !op.IsEscaped {
		switch op.Operator {
		case 11: // return
			return state.Return() // do not clear the arg stack
		case 14: // endchar
			if state.ArgStack.Top > 0 { // width is optional
				met.width = met.nominalWidthX + state.ArgStack.Vals[0]
			}
			met.cs.ClosePath()
			return ps.ErrInterrupt
		case 10: // callsubr
			return ps.LocalSubr(state) // do not clear the arg stack
		case 29: // callgsubr
			return ps.GlobalSubr(state) // do not clear the arg stack
		case 21: // rmoveto
			if state.ArgStack.Top > 2 { // width is optional
				met.width = met.nominalWidthX + state.ArgStack.Vals[0]
			}
			err = met.cs.Rmoveto(state)
		case 22: // hmoveto
			if state.ArgStack.Top > 1 { // width is optional
				met.width = met.nominalWidthX + state.ArgStack.Vals[0]
			}
			err = met.cs.Hmoveto(state)
		case 4: // vmoveto
			if state.ArgStack.Top > 1 { // width is optional
				met.width = met.nominalWidthX + state.ArgStack.Vals[0]
			}
			err = met.cs.Vmoveto(state)
		case 1, 18: // hstem, hstemhm
			met.cs.Hstem(state)
		case 3, 23: // vstem, vstemhm
			met.cs.Vstem(state)
		case 19, 20: // hintmask, cntrmask
			// variable number of arguments, but always even
			// for xxxmask, if there are arguments on the stack, then this is an impliied stem
			if state.ArgStack.Top&1 != 0 {
				met.width = met.nominalWidthX + state.ArgStack.Vals[0]
			}
			met.cs.Hintmask(state)
			// the stack is managed by the previous call
			return nil

		case 5: // rlineto
			met.cs.Rlineto(state)
		case 6: // hlineto
			met.cs.Hlineto(state)
		case 7: // vlineto
			met.cs.Vlineto(state)
		case 8: // rrcurveto
			met.cs.Rrcurveto(state)
		case 24: // rcurveline
			err = met.cs.Rcurveline(state)
		case 25: // rlinecurve
			err = met.cs.Rlinecurve(state)
		case 26: // vvcurveto
			met.cs.Vvcurveto(state)
		case 27: // hhcurveto
			met.cs.Hhcurveto(state)
		case 30: // vhcurveto
			met.cs.Vhcurveto(state)
		case 31: // hvcurveto
			met.cs.Hvcurveto(state)
		default:
			// no other operands are allowed before the ones handled above
			err = fmt.Errorf("invalid operator %s in charstring", op)
		}
	} else {
		switch op.Operator {
		case 34: // hflex
			err = met.cs.Hflex(state)
		case 35: // flex
			err = met.cs.Flex(state)
		case 36: // hflex1
			err = met.cs.Hflex1(state)
		case 37: // flex1
			err = met.cs.Flex1(state)
		default:
			// no other operands are allowed before the ones handled above
			err = fmt.Errorf("invalid operator %s in charstring", op)
		}
	}
	state.ArgStack.Clear()
	return err
}

// func (met *type2CharstringHandler) hstem(state *ps.Machine) {
// 	met.hstemCount += state.ArgStack.Top / 2
// }

// func (met *type2CharstringHandler) vstem(state *ps.Machine) {
// 	met.vstemCount += state.ArgStack.Top / 2
// }

// func (met *type2CharstringHandler) determineHintmaskSize(state *ps.Machine) {
// 	if !met.seenHintmask {
// 		met.vstemCount += state.ArgStack.Top / 2
// 		met.hintmaskSize = (met.hstemCount + met.vstemCount + 7) >> 3
// 		met.seenHintmask = true
// 	}
// }

// func (met *type2CharstringHandler) hintmask(state *ps.Machine) {
// 	met.determineHintmaskSize(state)
// 	state.SkipBytes(met.hintmaskSize)
// }

// // psType2CharstringsData contains fields specific to the Type 2 Charstrings
// // context.
// type psType2CharstringsData struct {
// 	f          *Font
// 	b          *Buffer
// 	x          int32
// 	y          int32
// 	firstX     int32
// 	firstY     int32
// 	hintBits   int32
// 	seenWidth  bool
// 	ended      bool
// 	glyphIndex GlyphIndex
// 	// fdSelectIndexPlusOne is the result of the Font Dict Select lookup, plus
// 	// one. That plus one lets us use the zero value to denote either unused
// 	// (for CFF fonts with a single Font Dict) or lazily evaluated.
// 	fdSelectIndexPlusOne int32
// }

// func (d *psType2CharstringsData) closePath() {
// 	if d.x != d.firstX || d.y != d.firstY {
// 		d.b.segments = append(d.b.segments, Segment{
// 			Op: SegmentOpLineTo,
// 			Args: [3]fixed.Point26_6{{
// 				X: fixed.Int26_6(d.firstX),
// 				Y: fixed.Int26_6(d.firstY),
// 			}},
// 		})
// 	}
// }

// func (d *psType2CharstringsData) moveTo(dx, dy int32) {
// 	d.closePath()
// 	d.x += dx
// 	d.y += dy
// 	d.b.segments = append(d.b.segments, Segment{
// 		Op: SegmentOpMoveTo,
// 		Args: [3]fixed.Point26_6{{
// 			X: fixed.Int26_6(d.x),
// 			Y: fixed.Int26_6(d.y),
// 		}},
// 	})
// 	d.firstX = d.x
// 	d.firstY = d.y
// }

// func (d *psType2CharstringsData) lineTo(dx, dy int32) {
// 	d.x += dx
// 	d.y += dy
// 	d.b.segments = append(d.b.segments, Segment{
// 		Op: SegmentOpLineTo,
// 		Args: [3]fixed.Point26_6{{
// 			X: fixed.Int26_6(d.x),
// 			Y: fixed.Int26_6(d.y),
// 		}},
// 	})
// }

// func (d *psType2CharstringsData) cubeTo(dxa, dya, dxb, dyb, dxc, dyc int32) {
// 	d.x += dxa
// 	d.y += dya
// 	xa := fixed.Int26_6(d.x)
// 	ya := fixed.Int26_6(d.y)
// 	d.x += dxb
// 	d.y += dyb
// 	xb := fixed.Int26_6(d.x)
// 	yb := fixed.Int26_6(d.y)
// 	d.x += dxc
// 	d.y += dyc
// 	xc := fixed.Int26_6(d.x)
// 	yc := fixed.Int26_6(d.y)
// 	d.b.segments = append(d.b.segments, Segment{
// 		Op:   SegmentOpCubeTo,
// 		Args: [3]fixed.Point26_6{{X: xa, Y: ya}, {X: xb, Y: yb}, {X: xc, Y: yc}},
// 	})
// }

type psInterpreter struct{}

type psOperator struct {
	// run is the function that implements the operator. Nil means that we
	// ignore the operator, other than popping its arguments off the stack.
	run func(*psInterpreter) error
	// name is the operator name. An empty name (i.e. the zero value for the
	// struct overall) means an unrecognized 1-byte operator.
	name string
	// numPop is the number of stack values to pop. -1 means "array" and -2
	// means "delta" as per 5176.CFF.pdf Table 6 "Operand Types".
	numPop int32
}

// psOperators holds the 1-byte and 2-byte operators for PostScript interpreter
// contexts.
var psOperators = [...][2][]psOperator{
	// // The Type 2 Charstring operators are defined by 5177.Type2.pdf Appendix A
	// // "Type 2 Charstring Command Codes".
	// psContextType2Charstring: {{
	// 	// 1-byte operators.
	// 	0:  {}, // Reserved.
	// 	2:  {}, // Reserved.
	// 	1:  {-1, "hstem", t2CStem},
	// 	3:  {-1, "vstem", t2CStem},
	// 	18: {-1, "hstemhm", t2CStem},
	// 	23: {-1, "vstemhm", t2CStem},
	// 	5:  {-1, "rlineto", t2CRlineto},
	// 	6:  {-1, "hlineto", t2CHlineto},
	// 	7:  {-1, "vlineto", t2CVlineto},
	// 	8:  {-1, "rrcurveto", t2CRrcurveto},
	// 	9:  {}, // Reserved.
	// 	10: {+1, "callsubr", t2CCallsubr},
	// 	11: {+0, "return", t2CReturn},
	// 	12: {}, // escape.
	// 	13: {}, // Reserved.
	// 	14: {-1, "endchar", t2CEndchar},
	// 	15: {}, // Reserved.
	// 	16: {}, // Reserved.
	// 	17: {}, // Reserved.
	// 	19: {-1, "hintmask", t2CMask},
	// 	20: {-1, "cntrmask", t2CMask},
	// 	4:  {-1, "vmoveto", t2CVmoveto},
	// 	21: {-1, "rmoveto", t2CRmoveto},
	// 	22: {-1, "hmoveto", t2CHmoveto},
	// 	24: {-1, "rcurveline", t2CRcurveline},
	// 	25: {-1, "rlinecurve", t2CRlinecurve},
	// 	26: {-1, "vvcurveto", t2CVvcurveto},
	// 	27: {-1, "hhcurveto", t2CHhcurveto},
	// 	28: {}, // shortint.
	// 	29: {+1, "callgsubr", t2CCallgsubr},
	// 	30: {-1, "vhcurveto", t2CVhcurveto},
	// 	31: {-1, "hvcurveto", t2CHvcurveto},
	// }, {
	// 	// 2-byte operators. The first byte is the escape byte.
	// 	34: {+7, "hflex", t2CHflex},
	// 	36: {+9, "hflex1", t2CHflex1},
	// 	// TODO: more operators.
	// }},
}
