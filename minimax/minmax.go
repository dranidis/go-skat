package minimax

import (
	"math"
	"fmt"
)

var DEBUG = false
var MAXDEPTH = math.MaxInt32

func debugLog(format string, a ...interface{}) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}

type Node struct {
	state State
	score float64
	children []*Node
}

type State interface {
	FindLegals() []Action
	FindNextState(Action) State
	IsTerminal() bool
	FindReward() float64
	IsOpponentTurn() bool
	Heuristic() float64
}

type Action interface {
}

func Minimax(state State) *Action {
	action, _ := minimaxAlg(state, MAXDEPTH, "")
	return action
}

func minimaxAlg(state State, depth int, tab string) (*Action, float64) {
	if depth == 0  || state.IsTerminal() {
		return nil, state.Heuristic()
	}
	var bestValue float64
	var bestAction Action

	if !state.IsOpponentTurn() {
		bestValue = float64(math.MinInt32)
	} else {
		bestValue = float64(math.MaxInt32)
	}

	for _, action := range state.FindLegals() {
		nextState := state.FindNextState(action)
		debugLog("%saction %v :nextstate %v\n", tab, action, nextState)
		_, value := minimaxAlg(nextState, depth - 1, tab + "....")
		debugLog("%sVALUE of action %v : %.2f at state %v\n", tab, action, value, state)
		if !state.IsOpponentTurn() { // MAX
			if value > bestValue {
				bestValue, bestAction = value, action
				debugLog("%sbestValue %.2f, bestAction so far %s\n", tab, bestValue, bestAction)
			}
		} else { // MIN
			if value < bestValue {
				bestValue, bestAction = value, action
			}
		}
	}
	debugLog("%s action %s : %.2f at state %v\n", tab, bestAction, bestValue, state)
	return &bestAction, bestValue
}


