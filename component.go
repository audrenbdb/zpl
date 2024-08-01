package zpl

import (
	"fmt"
	"image"
	"strconv"
	"strings"
)

type Component fmt.Stringer

type BarCodeConfig struct {
	// Width of the bar code module, in dots.
	// Any number between 1 and 100 may be used.
	//
	// The default value is 2.
	Width int
	// WidthRatio between wide bars and narrow bars. Any decimal number between 2 and 3 may be used.
	//
	// The number must be a multiple of 0.1 (i.e. 2.0, 2.1, 2.2, 2.3, ... , 2.9, 3.0).
	//
	// Larger numbers generally result in fewer bar code scan failures.
	// The default value is 3.
	WidthRatio int
	// The default bar code height, in dots.
	// Any positive number may be used.
	//
	// The default value is 10.
	Height int
}

func (bc BarCodeConfig) String() string {
	var sb strings.Builder

	sb.WriteString("^BY")
	sb.WriteString(strconv.Itoa(bc.Width))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(bc.WidthRatio))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(bc.Height))

	return sb.String()
}

func NewBarCodeConfig() BarCodeConfig {
	return BarCodeConfig{
		Width:      2,
		WidthRatio: 3,
		Height:     10,
	}
}

func (x BarCodeConfig) WithWidth(width int) BarCodeConfig {
	x.Width = width
	return x
}

func (x BarCodeConfig) WithWidthRatio(widthRatio int) BarCodeConfig {
	x.WidthRatio = widthRatio
	return x
}

func (x BarCodeConfig) WithHeight(height int) BarCodeConfig {
	x.Height = height
	return x
}

// Coordinates of top left corner of the current field.
type Coordinates struct {
	// X represents the field position x-coordinate, in dots.
	X int
	// Y represents the field position y-coordinate, in dots.
	Y int
}

func (c Coordinates) String() string {
	var sb strings.Builder

	sb.WriteString("^FO")
	sb.WriteString(strconv.Itoa(c.X))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(c.Y))

	return sb.String()
}

type Font struct {
	Height int
	Width  int
}

func (f Font) String() string {
	var sb strings.Builder

	sb.WriteString("^A0N,")
	sb.WriteString(strconv.Itoa(f.Height))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(f.Width))

	return sb.String()
}

func NewFont(height, width int) Font {
	return Font{
		Height: height,
		Width:  width,
	}
}

// BarCode128 builds a field as a Code 128 bar code.
type BarCode128 struct {
	Coordinates
	// Code is the code to display as an EAN 128 bar code.
	Code   string
	Height int
}

func NewBarCode(x, y int, code string) BarCode128 {
	return BarCode128{
		Code: code,
		Coordinates: Coordinates{
			X: x,
			Y: y,
		},
	}
}

func (bc BarCode128) WithHeight(height int) BarCode128 {
	bc.Height = height
	return bc
}

func (bc BarCode128) String() string {
	var sb strings.Builder

	// start field
	sb.WriteString(bc.Coordinates.String())

	// bar code
	sb.WriteString("^BC,")

	if bc.Height > 0 {
		sb.WriteString(strconv.Itoa(bc.Height))
	}

	// mode
	sb.WriteString(",,,,A")

	sb.WriteString(`^FH\^FD`)
	sb.WriteString(bc.Code)
	sb.WriteString("^FS")

	return sb.String()
}

// Line builds a field to represent a text.
type Line struct {
	Text string
	Bold bool
	Coordinates
	Font     Font
	Reversed bool
}

func (l Line) WithBold() Line {
	l.Bold = true
	return l
}

func (l Line) WithFontSize(size int) Line {
	l.Font.Height = size
	l.Font.Width = size * 90 / 100

	return l
}

func (l Line) WithReversed() Line {
	l.Reversed = true
	return l
}

func NewLine(x, y int, text string) Line {
	return Line{
		Text: text,
		Coordinates: Coordinates{
			X: x,
			Y: y,
		},
		Font: Font{
			Height: 14,
			Width:  14,
		},
	}
}

func (l Line) String() string {
	var sb strings.Builder

	if l.Bold {
		l.Bold = false

		l.Coordinates.X++
		sb.WriteString(l.String())

		l.Coordinates.X++
		sb.WriteString(l.String())
	}

	// start field
	sb.WriteString(l.Coordinates.String())
	sb.WriteString(l.Font.String())

	if l.Reversed {
		sb.WriteString("^FR")
	}

	// text
	sb.WriteString(`^FH\^FD`)
	sb.WriteString(escape(l.Text))

	// end field
	sb.WriteString("^FS")

	return sb.String()
}

type TextBlockAlignment string

const (
	TextBlockAlignmentLeft      = "L"
	TextBlockAlignmentCenter    = "C"
	TextBlockAlignmentRight     = "R"
	TextBlockAlignmentJustified = "J"
)

type TextBlock struct {
	Text string
	Font
	Coordinates

	Width       int
	MaxLines    int
	LineSpacing int

	Reversed  bool
	Alignment TextBlockAlignment
}

func NewTextBlock(x, y, width int, text string) TextBlock {
	return TextBlock{
		Text: text,
		Coordinates: Coordinates{
			X: x,
			Y: y,
		},
		Font: Font{
			Height: 14,
			Width:  14,
		},
		Width:       width,
		MaxLines:    5,
		LineSpacing: 0,
		Alignment:   TextBlockAlignmentLeft,
	}
}

func (tb TextBlock) WithMaxLines(maxLines int) TextBlock {
	tb.MaxLines = maxLines
	return tb
}

func (tb TextBlock) WithLineSpacing(lineSpacing int) TextBlock {
	tb.LineSpacing = lineSpacing
	return tb
}

func (tb TextBlock) WithReversed(reversed bool) TextBlock {
	tb.Reversed = reversed
	return tb
}

func (tb TextBlock) WithFontSize(size int) TextBlock {
	tb.Font.Height = size
	tb.Font.Width = size

	return tb
}

func (tb TextBlock) WithAlignment(alignment TextBlockAlignment) TextBlock {
	tb.Alignment = alignment
	return tb
}

func (tb TextBlock) String() string {
	var sb strings.Builder

	// start field
	sb.WriteString(tb.Coordinates.String())
	sb.WriteString(tb.Font.String())

	// text
	sb.WriteString("^FB")
	sb.WriteString(strconv.Itoa(tb.Width))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(tb.MaxLines))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(tb.LineSpacing))
	sb.WriteString("," + string(tb.Alignment))

	if tb.Reversed {
		sb.WriteString("^FR")
	}

	sb.WriteString(`^FH\^FD`)
	sb.WriteString(escape(tb.Text))

	// end field
	sb.WriteString("^FS")

	return sb.String()
}

type ImageField struct {
	Image image.Image
	Coordinates
}

func NewImageField(x, y int, img image.Image) ImageField {
	return ImageField{
		Image: img,
		Coordinates: Coordinates{
			X: x,
			Y: y,
		},
	}
}

func (b ImageField) String() string {
	var sb strings.Builder

	sb.WriteString(b.Coordinates.String())
	sb.WriteString(convertFromImage(b.Image))

	// end field
	sb.WriteString("^FS")

	return sb.String()
}

type Square struct {
	Coordinates
	Width              int
	Height             int
	Thickness          int
	TexturedBackground bool
}

func NewVerticalLine(x, y, height, thickness int) Square {
	return Square{
		Coordinates: Coordinates{
			X: x,
			Y: y,
		},
		Width:     5,
		Height:    height,
		Thickness: thickness,
	}
}

func NewHorizontalLine(x, y, width, thickness int) Square {
	return Square{
		Coordinates: Coordinates{
			X: x,
			Y: y,
		},
		Width:     width,
		Height:    5,
		Thickness: thickness,
	}
}

func NewSquare(x, y, width, height int) Square {
	return Square{
		Coordinates: Coordinates{
			X: x,
			Y: y,
		},
		Width: width,
		// 5 is the minimum value for this component.
		Height:    max(height, 5),
		Thickness: 5,
	}
}

func (s Square) WithTexturedBackground() Square {
	s.TexturedBackground = true
	return s
}

func (s Square) WithPlainBackground() Square {
	s.TexturedBackground = false
	s.Thickness = s.Width

	if s.Width > s.Height {
		s.Thickness = s.Height
	}

	return s
}

func (s Square) String() string {
	var sb strings.Builder

	sb.WriteString(s.Coordinates.String())

	sb.WriteString("^GB")
	sb.WriteString(strconv.Itoa(s.Width))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(s.Height))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(s.Thickness))
	sb.WriteString("^FS")

	if s.TexturedBackground {
		dotSize := 3
		dotPadding := 14

		// how many dots can fit in the square.
		widthCapacity := (s.Width - s.Thickness*2) / (dotSize + dotPadding*2)
		heightCapacity := (s.Height - s.Thickness*2) / (dotSize + dotPadding*2)

		// calculate initial dot position based of the capacity so that everything is centered.
		startX := s.X + (s.Width-widthCapacity*(dotSize+dotPadding*2))/2
		startY := s.Y + (s.Height-heightCapacity*(dotSize+dotPadding*2))/2

		endX := startX + widthCapacity*(dotSize+dotPadding*2)
		endY := startY + heightCapacity*(dotSize+dotPadding*2)

		for x := startX; x < endX; x += dotSize + dotPadding/2 {
			for y := startY; y < endY; y += dotSize + dotPadding/2 {
				sb.WriteString("^FO")
				sb.WriteString(strconv.Itoa(x))
				sb.WriteString(",")
				sb.WriteString(strconv.Itoa(y))
				sb.WriteString("^GB")
				sb.WriteString(strconv.Itoa(dotSize))
				sb.WriteString(",")
				sb.WriteString(strconv.Itoa(dotSize))
				sb.WriteString(",")
				sb.WriteString(strconv.Itoa(dotSize))
				sb.WriteString("^FS")
			}
		}
	}

	return sb.String()
}

type QRBarCode struct {
	Coordinates
	URI           string
	Magnification int
}

func NewQRBarCode(x, y int, uri string) QRBarCode {
	return QRBarCode{
		URI: uri,
		Coordinates: Coordinates{
			X: x,
			Y: y,
		},
		Magnification: 2,
	}
}

func (q QRBarCode) WithMagnification(magnification int) QRBarCode {
	q.Magnification = magnification
	return q
}

func (q QRBarCode) String() string {
	var sb strings.Builder

	sb.WriteString(q.Coordinates.String())

	sb.WriteString("^BQN,2,")
	sb.WriteString(strconv.Itoa(q.Magnification))
	sb.WriteString(`^FH\^FDLA,`)
	sb.WriteString(q.URI)
	sb.WriteString("^FS")

	return sb.String()
}

// escape takes any characters that ZPL-reserved, such as ~
// and replace it with the HEX representation.
func escape(in string) string {
	out := strings.ReplaceAll(in, `\`, `\1F`)
	out = strings.ReplaceAll(in, "~", `\7E`)

	return strings.ReplaceAll(out, "^", `\5E`)
}

// Raw is a flexible type where the consumer provide valid ZPL string.
type Raw string

func (x Raw) String() string {
	return string(x)
}

func NewRaw(raw string) Raw {
	return Raw(raw)
}
