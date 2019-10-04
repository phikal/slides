package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var opts = map[string]Option{
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
	"fill":    &Bool{&fill, fill},
	"image":   &Image{},
	"title": &Aggregated{
		"center": "t",
		"style":  "bold",
		"size":   "large",
	},
}

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
	*s.ref = strings.ToLower(val)
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
