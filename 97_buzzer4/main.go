package main

import (
	"machine"
	"machine/usb/adc/midi"
	"time"

	"tinygo.org/x/drivers/tone"
)

var pinToPWM = map[machine.Pin]tone.PWM{
	machine.GPIO14: machine.PWM7, // for EX01
	machine.GPIO15: machine.PWM7, // for EX02
	machine.GPIO26: machine.PWM5, // for EX01
	machine.GPIO27: machine.PWM5, // for EX01
}

func main() {
	btn := machine.GPIO2
	btn.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

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

	bzrPins := []machine.Pin{machine.GPIO14, machine.GPIO26, machine.GPIO15, machine.GPIO27}
	speakers := []tone.Speaker{}

	for _, bzrPin := range bzrPins {
		pwm := pinToPWM[bzrPin]
		speaker, err := tone.New(pwm, bzrPin)
		if err != nil {
			println("failed to configure PWM")
			return
		}
		speakers = append(speakers, speaker)
	}

	notes := []tone.Note{
		tone.D5,
		tone.G4,
		tone.C4,

		tone.E5,
		tone.A4,
		tone.D4,

		tone.F5,
		tone.B4,
		tone.E4,

		tone.G5,
		tone.C5,
		tone.F4,
	}

	pressed := make([]Pressed, 0, 4)
	for {
		for i, s := range getKeys(colPins, rowPins) {
			note := notes[i]
			switch s {
			case off2on:
				//m.NoteOn(cable, channel, note, velocity)
				if len(pressed) < 4 {
					speakers[len(pressed)].SetNote(note)
					pressed = append(pressed, Pressed{
						Index: len(pressed),
						Note:  note,
					})
				}
				time.Sleep(1 * time.Millisecond)
			case on2off:
				num := -1
				for i, p := range pressed {
					if p.Note == note {
						num = i
						speakers[p.Index].Stop()
					}
				}
				if num >= 0 {
					pressed = append(pressed[:num], pressed[num+1:]...)
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
}

type Pressed struct {
	Index int
	Note  tone.Note
}

var States [12]State

type State int8

const (
	off State = iota
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

func getKeys(colPins, rowPins []machine.Pin) []State {
	// COL1
	colPins[0].High()
	colPins[1].Low()
	colPins[2].Low()
	colPins[3].Low()
	time.Sleep(1 * time.Millisecond)

	States[0] = updateState(States[0], rowPins[0].Get())
	States[1] = updateState(States[1], rowPins[1].Get())
	States[2] = updateState(States[2], rowPins[2].Get())

	// COL2
	colPins[0].Low()
	colPins[1].High()
	colPins[2].Low()
	colPins[3].Low()
	time.Sleep(1 * time.Millisecond)

	States[3] = updateState(States[3], rowPins[0].Get())
	States[4] = updateState(States[4], rowPins[1].Get())
	States[5] = updateState(States[5], rowPins[2].Get())

	// COL3
	colPins[0].Low()
	colPins[1].Low()
	colPins[2].High()
	colPins[3].Low()
	time.Sleep(1 * time.Millisecond)

	States[6] = updateState(States[6], rowPins[0].Get())
	States[7] = updateState(States[7], rowPins[1].Get())
	States[8] = updateState(States[8], rowPins[2].Get())

	// COL4
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

func updateState(s State, btn bool) State {
	ret := s
	switch s {
	case off:
		if btn {
			ret = off2on
		}
	case on:
		if !btn {
			ret = on2off
		}
	case on2offX:
		ret = off
	default:
		ret = s + 1
	}
	return ret
}

var pbuf [4]byte

func programChange(cable, channel uint8, patch uint8) []byte {
	pbuf[0], pbuf[1], pbuf[2], pbuf[3] = ((cable&0xf)<<4)|midi.CINProgramChange, midi.MsgProgramChange|((channel-1)&0xf), patch&0x7f, 0x00
	return pbuf[:4]
}
