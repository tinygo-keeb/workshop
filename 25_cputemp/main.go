package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/shnm"
)

var (
	white = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	black = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
)

type celsius float32

func (c celsius) String() string {
	return fmt.Sprintf("温度:%4.1f℃", c)
}

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
	display.SetRotation(drivers.Rotation180)
	display.ClearDisplay()
	time.Sleep(50 * time.Millisecond)

	for {
		temp := celsius(float32(machine.ReadTemperature()) / 1000)
		println(temp.String())

		display.ClearBuffer()
		tinyfont.WriteLine(display, &shnm.Shnmk12, 5, 10, temp.String(), white)
		display.Display()
		time.Sleep(time.Second)
	}
}
