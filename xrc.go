package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"os"
	"text/template"
)

var tfile = flag.String("t", "", "path to template file")
var pfile = flag.String("p", "", "path to palette file")

func errh(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "xrc: %s\n", err)
		os.Exit(1)
	}
}

var colors = [9]struct{ hi, lo string }{
	{"black", "lblack"},
	{"white", "lwhite"},
	{"red", "lred"},
	{"green", "lgreen"},
	{"yellow", "lyellow"},
	{"blue", "lblue"},
	{"magenta", "lmagenta"},
	{"cyan", "lcyan"},
	{"background", "foreground"},
}

var c256 = [...]struct{ hi, lo string }{
	{"black", "lblack"},
	{"white", "lwhite"},
	{"red", "lred"},
	{"green", "lgreen"},
	{"yellow", "lyellow"},
	{"blue", "lblue"},
	{"magenta", "lmagenta"},
	{"cyan", "lcyan"},
	{"background", "foreground"},
}

type co struct{ color.Color }

func hex(c co) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

func (c co) String() string {
	return "#" + hex(c)
}

func x(c co) string {
	return "0x" + hex(c)
}

var funcs = map[string]interface{}{
	"r": hex,
	"x": x,
}

func main() {
	flag.Parse()
	if *pfile == "" {
		errh(errors.New("missing palette file"))
	}
	if *tfile == "" {
		errh(errors.New("missing template file"))
	}
	rd, err := os.Open(*pfile)
	errh(err)
	img, _, err := image.Decode(rd)
	errh(err)
	b, err := ioutil.ReadFile(*tfile)
	errh(err)
	t, err := template.New("aaaa").Funcs(funcs).Parse(string(b))
	errh(err)
	t = t.Funcs(funcs)
	var colorMap = make(map[string]co)
	for i, c := range colors {
		hi, lo := co{img.At(i*2+1, 1)}, co{img.At(i*2+1, 3)}
		colorMap[c.hi] = hi
		colorMap[c.lo] = lo
	}
	errh(t.Execute(os.Stdout, colorMap))
}