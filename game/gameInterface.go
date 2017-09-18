package game

type State interface {
	FindLegals() []Action
	FindNextState(Action) State
	IsTerminal() bool
	IsOpponentTurn() bool
	Heuristic() float64
	GetTacticsMove() Action
	FindReward() float64
}

type Action interface {
}
