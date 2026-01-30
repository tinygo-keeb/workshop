package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/tinyfont/gophers"
)

func main() {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 2.8 * machine.MHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})
	display.ClearDisplay()
	time.Sleep(50 * time.Millisecond)

	rotDisplay := RotatedDisplay{display}

	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	tinyfont.WriteLine(&rotDisplay, &freemono.Bold9pt7b, 5, 10, "hello", white)
	tinyfont.WriteLine(&rotDisplay, &gophers.Regular58pt, 10, 70, "B", white)
	tinyfont.WriteLine(&rotDisplay, &gophers.Regular58pt, 10, 110, "H", white)
	display.Display()
}

type RotatedDisplay struct {
	drivers.Displayer
}

func (d *RotatedDisplay) Size() (x, y int16) {
	return y, x
}

func (d *RotatedDisplay) SetPixel(x, y int16, c color.RGBA) {
	_, sy := d.Displayer.Size()
	d.Displayer.SetPixel(y, sy-x, c)
}
