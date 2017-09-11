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

func Minimax(state State) (Action, float64) {
	action, value := minimaxAlg(state, MAXDEPTH, "")
	return *action, value
}

func minimaxAlg(state State, depth int, tab string) (*Action, float64) {
	treedepth := MAXDEPTH - depth

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

	debugStr := "MAX"
	if state.IsOpponentTurn() {
		debugStr = "MIN"
	}
	for _, action := range state.FindLegals() {
		nextState := state.FindNextState(action)
		debugLog("%4d %s(%s) Action %v :nextstate %v\n", treedepth, tab, debugStr, action, nextState)
		_, value := minimaxAlg(nextState, depth - 1, tab + "....")
		debugLog("%4d %s(%s) VALUE of action %v : %.2f at state %v\n", treedepth, tab, debugStr, action, value, state)
		if !state.IsOpponentTurn() { // MAX
			if value > bestValue {
				bestValue, bestAction = value, action
				debugLog("%4d %s(%s) bestValue %.2f, bestAction so far %s\n", treedepth, tab, debugStr, bestValue, bestAction)
			}
		} else { // MIN
			if value < bestValue {
				bestValue, bestAction = value, action
				debugLog("%4d %s(%s) bestValue %.2f, bestAction so far %s\n", treedepth, tab, debugStr, bestValue, bestAction)
			}
		}
	}

	debugLog("%4d %s(%s) Best action %s : %.2f at state %v\n", treedepth, tab, debugStr, bestAction, bestValue, state)
	return &bestAction, bestValue
}

func AlphaBeta(state State) (Action, float64) {
	alpha := float64(math.MinInt32)
	beta := float64(math.MaxInt32)
	action, value := alphaBetaAlg(state, alpha, beta, MAXDEPTH, "")
	return *action, value
}


func alphaBetaAlg(state State, alpha, beta float64, depth int, tab string) (*Action, float64) {
	treedepth := MAXDEPTH - depth

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

	debugStr := "MAX"
	if state.IsOpponentTurn() {
		debugStr = "MIN"
	}

	for _, action := range state.FindLegals() {
		nextState := state.FindNextState(action)
		debugLog("%4d %s(%s) Action %v :nextstate %v\n", treedepth, tab, debugStr, action, nextState)
		_, value := alphaBetaAlg(nextState, alpha, beta, depth - 1, tab + "....")
		debugLog("%4d %s(%s) VALUE of action %v : %.2f at state %v\n", treedepth, tab, debugStr, action, value, state)
		if !state.IsOpponentTurn() { // MAX
			if value > bestValue {
				bestValue, bestAction = value, action
				debugLog("%4d %s(%s) bestValue %.2f, bestAction so far %s\n", treedepth, tab, debugStr, bestValue, bestAction)
			}
			if value > alpha {
				alpha = value
			}
			if beta <= alpha {
				debugLog("%4d %s(%s) Pruning at state %v and action %v\n", treedepth, tab, state, action)
				break
			}
		} else { // MIN
			if value < bestValue {
				bestValue, bestAction = value, action
				debugLog("%4d %s(%s) bestValue %.2f, bestAction so far %s\n", treedepth, tab, debugStr, bestValue, bestAction)
			}
			if value < beta {
				beta = value
			}
			if beta <= alpha {
				debugLog("%4d %s(%s) Pruning at state %v and action %v\n", treedepth, tab, state, action)
				break
			}
		}
	}
	debugLog("%4d %s(%s) Best action %s : %.2f at state %v\n", treedepth, tab, debugStr, bestAction, bestValue, state)
	return &bestAction, bestValue
}