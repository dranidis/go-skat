package minimax

import (
	"github.com/dranidis/go-skat/game"
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
	state game.State
	score float64
	children []*Node
}


func Minimax(state game.State) (game.Action, float64) {
	action, value := minimaxAlg(state, MAXDEPTH, "")
	return *action, value
}

func minimaxAlg(state game.State, depth int, tab string) (*game.Action, float64) {
	treedepth := MAXDEPTH - depth

	if depth == 0  || state.IsTerminal() {
		return nil, state.Heuristic()
	}
	var bestValue float64
	var bestAction game.Action

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
		debugLog("%4d %s(%s) game.Action %v :nextstate %v\n", treedepth, tab, debugStr, action, nextState)
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

func AlphaBeta(state game.State) (game.Action, float64) {
	alpha := float64(math.MinInt32)
	beta := float64(math.MaxInt32)
	action, value := alphaBetaAlg(state, alpha, beta, MAXDEPTH, "")
	return *action, value
}


func alphaBetaAlg(state game.State, alpha, beta float64, depth int, tab string) (*game.Action, float64) {
	treedepth := MAXDEPTH - depth

	if depth == 0  || state.IsTerminal() {
		return nil, state.Heuristic()
	}
	var bestValue float64
	var bestAction game.Action

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
		debugLog("%4d %s(%s) game.Action %v :nextstate %v\n", treedepth, tab, debugStr, action, nextState)
		_, value := alphaBetaAlg(nextState, alpha, beta, depth - 1, tab + "....")
		debugLog("%4d %s(%s) VALUE of action %v : %.2f at state %v\n", treedepth, tab, debugStr, action, value, state)
		if !state.IsOpponentTurn() { // MAX
			if value > bestValue {
				bestValue, bestAction = value, action
				debugLog("%4d %s(%s) Best Value so far: %.2f, Best game.Action so far: %s\n", treedepth, tab, debugStr, bestValue, bestAction)
			}
			if value > alpha {
				alpha = value
			}
			if beta <= alpha {
				debugLog("%4d %s(%s) Pruning at state %v and action %s\n", treedepth, tab, debugStr, state, action)
				break
			}
		} else { // MIN
			if value < bestValue {
				bestValue, bestAction = value, action
				debugLog("%4d %s(%s) Best Value so far: %.2f, Best game.Action so far: %s\n", treedepth, tab, debugStr, bestValue, bestAction)
			}
			if value < beta {
				beta = value
			}
			if beta <= alpha {
				debugLog("%4d %s(%s) Pruning at state %v and action %s\n", treedepth, tab, debugStr, state, action)
				break
			}
		}
	}
	debugLog("%4d %s(%s) Best action %s : %.2f at state %v\n", treedepth, tab, debugStr, bestAction, bestValue, state)
	return &bestAction, bestValue
}



func ZeroWindowAlg(state game.State, beta float64) (game.Action, float64) {
	action, value := alphaBetaAlg(state, beta - 1.0, beta, MAXDEPTH, "")
	return *action, value
}

func ABWindowAlg(state game.State, alpha, beta float64) (game.Action, float64) {
	action, value := alphaBetaAlg(state, alpha, beta, MAXDEPTH, "")
	return *action, value
}

func AlphaBetaTactics(state game.State) (game.Action, float64) {
	alpha := float64(math.MinInt32)
	beta := float64(math.MaxInt32)
	action, value := alphaBetaTacticsAlg(state, alpha, beta, MAXDEPTH, "")
	return *action, value
}

func alphaBetaTacticsAlg(state game.State, alpha, beta float64, depth int, tab string) (*game.Action, float64) {
	treedepth := MAXDEPTH - depth

	if depth == 0  || state.IsTerminal() {
		return nil, state.Heuristic()
	}
	var bestValue float64
	var bestAction game.Action

	if !state.IsOpponentTurn() {
		bestValue = float64(math.MinInt32)
	} else {
		bestValue = float64(math.MaxInt32)
	}

	debugStr := "MAX"
	if state.IsOpponentTurn() {
		debugStr = "TAC"
	}

	if !state.IsOpponentTurn() { // MAX
		for _, action := range state.FindLegals() {
			nextState := state.FindNextState(action)
			debugLog("%4d %s(%s) game.Action %v :nextstate %v\n", treedepth, tab, debugStr, action, nextState)
			_, value := alphaBetaTacticsAlg(nextState, alpha, beta, depth - 1, tab + "....")
			debugLog("%4d %s(%s) VALUE of action %v : %.2f at state %v\n", treedepth, tab, debugStr, action, value, state)
			if value > bestValue {
				bestValue, bestAction = value, action
				debugLog("%4d %s(%s) Best Value so far: %.2f, Best game.Action so far: %s\n", treedepth, tab, debugStr, bestValue, bestAction)
			}
			if value > alpha {
				alpha = value
			}
			if beta <= alpha {
				debugLog("%4d %s(%s) Pruning at state %v and action %s\n", treedepth, tab, debugStr, state, action)
				break
			}
		}
	} else { // Tactics
		action := state.GetTacticsMove()
		nextState := state.FindNextState(action)
		_, value := alphaBetaTacticsAlg(nextState, alpha, beta, depth - 1, tab + "....")

		bestValue, bestAction = value, action
	}
	debugLog("%4d %s(%s) Best action %s : %.2f at state %v\n", treedepth, tab, debugStr, bestAction, bestValue, state)
	return &bestAction, bestValue

}

