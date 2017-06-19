package mixer

import ()

type WipeState struct {
	MattSat    int
	MattLum    int
	MattHue    int
	Border     int
	BorderType int
}

func (state *WipeState) copy() *WipeState {
	return &WipeState{
		MattSat:    state.MattSat,
		MattLum:    state.MattLum,
		MattHue:    state.MattHue,
		Border:     state.Border,
		BorderType: state.BorderType,
	}
}
