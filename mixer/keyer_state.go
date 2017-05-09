package mixer

import ()

type KeyerState struct {
	MattSat      int
	MattLum      int
	MattHue      int
	Edge         int
	EdgeSat      int
	EdgeLum      int
	EdgeHue      int
	Transparency int
	ShadowX      int
	ShadowY      int
}

func (state *KeyerState) copy() *KeyerState {
	return &KeyerState{
		MattSat:      state.MattSat,
		MattLum:      state.MattLum,
		MattHue:      state.MattHue,
		Edge:         state.Edge,
		EdgeSat:      state.EdgeSat,
		EdgeLum:      state.EdgeLum,
		EdgeHue:      state.EdgeHue,
		Transparency: state.Transparency,
		ShadowX:      state.ShadowX,
		ShadowY:      state.ShadowY,
	}
}
