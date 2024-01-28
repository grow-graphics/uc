// Package uc provides a useful color representation suitable for blending and graphics.
package uc

import "math"

/*
Color represented in RGBA format by a red (r), green (g), blue (b), and alpha (a) component. Each component is a
32-bit floating-point value, usually ranging from 0.0 to 1.0. May support values greater than 1.0, for overbright
or HDR (High Dynamic Range) colors.
*/
type Color [4]float32

// NewColor constructs a Color from RGBA values, typically between 0.0 and 1.0.
//
//	var color = RGBA(0.2, 1.0, 0.7, 0.8) // Similar to `RGBA8(51, 255, 178, 204)`
func NewColor(r, g, b, a float64) Color {
	return Color{float32(r), float32(g), float32(b), float32(a)}
}

const (
	r = iota
	g
	b
	a
)

// ʕ is a little ternary operator for porting C code.
func ʕ[T any](q bool, a T, b T) T {
	if q {
		return a
	}
	return b
}

// Blend returns a new color resulting from overlaying this color over the given color. In a painting program,
// you can imagine it as the over color painted over this color (including alpha).
//
//	var bg = Color{0,1,0,0.5} // Green with alpha of 50%
//	var fg = Color{1,0,0,0.5} // Red with alpha of 50%
//	var blended_color = bg.Blend(fg) // Brown with alpha of 75%
func (c Color) Blend(over Color) Color { // Color.blend
	var res Color
	var sa = 1.0 - over[a]
	res[a] = a*sa + over[a]
	if res[a] == 0 {
		return Color{}
	} else {
		res[r] = (r*a*sa + over[r]*over[a]) / res[a]
		res[g] = (g*a*sa + over[g]*over[a]) / res[a]
		res[b] = (b*a*sa + over[b]*over[a]) / res[a]
	}
	return res
}

// Clamp returns a new color with all components clamped between the components of minimum and maximum.
func (c Color) Clamp(minimum, maximum Color) Color { // Color.clamp
	return Color{
		max(minimum[r], min(maximum[r], c[r])),
		max(minimum[g], min(maximum[g], c[g])),
		max(minimum[b], min(maximum[b], c[b])),
		max(minimum[a], min(maximum[a], c[a])),
	}
}

// Darkened returns a new color resulting from making this color darker by the specified amount
// (ratio from 0.0 to 1.0). See also [Color.Lightened].
//
//	var green = Color{0,1,0,1}
//	var darkgreen = green.Darkened(0.2) // 20% darker than regular green
func (c Color) Darkened(amount float64) Color { // Color.darkened
	return Color{
		float32(float64(c[r]) * (1.0 - amount)),
		float32(float64(c[g]) * (1.0 - amount)),
		float32(float64(c[b]) * (1.0 - amount)),
		c[a],
	}
}

// Luminance returns the light intensity of the color, as a value between 0.0 and 1.0 (inclusive).
// This is useful when determining light or dark color. Colors with a luminance smaller than 0.5
// can be generally considered dark.
//
// Note: [Color.Luminance] relies on the color being in the linear color space to return an accurate
// relative luminance value. If the color is in the sRGB color space, use srgb_to_linear to convert
// it to the linear color space first.
func (c Color) Luminance() float64 { // Color.luminance
	return 0.2126*float64(c[r]) + 0.7152*float64(c[g]) + 0.0722*float64(c[b])
}

// Inverted returns the color with its r, g, and b components inverted ((1 - r, 1 - g, 1 - b, a)).
//
//	var black = X11.White.Inverted()
//	var color = Color{0.3, 0.4, 0.9, 1.0}
//	var inverted_color = color.Inverted() // Equivalent to `Color{(0.7, 0.6, 0.1, 1.0}`
func (c Color) Inverted() Color { return Color{1 - c[r], 1 - c[g], 1 - c[b], c[a]} } // Color.inverted

// IsApproximatelyEqual returns true if this color and to are approximately equal.
func (c Color) IsApproximatelyEqual(to Color) bool { // Color.is_equal_approx
	return isApproximatelyEqual(c[r], to[r]) &&
		isApproximatelyEqual(c[g], to[g]) &&
		isApproximatelyEqual(c[b], to[b]) &&
		isApproximatelyEqual(c[a], to[a])
}

func isApproximatelyEqual(a, b float32) bool {
	const cmpEpsilon = 0.00001
	// Check for exact equality first, required to handle "infinity" values.
	if a == b {
		return true
	}
	// Then check for approximate equality.
	tolerance := cmpEpsilon * float32(math.Abs(float64(a)))
	if tolerance < cmpEpsilon {
		tolerance = cmpEpsilon
	}
	return float32(math.Abs(float64(a-b))) < tolerance
}

// Lerp returns the linear interpolation between this color's components and to's components.
// The interpolation factor weight should be between 0.0 and 1.0 (inclusive).
//
//	var red = Color{1.0, 0.0, 0.0, 1.0}
//	var aqua = Color{0.0, 1.0, 0.8, 1.0}
//
//	red.Lerp(aqua, 0.2) // Returns Color{0.8, 0.2, 0.16, 1}
//	red.Lerp(aqua, 0.5) // Returns Color{0.5, 0.5, 0.4, 1}
//	red.Lerp(aqua, 1.0) // Returns Color{0.0, 1.0, 0.8, 1}
func (c Color) Lerp(to Color, weight float64) Color { // Color.lerp
	return Color{
		float32(lerpf(float64(c[r]), float64(to[r]), weight)),
		float32(lerpf(float64(c[g]), float64(to[g]), weight)),
		float32(lerpf(float64(c[b]), float64(to[b]), weight)),
		float32(lerpf(float64(c[a]), float64(to[a]), weight)),
	}
}

func lerpf(from, to, weight float64) float64 { return from + (to-from)*weight }

// Lightened returns a new color resulting from making this color lighter by the specified amount
// which should be a ratio from 0.0 to 1.0. See also [Color.Darkened].
//
//	var green = Color{0.0, 1.0, 0.0, 1.0}
//	var light_green = green.Lightened(0.2) // 20% lighter than regular green
func (c Color) Lightened(amount float64) Color { // Color.lightened
	return Color{
		float32(float64(c[r]) + (1.0-float64(c[r]))*amount),
		float32(float64(c[g]) + (1.0-float64(c[g]))*amount),
		float32(float64(c[b]) + (1.0-float64(c[b]))*amount),
		c[a],
	}
}

// SRGB returns the color converted to the sRGB color space. This method assumes the original color is
// in the linear color space. See also [Color.Linear] which performs the opposite operation.
func (c Color) SRGB() Color { // Color.linear_to_srgb
	return Color{
		float32(ʕ(r < 0.0031308, 12.92*r, (1.0+0.055)*math.Pow(r, 1.0/2.4)-0.055)),
		float32(ʕ(g < 0.0031308, 12.92*g, (1.0+0.055)*math.Pow(g, 1.0/2.4)-0.055)),
		float32(ʕ(b < 0.0031308, 12.92*b, (1.0+0.055)*math.Pow(b, 1.0/2.4)-0.055)),
		a,
	}
}

// Linear returns the color converted to the linear color space. This method assumes the original color
// already is in the sRGB color space. See also [Color.SRGB] which performs the opposite operation.
func (c Color) Linear() Color { // Color.srgb_to_linear
	return Color{
		float32(ʕ(r <= 0.04045, r/12.92, math.Pow((r+0.055)/1.055, 2.4))),
		float32(ʕ(g <= 0.04045, g/12.92, math.Pow((g+0.055)/1.055, 2.4))),
		float32(ʕ(b <= 0.04045, b/12.92, math.Pow((b+0.055)/1.055, 2.4))),
		a,
	}
}

// ABGR32 returns the color converted to a 32-bit integer in ABGR format (each component is 8 bits).
// ABGR is the reversed version of the default RGBA format.
//
//	var color = Color{1, 0.5, 0.2, 1}
//	print(color.ABGR32()) // Prints 4281565439
func (c Color) ABGR32() uint32 { // Color.to_abgr32
	var u32 = uint32(math.Round(float64(c[a]) * 255))
	u32 <<= 8
	u32 |= uint32(math.Round(float64(c[b]) * 255.0))
	u32 <<= 8
	u32 |= uint32(math.Round(float64(c[g]) * 255.0))
	u32 <<= 8
	u32 |= uint32(math.Round(float64(c[r]) * 255.0))
	return u32
}

// ABGR64 returns the color converted to a 64-bit integer in ABGR format (each component is 16 bits).
// ABGR is the reversed version of the default RGBA format.
func (c Color) ABGR64() int64 { // Color.to_abgr64
	var u64 = int64(math.Round(float64(c[a]) * 65535))
	u64 <<= 16
	u64 |= int64(math.Round(float64(c[b]) * 65535))
	u64 <<= 16
	u64 |= int64(math.Round(float64(c[g]) * 65535))
	u64 <<= 16
	u64 |= int64(math.Round(float64(c[r]) * 65535))
	return u64
}

// ARGB32 returns the color converted to a 32-bit integer in ARGB format (each component is 8 bits).
// ARGB is more compatible with DirectX.
//
//	var color = Color{1, 0.5, 0.2, 1}
//	print(color.ARGB32()) // Prints 4294934323
func (c Color) ARGB32() uint32 { // Color.to_argb32
	var u32 = uint32(math.Round(float64(c[a]) * 255))
	u32 <<= 8
	u32 |= uint32(math.Round(float64(c[r]) * 255.0))
	u32 <<= 8
	u32 |= uint32(math.Round(float64(c[g]) * 255.0))
	u32 <<= 8
	u32 |= uint32(math.Round(float64(c[b]) * 255.0))
	return u32
}

// ARGB64 returns the color converted to a 64-bit integer in ARGB format (each component is 16 bits).
// ARGB is more compatible with DirectX.
//
//	var color = Color{1, 0.5, 0.2, 1}
//	print(color.ARGB64()) // Prints -2147470541
func (c Color) ARGB64() int64 { // Color.to_argb64
	var u64 = int64(math.Round(float64(c[a]) * 65535))
	u64 <<= 16
	u64 |= int64(math.Round(float64(c[r]) * 65535))
	u64 <<= 16
	u64 |= int64(math.Round(float64(c[g]) * 65535))
	u64 <<= 16
	u64 |= int64(math.Round(float64(c[b]) * 65535))
	return u64
}

// HTML returns the color converted to an HTML hexadecimal color String in RGBA format
// without the hash (#) prefix.
//
// Setting with_alpha to false, excludes alpha from the hexadecimal string, using RGB
// format instead of RGBA format.
//
//	var white = Color{1, 1, 1, 0.5}
//	var with_alpha = white.HTML(true) // Returns "ffffff7f"
//	var without_alpha = white.HTML(false) // Returns "ffffff"
func (c Color) HTML(with_alpha bool) string { // Color.to_html
	var txt string
	txt += _to_hex(c[r])
	txt += _to_hex(c[g])
	txt += _to_hex(c[b])
	if with_alpha {
		txt += _to_hex(c[a])
	}
	return txt
}

func _to_hex(val float32) string {
	v := rune(min(255, max(0, math.Round(float64(val*255)))))
	var ret string
	for i := 0; i < 2; i++ {
		var c = [2]rune{0, 0}
		var lv = v & 0xF
		if lv < 10 {
			c[0] = '0' + lv
		} else {
			c[0] = 'a' + lv - 10
		}
		v >>= 4
		var cs = string(c[:])
		ret = cs + ret
	}
	return ret
}

// RGBA32 returns the color converted to a 32-bit integer in RGBA format
// (each component is 8 bits).
//
//	var color = Color{1, 0.5, 0.2, 1.0}
//	print(color.RGBA32()) // Prints 4286526463
func (c Color) RGBA32() uint32 { // Color.to_rgba32
	var u32 = uint32(math.Round(float64(c[r]) * 255))
	u32 <<= 8
	u32 |= uint32(math.Round(float64(c[g]) * 255.0))
	u32 <<= 8
	u32 |= uint32(math.Round(float64(c[b]) * 255.0))
	u32 <<= 8
	u32 |= uint32(math.Round(float64(c[a]) * 255.0))
	return u32
}

// RGBA64 returns the color converted to a 64-bit integer in RGBA format
// (each component is 16 bits).
//
//	var color = Color{1, 0.5, 0.2, 1.0}
//	print(color.RGBA64()) // Prints -140736629309441
func (c Color) RGBA64() int64 { // Color.to_rgba64
	var u64 = int64(math.Round(float64(c[r]) * 65535))
	u64 <<= 16
	u64 |= int64(math.Round(float64(c[g]) * 65535))
	u64 <<= 16
	u64 |= int64(math.Round(float64(c[b]) * 65535))
	u64 <<= 16
	u64 |= int64(math.Round(float64(c[a]) * 65535))
	return u64
}

func (c Color) Add(o Color) Color {
	return Color{c[r] + o[r], c[g] + o[g], c[b] + o[b], c[a] + o[a]}
}
func (c Color) Sub(o Color) Color {
	return Color{c[r] - o[r], c[g] - o[g], c[b] - o[b], c[a] - o[a]}
}
func (c Color) Div(o Color) Color {
	return Color{c[r] / o[r], c[g] / o[g], c[b] / o[b], c[a] / o[a]}
}
func (c Color) Mul(o Color) Color {
	return Color{c[r] * o[r], c[g] * o[g], c[b] * o[b], c[a] * o[a]}
}
func (c Color) Neg() Color {
	return Color{-c[r], -c[g], -c[b], -c[a]}
}
func (c Color) Addf(f float64) Color {
	return Color{c[r] + float32(f), c[g] + float32(f), c[b] + float32(f), c[a] + float32(f)}
}
func (c Color) Subf(f float64) Color {
	return Color{c[r] - float32(f), c[g] - float32(f), c[b] - float32(f), c[a] - float32(f)}
}
func (c Color) Divf(f float64) Color {
	return Color{c[r] / float32(f), c[g] / float32(f), c[b] / float32(f), c[a] / float32(f)}
}
func (c Color) Mulf(f float64) Color {
	return Color{c[r] * float32(f), c[g] * float32(f), c[b] * float32(f), c[a] * float32(f)}
}
