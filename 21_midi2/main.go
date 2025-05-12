package main

import (
	"image/color"
	"machine"
	"machine/usb/adc/midi"
	"time"

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/shnm"
)

const (
	cable    = 0
	channel  = 1  // ピアノチャンネル
	drumCh   = 10 // ドラムチャンネル (MIDI仕様では10チャンネル、0ベースでは9)
	velocity = 0x7F
	bpm      = 100 // リズムパターンのテンポ

	BassDrum      = 36 // バスドラム
	SideStick     = 37 // サイドスティック/リムショット
	SnareDrum     = 38 // スネアドラム
	HandClap      = 39 // ハンドクラップ
	ElectricSnare = 40 // エレクトリックスネア
	LowFloorTom   = 41 // ローフロアタム
	ClosedHiHat   = 42 // クローズドハイハット
	HighFloorTom  = 43 // ハイフロアタム
	PedalHiHat    = 44 // ペダルハイハット
	LowTom        = 45 // ロータム
	OpenHiHat     = 46 // オープンハイハット
	LowMidTom     = 47 // ロー・ミッド・タム
	HighMidTom    = 48 // ハイミッドタム
	CrashCymbal1  = 49 // クラッシュシンバル1
	HighTom       = 50 // ハイトム
	RideCymbal1   = 51 // ライドシンバル1
	ChineseCymbal = 52 // チャイニーズシンバル
	RideBell      = 53 // ライドベル
	Tambourine    = 54 // タンバリン
	SplashCymbal  = 55 // スプラッシュシンバル
	Cowbell       = 56 // カウベル
	CrashCymbal2  = 57 // クラッシュシンバル2
	Vibraslap     = 58 // ビブラスラップ
	RideCymbal2   = 59 // ライドシンバル2
	Maracas       = 70 // マラカス
	Claves        = 75 // クラベス

	// LED (GRB)
	white  = 0x3F3F3FFF // 白色
	red    = 0x00FF00FF // 赤色
	green  = 0xFF0000FF // 緑色
	blue   = 0x0000FFFF // 青色
	yellow = 0xFFFF00FF // 黄色
	purple = 0x008080FF // 紫色
	pink   = 0x69FFB4FF // ピンク色
	cyan   = 0xFF00FFFF // シアン
	orange = 0x80FF00FF // オレンジ
	black  = 0x000000FF // 黒色
)

// ディスプレイ用の色定義
var (
	displayWhite = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	displayBlack = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
)

// Not番号から音名のマッピング
var noteNamesKatakana = map[midi.Note]string{
	midi.C3: "ド3", midi.C4: "ド4", midi.C5: "ド5", midi.C6: "ド6",
	midi.D3: "レ3", midi.D4: "レ4", midi.D5: "レ5", midi.D6: "レ6",
	midi.E3: "ミ3", midi.E4: "ミ4", midi.E5: "ミ5", midi.E6: "ミ6",
	midi.F3: "ファ3", midi.F4: "ファ4", midi.F5: "ファ5", midi.F6: "ファ6",
	midi.G3: "ソ3", midi.G4: "ソ4", midi.G5: "ソ5", midi.G6: "ソ6",
	midi.A3: "ラ3", midi.A4: "ラ4", midi.A5: "ラ5", midi.A6: "ラ6",
	midi.B3: "シ3", midi.B4: "シ4", midi.B5: "シ5", midi.B6: "シ6",
}

// ドラム音のマッピング
var drumNames = map[uint8]string{
	BassDrum:      "バスドラム",
	SideStick:     "リムショット",
	SnareDrum:     "スネア",
	HandClap:      "クラップ",
	ElectricSnare: "Eスネア",
	LowFloorTom:   "ロータム",
	ClosedHiHat:   "ハイハット(閉)",
	HighFloorTom:  "ハイタム",
	PedalHiHat:    "ペダルHH",
	LowTom:        "ロータム",
	OpenHiHat:     "ハイハット(開)",
	CrashCymbal1:  "クラッシュ",
	RideCymbal1:   "ライド",
	Cowbell:       "カウベル",
	Claves:        "クラベス",
	Maracas:       "マラカス",
}

type DrumPattern struct {
	Name    string    // パターン名
	Steps   [][]uint8 // 各ステップで鳴るドラム音のリスト
	StepLen int       // 1ステップの長さ (ミリ秒)
}

var metronomeBeat = DrumPattern{
	Name: "Metronome",
	Steps: [][]uint8{
		{SnareDrum}, // 1拍目: スネア
		{},          // 2
		{Claves},    // 3
		{},          // 4
		{Claves},    // 5
		{},          // 6
		{Claves},    // 7
		{},          // 8
		{Claves},    // 9
		{},          // 10
		{Claves},    // 11
		{},          // 12
		{Claves},    // 13
		{},          // 14
		{Claves},    // 15
		{},          // 16
	},
	StepLen: 60000 / bpm / 4, // 16分音符の長さ
}

var basicBeat = DrumPattern{
	Name: "8Beat",
	Steps: [][]uint8{
		{BassDrum, ClosedHiHat},  // 1拍目
		{ClosedHiHat},            // 2
		{SnareDrum, ClosedHiHat}, // 3
		{OpenHiHat},              // 4
		{BassDrum, ClosedHiHat},  // 5
		{ClosedHiHat},            // 6
		{SnareDrum, ClosedHiHat}, // 7
		{ClosedHiHat},            // 8
		{BassDrum, ClosedHiHat},  // 9
		{ClosedHiHat},            // 10
		{SnareDrum, ClosedHiHat}, // 11
		{ClosedHiHat},            // 12
		{BassDrum, ClosedHiHat},  // 13
		{ClosedHiHat},            // 14
		{SnareDrum, ClosedHiHat}, // 15
		{CrashCymbal1},           // 16
	},
	StepLen: 60000 / 100 / 4,
}

var latinBeatBrazil = DrumPattern{
	Name: "LatinBrazil",
	Steps: [][]uint8{
		{BassDrum, Maracas},  // 1
		{Claves},             // 2
		{SideStick, Maracas}, // 3
		{Cowbell},            // 4

		{BassDrum, Maracas},  // 5
		{Claves},             // 6
		{SideStick, Maracas}, // 7
		{Cowbell},            // 8

		{BassDrum, Maracas},  // 9
		{Claves},             // 10
		{SideStick, Maracas}, // 11
		{Cowbell},            // 12

		{SideStick, Maracas}, // 13
		{Claves},             // 14
		{Maracas, Cowbell},   // 15
		{BassDrum, Claves},   // 16
	},
	StepLen: 60000 / 115 / 4,
}

var drumBassBeat1 = DrumPattern{
	Name: "DnB 1",
	Steps: [][]uint8{
		{BassDrum, ClosedHiHat},  // 1
		{ClosedHiHat},            // 2
		{ClosedHiHat},            // 3
		{BassDrum},               // 4
		{SnareDrum, ClosedHiHat}, // 5
		{ClosedHiHat},            // 6
		{ClosedHiHat},            // 7
		{BassDrum},               // 8
		{ClosedHiHat},            // 9
		{ClosedHiHat},            // 10
		{ClosedHiHat},            // 11
		{BassDrum},               // 12
		{SnareDrum, ClosedHiHat}, // 13
		{ClosedHiHat},            // 14
		{ClosedHiHat},            // 15
		{ClosedHiHat},            // 16
	},
	StepLen: 60000 / 174 / 4,
}

var drumBassBeat2 = DrumPattern{
	Name: "DnB 2",
	Steps: [][]uint8{
		{BassDrum, ClosedHiHat}, // 1
		{ClosedHiHat},           // 2
		{ClosedHiHat},           // 3
		{BassDrum},              // 4

		{SnareDrum, ClosedHiHat}, // 5
		{ClosedHiHat},            // 6
		{ClosedHiHat, BassDrum},  // 7
		{ClosedHiHat},            // 8

		{BassDrum, ClosedHiHat}, // 9
		{ClosedHiHat},           // 10
		{SnareDrum},             // 11
		{OpenHiHat},             // 12

		{SnareDrum, ClosedHiHat}, // 13
		{ClosedHiHat},            // 14
		{SideStick, BassDrum},    // 15
		{ClosedHiHat},            // 16
	},
	StepLen: 60000 / 176 / 4,
}

var drumBassJungle = DrumPattern{
	Name: "DnB Jungle",
	Steps: [][]uint8{
		{BassDrum, ClosedHiHat, CrashCymbal1},
		{ClosedHiHat},
		{BassDrum, ClosedHiHat},
		{ClosedHiHat},

		{SnareDrum, ClosedHiHat},
		{ClosedHiHat},
		{BassDrum, ClosedHiHat},
		{ClosedHiHat},

		{BassDrum, ClosedHiHat},
		{BassDrum, ClosedHiHat},
		{ClosedHiHat},
		{OpenHiHat, SideStick},

		{SnareDrum, ClosedHiHat},
		{BassDrum, ClosedHiHat},
		{ClosedHiHat},
		{ClosedHiHat},
	},
	StepLen: 60000 / 185 / 4,
}

var drumBassRollin = DrumPattern{
	Name: "DnB Rollin",
	Steps: [][]uint8{
		{BassDrum, ClosedHiHat},
		{ClosedHiHat},
		{BassDrum},
		{ClosedHiHat},
		{SnareDrum, ClosedHiHat},
		{ClosedHiHat},
		{BassDrum, ClosedHiHat},
		{ClosedHiHat},
		{BassDrum, ClosedHiHat},
		{ClosedHiHat},
		{SnareDrum, ClosedHiHat},
		{ClosedHiHat},
		{BassDrum, OpenHiHat},
		{ClosedHiHat},
		{SideStick, ClosedHiHat},
		{ClosedHiHat},
	},
	StepLen: 60000 / 180 / 4,
}

var drumBassStepJump = DrumPattern{
	Name: "DnB StepJump",
	Steps: [][]uint8{
		{BassDrum, ClosedHiHat, CrashCymbal1},
		{},
		{ClosedHiHat},
		{BassDrum, ClosedHiHat},
		{},
		{SnareDrum, OpenHiHat},
		{BassDrum},
		{ClosedHiHat},
		{BassDrum},
		{},
		{ClosedHiHat},
		{SnareDrum, ClosedHiHat},
		{BassDrum},
		{ClosedHiHat},
		{SideStick},
		{OpenHiHat},
	},
	StepLen: 60000 / 174 / 4,
}

var drumBassTechstep = DrumPattern{
	Name: "DnB TechStep",
	Steps: [][]uint8{
		{BassDrum, ClosedHiHat},
		{},
		{ClosedHiHat},
		{BassDrum, ClosedHiHat},
		{SnareDrum},
		{},
		{BassDrum, ClosedHiHat},
		{ClosedHiHat},
		{BassDrum},
		{ClosedHiHat},
		{SnareDrum},
		{ClosedHiHat},
		{ClosedHiHat},
		{},
		{BassDrum, SideStick},
		{ClosedHiHat},
	},
	StepLen: 60000 / 176 / 4,
}

var drumBassNeuro = DrumPattern{
	Name: "DnB Neuro",
	Steps: [][]uint8{
		{BassDrum, ClosedHiHat, CrashCymbal1}, // 1
		{},                                    // 2
		{ClosedHiHat},                         // 3
		{SnareDrum, ClosedHiHat},              // 4
		{ClosedHiHat},                         // 5
		{BassDrum, ClosedHiHat},               // 6
		{ClosedHiHat},                         // 7
		{ClosedHiHat},                         // 8
		{BassDrum, ClosedHiHat},               // 9
		{},                                    // 10
		{ClosedHiHat, SnareDrum},              // 11
		{OpenHiHat},                           // 12
		{BassDrum, ClosedHiHat},               // 13
		{ClosedHiHat},                         // 14
		{SideStick, ClosedHiHat},              // 15
		{ClosedHiHat},                         // 16
	},
	StepLen: 60000 / 178 / 4,
}

var drumBassLiquid = DrumPattern{
	Name: "DnB Liquid",
	Steps: [][]uint8{
		{BassDrum, ClosedHiHat},  // 1
		{ClosedHiHat},            // 2
		{ClosedHiHat},            // 3
		{SnareDrum, ClosedHiHat}, // 4
		{ClosedHiHat},            // 5
		{BassDrum, ClosedHiHat},  // 6
		{ClosedHiHat},            // 7
		{ClosedHiHat},            // 8
		{BassDrum, ClosedHiHat},  // 9
		{ClosedHiHat},            // 10
		{SnareDrum, ClosedHiHat}, // 11
		{ClosedHiHat},            // 12
		{ClosedHiHat},            // 13
		{ClosedHiHat},            // 14
		{ClosedHiHat},            // 15
		{ClosedHiHat},            // 16
	},
	StepLen: 60000 / 172 / 4,
}

var drumPatterns = []DrumPattern{
	metronomeBeat,
	basicBeat,
	latinBeatBrazil,
	drumBassBeat1,
	drumBassBeat2,
	drumBassJungle,
	drumBassRollin,
	drumBassStepJump,
	drumBassTechstep,
	drumBassNeuro,
	drumBassLiquid,
}

type State struct {
	Up               bool
	Down             bool
	Left             bool
	Right            bool
	Center           bool
	RotaryButton     bool
	RotaryLeft       bool
	RotaryRight      bool
	Keys             [12]bool
	ActiveNotes      [12]string // 押されているキーに対応する音名
	DrumPlaying      bool       // ドラムが再生中かどうか
	DrumPatternIndex int        // 現在のドラムパターン
}

// WS2812B LED
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

func (ws *WS2812B) WriteRaw(rawGRB []uint32) error {
	return ws.ws.WriteRaw(rawGRB)
}

func main() {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 2.8 * machine.MHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	// ディスプレイ初期化
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})
	display.SetRotation(drivers.Rotation180)
	display.ClearDisplay()
	time.Sleep(50 * time.Millisecond)

	m := midi.Port()

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

	for _, c := range colPins {
		c.Configure(machine.PinConfig{Mode: machine.PinOutput})
		c.Low()
	}

	for _, c := range rowPins {
		c.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	// ピアノ
	notes := []midi.Note{
		midi.D5,
		midi.G4,
		midi.C4,

		midi.E5,
		midi.A4,
		midi.D4,

		midi.F5,
		midi.B4,
		midi.E4,

		midi.G5,
		midi.C5,
		midi.F4,
	}

	// Note Colors
	noteColors := map[midi.Note]uint32{
		midi.C4: white,
		midi.C5: white,
		midi.D4: red,
		midi.D5: red,
		midi.E4: green,
		midi.E5: green,
		midi.F4: blue,
		midi.F5: blue,
		midi.G4: yellow,
		midi.G5: yellow,
		midi.A4: purple,
		midi.B4: pink,
	}

	state := State{
		DrumPlaying:      false,
		DrumPatternIndex: 0, // 最初のドラムパターンを選択
	}

	// LED
	colors := make([]uint32, 12)
	for i := range colors {
		colors[i] = black
	}

	ws := NewWS2812B(machine.GPIO1)

	// ロータリーエンコーダーボタン
	rotaryButton := machine.GPIO2
	rotaryButton.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	prevRotaryButton := rotaryButton.Get()

	// ジョイスティック
	machine.InitADC()

	ax := machine.ADC{Pin: machine.GPIO29}
	ax.Configure(machine.ADCConfig{})
	ay := machine.ADC{Pin: machine.GPIO28}
	ay.Configure(machine.ADCConfig{})

	// ロータリーエンコーダー
	rotaryEncoder := encoders.NewQuadratureViaInterrupt(
		machine.GPIO3,
		machine.GPIO4,
	)
	rotaryEncoder.Configure(encoders.QuadratureConfig{
		Precision: 4,
	})
	encOldValue := 0

	// 初期化待ち
	time.Sleep(1 * time.Second)

	// 初期音色
	pcOfs := 0x00 // Piano
	m.Write(programChange(cable, channel, uint8(pcOfs)))

	prevX := uint16(0)
	prevY := uint16(0)

	// 初期表示
	redraw(display, state)

	var lastDrumTime time.Time
	currentStep := 0

	xxx := time.Time{}
	currentPattern := DrumPattern{}

	cnt := 0
	for {
		// ジョイスティック X 軸処理
		{
			x := ax.Get()
			if 0x7000 <= x && x <= 0x9000 {
				x = 0x8000
			}
			if prevX != x {
				m.PitchBend(cable, channel, x>>2)
				prevX = x
			}
		}

		// ジョイスティック Y 軸処理
		{
			y := ay.Get()
			if 0x7000 <= y && y <= 0x9000 {
				y = 0x8000
			}
			if y >= 0x8000 {
				if prevY != y {
					m.ControlChange(cable, channel, midi.CCModulationWheel, byte((y-0x8000)>>8))
					prevY = y
				}
			}
		}

		// ロータリーエンコーダー位置変更時の処理
		if newValue := rotaryEncoder.Position(); newValue != encOldValue {
			// エンコーダーの変化方向を検出
			if newValue > encOldValue {
				// 右回転 - ドラムパターンを次へ
				state.DrumPatternIndex = (state.DrumPatternIndex + 1) % len(drumPatterns)
			} else {
				// 左回転 - ドラムパターンを前へ
				state.DrumPatternIndex = (state.DrumPatternIndex - 1 + len(drumPatterns)) % len(drumPatterns)
			}

			// ディスプレイ更新
			redraw(display, state)

			encOldValue = newValue
		}

		// ロータリーエンコーダーボタン処理
		currentRotaryButton := rotaryButton.Get()
		if prevRotaryButton && !currentRotaryButton {
			// ボタンが押された
			state.DrumPlaying = !state.DrumPlaying
			// ディスプレイ更新
			redraw(display, state)
		}
		prevRotaryButton = currentRotaryButton

		// ドラムパターン再生処理
		if state.DrumPlaying && state.DrumPatternIndex >= 0 && state.DrumPatternIndex < len(drumPatterns) {
			currentPattern = drumPatterns[state.DrumPatternIndex]

			// 次のステップを再生する時間になったか確認
			if lastDrumTime.IsZero() || time.Since(lastDrumTime) >= time.Duration(currentPattern.StepLen)*time.Millisecond {
				for _, note := range currentPattern.Steps[currentStep] {
					m.NoteOn(cable, drumCh, midi.Note(note), velocity)
				}
				//time.Sleep(40 * time.Millisecond)
				xxx = time.Now()

				// 次のステップに進む
				currentStep = (currentStep + 1) % len(currentPattern.Steps)
				lastDrumTime = time.Now()
			}
		}

		if !xxx.IsZero() && time.Since(xxx) >= 40*time.Millisecond {
			for _, note := range currentPattern.Steps[currentStep] {
				m.NoteOff(cable, drumCh, midi.Note(note), 0)
			}

			xxx = time.Time{}
		}

		// キーの状態更新と処理
		for i, s := range getKeys(colPins, rowPins) {
			note := notes[i]
			switch s {
			case off2on:
				m.NoteOn(cable, channel, note, velocity)

				// 対応する色をLEDに設定
				if color, exists := noteColors[note]; exists {
					colors[i] = color
				} else {
					colors[i] = white // マッピングがない場合は白色を使用
				}

				// カタカナ音名を保存
				if name, exists := noteNamesKatakana[note]; exists {
					state.ActiveNotes[i] = name
				} else {
					state.ActiveNotes[i] = ""
				}

				state.Keys[i] = true

			case on2off:
				m.NoteOff(cable, channel, note, velocity)

				// LED の色をリセット
				colors[i] = black

				// 音名をクリア
				state.ActiveNotes[i] = ""
				state.Keys[i] = false
			}
		}

		if cnt%10 == 0 {
			// LED に色を反映
			ws.WriteRaw(colors)

			// 画面を更新
			redraw(display, state)
		}

		cnt = (cnt + 1) % 10
		time.Sleep(1 * time.Millisecond)
	}
}

func redraw(d ssd1306.Device, state State) {
	d.ClearBuffer()

	sz := int16(8)

	// ドラムパターン情報を表示
	patternName := "なし"
	if state.DrumPatternIndex >= 0 && state.DrumPatternIndex < len(drumPatterns) {
		patternName = drumPatterns[state.DrumPatternIndex].Name
	}
	tinyfont.WriteLine(&d, &shnm.Shnmk12, 5, 12, patternName, displayWhite)

	if state.DrumPlaying {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, 5, 24, "State: Playing", displayWhite)
	} else {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, 5, 24, "State: Pausing", displayWhite)
	}

	// キーボード表示
	x := 128/2 - (sz+2)*2

	fontKeyNameX := x + (sz+2)*5 - 4
	fontKeyNameY := (sz+2)*(3+0) + sz + 8

	// 1行目のキー
	Rectangle(state.Keys[0], &d, x+(sz+2)*0, (sz+2)*(3+0), sz, sz, displayWhite)
	Rectangle(state.Keys[3], &d, x+(sz+2)*1, (sz+2)*(3+0), sz, sz, displayWhite)
	Rectangle(state.Keys[6], &d, x+(sz+2)*2, (sz+2)*(3+0), sz, sz, displayWhite)
	Rectangle(state.Keys[9], &d, x+(sz+2)*3, (sz+2)*(3+0), sz, sz, displayWhite)

	// 音名表示
	if state.ActiveNotes[0] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[0], displayWhite)
	}
	if state.ActiveNotes[3] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[3], displayWhite)
	}
	if state.ActiveNotes[6] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[6], displayWhite)
	}
	if state.ActiveNotes[9] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[9], displayWhite)
	}

	// 2行目のキー
	Rectangle(state.Keys[1], &d, x+(sz+2)*0, (sz+2)*(3+1), sz, sz, displayWhite)
	Rectangle(state.Keys[4], &d, x+(sz+2)*1, (sz+2)*(3+1), sz, sz, displayWhite)
	Rectangle(state.Keys[7], &d, x+(sz+2)*2, (sz+2)*(3+1), sz, sz, displayWhite)
	Rectangle(state.Keys[10], &d, x+(sz+2)*3, (sz+2)*(3+1), sz, sz, displayWhite)

	// 音名表示
	if state.ActiveNotes[1] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[1], displayWhite)
	}
	if state.ActiveNotes[4] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[4], displayWhite)
	}
	if state.ActiveNotes[7] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[7], displayWhite)
	}
	if state.ActiveNotes[10] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[10], displayWhite)
	}

	// 3行目のキー
	Rectangle(state.Keys[2], &d, x+(sz+2)*0, (sz+2)*(3+2), sz, sz, displayWhite)
	Rectangle(state.Keys[5], &d, x+(sz+2)*1, (sz+2)*(3+2), sz, sz, displayWhite)
	Rectangle(state.Keys[8], &d, x+(sz+2)*2, (sz+2)*(3+2), sz, sz, displayWhite)
	Rectangle(state.Keys[11], &d, x+(sz+2)*3, (sz+2)*(3+2), sz, sz, displayWhite)

	// 音名表示
	if state.ActiveNotes[2] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[2], displayWhite)
	}
	if state.ActiveNotes[5] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[5], displayWhite)
	}
	if state.ActiveNotes[8] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[8], displayWhite)
	}
	if state.ActiveNotes[11] != "" {
		tinyfont.WriteLine(&d, &shnm.Shnmk12, fontKeyNameX, fontKeyNameY, state.ActiveNotes[11], displayWhite)
	}

	d.Display()
}

func Rectangle(b bool, d drivers.Displayer, x int16, y int16, w int16, h int16, c color.RGBA) error {
	if b {
		tinydraw.FilledRectangle(d, x, y, w, h, c)
	} else {
		tinydraw.Rectangle(d, x, y, w, h, c)
	}
	return nil
}

var States [12]state

type state int8

const (
	off state = iota
	off2on
	off2on2
	off2on3
	off2on4
	off2onX
	on
	on2off
	on2off2
	on2off3
	on2off4
	on2offX
)

func getKeys(colPins, rowPins []machine.Pin) []state {
	colPins[0].High()
	colPins[1].Low()
	colPins[2].Low()
	colPins[3].Low()
	time.Sleep(1 * time.Millisecond)

	States[0] = updateState(States[0], rowPins[0].Get())
	States[1] = updateState(States[1], rowPins[1].Get())
	States[2] = updateState(States[2], rowPins[2].Get())

	colPins[0].Low()
	colPins[1].High()
	colPins[2].Low()
	colPins[3].Low()
	time.Sleep(1 * time.Millisecond)

	States[3] = updateState(States[3], rowPins[0].Get())
	States[4] = updateState(States[4], rowPins[1].Get())
	States[5] = updateState(States[5], rowPins[2].Get())

	colPins[0].Low()
	colPins[1].Low()
	colPins[2].High()
	colPins[3].Low()
	time.Sleep(1 * time.Millisecond)

	States[6] = updateState(States[6], rowPins[0].Get())
	States[7] = updateState(States[7], rowPins[1].Get())
	States[8] = updateState(States[8], rowPins[2].Get())

	colPins[0].Low()
	colPins[1].Low()
	colPins[2].Low()
	colPins[3].High()
	time.Sleep(1 * time.Millisecond)

	States[9] = updateState(States[9], rowPins[0].Get())
	States[10] = updateState(States[10], rowPins[1].Get())
	States[11] = updateState(States[11], rowPins[2].Get())

	return States[:]
}

func updateState(s state, btn bool) state {
	ret := s
	switch s {
	case off:
		if btn {
			ret = off2on
		}
	case off2on:
		ret = off2on2
	case off2on2:
		ret = off2on3
	case off2on3:
		ret = off2on4
	case off2on4:
		ret = off2onX
	case off2onX:
		ret = on
	case on:
		if !btn {
			ret = on2off
		}
	case on2off:
		ret = on2off2
	case on2off2:
		ret = on2off3
	case on2off3:
		ret = on2off4
	case on2off4:
		ret = on2offX
	case on2offX:
		ret = off
	}
	return ret
}

func programChange(cable, channel uint8, patch uint8) []byte {
	var pbuf [4]byte
	pbuf[0], pbuf[1], pbuf[2], pbuf[3] = ((cable&0xf)<<4)|midi.CINProgramChange, midi.MsgProgramChange|((channel-1)&0xf), patch&0x7f, 0x00
	return pbuf[:4]
}
