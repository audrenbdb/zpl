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
	Components    []Component
}

func NewBuilder() *Builder {
	return &Builder{
		Darkness:      15,
		Width:         812,
		ShiftDistance: 0,
		Components:    []Component{},
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

func (b *Builder) WithComponents(components ...Component) *Builder {
	b.Components = append(b.Components, components...)
	return b
}

func (b *Builder) String() string {
	var sb strings.Builder

	sb.WriteString("^XA\n")                                      // Start of label.
	sb.WriteString("~SD" + strconv.Itoa(b.Darkness) + "\n")      // Set darkness.
	sb.WriteString("^LRN\n")                                     // Disable reverse printing.
	sb.WriteString("^CI28\n")                                    // UTF-8 encoding.
	sb.WriteString("^MMT\n")                                     // Post print action tear off.
	sb.WriteString("^PW" + strconv.Itoa(b.Width) + "\n")         // Label width.
	sb.WriteString("^LS" + strconv.Itoa(b.ShiftDistance) + "\n") // Shift distance (can be negative to shift left).
	sb.WriteString("^PQ1,0,1,Y" + "\n")                          // Print quantity, pause, replicate, and tear off.

	for _, c := range b.Components {
		sb.WriteString(c.String())
		sb.WriteString("\n")
	}

	sb.WriteString("^XZ") // End of label.

	return sb.String()
}
