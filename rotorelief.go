package main

import "math"

// RGB es un color de 24 bits.
type RGB struct{ R, G, B uint8 }

// Palette define los dos tonos alternos de los anillos.
type Palette struct {
	A, B RGB
}

var palettes = map[string]Palette{
	"classic": {A: RGB{245, 245, 240}, B: RGB{15, 15, 15}}, // blanco y negro, como el original
	"opart":   {A: RGB{230, 30, 40}, B: RGB{20, 40, 200}},  // contraste extremo
	"retro":   {A: RGB{80, 255, 120}, B: RGB{5, 25, 10}},   // fósforo verde CRT
}

// Cuánto se desplaza el centro de los anillos respecto al del disco
// (el radio del disco es 1).
const offset = 0.30

// Pixel devuelve el color del punto (x, y) del disco en el ángulo t.
// inside es false si el punto queda fuera del disco.
//
// El patrón clásico reproduce los círculos alternos desplazados de la
// referencia de Anémic Cinéma. Los demás patrones usan bandas calculadas
// alrededor de centros excéntricos.
func Pixel(x, y, t float64, pal Palette, p Params) (c RGB, inside bool) {
	if math.Hypot(x, y) > 1 {
		return RGB{}, false
	}

	// El patrón principal sigue la construcción de Anémic Cinéma: círculos
	// alternos, cada uno desplazado siguiendo una trayectoria angular. La
	// animación gira la composición completa alrededor del centro del disco.
	if p.Pattern == 0 {
		ct, st := math.Cos(t), math.Sin(t)
		rx, ry := x*ct+y*st, -x*st+y*ct
		const circles = 33
		const turns = 2.0
		const smaller = 0.8
		startRadius := 0.95
		step := startRadius / circles
		angleStep := -2 * math.Pi * turns / (circles - 1)
		cx, cy, angle := 0.0, 0.0, 0.0
		color := pal.B // fondo oscuro dentro del límite del disco
		for i := 0; i < circles; i++ {
			radius := startRadius - float64(i)*step
			dx, dy := rx-cx, ry-cy
			if dx*dx+dy*dy <= radius*radius {
				if i%2 == 0 {
					color = pal.A
				} else {
					color = pal.B
				}
			}
			cx += smaller * step * math.Cos(angle)
			cy += smaller * step * math.Sin(angle)
			angle += angleStep
		}
		return color, true
	}

	// Temperatura: deformación ondulada que se mueve con la rotación.
	if p.Warp > 0 {
		x, y = x+p.Warp*math.Sin(5*y+t), y+p.Warp*math.Cos(5*x+t)
	}

	cx := offset * math.Cos(t)
	cy := offset * math.Sin(t)
	dx, dy := x-cx, y-cy
	d := math.Hypot(dx, dy)

	var band int
	switch p.Pattern {
	case 1: // espiral de dos brazos alrededor del centro excéntrico
		band = int(math.Floor(d/p.Spacing + 2*math.Atan2(dy, dx)/math.Pi))
	case 2: // doble centro: dos juegos de anillos que interfieren (moiré)
		d2 := math.Hypot(x+cx, y+cy)
		band = int(math.Floor(d/p.Spacing)) + int(math.Floor(d2/p.Spacing))
	case 3: // flor: los anillos se ondulan con 6 lóbulos
		lobed := d * (1 + 0.12*math.Cos(6*math.Atan2(dy, dx)+t))
		band = int(math.Floor(lobed / p.Spacing))
	default: // respaldo para patrones no definidos
		band = int(math.Floor(d / p.Spacing))
	}

	if (band%2+2)%2 == 0 { // el doble módulo evita problemas con negativos
		return pal.A, true
	}
	return pal.B, true
}
