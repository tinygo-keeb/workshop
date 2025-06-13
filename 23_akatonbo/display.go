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

// 画面レイアウト定数
const (
	STATUS_Y         = 12 // ステータス行のテキスト表示位置
	STATUS_AREA_END  = 19 // ステータス行エリアの終端
	SEPARATOR_Y      = 16 // 区切り線の位置
	SCROLL_START_Y   = 30 // スクロール領域の開始位置
	SCROLL_CLEAR_Y   = 20 // スクロール領域のクリア開始位置
	LINE_HEIGHT      = 12 // 行の高さ
	SCROLL_MAX_Y     = 60 // スクロール領域の最大Y位置
	MAX_SCROLL_LINES = 3  // スクロール表示可能行数
)

var white = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
var black = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}

// Display 構造体 - 画面制御とスクロール管理
type Display struct {
	device      *ssd1306.Device
	currentLine int                      // 現在の表示行（自動スクロール用）
	lines       [MAX_SCROLL_LINES]string // 表示中の行を保存
}

// InitDisplay ディスプレイを初期化する関数
func InitDisplay() *Display {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 2.8 * machine.MHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	device := ssd1306.NewI2C(machine.I2C0)
	device.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})

	device.SetRotation(drivers.Rotation180)
	device.ClearDisplay()
	time.Sleep(50 * time.Millisecond)

	return &Display{
		device:      &device,
		currentLine: 0,
		lines:       [MAX_SCROLL_LINES]string{}, // 空文字列で初期化
	}
}

// UpdateStatus ステータス行を更新する関数
func (d *Display) UpdateStatus(status string) {
	// ステータス行のエリアをクリア
	for y := 0; y < STATUS_AREA_END; y++ {
		for x := 0; x < 128; x++ {
			d.device.SetPixel(int16(x), int16(y), black)
		}
	}

	// ステータス行に表示
	tinyfont.WriteLine(d.device, &shnm.Shnmk12, 0, STATUS_Y, status, white)

	// 区切り線を描画
	for x := 0; x < 128; x++ {
		d.device.SetPixel(int16(x), SEPARATOR_Y, white)
	}

	d.device.Display()
}

// ClearScrollArea ステータス行以外をクリアする関数
func (d *Display) ClearScrollArea() {
	for y := SCROLL_CLEAR_Y; y < 64; y++ {
		for x := 0; x < 128; x++ {
			d.device.SetPixel(int16(x), int16(y), black)
		}
	}
	d.currentLine = 0 // 表示行をリセット
	// 行バッファもクリア
	for i := range d.lines {
		d.lines[i] = ""
	}
	d.device.Display()
}

// PrintLine 一行表示する関数（必要に応じて自動スクロール）
func (d *Display) PrintLine(message string) {
	// デバッグ情報をコンソールに出力
	fmt.Printf("PrintLine: currentLine=%d, message='%s'\n", d.currentLine, message)

	// 表示可能行数を超えた場合は画面をスクロール
	if d.currentLine >= MAX_SCROLL_LINES {
		fmt.Println("スクロール実行中...")
		d.scrollUp()
		d.currentLine = MAX_SCROLL_LINES - 1
		fmt.Printf("スクロール後: currentLine=%d\n", d.currentLine)
	}

	// 新しいメッセージを配列に保存
	d.lines[d.currentLine] = message
	fmt.Printf("行[%d]に保存: '%s'\n", d.currentLine, message)

	// 画面を再描画（スクロール領域のみ）
	d.redrawScrollArea()
	d.currentLine++
	fmt.Printf("currentLine更新: %d\n", d.currentLine)
}

// scrollUp 画面を1行上にスクロールする（内部関数）
func (d *Display) scrollUp() {
	fmt.Println("scrollUp開始")
	// デバッグ: スクロール前の状態を表示
	fmt.Print("スクロール前の行: ")
	for i := 0; i < MAX_SCROLL_LINES; i++ {
		fmt.Printf("[%d]='%s' ", i, d.lines[i])
	}
	fmt.Println()

	// 行を1つずつ上に移動
	for i := 0; i < MAX_SCROLL_LINES-1; i++ {
		d.lines[i] = d.lines[i+1]
	}
	// 最後の行をクリア
	d.lines[MAX_SCROLL_LINES-1] = ""

	// デバッグ: スクロール後の状態を表示
	fmt.Print("スクロール後の行: ")
	for i := 0; i < MAX_SCROLL_LINES; i++ {
		fmt.Printf("[%d]='%s' ", i, d.lines[i])
	}
	fmt.Println()
}

// redrawScrollArea スクロール領域を再描画する（内部関数）
func (d *Display) redrawScrollArea() {
	fmt.Println("redrawScrollArea開始")

	// スクロール領域をクリア
	for y := SCROLL_CLEAR_Y; y < 64; y++ {
		for x := 0; x < 128; x++ {
			d.device.SetPixel(int16(x), int16(y), black)
		}
	}

	// 保存されている全ての行を再描画
	for i := 0; i < MAX_SCROLL_LINES; i++ {
		if d.lines[i] != "" {
			y := int16(i*LINE_HEIGHT + SCROLL_START_Y)
			fmt.Printf("描画: 行[%d] Y=%d '%s'\n", i, y, d.lines[i])
			if y >= SCROLL_START_Y && y <= SCROLL_MAX_Y {
				tinyfont.WriteLine(d.device, &shnm.Shnmk12, 0, y, d.lines[i], white)
			} else {
				fmt.Printf("描画範囲外: Y=%d (範囲: %d-%d)\n", y, SCROLL_START_Y, SCROLL_MAX_Y)
			}
		}
	}

	d.device.Display()
	fmt.Println("redrawScrollArea完了")
}

// GetDevice デバイスへの直接アクセス（必要な場合）
func (d *Display) GetDevice() *ssd1306.Device {
	return d.device
}
