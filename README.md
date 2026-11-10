# duchampfetch

`duchampfetch` es una herramienta de información del sistema para terminal,
inspirada en los *Rotoreliefs* de Marcel Duchamp y en la película *Anémic
Cinéma* (1926). Muestra un disco óptico animado junto a datos del ordenador.
El dibujo ASCII cambia con las métricas disponibles: uso de CPU, memoria RAM
y temperatura.

## Requisitos

- Go 1.27.3 o posterior para ejecutarlo desde el código.
- Linux, macOS o Windows. La disponibilidad de la temperatura depende de los
  sensores y de la plataforma; si no se detecta, las demás funciones siguen
  disponibles.
- Una terminal que admita secuencias ANSI para mostrar los colores y animar
  el disco.

## Instalación y ejecución

Clona el repositorio y entra en su carpeta. Para iniciarlo directamente:

```sh
go run .
```

Para compilar un ejecutable local:

```sh
go build -o duchampfetch .
./duchampfetch
```

En Windows, ejecuta `duchampfetch.exe`.

## Opciones

```text
-static             dibuja un fotograma y termina
-size N             alto del disco en filas (por defecto: 24)
-duration 3s        duración de la animación
-palette classic    paleta: classic, opart o retro
-no-react           desactiva la respuesta a las métricas del sistema
-cpu 0..100         simula el uso de CPU
-ram 0..100         simula el uso de RAM
-temp grados        simula la temperatura
-uptime 13h         simula el tiempo encendido
```

El patrón clásico usa círculos alternos cuyos centros se desplazan en
secuencia, mientras el conjunto gira. La actividad de CPU selecciona entre
los patrones clásico, espiral, moiré y flor. La RAM ajusta la separación de
las bandas; la temperatura puede añadir deformación y aumentar la velocidad.

## Ejemplos

Mostrar un fotograma estático:

```sh
go run . -static
```

Aumentar el tamaño y usar la paleta op-art:

```sh
go run . -size 32 -palette opart -duration 5s
```

Probar valores simulados:

```sh
go run . -cpu 85 -ram 70 -temp 90 -uptime 13h
```

Desactivar la reacción a las métricas y mantener el patrón clásico:

```sh
go run . -no-react
```
