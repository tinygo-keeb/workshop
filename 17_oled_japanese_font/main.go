package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/gophers"
	"tinygo.org/x/tinyfont/shnm"
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
		//Rotation: drivers.Rotation180,
	})
	display.SetRotation(drivers.Rotation180)
	display.ClearDisplay()
	time.Sleep(50 * time.Millisecond)

	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	tinyfont.WriteLine(display, &shnm.Shnmk12, 5, 10, "こんにちは世界", white)
	tinyfont.WriteLine(display, &gophers.Regular32pt, 5, 50, "ABCEF", white)
	display.Display()
}
