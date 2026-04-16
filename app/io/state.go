package io

type IOState struct {
	IF uint8
}

func (i *IO) GetState() IOState {
	return IOState{
		IF: i.IF,
	}
}

func (i *IO) LoadState(state IOState) {
	i.IF = state.IF
}
