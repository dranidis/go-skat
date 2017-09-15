package game

type State interface {
	FindLegals() []Action
	FindNextState(Action) State
	IsTerminal() bool
	FindReward() float64
	IsOpponentTurn() bool
	Heuristic() float64
	GetTacticsMove() Action
}

type Action interface {
}
