package main

import (
	"image/color"
	"machine"
	"time"

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
)

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

	state := State{}
	redraw(display, state)

	colors := []uint32{
		0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF,
		0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF,
		0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF,
	}

	ws := NewWS2812B(machine.GPIO1)

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

	rotaryButton := machine.GPIO2
	rotaryButton.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	machine.InitADC()

	joystickX := machine.ADC{Pin: machine.GPIO29}
	joystickX.Configure(machine.ADCConfig{})
	joystickY := machine.ADC{Pin: machine.GPIO28}
	joystickY.Configure(machine.ADCConfig{})

	joystickButton := machine.GPIO0
	joystickButton.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	rotaryEncoder := encoders.NewQuadratureViaInterrupt(
		machine.GPIO4,
		machine.GPIO3,
	)
	rotaryEncoder.Configure(encoders.QuadratureConfig{
		Precision: 4,
	})
	rotaryEncoderOldValue := 0
	rotaryEncoderTimer := 0

	for {
		// COL1
		colPins[0].High()
		colPins[1].Low()
		colPins[2].Low()
		colPins[3].Low()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			colors[0] = 0x00000000
		} else {
			colors[0] = 0xFFFFFFFF
		}
		if rowPins[1].Get() {
			colors[1] = 0x00000000
		} else {
			colors[1] = 0xFFFFFFFF
		}
		if rowPins[2].Get() {
			colors[2] = 0x00000000
		} else {
			colors[2] = 0xFFFFFFFF
		}

		// COL2
		colPins[0].Low()
		colPins[1].High()
		colPins[2].Low()
		colPins[3].Low()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			colors[3] = 0x00000000
		} else {
			colors[3] = 0xFFFFFFFF
		}
		if rowPins[1].Get() {
			colors[4] = 0x00000000
		} else {
			colors[4] = 0xFFFFFFFF
		}
		if rowPins[2].Get() {
			colors[5] = 0x00000000
		} else {
			colors[5] = 0xFFFFFFFF
		}

		// COL3
		colPins[0].Low()
		colPins[1].Low()
		colPins[2].High()
		colPins[3].Low()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			colors[6] = 0x00000000
		} else {
			colors[6] = 0xFFFFFFFF
		}
		if rowPins[1].Get() {
			colors[7] = 0x00000000
		} else {
			colors[7] = 0xFFFFFFFF
		}
		if rowPins[2].Get() {
			colors[8] = 0x00000000
		} else {
			colors[8] = 0xFFFFFFFF
		}

		// COL4
		colPins[0].Low()
		colPins[1].Low()
		colPins[2].Low()
		colPins[3].High()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			colors[9] = 0x00000000
		} else {
			colors[9] = 0xFFFFFFFF
		}
		if rowPins[1].Get() {
			colors[10] = 0x00000000
		} else {
			colors[10] = 0xFFFFFFFF
		}
		if rowPins[2].Get() {
			colors[11] = 0x00000000
		} else {
			colors[11] = 0xFFFFFFFF
		}

		if !rotaryButton.Get() {
			colors[10] = 0x0000FFFF
		}
		if !joystickButton.Get() {
			colors[4] = 0x0000FFFF
		}
		if joystickX.Get() < 0x6000 {
			colors[1] = 0x0000FFFF
		}
		if 0xA000 < joystickX.Get() {
			colors[7] = 0x0000FFFF
		}
		if 0xA000 < joystickY.Get() {
			colors[3] = 0x0000FFFF
		}
		if joystickY.Get() < 0x6000 {
			colors[5] = 0x0000FFFF
		}

		for i := 0; i < 12; i++ {
			state.Keys[i] = colors[i] == 0x00000000
		}
		state.RotaryButton = !rotaryButton.Get()
		state.Center = !joystickButton.Get()
		state.Left = joystickX.Get() < 0x6000
		state.Right = 0xA000 < joystickX.Get()
		state.Up = 0xA000 < joystickY.Get()
		state.Down = joystickY.Get() < 0x6000

		if newValue := rotaryEncoder.Position(); newValue != rotaryEncoderOldValue {
			if newValue < rotaryEncoderOldValue {
				state.RotaryRight = true
				rotaryEncoderTimer = 5
				colors[rkIndex(rotaryEncoderOldValue)] = 0xFF0000FF
			} else {
				state.RotaryLeft = true
				rotaryEncoderTimer = 5
				colors[rkIndex(rotaryEncoderOldValue)] = 0x00FF00FF
			}
			rotaryEncoderOldValue = newValue
		} else {
			if rotaryEncoderTimer > 0 {
				rotaryEncoderTimer--
			} else {
				state.RotaryLeft = false
				state.RotaryRight = false
			}
		}

		ws.WriteRaw(colors)
		redraw(display, state)
		time.Sleep(32 * time.Millisecond)
	}
}

var (
	white = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	black = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
)

type State struct {
	Up           bool
	Down         bool
	Left         bool
	Right        bool
	Center       bool
	RotaryButton bool
	RotaryLeft   bool
	RotaryRight  bool
	Keys         [12]bool
}

func redraw(d *ssd1306.Device, state State) {
	d.ClearBuffer()

	sz := int16(8)

	// joystick
	Rectangle(state.Up, d, (sz+2)*1, (sz+2)*0, sz, sz, white)
	Rectangle(state.Left, d, (sz+2)*0, (sz+2)*1, sz, sz, white)
	Rectangle(state.Right, d, (sz+2)*2, (sz+2)*1, sz, sz, white)
	Rectangle(state.Down, d, (sz+2)*1, (sz+2)*2, sz, sz, white)
	Rectangle(state.Center, d, (sz+2)*1, (sz+2)*1, sz, sz, white)

	// rotary encoder
	x := 128 - sz*3
	Triangle(state.RotaryLeft, d, x-sz, sz*2+sz/2, x-sz, sz*2-sz/2, x-2*sz, sz*2, white)
	Circle(state.RotaryButton, d, x, sz*2, sz-2, white)
	Triangle(state.RotaryRight, d, x+sz, sz*2+sz/2, x+sz, sz*2-sz/2, x+2*sz, sz*2, white)

	// Keys
	x = 128/2 - (sz+2)*2
	Rectangle(state.Keys[0], d, x+(sz+2)*0, (sz+2)*(3+0), sz, sz, white)
	Rectangle(state.Keys[1], d, x+(sz+2)*0, (sz+2)*(3+1), sz, sz, white)
	Rectangle(state.Keys[2], d, x+(sz+2)*0, (sz+2)*(3+2), sz, sz, white)

	Rectangle(state.Keys[3], d, x+(sz+2)*1, (sz+2)*(3+0), sz, sz, white)
	Rectangle(state.Keys[4], d, x+(sz+2)*1, (sz+2)*(3+1), sz, sz, white)
	Rectangle(state.Keys[5], d, x+(sz+2)*1, (sz+2)*(3+2), sz, sz, white)

	Rectangle(state.Keys[6], d, x+(sz+2)*2, (sz+2)*(3+0), sz, sz, white)
	Rectangle(state.Keys[7], d, x+(sz+2)*2, (sz+2)*(3+1), sz, sz, white)
	Rectangle(state.Keys[8], d, x+(sz+2)*2, (sz+2)*(3+2), sz, sz, white)

	Rectangle(state.Keys[9], d, x+(sz+2)*3, (sz+2)*(3+0), sz, sz, white)
	Rectangle(state.Keys[10], d, x+(sz+2)*3, (sz+2)*(3+1), sz, sz, white)
	Rectangle(state.Keys[11], d, x+(sz+2)*3, (sz+2)*(3+2), sz, sz, white)

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

func Circle(b bool, d drivers.Displayer, x0 int16, y0 int16, r int16, c color.RGBA) {
	if b {
		tinydraw.FilledCircle(d, x0, y0, r, c)
	} else {
		tinydraw.Circle(d, x0, y0, r, c)
	}
}

func Triangle(b bool, d drivers.Displayer, x0 int16, y0 int16, x1 int16, y1 int16, x2 int16, y2 int16, c color.RGBA) {
	if b {
		tinydraw.FilledTriangle(d, x0, y0, x1, y1, x2, y2, c)
	} else {
		tinydraw.Triangle(d, x0, y0, x1, y1, x2, y2, c)
	}
}

func rkIndex(idx int) int {
	for idx < 0 {
		idx += 10
	}
	switch idx % 10 {
	case 0:
		idx = 0
	case 1:
		idx = 1
	case 2:
		idx = 2
	case 3:
		idx = 5
	case 4:
		idx = 8
	case 5:
		idx = 11
	case 6:
		idx = 10
	case 7:
		idx = 9
	case 8:
		idx = 6
	case 9:
		idx = 3
	}
	return idx
}
