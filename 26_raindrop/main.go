package main

import (
	"machine"
	"time"

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
)

func NewWS2812B(pin machine.Pin) *piolib.WS2812B {
	s, _ := pio.PIO0.ClaimStateMachine()
	ws, _ := piolib.NewWS2812B(s, pin)
	ws.EnableDMA(true)
	return ws
}

func main() {
	// 16 random colors
	randCol := 0
	randColors := [16]uint32{
		0xFF00FF00, 0xFF0000FF, 0xFFFF0000, 0xFFFFFF00,
		0xFFFF00FF, 0xFF00FFFF, 0xFF800000, 0xFF008000,
		0xFF000080, 0xFF808000, 0xFF800080, 0xFF008080,
		0xFFC00000, 0xFF00C000, 0xFF0000C0, 0xFFC000C0,
	}

	colors := []uint32{
		0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF,
		0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF,
		0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF,
	}

	ws := NewWS2812B(machine.GPIO1)

	colPins := []machine.Pin{
		machine.GPIO5,
		machine.GPIO6,
		machine.GPIO7,
		machine.GPIO8,
	}

	rowPins := []machine.Pin{
		machine.GPIO9,
		machine.GPIO10,
		machine.GPIO11,
	}

	for _, c := range colPins {
		c.Configure(machine.PinConfig{Mode: machine.PinOutput})
		c.Low()
	}

	for _, c := range rowPins {
		c.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	input := make(chan int, 12)

	update := time.NewTicker(10 * time.Millisecond)
	decay := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case key := <-input:
			randCol++
			colors[key] = randColors[randCol%16]
		case <-decay.C:
			for i := range colors {
				g := ((colors[i] >> 24) & 0xFF) >> 1
				r := ((colors[i] >> 16) & 0xFF) >> 1
				b := ((colors[i] >> 8) & 0xFF) >> 1
				colors[i] = (g << 24) | (r << 16) | (b << 8) | 0
			}
		case <-update.C:
			ws.WriteRaw(colors)
		default:
			// COL1
			colPins[0].High()
			colPins[1].Low()
			colPins[2].Low()
			colPins[3].Low()
			time.Sleep(1 * time.Millisecond)

			if rowPins[0].Get() {
				input <- 0
			}
			if rowPins[1].Get() {
				input <- 1
			}
			if rowPins[2].Get() {
				input <- 2
			}

			// COL2
			colPins[0].Low()
			colPins[1].High()
			colPins[2].Low()
			colPins[3].Low()
			time.Sleep(1 * time.Millisecond)

			if rowPins[0].Get() {
				input <- 3
			}
			if rowPins[1].Get() {
				input <- 4
			}
			if rowPins[2].Get() {
				input <- 5
			}

			// COL3
			colPins[0].Low()
			colPins[1].Low()
			colPins[2].High()
			colPins[3].Low()
			time.Sleep(1 * time.Millisecond)

			if rowPins[0].Get() {
				input <- 6
			}
			if rowPins[1].Get() {
				input <- 7
			}
			if rowPins[2].Get() {
				input <- 8
			}

			// COL4
			colPins[0].Low()
			colPins[1].Low()
			colPins[2].Low()
			colPins[3].High()
			time.Sleep(1 * time.Millisecond)

			if rowPins[0].Get() {
				input <- 9
			}
			if rowPins[1].Get() {
				input <- 10
			}
			if rowPins[2].Get() {
				input <- 11
			}
		}
	}
}
