package main

import "math"

// Params controla la geometría y el movimiento del rotorelief.
type Params struct {
	Pattern int
	Spacing float64
	Warp    float64
	Speed   float64
}

var patternNames = map[int]string{
	0: "clásico",
	1: "espiral",
	2: "moiré",
	3: "flor",
}

// DefaultParams devuelve el rotorelief clásico con valores seguros.
func DefaultParams() Params {
	return Params{Pattern: 0, Spacing: 0.065, Speed: 0.35}
}

// ParamsFromMetrics convierte las métricas del sistema en cambios visuales.
func ParamsFromMetrics(m Metrics, rows int) Params {
	p := DefaultParams()
	// El uso de CPU selecciona uno de los cuatro patrones.
	p.Pattern = int(math.Min(3, math.Floor(m.CPU/25)))
	// La RAM controla la separación de los anillos; más uso produce más anillos.
	p.Spacing = 0.09 - 0.04*clampPercent(m.RAM)/100
	// La temperatura añade deformación y acelera ligeramente el giro.
	if m.HasTemp {
		temp := math.Max(0, math.Min(100, (m.TempC-35)/55*100))
		p.Warp = 0.10 * temp / 100
		p.Speed += 0.06 * temp / 100
	}
	// Mantener el detalle legible en tamaños pequeños.
	if rows > 0 && rows < 8 {
		p.Spacing = math.Max(p.Spacing, 0.075)
	}
	return p
}

func clampPercent(v float64) float64 {
	return math.Max(0, math.Min(100, v))
}
