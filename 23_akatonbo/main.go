package main

// Please connect a piezo buzzer to the 3V3 and EX01 pins on the back terminal.
//
// | EX01 | EX03 | 3V3 | SDA0 | 3V3 | 3V3 |     |        GROVE            |
// | EX02 | EX04 | GND | SCL0 | GND | GND | - - | GND | 3V3 | SDA0 | SCL0 |

import (
	"fmt"
	"machine"
	"time"
	"tinygo.org/x/drivers/tone"
)

// タクトスイッチのピン定義
const BUTTON_PIN = machine.GPIO2

// ボタンの初期化
func initButton() {
	BUTTON_PIN.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
}

// ボタンが押されたかチェック（チャタリング対策付き）
func isButtonPressed() bool {
	if !BUTTON_PIN.Get() { // プルアップなので押されると false
		time.Sleep(50 * time.Millisecond) // チャタリング対策
		if !BUTTON_PIN.Get() {            // 再確認
			// ボタンが離されるまで待つ
			for !BUTTON_PIN.Get() {
				time.Sleep(10 * time.Millisecond)
			}
			time.Sleep(50 * time.Millisecond) // 離した後のチャタリング対策
			return true
		}
	}
	return false
}

type NoteWithDuration struct {
	Note     tone.Note
	Duration time.Duration
}

var pinToPWM = map[machine.Pin]tone.PWM{
	machine.GPIO14: machine.PWM7, // for EX01
}

// ブザーを初期化する関数
func initBuzzer() (tone.Speaker, error) {
	bzrPin := machine.GPIO14
	pwm := pinToPWM[bzrPin]
	speaker, err := tone.New(pwm, bzrPin)
	if err != nil {
		return tone.Speaker{}, err
	}
	return speaker, nil
}

func getSong() []interface{} {
	const bpm = 100
	beat := time.Minute / bpm
	eighth := beat / 2
	quarter := beat
	dottedQuarter := beat * 3 / 2
	dottedHalf := beat * 3

	return []interface{}{
		"夕焼け",
		NoteWithDuration{tone.G4, eighth},
		NoteWithDuration{tone.C5, eighth},
		NoteWithDuration{tone.C5, dottedQuarter},
		NoteWithDuration{tone.D5, eighth},

		"小焼けの",
		NoteWithDuration{tone.E5, eighth},
		NoteWithDuration{tone.G5, eighth},
		NoteWithDuration{tone.C6, eighth},
		NoteWithDuration{tone.A5, eighth},
		NoteWithDuration{tone.G5, quarter},

		"赤とんぼ",
		NoteWithDuration{tone.A5, eighth},
		NoteWithDuration{tone.C5, eighth},
		NoteWithDuration{tone.C5, quarter},
		NoteWithDuration{tone.D5, quarter},
		NoteWithDuration{tone.E5, dottedHalf},

		"負われて",
		NoteWithDuration{tone.E5, eighth},
		NoteWithDuration{tone.A5, eighth},
		NoteWithDuration{tone.G5, dottedQuarter},
		NoteWithDuration{tone.A5, eighth},

		"見たのは",
		NoteWithDuration{tone.C6, eighth},
		NoteWithDuration{tone.A5, eighth},
		NoteWithDuration{tone.G5, eighth},
		NoteWithDuration{tone.A5, eighth},
		NoteWithDuration{tone.G5, eighth},
		NoteWithDuration{tone.E5, eighth},

		"何時の日か",
		NoteWithDuration{tone.G5, eighth},
		NoteWithDuration{tone.E5, eighth},
		NoteWithDuration{tone.C5, eighth},
		NoteWithDuration{tone.E5, eighth},
		NoteWithDuration{tone.D5, eighth},
		NoteWithDuration{tone.C5, eighth},
		NoteWithDuration{tone.C5, dottedHalf},
	}
}

// 音符名と周波数を取得する関数（全オクターブ対応）
func getNoteName(note tone.Note) string {
	switch note {
	// 3オクターブ
	case tone.C3:
		return "ド(C3) 131Hz"
	case tone.CS3:
		return "ド#(C#3) 139Hz"
	case tone.D3:
		return "レ(D3) 147Hz"
	case tone.DS3:
		return "レ#(D#3) 156Hz"
	case tone.E3:
		return "ミ(E3) 165Hz"
	case tone.F3:
		return "ファ(F3) 175Hz"
	case tone.FS3:
		return "ファ#(F#3) 185Hz"
	case tone.G3:
		return "ソ(G3) 196Hz"
	case tone.GS3:
		return "ソ#(G#3) 208Hz"
	case tone.A3:
		return "ラ(A3) 220Hz"
	case tone.AS3:
		return "ラ#(A#3) 233Hz"
	case tone.B3:
		return "シ(B3) 247Hz"

	// 4オクターブ
	case tone.C4:
		return "ド(C4) 262Hz"
	case tone.CS4:
		return "ド#(C#4) 277Hz"
	case tone.D4:
		return "レ(D4) 294Hz"
	case tone.DS4:
		return "レ#(D#4) 311Hz"
	case tone.E4:
		return "ミ(E4) 330Hz"
	case tone.F4:
		return "ファ(F4) 349Hz"
	case tone.FS4:
		return "ファ#(F#4) 370Hz"
	case tone.G4:
		return "ソ(G4) 392Hz"
	case tone.GS4:
		return "ソ#(G#4) 415Hz"
	case tone.A4:
		return "ラ(A4) 440Hz"
	case tone.AS4:
		return "ラ#(A#4) 466Hz"
	case tone.B4:
		return "シ(B4) 494Hz"

	// 5オクターブ
	case tone.C5:
		return "ド(C5) 523Hz"
	case tone.CS5:
		return "ド#(C#5) 554Hz"
	case tone.D5:
		return "レ(D5) 587Hz"
	case tone.DS5:
		return "レ#(D#5) 622Hz"
	case tone.E5:
		return "ミ(E5) 659Hz"
	case tone.F5:
		return "ファ(F5) 698Hz"
	case tone.FS5:
		return "ファ#(F#5) 740Hz"
	case tone.G5:
		return "ソ(G5) 784Hz"
	case tone.GS5:
		return "ソ#(G#5) 831Hz"
	case tone.A5:
		return "ラ(A5) 880Hz"
	case tone.AS5:
		return "ラ#(A#5) 932Hz"
	case tone.B5:
		return "シ(B5) 988Hz"

	// 6オクターブ
	case tone.C6:
		return "ド(C6) 1047Hz"
	case tone.CS6:
		return "ド#(C#6) 1109Hz"
	case tone.D6:
		return "レ(D6) 1175Hz"
	case tone.DS6:
		return "レ#(D#6) 1245Hz"
	case tone.E6:
		return "ミ(E6) 1319Hz"
	case tone.F6:
		return "ファ(F6) 1397Hz"
	case tone.FS6:
		return "ファ#(F#6) 1480Hz"
	case tone.G6:
		return "ソ(G6) 1568Hz"
	case tone.GS6:
		return "ソ#(G#6) 1661Hz"
	case tone.A6:
		return "ラ(A6) 1760Hz"
	case tone.AS6:
		return "ラ#(A#6) 1865Hz"
	case tone.B6:
		return "シ(B6) 1976Hz"

	default:
		return "不明な音"
	}
}

// 楽曲を演奏する関数（画面表示付き）
func playSong(speaker tone.Speaker, song []interface{}, display *Display) {
	noteIndex := 0
	for _, element := range song {
		switch v := element.(type) {
		case string:
			display.UpdateStatus(v)
		case NoteWithDuration:
			noteIndex++
			noteName := getNoteName(v.Note)
			display.PrintLine(fmt.Sprintf("%d: %s", noteIndex, noteName))

			speaker.SetNote(v.Note)
			time.Sleep(v.Duration)
			speaker.Stop()
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// デモ用のスクロール表示（削除）

func main() {
	fmt.Println("プログラム開始")

	// ディスプレイの初期化
	display := InitDisplay()
	fmt.Println("ディスプレイ初期化完了")
	display.PrintLine("ディスプレイ初期化完了")

	// ボタンの初期化
	initButton()
	fmt.Println("ボタン初期化完了")
	display.PrintLine("ボタン初期化完了")

	// ステータス行に「テスト音楽」を表示
	display.UpdateStatus("テスト音楽")
	fmt.Println("ステータス行表示完了")
	display.PrintLine("ステータス行表示完了")

	// ブザーの初期化
	speaker, err := initBuzzer()
	if err != nil {
		fmt.Println("failed to configure PWM")
		display.PrintLine("PWM設定エラー")
		return
	}
	fmt.Println("ブザー初期化完了")
	display.PrintLine("ブザー初期化完了")

	// 楽曲データの取得
	song := getSong()
	fmt.Println("楽曲データ取得完了")
	display.PrintLine("楽曲データ取得完了")

	// 演奏回数カウンター
	playCount := 0

	// メインループ
	display.PrintLine("ボタンを押す")
	fmt.Println("ボタン待機中...")
	display.UpdateStatus("待機中")

	for {
		if isButtonPressed() {
			playCount++
			fmt.Printf("ボタンが押されました - 演奏回数: %d\n", playCount)
			display.PrintLine(fmt.Sprintf("演奏開始 (%d回目)", playCount))
			display.UpdateStatus("演奏中...")

			playSong(speaker, song, display)

			fmt.Printf("演奏完了 - %d回目\n", playCount)
			display.PrintLine(fmt.Sprintf("演奏完了 (%d回目)", playCount))
			display.PrintLine("ボタンを押す")
			display.UpdateStatus("待機中")
		}
		time.Sleep(100 * time.Millisecond) // CPU負荷軽減
	}
}
