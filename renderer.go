package main

import (
	"fmt"
	"strings"
)

const reset = "\x1b[0m"

func fg(c RGB) string { return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B) }
func bg(c RGB) string { return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B) }

// RenderFrame dibuja el disco con caracteres ASCII. Usa el doble de columnas
// que filas para compensar la proporción habitual de una celda de terminal.
func RenderFrame(rows int, t float64, pal Palette, p Params) []string {
	// Un carácter ASCII suele ser aproximadamente la mitad de alto que ancho;
	// dos columnas por fila conservan la forma circular en la terminal.
	cols := 2 * rows
	lines := make([]string, rows)

	toDisc := func(pos, count int) float64 { return (float64(pos)+0.5)/float64(count)*2 - 1 }

	for row := 0; row < rows; row++ {
		var sb strings.Builder
		for col := 0; col < cols; col++ {
			x, y := toDisc(col, cols), toDisc(row, rows)
			_, inside := Pixel(x, y, t, pal, p)
			if !inside {
				sb.WriteByte(' ')
				continue
			}
			coverage := 0.0
			for sy := -1; sy <= 1; sy++ {
				for sx := -1; sx <= 1; sx++ {
					nx := x + float64(sx)/(2*float64(rows))
					ny := y + float64(sy)/float64(rows)
					neighbor, ok := Pixel(nx, ny, t, pal, p)
					if ok && neighbor == pal.A {
						coverage++
					}
				}
			}
			glyph := asciiCoverage(coverage / 9)
			if glyph == ' ' {
				sb.WriteByte(' ')
			} else {
				sb.WriteString(fg(pal.A) + string(glyph) + reset)
			}
		}
		lines[row] = sb.String()
	}
	return lines
}

// asciiCoverage convierte la proporción de píxeles claros vecinos en un
// glifo de densidad equivalente, suavizando los bordes y afinando las líneas.
func asciiCoverage(coverage float64) rune {
	const ramp = " .,:;irsXA253hMHGS#9B&@"
	level := int(coverage*float64(len(ramp)-1) + 0.5)
	if level < 0 {
		level = 0
	} else if level >= len(ramp) {
		level = len(ramp) - 1
	}
	return rune(ramp[level])
}

// InfoLines da formato a la información del sistema (clave en negrita).
func InfoLines(info []Info, accent RGB) []string {
	lines := make([]string, 0, len(info))
	for _, i := range info {
		lines = append(lines, fmt.Sprintf("%s\x1b[1m%-9s%s %s", fg(accent), i.Key, reset, i.Value))
	}
	return lines
}

// Compose coloca la info a la derecha del disco, centrada en vertical.
func Compose(frame, info []string) []string {
	const gap = "   "
	out := make([]string, len(frame))
	start := (len(frame) - len(info)) / 2
	if start < 0 {
		start = 0
	}
	for i, line := range frame {
		out[i] = line
		if j := i - start; j >= 0 && j < len(info) {
			out[i] += gap + info[j]
		}
	}
	return out
}
