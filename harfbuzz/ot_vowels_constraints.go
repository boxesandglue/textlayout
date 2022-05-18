package harfbuzz

import "github.com/speedata/textlayout/language"

// Code generated by unicodedata/generate/main.go DO NOT EDIT.

func outputDottedCircle(buffer *Buffer) {
	buffer.outputRune(0x25CC)
	buffer.prev().resetContinutation()
}

func outputWithDottedCircle(buffer *Buffer) {
	outputDottedCircle(buffer)
	buffer.nextGlyph()
}

func preprocessTextVowelConstraints(buffer *Buffer) {
	if (buffer.Flags & DoNotinsertDottedCircle) != 0 {
		return
	}

	/* UGLY UGLY UGLY business of adding dotted-circle in the middle of
	* vowel-sequences that look like another vowel. Data for each script
	* collected from the USE script development spec.
	*
	* https://github.com/harfbuzz/harfbuzz/issues/1019
	 */
	buffer.clearOutput()
	count := len(buffer.Info)
	switch buffer.Props.Script {

	case language.Devanagari:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x0905:
				switch buffer.cur(1).codepoint {
				case 0x093A, 0x093B, 0x093E, 0x0945, 0x0946, 0x0949, 0x094A, 0x094B, 0x094C, 0x094F, 0x0956, 0x0957:
					matched = true
				}
			case 0x0906:
				switch buffer.cur(1).codepoint {
				case 0x093A, 0x0945, 0x0946, 0x0947, 0x0948:
					matched = true
				}
			case 0x0909:
				matched = 0x0941 == buffer.cur(1).codepoint
			case 0x090F:
				switch buffer.cur(1).codepoint {
				case 0x0945, 0x0946, 0x0947:
					matched = true
				}
			case 0x0930:
				if 0x094D == buffer.cur(1).codepoint &&
					buffer.idx+2 < count &&
					0x0907 == buffer.cur(2).codepoint {
					buffer.nextGlyph()
					matched = true
				}
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Bengali:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x0985:
				matched = 0x09BE == buffer.cur(1).codepoint
			case 0x098B:
				matched = 0x09C3 == buffer.cur(1).codepoint
			case 0x098C:
				matched = 0x09E2 == buffer.cur(1).codepoint
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Gurmukhi:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x0A05:
				switch buffer.cur(1).codepoint {
				case 0x0A3E, 0x0A48, 0x0A4C:
					matched = true
				}
			case 0x0A72:
				switch buffer.cur(1).codepoint {
				case 0x0A3F, 0x0A40, 0x0A47:
					matched = true
				}
			case 0x0A73:
				switch buffer.cur(1).codepoint {
				case 0x0A41, 0x0A42, 0x0A4B:
					matched = true
				}
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Gujarati:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x0A85:
				switch buffer.cur(1).codepoint {
				case 0x0ABE, 0x0AC5, 0x0AC7, 0x0AC8, 0x0AC9, 0x0ACB, 0x0ACC:
					matched = true
				}
			case 0x0AC5:
				matched = 0x0ABE == buffer.cur(1).codepoint
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Oriya:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x0B05:
				matched = 0x0B3E == buffer.cur(1).codepoint
			case 0x0B0F, 0x0B13:
				matched = 0x0B57 == buffer.cur(1).codepoint
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Tamil:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			if 0x0B85 == buffer.cur(0).codepoint &&
				0x0BC2 == buffer.cur(1).codepoint {
				matched = true
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Telugu:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x0C12:
				switch buffer.cur(1).codepoint {
				case 0x0C4C, 0x0C55:
					matched = true
				}
			case 0x0C3F, 0x0C46, 0x0C4A:
				matched = 0x0C55 == buffer.cur(1).codepoint
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Kannada:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x0C89, 0x0C8B:
				matched = 0x0CBE == buffer.cur(1).codepoint
			case 0x0C92:
				matched = 0x0CCC == buffer.cur(1).codepoint
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Malayalam:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x0D07, 0x0D09:
				matched = 0x0D57 == buffer.cur(1).codepoint
			case 0x0D0E:
				matched = 0x0D46 == buffer.cur(1).codepoint
			case 0x0D12:
				switch buffer.cur(1).codepoint {
				case 0x0D3E, 0x0D57:
					matched = true
				}
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Sinhala:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x0D85:
				switch buffer.cur(1).codepoint {
				case 0x0DCF, 0x0DD0, 0x0DD1:
					matched = true
				}
			case 0x0D8B, 0x0D8F, 0x0D94:
				matched = 0x0DDF == buffer.cur(1).codepoint
			case 0x0D8D:
				matched = 0x0DD8 == buffer.cur(1).codepoint
			case 0x0D91:
				switch buffer.cur(1).codepoint {
				case 0x0DCA, 0x0DD9, 0x0DDA, 0x0DDC, 0x0DDD, 0x0DDE:
					matched = true
				}
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Brahmi:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x11005:
				matched = 0x11038 == buffer.cur(1).codepoint
			case 0x1100B:
				matched = 0x1103E == buffer.cur(1).codepoint
			case 0x1100F:
				matched = 0x11042 == buffer.cur(1).codepoint
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Khudawadi:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x112B0:
				switch buffer.cur(1).codepoint {
				case 0x112E0, 0x112E5, 0x112E6, 0x112E7, 0x112E8:
					matched = true
				}
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Tirhuta:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x11481:
				matched = 0x114B0 == buffer.cur(1).codepoint
			case 0x1148B, 0x1148D:
				matched = 0x114BA == buffer.cur(1).codepoint
			case 0x114AA:
				switch buffer.cur(1).codepoint {
				case 0x114B5, 0x114B6:
					matched = true
				}
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Modi:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x11600, 0x11601:
				switch buffer.cur(1).codepoint {
				case 0x11639, 0x1163A:
					matched = true
				}
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	case language.Takri:
		for buffer.idx = 0; buffer.idx+1 < count; {
			matched := false
			switch buffer.cur(0).codepoint {
			case 0x11680:
				switch buffer.cur(1).codepoint {
				case 0x116AD, 0x116B4, 0x116B5:
					matched = true
				}
			case 0x11686:
				matched = 0x116B2 == buffer.cur(1).codepoint
			}

			buffer.nextGlyph()
			if matched {
				outputWithDottedCircle(buffer)
			}
		}
	}
	buffer.swapBuffers()
}
