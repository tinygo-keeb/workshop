package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/gophers"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const alphabetWidth = 598 // 実際に描画してみて、適切な値を設定する

func main() {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 2.8 * machine.MHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})
	enc := encoders.NewQuadratureViaInterrupt(
		machine.GPIO3,
		machine.GPIO4,
	)
	enc.Configure(encoders.QuadratureConfig{
		Precision: 4,
	})

	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address:  0x3C,
		Width:    128,
		Height:   64,
		Rotation: drivers.Rotation180,
	})
	display.ClearDisplay()
	time.Sleep(50 * time.Millisecond)

	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	tinyfont.WriteLine(&display, &gophers.Regular32pt, 0, 50, alphabet, white)
	display.Display()

	oldValue := int16(0)
	for {
		if newValue := int16(enc.Position()) % alphabetWidth; newValue != oldValue {
			display.ClearBuffer()
			oldValue = newValue

			println("value: ", newValue)

			tinyfont.WriteLine(&display, &gophers.Regular32pt, oldValue-alphabetWidth, 50, alphabet, white)
			tinyfont.WriteLine(&display, &gophers.Regular32pt, oldValue, 50, alphabet, white)
			tinyfont.WriteLine(&display, &gophers.Regular32pt, oldValue+alphabetWidth, 50, alphabet, white)
			display.Display()
		}
		time.Sleep(10 * time.Millisecond)
	}
}
