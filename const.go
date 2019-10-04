package main

import "regexp"

const header = `%PS
/height 300 def
/width 400 def
/newpage {
    << /PageSize [height width] /Orientation 3 >> setpagedevice
    90 rotate
    0 -300 translate
} def`

var (
	cmdRe  = regexp.MustCompile(`^#\+([[:alpha:]]*)(!)?:?[[:space:]]*(.*?)[[:space:]]*$`)
	glyphs = map[rune]string{
		'â': "acircumflex",
		'ä': "adieresis",
		'Ä': "Adieresis",
		'à': "agrave",
		'ç': "ccedilla",
		'Ç': "Ccedilla",
		'é': "eacute",
		'É': "Eacute",
		'ê': "ecircumflex",
		'ë': "edieresis",
		'è': "egrave",
		'€': "Euro",
		'ï': "idieresis",
		'ô': "ocircumflex",
		'ö': "odieresis",
		'Ö': "Odieresis",
		'ß': "germandbls",
		'ü': "udieresis",
		'Ü': "Udieresis",
	}
)
