package joypad

type Button uint8

// TODO - technically doesn't implement the joypad interrupt
// but that mostly handles waking up from stop so not technically necessary

const (
	DPAD_RIGHT Button = iota
	DPAD_LEFT
	DPAD_UP
	DPAD_DOWN

	BUTTON_A
	BUTTON_B
	BUTTON_SELECT
	BUTTON_START
)

type Joypad struct {
	dpad    uint8 // low nibble only
	buttons uint8 // low nibble only
	sel     uint8 // bits 4/5
}

func CreateJoypad() *Joypad {
	return &Joypad{
		dpad:    0x0F, // Initialize to all deactivated
		buttons: 0x0F, // Same here
		sel:     0x00,
	}
}

func (j *Joypad) PressButton(button Button) {
	if button < BUTTON_A {
		j.dpad &^= (1 << button)
	} else {
		j.buttons &^= (1 << (button - BUTTON_A))
	}
}

func (j *Joypad) ReleaseButton(button Button) {
	if button < BUTTON_A {
		j.dpad |= (1 << button)
	} else {
		j.buttons |= (1 << (button - BUTTON_A))
	}
}

func (j *Joypad) Write(val uint8) {
	j.sel = val & 0x30 // only bits 4/5 matter
}

func (j *Joypad) Read() uint8 {
	res := uint8(0x0F)

	if j.sel&0x10 == 0 { // dpad selected
		res &= j.dpad
	}

	if j.sel&0x20 == 0 { // buttons selected
		res &= j.buttons
	}

	// 0xC0 is first bits set to 1
	return 0xC0 | j.sel | res
}
