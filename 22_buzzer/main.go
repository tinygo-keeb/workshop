package main

// Please connect a piezo buzzer to the 3V3 and EX01 pins on the back terminal.
//
// | EX01 | EX03 | 3V3 | SDA0 | 3V3 | 3V3 |     |        GROVE            |
// | EX02 | EX04 | GND | SCL0 | GND | GND | - - | GND | 3V3 | SDA0 | SCL0 |

import (
	"machine"
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
	bzrPin := machine.GPIO14
	pwm := pinToPWM[bzrPin]
	speaker, err := tone.New(pwm, bzrPin)
	if err != nil {
		println("failed to configure PWM")
		return
	}

	song := []tone.Note{
		tone.C5,
		tone.D5,
		tone.E5,
		tone.F5,
		tone.G5,
		tone.A5,
		tone.B5,
		tone.C6,
		tone.C6,
		tone.B5,
		tone.A5,
		tone.G5,
		tone.F5,
		tone.E5,
		tone.D5,
		tone.C5,
	}

	for {
		for _, val := range song {
			speaker.SetNote(val)
			time.Sleep(time.Second / 2)
		}
	}
}
