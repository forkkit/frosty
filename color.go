package main

import (
	"fmt"
	"strconv"
	"strings"
)

// A Color is some color with unbounded intensity in each channel.
// This is unlike, for example, color.RGBA and friends, which are traditional
// display colors with some fixed range per channel. Color satisfies the
// color.Color interface by clamping the range to [0, 1] and scaling
// appropriately. The entire image should be scaled together (tone mapped?)
// before rendering to an image.Image.
type Color struct {
	R, G, B float64
}

var (
	Black  = Color{0, 0, 0}
	Pink   = Color{1, 0, 0.5}
	Yellow = Color{0.5, 0.5, 0}
)

// Add adds two colors
func (c1 Color) Add(c2 Color) Color {
	return Color{
		R: c1.R + c2.R,
		G: c1.G + c2.G,
		B: c1.B + c2.B,
	}
}

// Mul multiplies two colors component-wise.
func (c1 Color) Mul(c2 Color) Color {
	return Color{
		R: c1.R * c2.R,
		G: c1.G * c2.G,
		B: c1.B * c2.B,
	}
}

// MulS multiplies a color by a scalar.
func (c1 Color) MulS(f float64) Color {
	return Color{
		R: c1.R * f,
		G: c1.G * f,
		B: c1.B * f,
	}
}

func clamp(f float64) float64 {
	switch {
	case f < 0:
		return 0
	case f > 1:
		return 1
	}
	return f
}

func (c *Color) RGBA() (r, g, b, a uint32) {
	r = uint32(clamp(c.R) * 0xFF)
	g = uint32(clamp(c.G) * 0xFF)
	b = uint32(clamp(c.B) * 0xFF)
	a = 1
	return
}

func (c *Color) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("Bad color string: %s", string(b))
	}
	if b[0] == '"' {
		b = b[1:]
	}
	if b[len(b)-1] == '"' {
		b = b[:len(b)-1]
	}
	col, err := ParseHexColor(string(b))
	if err != nil {
		return err
	}
	*c = col
	return nil
}

// ParseHexColor parses a CSS-style hex color (e.g. #123abc, #fff) as a Color
// with values in [0, 1].
func ParseHexColor(c string) (Color, error) {
	if strings.HasPrefix(c, "#") {
		c = c[1:]
	}
	if len(c) == 3 {
		c = c[:1] + c[:1] + c[1:2] + c[1:2] + c[2:] + c[2:]
	}
	rgba := Color{}
	if len(c) == 6 {
		rgb, err := strconv.ParseUint(c, 16, 32)
		if err != nil {
			return rgba, err
		}
		rgba.R = float64(rgb>>16) / 0xFF
		rgba.G = float64((rgb>>8)&0xFF) / 0xFF
		rgba.B = float64(rgb&0xFF) / 0xFF
		return rgba, nil
	}
	return rgba, fmt.Errorf("Bad color string: %s", c)
}
