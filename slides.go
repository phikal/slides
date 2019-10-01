package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

type Font struct{ font string }

func (r *Font) Set(val string) {
	switch val {
	case "serif":
		font = "Times"
	case "mono":
		font = "Courier"
	default:
		font = "Helvetica"
	}
}
func (r *Font) Reset() { font = r.font }
func (r *Font) Push()  { r.font = font }

type Style struct{ style string }

func (r *Style) Set(val string) {
	switch val {
	case "bold":
		style = "Bold"
	case "italics":
		style = "Italics"
	default:
		style = ""
	}
}
func (r *Style) Reset() { style = r.style }
func (r *Style) Push()  { r.style = style }

type Size struct{ size int }

func (r *Size) Set(val string) {
	switch val {
	case "huge":
		size = height / 6
	case "large":
		size = height / 10
	case "small":
		size = height / 25
	case "tiny":
		size = height / 35
	default:
		size = height / 15
	}
}
func (r *Size) Reset() { size = r.size }
func (r *Size) Push()  { r.size = size }

type Indent struct{ indent bool }

func (r *Indent) Set(val string) { indent = (val != "") }
func (r *Indent) Reset()         { indent = r.indent }
func (r *Indent) Push()          { r.indent = indent }

type Height struct{ height int }

func (r *Height) Set(val string) {
	i, err := strconv.Atoi(val)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid height value %q in line %d\n", val, linum)
		return
	}
	height = i
}
func (r *Height) Reset() { height = r.height }
func (r *Height) Push()  { r.height = height }

type Width struct{ width int }

func (r *Width) Set(val string) {
	i, err := strconv.Atoi(val)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid height value %q in line %d\n", val, linum)
		return
	}
	width = i
}
func (r *Width) Reset() { width = r.width }
func (r *Width) Push()  { r.width = width }

type Image struct{}

func (r *Image) Set(val string) {
	scale := fmt.Sprintf("%dx%d>", (width-width/8)*14, (height-height/8)*14)
	f, err := ioutil.TempFile("", "slides-*.eps")
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't create temporary file: ", err)
		return
	}
	cmd1 := exec.Command("convert", val, "-resize", scale, f.Name())
	err = cmd1.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while converting file: ", err)
		return
	}

	cmd2 := exec.Command("identify", "-format", "%w %h", f.Name())
	output, err := cmd2.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading `identify' output: ", err)
		return
	}
	fmt.Sscanf(string(output), "%d %d", &iwidth, &iheight)

	images = append(images, f.Name())
	image = true
}
func (r *Image) Reset() { image = false }
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
	image   = false
	iwidth  = 0
	iheight = 0

	lines   []string
	images  []string
	linum   = 1
	newpage = true
	opts    = map[string]Option{
		"font":   &Font{font},
		"style":  &Style{style},
		"size":   &Size{size},
		"indent": &Indent{indent},
		"height": &Height{height},
		"width":  &Width{width},
		"image":  &Image{},
	}
)

func printLine(line string) {
	escaped := strings.NewReplacer("(", "\\(", ")", "\\)").Replace(line)
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
	if !newpage || lines == nil {
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

	if image {
		fmt.Printf("%d %d translate\n", (width-iwidth)/2, (height-iheight)/2)
		fmt.Printf("(%s) run\n", images[0])
		return
	}

	fmt.Printf("/%s%s %d selectfont\n", font, style, size)

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
