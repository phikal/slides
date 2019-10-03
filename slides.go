package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const header = `%PS
/height 300 def
/width 400 def
/newpage {
    << /PageSize [height width] /Orientation 3 >> setpagedevice
    90 rotate
    0 -300 translate
} def`

type Option interface {
	Set(val string)
	Reset()
	Push()
}

type Wrapped struct {
	opt   Option
	parse func(val string)
}

func (w *Wrapped) Set(val string) { w.parse(val) }
func (w *Wrapped) Reset()         { w.opt.Reset() }
func (w *Wrapped) Push()          { w.opt.Push() }

func Map(ref *string, dfl string, tbl ...string) Option {
	mapping := make(map[string]string)
	for i := 0; i < len(tbl); i += 2 {
		mapping[tbl[i]] = tbl[i+1]
	}
	opt := &String{ref, *ref}
	return &Wrapped{
		opt,
		func(val string) {
			str, ok := mapping[val]
			if ok {
				opt.Set(str)
			} else {
				opt.Set(dfl)
			}
		},
	}
}

type Aggregated map[string]string

func (a *Aggregated) Set(val string) {
	for name, val := range *a {
		opts[name].Set(val)
	}
}
func (a *Aggregated) Reset() {
	for name := range *a {
		opts[name].Reset()
	}
}
func (a *Aggregated) Push() {}

type String struct {
	ref *string
	val string
}

func (s *String) Set(val string) {
	*s.ref = val
}
func (s *String) Reset() { *s.ref = s.val }
func (s *String) Push()  { s.val = *s.ref }

type Int struct {
	ref *int
	val int
}

func (i *Int) Set(val string) {
	j, err := strconv.Atoi(val)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid int value %q in line %d\n", val, linum)
		return
	}
	*i.ref = j
}
func (i *Int) Reset() { *i.ref = i.val }
func (i *Int) Push()  { i.val = *i.ref }

type Bool struct {
	ref *bool
	val bool
}

func (b *Bool) Set(val string) { *b.ref = (val != "") }
func (b *Bool) Reset()         { *b.ref = b.val }
func (b *Bool) Push()          { b.val = *b.ref }

type Image struct{}

func (r *Image) Set(val string) {
	file, err := os.Open(val)
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't open image file: ", err)
		return
	}
	defer file.Close()

	img, _, err = image.Decode(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while decoding image: ", err)
	}
}
func (r *Image) Reset() { img = nil }
func (r *Image) Push()  {}

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

	font    = "Helvetica"
	style   = ""
	size    = 20
	indent  = false
	height  = 300
	width   = 400
	padding = 0
	center  = false

	lines   []string
	img     image.Image
	linum   = 1
	newpage = true
	opts    = map[string]Option{
		"font":  Map(&font, "Helvetica", "mono", "Courier", "serif", "Times"),
		"style": Map(&style, "", "bold", "Bold", "italic", "Italic"),
		"size": &Wrapped{
			&Int{&size, size},
			func(val string) {
				switch val {
				case "huge":
					size = height / 6
				case "large":
					size = height / 10
				case "small":
					size = height / 25
				case "tiny":
					size = height / 30
				default:
					size = height / 15
				}
			}},
		"center":  &Bool{&center, center},
		"indent":  &Bool{&indent, indent},
		"height":  &Int{&height, height},
		"width":   &Int{&width, width},
		"padding": &Int{&padding, padding},
		"image":   &Image{},
		"title": &Aggregated{
			"center": "t",
			"style":  "bold",
			"size":   "large",
		},
	}
)

func printLine(line string) {
	escaped := strings.NewReplacer("(", "\\(", ")", "\\)", "\\", "\\\\").Replace(line)
	fmt.Print("(")
	for _, r := range escaped {
		esc, ok := glyphs[r]
		if ok {
			fmt.Printf(") show /%s glyphshow (", esc)
		} else {
			fmt.Printf("%c", r)
		}
	}
	fmt.Println(") show")
}

func printPage() {
	if img == nil && (!newpage || lines == nil) {
		return
	}
	newpage = false
	defer func() {
		for _, opt := range opts {
			opt.Reset()
		}
		lines = nil
		newpage = true
	}()

	fmt.Printf("/width %d def\n", width)
	fmt.Printf("/height %d def\n", height)
	fmt.Println("newpage")

	if img != nil {
		var scale, xoff, yoff float64
		bounds := img.Bounds()
		iheight := bounds.Max.Y - bounds.Min.Y
		iwidth := bounds.Max.X - bounds.Min.X
		if iwidth/width > iheight/height {
			scale = float64(width-2*padding) / float64(iwidth)
			yoff += (float64(iheight)*scale - float64(height-2*padding)) / 2
		} else {
			scale = float64(height-2*padding) / float64(iheight)
			xoff += (float64(iwidth)*scale - float64(width-2*padding)) / 2
		}
		fmt.Printf("%d %d 8 [%g 0 0 %g %g %g] {<",
			iwidth, iheight, 1/scale, 1/scale,
			(-float64(padding)+xoff)/scale,
			(-float64(padding)+yoff)/scale)
		for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				fmt.Printf("%02x%02x%02x", r/0x101, g/0x101, b/0x0101)
			}
		}
		fmt.Println(">} false 3 colorimage showpage")
		return
	}

	fmt.Printf("/%s%s %d selectfont\n", font, style, size)

	base := height / 2
	count := len(lines)
	if count&1 != 1 {
		base -= size / 2
	}
	for i, line := range lines {
		if line == "" {
			continue
		}

		y, x := base-size*(i-count/2), width/20
		if center {
			fmt.Println("gsave 0 -1000 moveto ")
			printLine(line)
			fmt.Println("currentpoint pop 400 exch sub 2 div grestore")
			fmt.Printf("%d moveto ", y)
		} else {
			if indent {
				x = width / 10
			}
			fmt.Printf("%d %d moveto ", x, y)
		}
		printLine(line)
	}
	fmt.Println("showpage")
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
			opt, ok := opts[strings.ToLower(sub[1])]
			if !ok {
				fmt.Fprintf(os.Stderr, "unknown command %q in line %d\n",
					sub[1], linum)
				break
			}
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
