package type1c

import (
	"github.com/boxesandglue/textlayout/fonts"
	"github.com/boxesandglue/textlayout/fonts/glyphsnames"
)

func (fnt *Font) GlyphName(g fonts.GID) string {
	if int(g) < len(fnt.charset) {
		sid := fnt.charset[int(g)]
		if int(sid) < len(fnt.global.strings) {
			return fnt.global.strings[sid]
		}
	}
	return ""
}

func (fnt *Font) NumGlyphs() int {
	return len(fnt.CharStrings)
}

// Type1 fonts have no natural notion of Unicode code points
// We use a glyph names table to identify the most commonly used runes
func (f *Font) synthesizeCmap() {
	f.cmap = make(map[rune]fonts.GID)
	for gid := range f.CharStrings {
		glyphName := f.GlyphName(fonts.GID(gid))
		r, _ := glyphsnames.GlyphToRune(glyphName)
		f.cmap[r] = fonts.GID(gid)
	}
}

func (f *Font) Cmap() (fonts.Cmap, fonts.CmapEncoding) {
	return f.cmap, fonts.EncUnicode
}
