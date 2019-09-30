package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const header = `%PS
/newpage {
    << /PageSize [300 400] /Orientation 3 >> setpagedevice
    90 rotate
    0 -300 translate
} def`

type Option interface {
	Set(val string)
	Reset()
	Push()
	Value() interface{}
}

type Font struct{ current, dflt string }

func (f *Font) Set(val string) {
	switch val {
	case "serif":
		f.current = "Times"
	case "mono":
		f.current = "Courier"
	default:
		f.current = "Helvetica"
	}
}
func (f *Font) Reset()             { f.current = f.dflt }
func (f *Font) Push()              { f.dflt = f.current }
func (f *Font) Value() interface{} { return f.current }

type Style struct{ current, dflt string }

func (f *Style) Set(val string) {
	switch val {
	case "bold":
		f.current = "Bold"
	case "italics":
		f.current = "Italics"
	default:
		f.current = ""
	}
}
func (f *Style) Reset()             { f.current = f.dflt }
func (f *Style) Push()              { f.dflt = f.current }
func (f *Style) Value() interface{} { return f.current }

type Size struct{ current, dflt int }

func (f *Size) Set(val string) {
	switch val {
	case "huge":
		f.current = 48
	case "large":
		f.current = 30
	case "small":
		f.current = 12
	case "tiny":
		f.current = 8
	default:
		f.current = 20
	}
}
func (f *Size) Reset()             { f.current = f.dflt }
func (f *Size) Push()              { f.dflt = f.current }
func (f *Size) Value() interface{} { return f.current }

type Indent struct{ current, dflt bool }

func (f *Indent) Set(val string)     { f.current = (val != "") }
func (f *Indent) Reset()             { f.current = f.dflt }
func (f *Indent) Push()              { f.dflt = f.current }
func (f *Indent) Value() interface{} { return f.current }

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

	lines   []string
	linum   = 1
	newpage = true
	opts    = map[string]Option{
		"font":   &Font{dflt: "Helvetica"},
		"style":  &Style{},
		"size":   &Size{dflt: 20},
		"indent": &Indent{},
	}
)

func psEscape(line string) string {
	var b strings.Builder
	escaped := strings.NewReplacer("(", "\\(", ")", "\\)").Replace(line)
	b.WriteString("(")
	for _, r := range escaped {
		esc, ok := glyphs[r]
		if ok {
			fmt.Fprintf(&b, ") show /%s glyphshow (", esc)
		} else {
			b.WriteRune(r)
		}
	}
	b.WriteString(") show")
	return b.String()
}

func printPage() {
	if !newpage || lines == nil {
		return
	}
	newpage = false

	size := opts["size"].Value().(int)
	font := opts["font"].Value().(string)
	style := opts["style"].Value().(string)
	indent := opts["indent"].Value().(bool)

	fmt.Printf("newpage /%s%s %d selectfont\n", font, style, size)

	base := 150
	count := len(lines)
	if count&1 != 1 {
		base = 140
	}
	for i, line := range lines {
		if line == "" {
			continue
		}

		y, x := base-size*(i-count/2), 20
		if indent {
			x = 40
		}
		fmt.Printf("%d %d moveto ", x, y)
		fmt.Println(psEscape(line))
	}
	fmt.Println("showpage")

	for _, opt := range opts {
		opt.Reset()
	}
	lines = nil
	newpage = true
}

func main() {
	in := os.Stdin
	if len(os.Args) > 1 {
		var err error
		in, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Println(header)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " \n\t\v")
		switch {
		case line == "":
			printPage()
		case strings.HasPrefix(line, "\\"):
			// escaped lines are inserted verbatim
			lines = append(lines, line[1:])
		case strings.HasPrefix(line, "#+"):
			// "command" comments are parsed seperatly
			sub := cmdRe.FindStringSubmatch(line)
			opt := opts[strings.ToLower(sub[1])]
			opt.Set(strings.ToLower(sub[3]))
			if sub[2] != "" {
				opt.Push()
			}
		case strings.HasPrefix(line, "#"):
			// regular comments are ignored
		default:
			lines = append(lines, line)
		}
		linum++
	}
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, scanner.Err().Error())
	}
	printPage()
}
