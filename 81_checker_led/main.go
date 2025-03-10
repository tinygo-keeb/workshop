//マイコンボード及びキーボードのLEDを点灯させます。
//ロータリーエンコーダのSWで色が変わります。

package main

import (
	"machine"
	"time"

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
)

type WS2812B struct {
	Pin machine.Pin
	ws  *piolib.WS2812B
}

func NewWS2812B(pin machine.Pin) *WS2812B {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	s, _ := pio.PIO0.ClaimStateMachine()
	ws, _ := piolib.NewWS2812B(s, pin)
	ws.EnableDMA(true)
	return &WS2812B{
		Pin: pin,
		ws:  ws,
	}
}

func (ws *WS2812B) WriteRaw(rawGRB []uint32) {
	ws.ws.WriteRaw(rawGRB)
}

func main() {
	const (
		white = 0xFFFFFFFF
		black = 0x00000000
		red   = 0x00FF0000
		green = 0xFF000000
		blue  = 0x0000FF00
	)

	// 色のプリセットパターン
	allColors := map[string][]uint32{
		"allWhite": {white, white, white, white, white, white, white, white, white, white, white, white},
		"allBlack": {black, black, black, black, black, black, black, black, black, black, black, black},
		"allRed":   {red, red, red, red, red, red, red, red, red, red, red, red},
		"allGreen": {green, green, green, green, green, green, green, green, green, green, green, green},
		"allBlue":  {blue, blue, blue, blue, blue, blue, blue, blue, blue, blue, blue, blue},
	}

	// LED初期化
	OnboardWs := NewWS2812B(machine.GPIO16)
	ws := NewWS2812B(machine.GPIO1)

	// GPIO 21 をボタン入力として設定 (プルアップ有効)
	bootBtn := machine.GPIO2
	bootBtn.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// 切り替え用
	patterns := []string{"allRed", "allGreen", "allBlue", "allWhite", "allBlack"}
	currentPattern := 0

	// メインループ
	for {
		// 色を設定
		color := allColors[patterns[currentPattern]]
		colorOne := []uint32{color[0]}
		ws.WriteRaw(color)
		OnboardWs.WriteRaw(colorOne)
		if bootBtn.Get() == false {
			time.Sleep(10 * time.Millisecond)
			if bootBtn.Get() == false {
				time.Sleep(200 * time.Millisecond)
				currentPattern = (currentPattern + 1) % len(patterns)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
