package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

const (
	cols = 4
	rows = 3
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

func main() {
	pin := machine.GPIO1
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(pin)

	var leds [cols * rows]color.RGBA
	for t := uint8(0); ; t++ {
		for i := range leds {
			// LEDs are numbered top to bottom, then left to right.
			col, row := i/rows, i%rows
			c := wheel(t + uint8((col+row)*256/(cols+rows-1)))
			leds[i] = color.RGBA{R: c.R / 8, G: c.G / 8, B: c.B / 8}
		}
		ws.WriteColors(leds[:])
		time.Sleep(20 * time.Millisecond)
	}
}
