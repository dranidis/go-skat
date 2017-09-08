package minimax

import (
	"math"
	// "log"
	"fmt"
)

var DEBUG = false

func debugLog(format string, a ...interface{}) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}

type Node struct {
	state State
	action *Action // last action moving to this state
	score float64
	children []*Node
}

type State interface {
	// players []Player
	FindLegals() []Action
	FindNextState(Action) State
	IsTerminal() bool
	FindReward() float64
	StateId() uint64
	IsOpponentTurn() bool
	Heuristic() float64
}

type Action interface {
	Equals(Action) bool
}

// func terminal(node Node) bool {
// 	return len(node.children) == 0
// }


// func children(node Node) []*Node {
// 	return node.children
// }

func Minimax(state State, depth int, maximizingPlayer bool) Action {
	if DEBUG {
		if maximizingPlayer {
			debugLog("MAXIMIZING")
		} else {
			debugLog("MINIMIZING")
		}
	}
	var action Action
	node := createMMTree(state, nil, depth)
	// if DEBUG {printTree(node, "")}
	best := MinimaxAux(node, depth, maximizingPlayer)
	_ = best

	if DEBUG {printTree(node, "", maximizingPlayer)}

	if maximizingPlayer {
		maxScore := float64(math.MinInt32)
		for _, child := range node.children {
			if child.score > maxScore {
				maxScore = child.score
				action = *child.action
			}
		}
	} else {
		minScore := float64(math.MaxInt32)
		for _, child := range node.children {
			if child.score < minScore {
				minScore = child.score
				action = *child.action
			}
		}

	}
	return action
}

func createMMTree(state State, action *Action, depth int) *Node {
	if depth == 0  || state.IsTerminal() {
		return &Node{state, action, 0, []*Node{}}
	}
	children := []*Node{}
	if depth != 0 {
		for _, action := range state.FindLegals() {
			nextState := state.FindNextState(action)
			nextNode := createMMTree(nextState, &action, depth - 1)
			children = append(children, nextNode)
		}
	}

	return &Node{state, action, 0, children}
}

func printTree(node *Node, indent string, maximizingPlayer bool) {
	debugLog("%sNode with score %.2f: %v, \n", indent, node.score, node.state)
	if maximizingPlayer {
		debugLog("%s\tMaximizing Children: \n", indent)
	} else {
		debugLog("%s\tMINimizing Children: \n", indent)
	}
	for _, child := range node.children {
		printTree(child, indent + "\t", !maximizingPlayer)
	}
}

func MinimaxAux(node *Node, depth int, maximizingPlayer bool) float64 {
	if depth == 0 || node.state.IsTerminal() {
		// debugLog("TERMINAL %v: returning %.2f\n", node.state, node.state.Heuristic())
		node.score = node.state.Heuristic()
		return node.state.Heuristic()
	}

	if depth > 0 {
		if maximizingPlayer {
			bestValue := float64(math.MinInt32)
			for _, child := range node.children {
				v := MinimaxAux(child, depth - 1, false)
				if v > bestValue {
					bestValue = v
				}
			}
			node.score = bestValue
			return bestValue
		} else {
			bestValue := float64(math.MaxInt32)
			for _, child := range node.children {
				v := MinimaxAux(child, depth - 1, true)
				if v < bestValue {
					bestValue = v
				}
			}
			node.score = bestValue
			return bestValue	
		}	
	}
	return 0.0
}