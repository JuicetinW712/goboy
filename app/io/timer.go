package io

// Reference https://github.com/Ashiepaws/GBEDG/blob/master/timers/index.md

// TODO - timers still tick during halt but stop during stop

// Div counter increments every 256 cycles
// Tima counter is a multiple of div that can be set by the program

type Timer struct {
	div uint16

	tima uint8 // timer counter
	tma  uint8 // timer modulo
	tac  uint8 // timer control

	// Callback for RequestInterrupt
	RequestInterrupt func(bit uint8)

	// Track the previous AND result for falling edge detection
	lastAndResult bool

	// TIMA overflow delay counter
	// On overflow, reload TIMA with TMA and request interrupt (after 4 cycles)
	// A write during the 4 cycles will cancel this
	overflowCycles  int
	overflowPending bool
}

// DIV is 16 bits and only the upper 8 are read. After 256
// cycles, the upper byte increments so it increments every 256 cycles
func (t *Timer) ReadDIV() uint8 {
	return uint8(t.div >> 8)
}

func (t *Timer) WriteDIV(_ uint8) {
	t.div = 0
	t.updateAndResult()
}

func (t *Timer) ReadTIMA() uint8 {
	return t.tima
}

func (t *Timer) WriteTIMA(val uint8) {
	if t.overflowPending {
		// Cancel reload and interrupt request
		t.overflowPending = false
		t.overflowCycles = 0
	}
	t.tima = val
}

func (t *Timer) ReadTMA() uint8 {
	return t.tma
}

func (t *Timer) WriteTMA(val uint8) {
	t.tma = val
}

func (t *Timer) ReadTAC() uint8 {
	return t.tac
}

func (t *Timer) WriteTAC(val uint8) {
	t.tac = val
	t.updateAndResult()
}

func (t *Timer) getTimerBit() int {
	switch t.tac & 0b11 {
	case 0:
		return 9 // 4096 Hz
	case 1:
		return 3 // 262144 Hz
	case 2:
		return 5 // 65536 Hz
	case 3:
		return 7 // 16384 Hz
	default:
		return 0
	}
}

func (t *Timer) andResult() bool {
	timerEnabled := (t.tac>>2)&1 == 1
	bitPos := t.getTimerBit()
	divBit := (t.div>>bitPos)&1 == 1
	return timerEnabled && divBit
}

func (t *Timer) updateAndResult() {
	current := t.andResult()
	if t.lastAndResult && !current {
		t.incrementTIMA()
	}
	t.lastAndResult = current
}

func (t *Timer) incrementTIMA() {
	if t.tima == 0xFF {
		// doesn't fire off interrupt or set to tima until 4 cycles later
		t.tima = 0
		t.overflowCycles = 4
		t.overflowPending = true
	} else {
		t.tima++
	}
}

func (t *Timer) Tick() {
	// Increment div counter
	t.div++

	// Check if there is a falling edge on AND result
	// If there is, increment TIMA
	t.updateAndResult()

	// Waiting for 4 cycles before interrupt fires
	if t.overflowPending {
		t.overflowCycles--
		if t.overflowCycles == 0 {
			// Reload TIMA with TMA and request interrupt
			t.tima = t.tma
			if t.RequestInterrupt != nil {
				t.RequestInterrupt(2)
			}
			t.overflowPending = false
		}
	}
}
