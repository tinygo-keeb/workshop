package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ws2812"
)

const (
	cols = 4
	rows = 3

	maxLevel = 16
)

// wheel returns a colour on the hue circle at position p.
func wheel(p uint8) color.RGBA {
	switch {
	case p < 85:
		return color.RGBA{R: 255 - p*3, G: p * 3}
	case p < 170:
		p -= 85
		return color.RGBA{G: 255 - p*3, B: p * 3}
	default:
		p -= 170
		return color.RGBA{R: p * 3, B: 255 - p*3}
	}
}

func scale(v uint8, level int) uint8 {
	return uint8(int(v) * level / maxLevel)
}

func main() {
	pin := machine.GPIO1
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(pin)

	enc := encoders.NewQuadratureViaInterrupt(machine.GPIO3, machine.GPIO4)
	enc.Configure(encoders.QuadratureConfig{Precision: 4})

	level := maxLevel / 2
	last := enc.Position()

	var leds [cols * rows]color.RGBA
	for t := uint8(0); ; t++ {
		pos := enc.Position()
		level = min(max(level+pos-last, 0), maxLevel)
		last = pos

		for i := range leds {
			// LEDs are numbered top to bottom, then left to right.
			col, row := i/rows, i%rows
			c := wheel(t + uint8((col+row)*256/(cols+rows-1)))
			leds[i] = color.RGBA{R: scale(c.R, level), G: scale(c.G, level), B: scale(c.B, level)}
		}
		ws.WriteColors(leds[:])
		time.Sleep(20 * time.Millisecond)
	}
}
