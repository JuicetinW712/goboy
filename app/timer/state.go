package timer

type TimerState struct {
	DIV  uint16
	TIMA uint8
	TMA  uint8
	TAC  uint8
}

func (t *Timer) GetState() TimerState {
	return TimerState{
		DIV:  t.div,
		TIMA: t.tima,
		TMA:  t.tma,
		TAC:  t.tac,
	}
}

func (t *Timer) LoadState(state TimerState) {
	t.div = state.DIV
	t.tima = state.TIMA
	t.tma = state.TMA
	t.tac = state.TAC
}
