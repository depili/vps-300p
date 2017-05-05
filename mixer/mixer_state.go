package mixer

type MixerState struct {
	PGM   byte // Active PGM selection
	PST   byte // Active PST selection
	Type  int  // Transition type (Mix or Wipe)
	Dir   bool // true == upwards, false == downwards
	Value int  // T-bar value 0-1023
}

func (state *MixerState) copy() MixerState {
	return MixerState{
		PGM:   state.PGM,
		PST:   state.PST,
		Type:  state.Type,
		Dir:   state.Dir,
		Value: state.Value,
	}
}

func (state *MixerState) cut() MixerState {
	s := state.copy()
	p := state.PGM
	s.PGM = state.PST
	s.PST = p
	return s
}
