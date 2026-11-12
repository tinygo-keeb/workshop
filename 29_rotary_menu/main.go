package main

// ロータリーエンコーダーでメニューを選び、ブザーで音を鳴らすサンプル
//
// ※ ブザーは EX01 と 3V3 の間につないでください (22_buzzer と同じ)
//
// 回す : カーソルを動かす (カチッと音が鳴る)
// 押す : 選んでいる音を鳴らす

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/tone"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
)

var (
	white = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	black = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
)

// メニューの項目と、選んだときに鳴らす音
type menuItem struct {
	name  string
	notes []tone.Note
}

var items = []menuItem{
	{"C", []tone.Note{tone.C5}},
	{"E", []tone.Note{tone.E5}},
	{"G", []tone.Note{tone.G5}},
	{"C E G C", []tone.Note{tone.C5, tone.E5, tone.G5, tone.C6}},
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

	enc := encoders.NewQuadratureViaInterrupt(
		machine.GPIO3,
		machine.GPIO4,
	)
	enc.Configure(encoders.QuadratureConfig{
		Precision: 4,
	})

	btn := machine.GPIO2
	btn.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	speaker, err := tone.New(machine.PWM7, machine.GPIO14)
	if err != nil {
		println("failed to configure PWM")
		return
	}

	cursor := 0
	oldValue := enc.Position()
	draw(display, cursor)

	for {
		// 回した分だけカーソルを動かす (端まで行ったら反対側に戻る)
		if newValue := enc.Position(); newValue != oldValue {
			cursor = ((cursor+newValue-oldValue)%len(items) + len(items)) % len(items)
			oldValue = newValue
			draw(display, cursor)
			beep(speaker, tone.C6, 20*time.Millisecond)
		}

		// 押したら選んでいる音を鳴らす
		if !btn.Get() {
			for _, note := range items[cursor].notes {
				beep(speaker, note, 200*time.Millisecond)
			}
			// 離すまで待つ
			for !btn.Get() {
				time.Sleep(10 * time.Millisecond)
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}

// beep は note を d の間だけ鳴らす
func beep(speaker tone.Speaker, note tone.Note, d time.Duration) {
	speaker.SetNote(note)
	time.Sleep(d)
	speaker.Stop()
}

// draw はメニューを表示する。選んでいる行は白黒を反転する
func draw(display *ssd1306.Device, cursor int) {
	display.ClearBuffer()
	for i, item := range items {
		y := int16(14 + i*16)
		c := white
		if i == cursor {
			tinydraw.FilledRectangle(display, 0, y-14, 128, 16, white)
			c = black
		}
		tinyfont.WriteLine(display, &freesans.Regular9pt7b, 8, y, item.name, c)
	}
	display.Display()
}
