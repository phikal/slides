package main

import (
	"fmt"
	"strings"
)

var (
	font    = "Helvetica"
	style   = ""
	size    = 20
	indent  = false
	height  = 300
	width   = 400
	padding = 0
	fill    = true
	center  = false
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
	if img == nil && lines == nil {
		return
	}
	defer func() {
		for _, opt := range opts {
			opt.Reset()
		}
		lines = nil
	}()

	fmt.Printf("/width %d def\n", width)
	fmt.Printf("/height %d def\n", height)
	fmt.Println("newpage")

	if img != nil {
		primtImage()
		return
	}

	fmt.Printf("/%s%s %d selectfont\n", font, style, size)
	base := height/2 - size/4
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
			fmt.Println("0 -10000 moveto ")
			printLine(line)
			fmt.Println("currentpoint pop 400 exch sub 2 div")
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
