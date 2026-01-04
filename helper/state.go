package helper

type GameState int

const (
	StatePlaying GameState = iota
	StatePaused
	StateQuit
)
