package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

const numLEDs = 10

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
	// GROVE pin 1 (SCL)
	pin := machine.GPIO13
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(pin)

	leds := make([]color.RGBA, numLEDs)
	for i := uint8(0); ; i++ {
		for j := range leds {
			c := wheel(i + uint8(j*256/numLEDs))
			leds[j] = color.RGBA{R: c.R / 8, G: c.G / 8, B: c.B / 8}
		}
		ws.WriteColors(leds)
		time.Sleep(20 * time.Millisecond)
	}
}
