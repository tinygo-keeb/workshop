package main

import (
	"image/color"
	"machine"
	"strconv"
	"time"

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/shnm"
)

// WS2812B LEDコントローラー構造体
type WS2812B struct {
	Pin machine.Pin
	ws  *piolib.WS2812B
}

func NewWS2812B(pin machine.Pin) *WS2812B {
	s, _ := pio.PIO0.ClaimStateMachine()
	ws, _ := piolib.NewWS2812B(s, pin)
	ws.EnableDMA(true)
	return &WS2812B{
		ws: ws,
	}
}

func (ws *WS2812B) PutColor(c color.Color) {
	ws.ws.PutColor(c)
}

func (ws *WS2812B) WriteRaw(rawGRB []uint32) error {
	return ws.ws.WriteRaw(rawGRB)
}

// 定数定義
const (
	testDuration          = 10 * time.Second
	ROWS                  = 3
	COLS                  = 4
	debounceTime          = 50 * time.Millisecond
	displayUpdateInterval = 50 * time.Millisecond
)

// キーの状態を保持する構造体
type KeyState struct {
	isPressed bool
	lastPress time.Time
}

// ディスプレイの状態を管理する構造体
type DisplayState struct {
	display    *ssd1306.Device // ポインタ型
	lastCount  int
	lastTime   int
	needUpdate bool
}

func NewDisplayState(display *ssd1306.Device) *DisplayState {
	return &DisplayState{
		display:    display,
		lastCount:  -1,
		lastTime:   -1,
		needUpdate: false,
	}
}

// ディスプレイ更新を最適化する関数
func (ds *DisplayState) updateDisplay(count int, remainingTime int, testStatus string) {
	if count != ds.lastCount || remainingTime != ds.lastTime || ds.needUpdate {
		ds.display.ClearBuffer()
		white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

		switch testStatus {
		case "waiting":
			tinyfont.WriteLine(ds.display, &shnm.Shnmk12, 5, 12, "赤いキーを押す！！", white)
		case "testing":
			tinyfont.WriteLine(ds.display, &shnm.Shnmk12, 5, 12, "残り時間: "+strconv.Itoa(remainingTime)+"秒", white)
			tinyfont.WriteLine(ds.display, &shnm.Shnmk12, 5, 30, "打鍵数: "+strconv.Itoa(count), white)
		case "result":
			tinyfont.WriteLine(ds.display, &shnm.Shnmk12, 5, 12, "テスト終了", white)
			tinyfont.WriteLine(ds.display, &shnm.Shnmk12, 5, 30, "合計打鍵数: "+strconv.Itoa(count), white)
			speed := float64(count) / 10.0
			speedStr := strconv.FormatFloat(speed, 'f', 1, 64)
			tinyfont.WriteLine(ds.display, &shnm.Shnmk12, 5, 48, "速度: "+speedStr+"回/秒", white)
		}

		ds.display.Display()
		ds.lastCount = count
		ds.lastTime = remainingTime
		ds.needUpdate = false
	}
}

func waitForSW12Key(colPins []machine.Pin, rowPins []machine.Pin) {
	wasPressed := false
	for {
		colPins[3].High()
		time.Sleep(time.Millisecond)

		if rowPins[2].Get() {
			if !wasPressed {
				wasPressed = true
			}
		} else if wasPressed {
			break
		}

		colPins[3].Low()
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
}

func scanKeys(colPins []machine.Pin, rowPins []machine.Pin, keyStates *[][]KeyState) int {
	keyCount := 0
	now := time.Now()

	for col := 0; col < len(colPins); col++ {
		// 現在の列をアクティブに
		for i, pin := range colPins {
			if i == col {
				pin.High()
			} else {
				pin.Low()
			}
		}

		time.Sleep(time.Millisecond)

		// 行をスキャン
		for row := 0; row < len(rowPins); row++ {
			currentState := rowPins[row].Get()

			if currentState {
				if !(*keyStates)[row][col].isPressed {
					(*keyStates)[row][col].isPressed = true
					(*keyStates)[row][col].lastPress = now
					keyCount++
				}
			} else {
				if (*keyStates)[row][col].isPressed {
					(*keyStates)[row][col].isPressed = false
					if now.Sub((*keyStates)[row][col].lastPress) > debounceTime {
						// デバウンス処理済み
					}
				}
			}
		}
	}

	return keyCount
}

func main() {
	// WS2812B LED初期化
	ws := NewWS2812B(machine.GPIO1)
	leds := make([]uint32, 12)
	leds[11] = 0x00FF0000 // 赤色LED点灯
	ws.WriteRaw(leds)

	// I2C初期化
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 2.8 * machine.MHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	// ディスプレイ初期化 - ここを修正
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})
	display.SetRotation(drivers.Rotation180)

	// ディスプレイ状態管理の初期化 - ポインタとして渡す
	displayState := NewDisplayState(display)

	// キーマトリックスピン設定
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

	// ピンの初期化
	for _, c := range colPins {
		c.Configure(machine.PinConfig{Mode: machine.PinOutput})
		c.Low()
	}
	for _, c := range rowPins {
		c.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	// キー状態の初期化
	keyStates := make([][]KeyState, ROWS)
	for i := range keyStates {
		keyStates[i] = make([]KeyState, COLS)
	}

	for {
		// 待機画面表示
		displayState.needUpdate = true
		displayState.updateDisplay(0, 0, "waiting")

		// キー待機
		waitForSW12Key(colPins, rowPins)

		// テスト開始
		keyCount := 0
		startTime := time.Now()
		endTime := startTime.Add(testDuration)

		// キー状態をリセット
		for i := range keyStates {
			for j := range keyStates[i] {
				keyStates[i][j].isPressed = false
			}
		}

		// 定期更新用のティッカー設定
		ticker := time.NewTicker(displayUpdateInterval)
		defer ticker.Stop()

		// メインテストループ
		for time.Now().Before(endTime) {
			keyCount += scanKeys(colPins, rowPins, &keyStates)
			remainingTime := int(endTime.Sub(time.Now()).Seconds())

			select {
			case <-ticker.C:
				displayState.updateDisplay(keyCount, remainingTime, "testing")
			default:
				// 継続
			}
		}

		// 結果表示
		displayState.needUpdate = true
		displayState.updateDisplay(keyCount, 0, "result")
		time.Sleep(3 * time.Second)
	}
}
