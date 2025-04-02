# tinygo-keeb/workshop

![](./images/zero-kb02.jpg)

This page is for the TinyGo Keeb Tour that started on 2024/08/04.
If you have any questions, please create in an issue in this repository or contact [twitter:sago35tk](https://x.com/sago35tk).

For hardware assembly, please refer to the following.


* [ビルドガイド](./buildguide.md)
* [build guide (english)](./buildguide_EN.md)

To this page QR code is here.  

![](./images/qr_tinygo_keeb_workshop_top_english.png)

# Environment Setup

## TinyGo Installation

The following dependencies are required.
For TinyGo, we have included the URL for v0.35.0, which is the latest version at the time of writting this page, but please use the latest version available.

* Git
    * https://git-scm.com/downloads
    * Not required for Go / TinyGo, but necessary for this workshop
* Go
    * https://go.dev/dl/
        * installation instructions : https://go.dev/doc/install
* TinyGo
    * https://github.com/tinygo-org/tinygo/releases/latest
        * installation instructions : https://tinygo.org/getting-started/install/

Note, that there is a version combination between Go and TinyGo.
Basically, TinyGo must be used with the latest and most recent version of Go.

| TinyGo | Compatible Go |
| ------ | ----------- |
| 0.35.0 | 1.23 - 1.22 |


You can check if the installation was successful or not at the following

```
$ tinygo version
tinygo version 0.35.0 windows/amd64 (using go version go1.23.6 and LLVM version 18.1.2)
```

```
$ tinygo build -o out.uf2 --target waveshare-rp2040-zero --size short examples/serial
   code    data     bss |   flash     ram
   7836     108    3152 |    7944    3260
```

```
$ tinygo flash --target waveshare-rp2040-zero --size short examples/serial
   code    data     bss |   flash     ram
   7932     108    3168 |    8040    3276

$ tinygo monitor --target waveshare-rp2040-zero
Connected to COM4. Press Ctrl-C to exit.
hello world!
hello world!
hello world!
```

### Windows + WSL2

You can use the linux version of TinyGo for Ubuntu on WSL2.
However, WSL2 cannot directly access USB devices connected to the Windows host.
Even when using WSL2, it is basically better to install the Windows version of TinyGo on the Windows path.
In this case, it is necessary to install the Windows version of Go as well.

If you really want to communicate with TinyGo on WSL2, you can use usbipd as follows.
However, it is not very comfortable because you need to attach usbipd every time you do tinygo flash.

* [Run tinygo monitor on raspberry pi pico with tinygo installed on WSL2 (japanese)](https://qiita.com/kn12abc/items/d6bfc172cf08d9be6e1a)

### Linux setup

To use `tinygo flash`, `tinygo monitor` or `Vial` on Linux, you need to configure udev rules.
Create `/etc/udev/rules.d/99-zero-kb02-udev.rules` with the following contents and restart.

```
# RP2040
# ref: https://docs.platformio.org/en/latest/core/installation/udev-rules.html
ATTRS{idVendor}=="2e8a", ATTRS{idProduct}=="[01]*", MODE:="0666", ENV{ID_MM_DEVICE_IGNORE}="1", ENV{ID_MM_PORT_IGNORE}="1"

# Vial
# ref: https://get.vial.today/manual/linux-udev.html
KERNEL=="hidraw*", SUBSYSTEM=="hidraw", ATTRS{serial}=="*vial:f64c2b3c*", MODE="0660", GROUP="users", TAG+="uaccess", TAG+="udev-acl"
```

Here is a file with the same contents as above.


* [./99-zero-kb02-udev.rules](./99-zero-kb02-udev.rules)

This file was created from the following documents.
Please refer to them for more details:

* https://docs.platformio.org/en/latest/core/installation/udev-rules.html
* https://get.vial.today/manual/linux-udev.html

### TinyGo dev branch version

If you want to use the latest version under development, download Artifact > release-double-zipped built by GitHub Actions.

* windows
    * https://github.com/tinygo-org/tinygo/actions/workflows/windows.yml?query=branch%3Adev
* linux
    * https://github.com/tinygo-org/tinygo/actions/workflows/linux.yml?query=branch%3Adev
* macos
    * https://github.com/tinygo-org/tinygo/actions/workflows/build-macos.yml?query=branch%3Adev


For more information see:

* [TinyGo の開発版のビルド方法と、ビルドせずに開発版バイナリを手に入れる方法](https://qiita.com/sago35/items/33e63ca5073f572ad69c#pr-%E5%86%85%E3%81%A7%E4%BD%9C%E6%88%90%E3%81%95%E3%82%8C%E3%81%9F%E3%83%90%E3%82%A4%E3%83%8A%E3%83%AA%E3%82%92%E4%BD%BF%E3%81%86)
* [Building from source TinyGo](https://tinygo.org/docs/guides/build/)

## LSP / gopls Support

TinyGo places packages like the machine package in GOROOT, so until you configure it, gopls will show errors and you won't be able to jump to definitions like machine.LED.
TinyGo has an unfamiliar package structure (even if you know Go well) and many build-tag branches, so it's better to set up LSP for TinyGo.

The official documentation is available here:

* https://tinygo.org/docs/guides/ide-integration/

For VSCode, it's good to install the TinyGo extension.
For Vim (+ vim-lsp), try `github.com/sago35/tinygo.vim`.

For information in Japanese, see:

* [How to integrate TinyGo + 'VSCode or Vim (or other LSP-compatible editors)' with gopls](https://qiita.com/sago35/items/c30cbce4a0a3e12d899c)
* [Settings for TinyGo + Vim with gopls](https://qiita.com/sago35/items/f0b058ed5c32b6446834)


# Development Target

TinyGo Keeb Tour will use a home-made keyboard/macro pad called zero-kb02.
Microcontroller is RP2040 (Cortex M0+) and microcontroller board is [Waveshare RP2040-Zero](https://www.waveshare.com/rp2040-zero.htm).

! [](. /images/zero-kb02.jpg)

Main features include.

* Waveshare RP2040-Zero.
    * https://www.waveshare.com/rp2040-zero.htm
    * https://www.waveshare.com/wiki/RP2040-Zero
* 12 keys with RGB LEDs
* 2-axis analog joystick for mouse cursor movement
* Rotary encoder
* OLED display - 128x64 monochrome
* [GROVE connector](https://lab.seeed.co.jp/entry/2019/10/25/120432)
* 2x6 pin socket on the back

Schematics, firmware, pinouts, etc. can be found at:


* https://github.com/sago35/keyboards
    * Schematics : [kicanvas](https://kicanvas.org/?github=https%3A%2F%2Fgithub.com%2Fsago35%2Fkeyboards%2Ftree%2Fmain%2Fzero-kb02%2Fzero-kb02)

## Assembly

For soldering and assembly instructions, please refer to the build guide:

* [Build Guide (English)](./buildguide_EN.md)
* [Build Guide (日本語)](./buildguide.md)


# TinyGo Basics

First, clone this repository somewhere.
From now on, we'll execute commands from the root of this repository.
If you want to modify the source code, please edit the local cod

```
$ git clone https://github.com/tinygo-keeb/workshop

$ cd workshop

# Launch VS Code or other editor
$ code .
```

The source code is in paths like `./00_basic` or `./12_matrix_basic`.

## build & flash (method 1)

You can build and flash from the command line with TinyGo, but we'll learn how to flash manually too.
Boards with RP2040 can boot into the bootloader by pressing the BOOT/BOOTSEL button while resetting (pressing the reset button or connecting to USB).
When in bootloader mode, the PC recognizes it as an external drive, so you can flash it by copying the binary file (`*.uf2`) to the newly created external drive.


Try flashing the following:

* [00_basic.uf2](https://github.com/tinygo-keeb/workshop/releases/download/0.1.0/00_basic.uf2)

If the LEDs on the key switches light up, the write was successful.

※This flashing method is also valid for uf2 files created outside of TinyGo. Putting the device into bootloader mode will help if for some reason the `flash` command does not work.

To create the 00_basic.uf2 yourself, execute the following command.
If no error messages are displayed and `00_basic.uf2` is created, it's successful.

```
$ tinygo build -o 00_basic.uf2 --target waveshare-rp2040-zero --size short ./00_basic/
   code    data     bss |   flash     ram
  20420     192    3240 |   20612    3432
```


## build & flash (method 2) + serial monitor

You can also build and flash at once using the tinygo flash command.
If no error messages are displayed, the flash has completed successfully.
If it fails in a Linux environment, check the udev rules settings mentioned earlier.

```
$ tinygo flash --target waveshare-rp2040-zero --size short examples/serial
   code    data     bss |   flash     ram
   7836     108    3152 |    7944    3260
```

The `examples/serial` written above is an example that displays `hello world!` to the serial output.
You can check if its working with the following:

```
$ tinygo monitor
Connected to COM7. Press Ctrl-C to exit.
hello world!
hello world!
hello world!
```


If you can't connect properly, check the port and add the --port option.
The waveshare-rp2040-zero uses the same USB VID/PID as other boards with the RP2040 microcontroller, so the Boards section might not display correctly, but don't worry about it.

```
$ tinygo ports
Port                 ID        Boards
COM7                 2E8A:0003 waveshare-rp2040-zero

$ tinygo monitor --port COM7
Connected to COM7. Press Ctrl-C to exit.
hello world!
hello world!
hello world!
```

There's also a way to run `tinygo flash --monitor` which combines `tinygo flash` and `tinygo monitor`.
However, depending on the environment, it may connect to the wrong port, so if it doesn't work well, run them separately as shown above.

```
$ tinygo flash --target waveshare-rp2040-zero --size short --monitor examples/serial
   code    data     bss |   flash     ram
   7836     108    3152 |    7944    3260
Connected to COM7. Press Ctrl-C to exit.
hello world!
hello world!
hello world!
```

### Troubleshooting tinygo flash doesn't work on macOS 15 Sequoia

Add `NO NAME` to `msd-volume-name` in `$TINYGOROOT/targets/rp2040.json`.  
You can find $TINYGOROOT with `tinygo env`.  

The modified JSON file is as follows:  

```json
{
    "inherits": ["cortex-m0plus"],
    "build-tags": ["rp2040", "rp"],
    "flash-1200-bps-reset": "true",
    "flash-method": "msd",
    "serial": "usb",
    "msd-volume-name": ["RPI-RP2", "NO NAME"],
    "msd-firmware-name": "firmware.uf2",
    "binary-format": "uf2",
    "uf2-family-id": "0xe48bff56",
    "rp2040-boot-patch": true,
    "extra-files": [
        "src/device/rp/rp2040.s"
    ],
    "linkerscript": "targets/rp2040.ld",
    "openocd-interface": "picoprobe",
    "openocd-transport": "swd",
    "openocd-target": "rp2040"
}
```

* https://github.com/tinygo-org/tinygo/issues/4519


### How to stop "Disk Not Ejected Properly" notifications from accumulating on macOS

Open a terminal and run the following, then restart:

```
$ sudo defaults write /Library/Preferences/SystemConfiguration/com.apple.DiskArbitration.diskarbitrationd.plist DADisableEjectNotification -bool YES && sudo pkill diskarbitrationd
```

To revert back, run the following and restart:

```
$ sudo defaults delete /Library/Preferences/SystemConfiguration/com.apple.DiskArbitration.diskarbitrationd.plist DADisableEjectNotification && sudo pkill diskarbitrationd
```

See: https://www.reddit.com/r/mac/comments/vsn1t6/how_to_disable_not_ejected_safely_notification_on/



### If you absolutely cannot flash to the device

The following possibilities exist:

* There's a problem with the USB cable
  * Check if it's recognized with `tinygo ports` or as a drive (try booting into bootloader mode)
* Writing to external drives is restricted
  * Company computers may restrict writing for security reasons
  * In this case, neither tinygo flash nor copying uf2 files will work



## LED Blink

Run the following:

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./01_blinky1/
```

If the RGB LED on the RP2040 Zero lights up successfully, try changing the color or blink speed by modifying the source code.
You can set a `color.Color` at the black or white locations below:

```go
// 01_blinky1/main.go
for {
    time.Sleep(time.Millisecond * 500)
    ws.PutColor(black)
    time.Sleep(time.Millisecond * 500)
    ws.PutColor(white)
}
```

Other color examples are as follows.
You can set any color by specifying RGBA.
You can (somewhat) reduce the brightness by making the 0xFF values smaller.

```
red     = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0x00}
green   = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0x00}
blue    = color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0x00}
yellow  = color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0x00}
cyan    = color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0x00}
magenta = color.RGBA{R: 0xFF, G: 0x00, B: 0xFF, A: 0x00}
```


## LED Blink (2)

Let's light up the keys.
The board has 12 SK2812MINI-E LEDs (WS2812B compatible).
They are mounted in the following positions/order:

```
 0  3  6  9
 1  4  7 10
 2  5  8 11
```

Run the following:
If it works successfully, try changing the color, blink speed, or blink pattern by modifying the source code.

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./02_blinky2/
```

It's almost the same as before, but `WriteRaw()` is used instead of `PutColor()`.
Here, `colors[0][:i+1]` is specified, but if `[:1]` is specified, only the first LED is set.
`[:4]` would change a total of 4 LEDs.

```go
// ./02_blinky2/main.go
ws.WriteRaw(colors[0][:i+1])
```

`WriteRaw()` allows you to specify colors as uint32.
The values are set as Green / Red / Blue, 8 bits each from the most significant bit.
For example:

```go
// ./02_blinky2/main.go
colors := []uint32{
    0xFFFFFFFF, // white
    0xFF0000FF, // green
    0x00FF00FF, // red
    0x0000FFFF, // blue
}
```

You can (somewhat) reduce the brightness by making the 0xFF values smaller.

## USB CDC Hello World

Let's also try USB CDC, which is useful for printf debugging and other purposes.
USB CDC stands for Universal Serial Bus Communications Device Class, and roughly speaking, it's for communication between a computer and a microcontroller throught the USB cable.
Rather than explaining, it's easier to understand by trying it, so first run the following:

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short examples/serial

$ tinygo monitor
```

On Windows, it will look like this:

```
$ tinygo flash --target waveshare-rp2040-zero --size short examples/serial
   code    data     bss |   flash     ram
   7836     108    3152 |    7944    3260

$ tinygo monitor
Connected to COM7. Press Ctrl-C to exit.
hello world!
hello world!
hello world!
(omitted)
```

examples/serial is a source ([./03_usbcdc-serial](./03_usbcdc-serial)) like the following.
It repeatedly displays `hello world!` and then waits for 1 second.
Try changing the wait time, display string, or using fmt.Printf() for writing.

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./03_usbcdc-serial/
```

Standard input can be handled with a code like the following: ([./04_usbcdc-echo/](./04_usbcdc-echo/)).
After pressing `Enter`/`Return`, you need to press `Ctrl-j` for a line break.

```go
// ./04_usbcdc-echo/main.go
package main

import (
        "bufio"
        "fmt"
        "os"
)

func main() {
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
                fmt.Printf("you typed : %s\n", scanner.Text())
        }
}
```



## Rotary Encoder

You can use encoders/quadrature interrupts from tinygo-org/drivers.

* https://github.com/tinygo-org/drivers/blob/release/examples/encoders/quadrature-interrupt/main.go

Here's the configuration adjusted for zero-kb02:

```
// ./05_rotary/main.go
enc := encoders.NewQuadratureViaInterrupt(
    machine.GPIO3,
    machine.GPIO4,
)
enc.Configure(encoders.QuadratureConfig{
    Precision: 4,
})
```

You can check and play with the example 05_rotary.
When you turn the rotary encoder, the value display will update.
It might be interesting to link it with the LEDs as an exercise.

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./05_rotary/
   code    data     bss |   flash     ram
   8276     108    3624 |    8384    3732

$ tinygo monitor
Connected to COM7. Press Ctrl-C to exit.
value:  -1
value:  -2
value:  -1
value:  0
value:  1
value:  2
(omitted)
```

Note that the rotary encoder can also be used as a button when pressed.
Getting the pressed state of the rotary encoder will be discussed in the next section.

## Getting the Pressed State of the Rotary Encoder

When the rotary encoder is pressed, it connects to GND and goes Low.
If you pull it up, it will be High when not pressed.

The basic code is as follows:

```go
// ./13_rotary_button/main.go
if !btn.Get() {
    println("pressed")
} else {
    println("released")
}
```

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./13_rotary_button/

$ tinygo monitor
```

When you press the rotary encoder, `pressed` will be output to the terminal running `tinygo monitor`.


## Analog Joystick

The analog joystick can be recognized as a digital value when pressed, and as analog values for the X and Y axes.
So it can be handled as follows:

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./06_joystick/
   code    data     bss |   flash     ram
  56792    1536    3176 |   58328    4712

$ tinygo monitor
Connected to COM7. Press Ctrl-C to exit.
7440 8000 false
7130 7F90 true
(omitted)
```

From the left, it shows X-axis value (voltage value), Y-axis value, and whether it's pressed.
When not doing anything, values close to 0x8000 are displayed.

## OLED

You can use ssd1306/i2c_128x64 from tinygo-org/drivers.

* https://github.com/tinygo-org/drivers/tree/release/examples/ssd1306/i2c_128x64/main.go

Here's the configuration adjusted for zero-kb02:

```go
// ./07_oled/main.go
machine.I2C0.Configure(machine.I2CConfig{
    Frequency: machine.TWI_FREQ_400KHZ,
    SDA:       machine.GPIO12,
    SCL:       machine.GPIO13,
})
display := ssd1306.NewI2C(machine.I2C0)
display.Configure(ssd1306.Config{
    Address: 0x3C,
    Width:   128,
    Height:  64,
})
```

You can write and check operation with the following command:

```go
$ tinygo flash --target waveshare-rp2040-zero --size short ./07_oled/
```

※As of 2024/08/04, OLED drawing sometimes stops (currently investigating)

### Drawing Shapes

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./08_oled_tinydraw/
```

### Drawing Text

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./09_oled_tinyfont/
```

### Rotating the Screen

On zero-kb02, the OLED is mounted upside down, so you need to rotate the screen somehow.
Here, let's try rotating it via hardware.

As shown below, you can rotate with Rotation in the Config.
For SSD1306, only `Rotation0` (no rotation) and `Rotation180` (inverted) can be used.

```go
// ./10_oled_rotated/main.go
display.Configure(ssd1306.Config{
    Address:  0x3C,
    Width:    128,
    Height:   64,
    Rotation: drivers.Rotation180,
})
```

You can also rotate with `SetRotation()` outside of Configure() time.

```go
display.SetRotation(drivers.Rotation180)
```

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./16_oled_inverted_hw/
```

### Rotating the Screen 90 Degrees

We have seen hwo to implement no rotation or inversion, however, in some cases, you might want to rotate the display 90 degrees to use it in portrait mode.
In this case, you need to rotate via software.
Screen drawing basically corresponds to the following Displayer interface, so we define a Displayer that can rotate the screen.

```go
// https://github.com/tinygo-org/drivers/blob/release/displayer.go
type Displayer interface {
    // Size returns the current size of the display.
    Size() (x, y int16)

    // SetPizel modifies the internal buffer.
    SetPixel(x, y int16, c color.RGBA)

    // Display sends the buffer (if any) to the screen.
    Display() error
}
```

Here, we've defined the following.
We embed Displayer in the struct and process the x and y values of Size and SetPixel.

```go
// ./16_oled_inverted_hw/main.go
type RotatedDisplay struct {
        drivers.Displayer
}

func (d *RotatedDisplay) Size() (x, y int16) {
        return y, x
}

func (d *RotatedDisplay) SetPixel(x, y int16, c color.RGBA) {
        _, sy := d.Displayer.Size()
        d.Displayer.SetPixel(y, sy-x, c)
}
```

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./10_oled_rotated/
```

### Animating

You can update the screen without flickering by using `display.ClearBuffer()` and `display.Display()`.

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./11_oled_animation/
```

### Displaying Japanese

Currently, either BDF or OTF/TTF fonts can be displayed.  
For small displays with 1-bit color like zero-kb02, BDF fonts are more suitable.  
You can use them as follows:  

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./17_oled_japanese_font/
```

## Getting Key Press States

zero-kb02 uses a wiring method called a matrix for its key connections.
As shown in the circuit below, 12 switches are connected using 7 pins in a COL:4 x ROW:3 configuration.

![](./images/matrix.png)

The reading process works as follows:

1. Set only COL1 to High, and set COL2 through COL4 to Low
2. Wait a moment
3. Read ROW1 through ROW3 in that state
    * If ROW1 is High, SW1 is pressed
    * If ROW2 is High, SW5 is pressed
    * If ROW3 is High, SW9 is pressed

Next, if only COL2 is set to High, SW2 / SW6 / SW10 can be read, and so on.

Matrix wiring is a widely used connection method in custom keyboards, so let's implement it.
A simple implementation of the above would look like this:

```go
// ./12_matrix_basic/main.go
colPins[0].High()
colPins[1].Low()
colPins[2].Low()
colPins[3].Low()
time.Sleep(1 * time.Millisecond)

if rowPins[0].Get() {
    fmt.Printf("sw1 pressed\n")
}
if rowPins[1].Get() {
    fmt.Printf("sw5 pressed\n")
}
if rowPins[2].Get() {
    fmt.Printf("sw9 pressed\n")
}
```

By implementing the same pattern for all columns, you can detect the press state of all keys.

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short --monitor ./12_matrix_basic/
```

By organizing the loops and making the number of keys variable, you can move closer to a keyboard firmware.

※ For those who want to learn more about matrix wiring, please see:
https://blog.ikejima.org/make/keyboard/2019/12/14/keyboard-circuit.html

## USB HID Keyboard Using Pin Input

Let's create a USB HID Keyboard using the rotary encoder's press state.
The following code allows you to link the press state with the `A` key.
In TinyGo, you can import `machine/usb/hid/keyboard` and call `keyboard.Port()`, so your device will act as a keyboard recognized by the computer.

```go
// ./14_hid_keyboard/main.go
kb := keyboard.Port()
for {
    if !btn.Get() {
        kb.Down(keyboard.KeyA)
    } else {
        kb.Up(keyboard.KeyA)
    }
}
```

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./14_hid_keyboard/
```

Press the rotary encoder to verify it's working properly.

## USB HID Mouse Using Pin Input

Now let's create a USB HID Mouse using the rotary encoder's press state.
With the following code, pressing the button becomes a mouse left click.

```go
// ./15_hid_mouse/main.go
m := mouse.Port()
for {
    if !btn.Get() {
        m.Press(mouse.Left)
    } else {
        m.Release(mouse.Left)
    }
}
```

```shell
$ tinygo flash --target waveshare-rp2040-zero --size short ./15_hid_mouse/
```

Press the rotary encoder to verify it's working properly.

## Using MIDI

TinyGo supports USB MIDI, so you can make MIDI sound sources or MIDI instruments.  
You can use the 12 keys and the rotary encoder press.  

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./18_midi/
```

After creation, you can test it on sites like:  

* https://midi.city/
* https://virtualpiano.eu/

In Windows environments, MIDI-OX is a good option:  

* http://www.midiox.com/

# Using sago35/tinygo-keyboard

The necessary elements for a custom keyboard vary from person to person.
However, some common requirements include:

* Layer functionality
* Ability to change settings without rebuilding
* Package-based methods for reading various switches

Creating these from scratch each time can be challenging and time consuming, so using some kind of library is common.
Here, let's create a custom keyboard using the package sago35/tinygo-keyboard.

With sago35/tinygo-keyboard, you can easily implement the following features:

* Support for various key input methods (matrix, GPIO, rotary encoder, etc.)
    * Possibility to write your own extensions
* Layer functionality
* Integration with mouse clicks and pointer movement
* Split keyboard support via TRRS cable
* Configuration changes through Vial via web browser
    * Keymaps
    * Layers
    * Matrix tester (key switch press test)
    * Macro functionality

The Vial integration is particularly important as it makes it easy to change settings according to individual preferences without the need to re-flash the firmware again.
Vial is available at the following URL and can be accessed from Edge/Chrome browsers that support the WebHID API:

* https://vial.rocks/

## Basic Usage

For detailed usage instructions, please refer to:

* [Creating a Custom Keyboard with sago35/tinygo-keyboard](https://qiita.com/sago35/items/b008ed03cd403742e7aa)
* [Create Your Own Keyboard with sago35/tinygo-keyboard](https://dev.to/sago35/create-your-own-keyboard-with-sago35tinygo-keyboard-4gbj)

## zero-kb02 firmware

Available here:

* https://github.com/sago35/keyboards

## Troubleshooting

- Cannot flash the program

Check if the microcontroller is recognized with the `tinygo ports` command. If recognized correctly, `waveshare-rp2040-zero` will be displayed.

```
$ tinygo ports
Port                 ID        Boards
COM7                 2E8A:0003 waveshare-rp2040-zero
```

If not recognized, disconnect the microcontroller from the PC and reconnect it. Try putting it in bootloader mode: hold the BOOT button on the back, press the RST button and release both.

## Examples

* https://x.com/ysaito8015/status/1827626098450166185
* https://x.com/ysaito8015/status/1827630059580231788
* https://x.com/sago35tk/status/1830208709471223966
* https://x.com/Ryu_07_29/status/1847921967070163377
* [./19_redkey/](./19_redkey/)
* [./20_rotary_gopher](./20_rotary_gopher/)
* https://github.com/conejoninja/midikeeb
* https://x.com/triring/status/1891448348818776323
* https://github.com/sago35/koebiten

# Other Tips

* When taking photos or videos, setting to 30 frames/second prevents LCD flickering.

# Announcements

I wrote a technical book "Learning TinyGo Embedded Development from the Basics" (released on 2022/11/12) using TinyGo 0.26 + Wio Terminal. Please check it along with this page.

* https://sago35.hatenablog.com/entry/2022/11/04/230919