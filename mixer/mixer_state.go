package mixer

type MixerState struct {
	PGM        byte // Active PGM selection
	PST        byte // Active PST selection
	Type       int  // Transition type (Mix or Wipe)
	Dir        bool // true == upwards, false == downwards
	Value      int  // T-bar value 0-1023
	Layers     byte // Layers to transition
	Key1       bool // Key1 led state
	Key2       bool // Key2 led state
	Pattern    int  // Pattern number
	PatternRev bool // Reverse pattern
	Key1State  *KeyerState
	Key2State  *KeyerState
}

func (state *MixerState) copy() MixerState {
	return MixerState{
		PGM:        state.PGM,
		PST:        state.PST,
		Type:       state.Type,
		Dir:        state.Dir,
		Value:      state.Value,
		Layers:     state.Layers,
		Key1:       state.Key1,
		Key2:       state.Key2,
		Pattern:    state.Pattern,
		PatternRev: state.PatternRev,
		Key1State:  state.Key1State.copy(),
		Key2State:  state.Key1State.copy(),
	}
}

func (state *MixerState) cut() MixerState {
	s := state.copy()
	p := state.PGM
	s.PGM = state.PST
	s.PST = p
	return s
}

func (state *MixerState) transComplete() MixerState {
	s := state.copy()
	if s.Layers&LayerBKGD != 0 {
		s = s.cut()
	}
	if s.Layers&LayerKey1 != 0 {
		s.Key1 = !s.Key1
	}
	if s.Layers&LayerKey2 != 0 {
		s.Key1 = !s.Key2
	}
	s.Dir = !s.Dir
	s.Value = 0

	return s
}
