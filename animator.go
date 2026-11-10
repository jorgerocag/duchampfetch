package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	hideCursor = "\x1b[?25l"
	showCursor = "\x1b[?25h"
)

// Animate reproduce la animación durante `duration` y se queda congelada en
// el último fotograma. La velocidad de giro sale de p.Speed (vueltas/s).
func Animate(rows int, pal Palette, p Params, info []string, duration time.Duration) {
	// Si nos interrumpen (Ctrl+C), devolvemos el cursor antes de salir.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	ticker := time.NewTicker(time.Second / 30) // ~30 fps
	defer ticker.Stop()

	start := time.Now()
	first := true

	for {
		elapsed := time.Since(start)
		done := elapsed >= duration
		if done {
			elapsed = duration
		}
		t := 2 * math.Pi * p.Speed * elapsed.Seconds()

		lines := Compose(RenderFrame(rows, t, pal, p), info)
		if !first {
			fmt.Printf("\x1b[%dA", len(lines)) // subir el cursor al inicio del bloque
		}
		first = false
		fmt.Print(strings.Join(lines, "\x1b[K\n") + "\x1b[K\n")

		if done {
			return
		}

		select {
		case <-ticker.C:
		case <-sig:
			return
		}
	}
}
