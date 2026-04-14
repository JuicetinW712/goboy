package timer

type TimerState struct {
	DIV  uint16
	TIMA uint8
	TMA  uint8
	TAC  uint8
}

func (t *Timer) GetTimerState() TimerState {
	return TimerState{
		DIV:  t.div,
		TIMA: t.tima,
		TMA:  t.tma,
		TAC:  t.tac,
	}
}
