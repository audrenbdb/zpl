package zpl

import (
	"fmt"
	"image"
	"strconv"
	"strings"
)

var mapCode = map[int]string{
	1: "G", 2: "H", 3: "I", 4: "J", 5: "K", 6: "L", 7: "M", 8: "N", 9: "O", 10: "P",
	11: "Q", 12: "R", 13: "S", 14: "T", 15: "U", 16: "V", 17: "W", 18: "X", 19: "Y", 20: "g",
	40: "h", 60: "i", 80: "j", 100: "k", 120: "l", 140: "m", 160: "n", 180: "o", 200: "p",
	220: "q", 240: "r", 260: "s", 280: "t", 300: "u", 320: "v", 340: "w", 360: "x", 380: "y", 400: "z",
}

func convertFromImage(img image.Image) string {
	var blackLimit = 380
	var total int
	var widthBytes int
	var index int

	var sb strings.Builder

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	auxBinaryChar := [8]byte{'0', '0', '0', '0', '0', '0', '0', '0'}
	widthBytes = width / 8

	if width%8 > 0 {
		widthBytes = (width / 8) + 1
	} else {
		widthBytes = width / 8
	}

	total = widthBytes * height

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			r, g, b, a := img.At(w, h).RGBA()
			red := int(r >> 8)
			green := int(g >> 8)
			blue := int(b >> 8)
			auxChar := '1'
			totalColor := red + green + blue
			if a == 0 {
				totalColor = blackLimit + 1
			}
			if totalColor > blackLimit {
				auxChar = '0'
			}
			auxBinaryChar[index] = byte(auxChar)
			index++
			if index == 8 || w == (width-1) {
				sb.WriteString(binaryToHex(auxBinaryChar))
				auxBinaryChar = [8]byte{'0', '0', '0', '0', '0', '0', '0', '0'}
				index = 0
			}
		}
		sb.WriteString("\n")
	}

	data := encodeHexAscii(sb.String(), widthBytes)

	sb.Reset()

	sb.WriteString("^GFA,")
	sb.WriteString(strconv.Itoa(total))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(total))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(widthBytes))
	sb.WriteString(",")
	sb.WriteString(data)

	return sb.String()
}

func binaryToHex(binary [8]byte) string {
	decimal := 0
	for i, bit := range binary {
		if bit == '1' {
			decimal += 1 << (7 - i)
		}
	}
	return fmt.Sprintf("%02X", decimal)
}

func encodeHexAscii(code string, widthBytes int) string {
	maxLine := widthBytes * 2
	var sbCode strings.Builder
	var sbLine strings.Builder
	var previousLine string
	counter := 1
	aux := code[0]
	firstChar := false

	for i := 1; i < len(code); i++ {
		if firstChar {
			aux = code[i]
			firstChar = false
			continue
		}

		if code[i] == '\n' {
			if counter >= maxLine && aux == '0' {
				sbLine.WriteString(",")
			} else if counter >= maxLine && aux == 'F' {
				sbLine.WriteString("!")
			} else if counter > 20 {
				multi20 := (counter / 20) * 20
				resto20 := (counter % 20)
				sbLine.WriteString(mapCode[multi20])
				if resto20 != 0 {
					sbLine.WriteString(mapCode[resto20] + string(aux))
				} else {
					sbLine.WriteString(string(aux))
				}
			} else {
				sbLine.WriteString(mapCode[counter] + string(aux))
			}

			counter = 1
			firstChar = true

			if sbLine.String() == previousLine {
				sbCode.WriteString(":")
			} else {
				sbCode.WriteString(sbLine.String())
			}

			previousLine = sbLine.String()
			sbLine.Reset()
			continue
		}

		if aux == code[i] {
			counter++
		} else {
			if counter > 20 {
				multi20 := (counter / 20) * 20
				resto20 := counter % 20
				sbLine.WriteString(mapCode[multi20])
				if resto20 != 0 {
					sbLine.WriteString(mapCode[resto20] + string(aux))
				} else {
					sbLine.WriteString(string(aux))
				}
			} else {
				sbLine.WriteString(mapCode[counter] + string(aux))
			}
			counter = 1
			aux = code[i]
		}
	}

	return sbCode.String()
}
