package uc

import "math"

// RGBA returns a Color constructed from red (r), green (g), blue (b) and alpha (a) integer channels, each divided
// by 255.0 for their final value. Using [RGBA] instead of the standard [New] constructor is useful when
// you need to match exact color values in an Image.
//
//	var red = RGBA8(255, 0, 0, 255)          // Same as NewColor(1, 0, 0, 1).
//	var dark_blue = RGBA8(0, 0, 51, 255)    // Same as NewColor(0, 0, 0.2, 1).
//	var my_color = RGBA8(306, 255, 0, 102)  // Same as NewColor(1.2, 1, 0, 0.4).
//
// Note: Due to the lower precision of [uint8] compared to the standard [float64] Color constructor, a color
// created with [RGBA] will generally not be equal to the same color created with the standard Color constructor.
// Use [Color.IsAproximatelyEqual] for comparisons to avoid issues with floating-point precision error.
func RGBA(r, g, b uint8) Color { // Color8
	return NewColor(
		float64(r)/255,
		float64(g)/255,
		float64(b)/255,
		1,
	)
}

// RGB returns a Color constructed from red (r), green (g) and blue (b) integer channels, each divided by 255.0
// for their final value. Using [RGB] instead of the standard [NewColor] constructor is useful when you need
// to match exact color values in an Image.
//
//	var red = RGB(255, 0, 0)       // Same as NewColor(1, 0, 0, 1).
//	var dark_blue =RGB(0, 0, 51)  // Same as NewColor(0, 0, 0.2, 1).
//
// Note: Due to the lower precision of [uint8] compared to the standard [float64] Color constructor, a color
// created with [RGB] will generally not be equal to the same color created with the standard Color constructor.
// Use [uc.Color.IsAproximatelyEqual] for comparisons to avoid issues with floating-point precision error.
func RGB(r, g, b uint8) Color { // Color8
	return NewColor(
		float64(r)/255,
		float64(g)/255,
		float64(b)/255,
		1,
	)
}

// RGBE9995 decodes a Color from a RGBE9995 format integer where the three color components have
// 9 bits of precision and all three share a single 5-bit exponent.
func RGBE9995(rgbe uint32) Color { // Color.from_rgbe9995
	var r = float64(rgbe & 0x1ff)
	var g = float64((rgbe >> 9) & 0x1ff)
	var b = float64((rgbe >> 18) & 0x1ff)
	var e = float64((rgbe >> 27))
	var m = math.Pow(2.0, e-15.0-9.0)
	var (
		rd = r * m
		gd = g * m
		bd = b * m
	)
	return NewColor(rd, gd, bd, 1.0)
}

// HSV constructs a color from an HSV profile. The hue (h), saturation (s), and value (v) are typically
// between 0.0 and 1.0.
func HSV(h, s, v float64) Color { // Color.from_hsv
	var (
		i, f, p, q, t float64
		a             float64 = 1
	)
	if s == 0.0 {
		return NewColor(v, v, v, a) // Achromatic (gray)
	}
	h *= 6.0
	h = math.Mod(h, 6)
	i = math.Floor(h)
	f = h - i
	p = v * (1.0 - s)
	q = v * (1.0 - s*f)
	t = v * (1.0 - s*(1.0-f))
	switch int(i) {
	case 0: // Red is the dominant color
		return NewColor(v, t, p, a)
	case 1: // Green is the dominant color
		return NewColor(q, v, p, a)
	case 2:
		return NewColor(p, v, t, a)
	case 3: // Blue is the dominant color
		return NewColor(p, q, v, a)
	case 4:
		return NewColor(t, p, v, a)
	default: // (5) Red is the dominant color
		return NewColor(v, p, q, a)
	}
}

// HSVA constructs a color from an HSV profile. The hue (h), saturation (s), and value (v) are typically
// between 0.0 and 1.0. Includes alpha.
func HSVA(h, s, v, a float64) Color { // Color.from_hsv
	c := HSV(h, s, v)
	c[3] = float32(a)
	return c
}

// Hex returns the Color associated with the provided hex integer in 32-bit RGBA format (8 bits per channel).
//
// The int is best visualized with hexadecimal notation ("0x" prefix, making it "0xRRGGBBAA").
//
//	var red = Hex(0xff0000ff)
//	var dark_cyan = Hex(0x008b8bff)
//	var my_color = Hex(0xbbefd2a4)
func Hex(hex uint32) Color { // Color.from_hex
	var a = float64(hex&0xFF) / 255
	hex >>= 8
	var b = float64(hex&0xFF) / 255
	hex >>= 8
	var g = float64(hex&0xFF) / 255
	hex >>= 8
	var r = float64(hex&0xFF) / 255
	return NewColor(r, g, b, a)
}

// Hex64 returns the Color associated with the provided hex integer in 64-bit RGBA format (16 bits per channel).
//
// The int is best visualized with hexadecimal notation ("0x" prefix, making it "0xRRRRGGGGBBBBAAAA").
func Hex64(hex int64) Color { // Color.from_hex64
	var a = float64(hex&0xFFFF) / 65535
	hex >>= 16
	var b = float64(hex&0xFFFF) / 65535
	hex >>= 16
	var g = float64(hex&0xFFFF) / 65535
	hex >>= 16
	var r = float64(hex&0xFFFF) / 65535
	return NewColor(r, g, b, a)
}

// HTML returns a new color from rgba, an HTML hexadecimal color string. rgba is not case-sensitive
// and may be prefixed by a hash sign (#).
//
// rgba must be a valid three-digit or six-digit hexadecimal color string, and may contain an alpha
// channel value. If rgba does not contain an alpha channel value, an alpha channel value of 1.0 is
// applied. If rgba is invalid, returns an empty color.
//
//	var blue = HTML("#0000ff") // blue is Color{0.0, 0.0, 1.0, 1.0}
//	var green = HTML("#0F0")   // green is Color{0.0, 1.0, 0.0, 1.0}
//	var col = HTML("663399cc") // col is Color{0.4, 0.2, 0.6, 0.8}
func HTML(rgba string) Color { // Color.from_html
	var color = rgba
	if len(color) == 0 {
		return Color{}
	}
	if color[0] == '#' {
		color = color[1:]
	}
	// If enabled, use 1 hex digit per channel instead of 2.
	// Other sizes aren't in the HTML/CSS spec but we could add them if desired.
	var (
		is_shorthand = len(color) < 5
		alpha        = false
	)
	if len(color) == 8 {
		alpha = true
	} else if len(color) == 6 {
		alpha = false
	} else if len(color) == 4 {
		alpha = true
	} else if len(color) == 3 {
		alpha = false
	}
	var r, g, b, a = 1.0, 1.0, 1.0, 1.0
	if is_shorthand {
		r = _parse_col4(color, 0) / 15
		g = _parse_col4(color, 1) / 15
		b = _parse_col4(color, 2) / 15
		if alpha {
			a = _parse_col4(color, 3) / 15
		}
	} else {
		r = _parse_col8(color, 0) / 255
		g = _parse_col8(color, 2) / 255
		b = _parse_col8(color, 4) / 255
		if alpha {
			a = _parse_col8(color, 6) / 255
		}
	}
	return NewColor(r, g, b, a)
}

func _parse_col4(s string, ofs int) float64 {
	var character = s[ofs]
	if character >= '0' && character <= '9' {
		return float64(character - '0')
	} else if character >= 'a' && character <= 'f' {
		return float64(character) + float64(10-'a')
	} else if character >= 'A' && character <= 'F' {
		return float64(character) + float64(10-'A')
	}
	return -1
}

func _parse_col8(s string, ofs int) float64 {
	return _parse_col4(s, ofs)*16 + _parse_col4(s, ofs+1)
}

// ValidHTML returns true if color is a valid HTML hexadecimal color string. The string must be a hexadecimal
// value (case-insensitive) of either 3, 4, 6 or 8 digits, and may be prefixed by a hash sign (#).
//
//	ValidHTML("#55aaFF")   // Returns true
//	ValidHTML("#55AAFF20") // Returns true
//	ValidHTML("55AAFF")    // Returns true
//	ValidHTML("#F2C")      // Returns true
//
//	ValidHTML("#AABBC")    // Returns false
//	ValidHTML("#55aaFF5")  // Returns false
func ValidHTML(color string) bool { // Color.html_is_valid
	if len(color) == 0 {
		return false
	}
	if color[0] == '#' {
		color = color[1:]
	}
	// Check if the amount of hex digits is valid.
	if !(len(color) == 3 || len(color) == 4 || len(color) == 6 || len(color) == 8) {
		return false
	}
	// Check if each hex digit is valid.
	for i := 0; i < len(color); i++ {
		if _parse_col4(color, i) == -1 {
			return false
		}
	}
	return true
}
