package main

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

type HSL struct {
	H, S, L float64
}

func main() {
}

func normalize(color string) (string, error) {

	// Remove leading '#'
	color = strings.TrimPrefix(color, "#")

	// Converting the passed hex to uppercase
	color = strings.ToUpper(color)

	i := len(color)
	if i == 8 {
		return color, nil
	}
	var buffer bytes.Buffer

	pad := func() {
		for _, i := range color {
			str := fmt.Sprintf("%c", i)
			buffer.WriteString(strings.Repeat(str, 2))
		}
	}

	prepend := func() {
		buffer.WriteString("FF")
	}
	switch i {
	case 3:
		prepend()
		pad()
	case 4:
		pad()
	case 6:
		prepend()
		buffer.WriteString(color)
	}

	str := buffer.String()
	if str == "" {
		return "", fmt.Errorf("#%v appears to be an invalid colorStr\n", color)
	}
	return str, nil
}

func rgbToHsl(rgba color.RGBA) HSL {
	r, g, b := float64(rgba.R), float64(rgba.G), float64(rgba.B)
	r /= 255
	g /= 255
	b /= 255
	min := math.Min(r, math.Min(g, b))
	max := math.Max(r, math.Max(g, b))
	delta := max - min

	l := (min + max) / 2

	var s float64
	if max != min {
		var divisor float64
		if l <= 0.5 {
			divisor = max + min
		} else {
			divisor = 2 - max - min
		}
		s = delta / divisor
	}

	var h float64

	if delta != 0 {
		var segment float64
		var shift float64
		switch max {
		case r:
			segment = (g - b) / delta
			if segment < 0 {
				shift = 360 / 60
			} else {
				shift = 0 / 60
			}
			break
		case g:
			segment = (b - r) / delta
			shift = 120 / 60
		case b:
			segment = (r - g) / delta
			shift = 240 / 60
		}
		h = segment + shift
	}
	return HSL{
		h * 60,
		s * 100,
		l * 100,
	}

}

func strToRGBA(str string) (color.RGBA, error) {
	rStr := fmt.Sprintf("0x%v", str[0:2])
	gStr := fmt.Sprintf("0x%v", str[2:4])
	bStr := fmt.Sprintf("0x%v", str[4:])

	r, err := strconv.ParseUint(rStr, 0, 8)
	if err != nil {
		return color.RGBA{}, err
	}

	g, err := strconv.ParseUint(gStr, 0, 8)
	if err != nil {
		return color.RGBA{}, err
	}

	b, err := strconv.ParseUint(bStr, 0, 8)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
	}, nil
}

func name(str string, rgb color.RGBA) (item, error) {
	var hsl = rgbToHsl(rgb)
	var h, s, l = hsl.H * 255, hsl.S * 255, hsl.L * 255
	var ndf, ndf1, ndf2 float64
	var cl = -1
	var df float64 = -1
	for i, v := range colorItems {
		if v.color == str {
			return v, nil
		}

		rbg2, _ := strToRGBA(v.color)
		hsl2 := rgbToHsl(rbg2)
		var h2, s2, l2 = hsl2.H * 255, hsl2.S * 255, hsl2.L * 255

		ndf1 = math.Pow(float64(rgb.R-rbg2.R), 2) +
			math.Pow(float64(rgb.G-rbg2.G), 2) +
			math.Pow(float64(rgb.B-rbg2.B), 2)

		ndf2 = math.Pow(h-h2, 2) +
			math.Pow(s-s2, 2) +
			math.Pow(l-l2, 2)

		ndf = ndf1 + (ndf2 * 2)
		if df < 0 || df > ndf {
			df = ndf
			cl = i
		}
	}

	if cl < 0 {
		return item{}, fmt.Errorf("#%s is an invalid color", str)
	}

	return colorItems[cl], nil
}
