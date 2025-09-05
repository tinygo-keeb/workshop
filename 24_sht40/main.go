package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/sht4x"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/shnm"
)

var (
	white = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	black = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
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
	display.SetRotation(drivers.Rotation180)
	display.ClearDisplay()
	time.Sleep(50 * time.Millisecond)

	sensor := sht4x.New(machine.I2C0)

	cnt := 0
	for {
		for x := int16(0); x < 128; x += 2 {
			temp, humidity, _ := sensor.ReadTemperatureHumidity()
			t := fmt.Sprintf("温度 %.2f ℃", float32(temp)/1000)
			h := fmt.Sprintf("湿度 %.2f ％", float32(humidity)/1000)

			display.ClearBuffer()
			tinyfont.WriteLine(&display, &shnm.Shnmk12, 5, 10, t, white)
			tinyfont.WriteLine(&display, &shnm.Shnmk12, 5, 30, h, white)
			display.Display()
			time.Sleep(1 * time.Second)
		}
		cnt++
	}
}
