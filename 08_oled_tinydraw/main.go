package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
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

	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	for {
		tinydraw.Rectangle(display, 10, 20, 30, 40, white)
		display.Display()
		time.Sleep(500 * time.Millisecond)

		tinydraw.FilledCircle(display, 60, 50, 10, white)
		display.Display()
		time.Sleep(500 * time.Millisecond)

		tinydraw.Triangle(display, 100, 10, 80, 40, 60, 10, white)
		display.Display()
		time.Sleep(500 * time.Millisecond)

		display.ClearDisplay()
		time.Sleep(500 * time.Millisecond)
	}
}
