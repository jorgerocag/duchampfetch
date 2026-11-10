// duchampfetch: neofetch con un Rotorelief de Duchamp animado que reacciona
// al estado de tu ordenador.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

func main() {
	static := flag.Bool("static", false, "dibuja un solo fotograma, sin animación")
	paletteName := flag.String("palette", "classic", "paleta: classic, opart o retro")
	rows := flag.Int("size", 24, "alto del disco en filas de terminal")
	duration := flag.Duration("duration", 3*time.Second, "duración de la animación")
	noReact := flag.Bool("no-react", false, "el disco no reacciona a las métricas del sistema")

	// Para previsualizar cómo se vería con otros valores (-1 = usar el real).
	cpuOver := flag.Float64("cpu", -1, "fuerza el uso de CPU (0-100)")
	ramOver := flag.Float64("ram", -1, "fuerza el uso de RAM (0-100)")
	tempOver := flag.Float64("temp", -1, "fuerza la temperatura en °C")
	uptimeOver := flag.Duration("uptime", 0, "fuerza el uptime (ej. 13h)")
	flag.Parse()

	pal, ok := palettes[*paletteName]
	if !ok {
		names := make([]string, 0, len(palettes))
		for n := range palettes {
			names = append(names, n)
		}
		sort.Strings(names)
		fmt.Fprintf(os.Stderr, "paleta desconocida %q (disponibles: %v)\n", *paletteName, names)
		os.Exit(1)
	}

	m := CollectMetrics()
	if *cpuOver >= 0 {
		m.CPU = *cpuOver
	}
	if *ramOver >= 0 {
		m.RAM = *ramOver
	}
	if *tempOver >= 0 {
		m.TempC, m.HasTemp = *tempOver, true
	}
	if *uptimeOver > 0 {
		m.Uptime = *uptimeOver
	}

	p := DefaultParams()
	if !*noReact {
		p = ParamsFromMetrics(m, *rows)
	}

	info := InfoLines(CollectInfo(m, p), pal.A)

	// Si la salida no es un terminal (pipe, redirección), forzamos modo estático.
	if st, err := os.Stdout.Stat(); err == nil && st.Mode()&os.ModeCharDevice == 0 {
		*static = true
	}

	if *static {
		for _, l := range Compose(RenderFrame(*rows, 0.8, pal, p), info) {
			fmt.Println(l)
		}
		return
	}

	Animate(*rows, pal, p, info, *duration)
}
