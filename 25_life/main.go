package main

import (
	"image/color"
	"machine"
	"math/rand"
	"time"

	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/ws2812"
)

const (
	DISPLAY_WIDTH  = 128
	DISPLAY_HEIGHT = 64
	CELL_SIZE      = 1
	WIDTH          = DISPLAY_WIDTH / CELL_SIZE
	HEIGHT         = DISPLAY_HEIGHT / CELL_SIZE
)

type Field [HEIGHT * WIDTH]uint8

var a, b Field

var display = ssd1306.NewI2C(machine.I2C0)

func main() {

	led := machine.GPIO16
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(led)

	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 2.8 * machine.MHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   DISPLAY_WIDTH,
		Height:  DISPLAY_HEIGHT,
	})
	display.ClearDisplay()

	enc := encoders.NewQuadratureViaInterrupt(
		machine.GPIO3,
		machine.GPIO4,
	)

	enc.Configure(encoders.QuadratureConfig{
		Precision: 4,
	})

	ch := make(chan *Field)
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	blue := color.RGBA{R: 0x00, G: 0x00, B: 0x80, A: 0x80}
	black := color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}

	go func() {
		for field := range ch {
			display.ClearBuffer()
			for y := range HEIGHT {
				for x := range WIDTH {
					cell := field[y*WIDTH+x]
					if cell > 0 {
						if CELL_SIZE > 1 {
							showRect(int16(x*CELL_SIZE), int16(y*CELL_SIZE), CELL_SIZE, CELL_SIZE, white)
						} else {
							display.SetPixel(int16(x*CELL_SIZE), int16(y*CELL_SIZE), white)
						}
					}
				}
			}
			display.Display()
		}
	}()

	field, next := &a, &b
	GenerateFirstRound(field, 4) // 1/4 of the cells will be alive
	for {
		ch <- field
		ws.WriteColors([]color.RGBA{blue})
		field.NextRound(next)
		ws.WriteColors([]color.RGBA{black})
		time.Sleep(time.Duration(enc.Position()) * 10 * time.Millisecond)
		field, next = next, field
	}
}

func showRect(x int16, y int16, w int16, h int16, c color.RGBA) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			display.SetPixel(i, j, c)
		}
	}
}

// setVitality sets the vitality (alive or dead) of the cell at the given coordinates.
// Coordinates are wrapped around the field dimensions.
func (field *Field) setVitality(x, y int, vitality uint8) {
	field[(y%HEIGHT)*WIDTH+(x%WIDTH)] = vitality
}

// getVitality returns the vitality (alive or dead) of the cell at the given coordinates.
// Coordinates are wrapped around the field dimensions.
func (field *Field) getVitality(x, y uint) uint8 {
	return field[(y%HEIGHT)*WIDTH+(x%WIDTH)]
}

// NextVitality determines and returns the vitality of the cell at (x, y) in the next round.
func (field *Field) NextVitality(x, y uint) uint8 {
	var neighbours uint8
	neighbours += field.getVitality(x-1, y-1)
	neighbours += field.getVitality(x-1, y)
	neighbours += field.getVitality(x-1, y+1)
	neighbours += field.getVitality(x, y-1)
	neighbours += field.getVitality(x, y+1)
	neighbours += field.getVitality(x+1, y-1)
	neighbours += field.getVitality(x+1, y)
	neighbours += field.getVitality(x+1, y+1)
	switch neighbours {
	case 3:
		// Cell is born
		return 1
	case 2:
		return field[(y%HEIGHT)*WIDTH+(x%WIDTH)] // Cell survives
	default:
		return 0 // Cell dies
	}

}

// NextRound calculates and returns the field state in the next round of the game.
func (field *Field) NextRound(next *Field) {
	for y := range HEIGHT {
		for x := range WIDTH {
			next.setVitality(x, y, field.NextVitality(uint(x), uint(y)))
		}
	}
}

// GenerateFirstRound generates a new field with a (pseudo) random seed
func GenerateFirstRound(field *Field, population int) {
	for i := 0; i < (WIDTH * HEIGHT / population); i++ {
		field.setVitality(rand.Intn(WIDTH), rand.Intn(HEIGHT), 1)
	}
}
