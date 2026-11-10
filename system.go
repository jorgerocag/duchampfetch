package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

// Info es una línea clave-valor de la información del sistema.
type Info struct {
	Key   string
	Value string
}

// Metrics son las cuatro métricas que gobiernan el aspecto del disco.
type Metrics struct {
	CPU      float64       // uso de CPU, 0-100
	RAM      float64       // RAM usada, 0-100
	MemUsed  uint64        // bytes
	MemTotal uint64        // bytes
	TempC    float64       // temperatura de la CPU en °C
	HasTemp  bool          // false si no se ha podido leer ningún sensor
	Uptime   time.Duration // tiempo encendido
}

// CollectMetrics lee las métricas del sistema. Tarda ~200 ms porque el uso de
// CPU se mide sobre una pequeña ventana de tiempo.
// Si algo falla, deja ese valor a cero.
func CollectMetrics() Metrics {
	var m Metrics

	if p, err := cpu.Percent(200*time.Millisecond, false); err == nil && len(p) > 0 {
		m.CPU = p[0]
	}
	if v, err := mem.VirtualMemory(); err == nil {
		m.RAM = v.UsedPercent
		m.MemUsed = v.Used
		m.MemTotal = v.Total
	}
	if u, err := host.Uptime(); err == nil {
		m.Uptime = time.Duration(u) * time.Second
	}
	m.TempC, m.HasTemp = readTemp()

	return m
}

// readTemp devuelve la temperatura de la CPU. Prefiere los sensores típicos de
// CPU (coretemp en Intel, k10temp/zenpower en AMD) y, si no hay, usa el más
// caliente de los que haya.
func readTemp() (float64, bool) {
	// En Linux gopsutil puede devolver datos parciales junto con un error
	// (warnings), así que solo nos importa si hay datos.
	stats, _ := sensors.SensorsTemperatures()

	cpuKeys := []string{"coretemp", "k10temp", "zenpower", "cpu"}
	var bestCPU, bestAny float64
	var haveCPU, haveAny bool

	for _, s := range stats {
		if s.Temperature < 1 || s.Temperature > 150 { // descarta lecturas absurdas
			continue
		}
		if !haveAny || s.Temperature > bestAny {
			bestAny, haveAny = s.Temperature, true
		}
		key := strings.ToLower(s.SensorKey)
		for _, k := range cpuKeys {
			if strings.Contains(key, k) {
				if !haveCPU || s.Temperature > bestCPU {
					bestCPU, haveCPU = s.Temperature, true
				}
				break
			}
		}
	}

	if haveCPU {
		return bestCPU, true
	}
	return bestAny, haveAny
}

// CollectInfo construye las líneas de información a partir de las métricas.
// Si algo falla, simplemente omite esa línea.
func CollectInfo(m Metrics, p Params) []Info {
	var out []Info

	if h, err := host.Info(); err == nil {
		if user := os.Getenv("USER"); user != "" {
			out = append(out, Info{"User", fmt.Sprintf("%s@%s", user, h.Hostname)})
		}
		os_ := h.Platform
		if os_ == "" {
			os_ = h.OS
		}
		out = append(out, Info{"OS", fmt.Sprintf("%s %s", os_, h.PlatformVersion)})
		out = append(out, Info{"Kernel", h.KernelVersion})
	}

	out = append(out, Info{"Uptime", formatUptime(m.Uptime)})

	if sh := os.Getenv("SHELL"); sh != "" {
		out = append(out, Info{"Shell", filepath.Base(sh)})
	}

	term := os.Getenv("TERM_PROGRAM")
	if term == "" {
		term = os.Getenv("TERM")
	}
	if term != "" {
		out = append(out, Info{"Terminal", term})
	}

	if c, err := cpu.Info(); err == nil && len(c) > 0 {
		name := c[0].ModelName
		if n, err := cpu.Counts(true); err == nil {
			name = fmt.Sprintf("%s (%d)", name, n)
		}
		out = append(out, Info{"CPU", fmt.Sprintf("%s · %.0f%%", name, m.CPU)})
	}

	if m.MemTotal > 0 {
		out = append(out, Info{"Memory", fmt.Sprintf("%.1f GiB / %.1f GiB (%.0f%%)",
			float64(m.MemUsed)/(1<<30), float64(m.MemTotal)/(1<<30), m.RAM)})
	}

	if m.HasTemp {
		out = append(out, Info{"Temp", fmt.Sprintf("%.0f°C", m.TempC)})
	}

	out = append(out, Info{"Pattern", patternNames[p.Pattern]})

	return out
}

func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, mins)
	default:
		return fmt.Sprintf("%dm", mins)
	}
}
