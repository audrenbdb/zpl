package zpl

import (
	"strconv"
	"strings"
)

const baseTemplate = `
^XA
~SD%d
^LRN
^CI28
^MM%s
^PW%d
^LL%d
^LS%d
^PQ1,0,1,Y
^XZ
`

type Builder struct {
	// Configures the print darkness level. The darkness level can also be modified using the ^MD command.
	Darkness int
	// Sets the label print width.
	Width int
	// Shifts all label content to the left or the right.
	ShiftDistance int
	BarCodeConfig BarCodeConfig
	Components    []Component
}

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

func NewBuilder() *Builder {
	return &Builder{
		Darkness:      15,
		Width:         812,
		ShiftDistance: 0,
		BarCodeConfig: BarCodeConfig{
			Width:      2,
			WidthRatio: 3,
			Height:     10,
		},
	}
}

func (b *Builder) WithWidth(width int) *Builder {
	b.Width = width
	return b
}

func (b *Builder) WithDarkness(darkness int) *Builder {
	b.Darkness = darkness
	return b
}

func (b *Builder) WithShiftDistance(shiftDistance int) *Builder {
	b.ShiftDistance = shiftDistance
	return b
}

func (b *Builder) WithBarCodeConfigHeight(height int) *Builder {
	b.BarCodeConfig.Height = height
	return b
}

func (b *Builder) WithComponents(components ...Component) *Builder {
	for _, component := range components {
		b.Components = append(b.Components, component)
	}

	return b
}

func (b *Builder) String() string {
	var sb strings.Builder

	sb.WriteString("^XA")                                 // Start of label.
	sb.WriteString("~SD" + strconv.Itoa(b.Darkness))      // Set darkness.
	sb.WriteString("^LRN")                                // Disable reverse printing.
	sb.WriteString("^CI28")                               // UTF-8 encoding.
	sb.WriteString("^MMT")                                // Post print action tear off.
	sb.WriteString("^PW" + strconv.Itoa(b.Width))         // Label width.
	sb.WriteString("^LS" + strconv.Itoa(b.ShiftDistance)) // Shift distance (can be negative to shift left).
	sb.WriteString(b.BarCodeConfig.String())              // Bar code field default.
	sb.WriteString("^PQ1,0,1,Y")                          // Print quantity, pause, replicate, and tear off.

	for _, c := range b.Components {
		sb.WriteString(c.String())
	}

	sb.WriteString("^XZ") // End of label.

	return sb.String()
}
