all: smoketest

smoketest:
	mkdir -p ./out
	tinygo build -o ./out/00_basic.uf2              --target waveshare-rp2040-zero --size short ./00_basic
	tinygo build -o ./out/01_blinky1.uf2            --target waveshare-rp2040-zero --size short ./01_blinky1
	tinygo build -o ./out/02_blinky2.uf2            --target waveshare-rp2040-zero --size short ./02_blinky2
	tinygo build -o ./out/03_usbcdc-serial.uf2      --target waveshare-rp2040-zero --size short ./03_usbcdc-serial
	tinygo build -o ./out/04_usbcdc-echo.uf2        --target waveshare-rp2040-zero --size short ./04_usbcdc-echo
	tinygo build -o ./out/05_rotary.uf2             --target waveshare-rp2040-zero --size short ./05_rotary
	tinygo build -o ./out/06_joystick.uf2           --target waveshare-rp2040-zero --size short ./06_joystick
	tinygo build -o ./out/07_oled.uf2               --target waveshare-rp2040-zero --size short ./07_oled
	tinygo build -o ./out/08_oled_tinydraw.uf2      --target waveshare-rp2040-zero --size short ./08_oled_tinydraw
	tinygo build -o ./out/09_oled_tinyfont.uf2      --target waveshare-rp2040-zero --size short ./09_oled_tinyfont
	tinygo build -o ./out/10_oled_rotated.uf2       --target waveshare-rp2040-zero --size short ./10_oled_rotated
	tinygo build -o ./out/11_oled_animation.uf2     --target waveshare-rp2040-zero --size short ./11_oled_animation
	tinygo build -o ./out/12_matrix_basic.uf2       --target waveshare-rp2040-zero --size short ./12_matrix_basic
	tinygo build -o ./out/13_rotary_button.uf2      --target waveshare-rp2040-zero --size short ./13_rotary_button
	tinygo build -o ./out/14_hid_keyboard.uf2       --target waveshare-rp2040-zero --size short ./14_hid_keyboard
	tinygo build -o ./out/15_hid_mouse.uf2          --target waveshare-rp2040-zero --size short ./15_hid_mouse
	tinygo build -o ./out/16_oled_inverted_hw.uf2   --target waveshare-rp2040-zero --size short ./16_oled_inverted_hw
	tinygo build -o ./out/17_oled_japanese_font.uf2 --target waveshare-rp2040-zero --size short ./17_oled_japanese_font
	tinygo build -o ./out/18_midi.uf2               --target waveshare-rp2040-zero --size short ./18_midi
	tinygo build -o ./out/19_redkey.uf2             --target waveshare-rp2040-zero --size short ./19_redkey
	tinygo build -o ./out/20_rotary_gopher.uf2      --target waveshare-rp2040-zero --size short ./20_rotary_gopher
	tinygo build -o ./out/21_midi2.uf2              --target waveshare-rp2040-zero --size short ./21_midi2
	tinygo build -o ./out/22_buzzer.uf2             --target waveshare-rp2040-zero --size short ./22_buzzer
	tinygo build -o ./out/80_checker.uf2            --target waveshare-rp2040-zero --size short ./80_checker
