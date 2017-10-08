package main

import (
	"math"
	// "math/rand"
)

type MCNode struct {
	state    State
	visits   int
	utility  int
	parent   *MCNode
	children []MCNode
}

type State struct {
}

type Action struct {
}

func selectNode(node MCNode) MCNode {
	if node.visits == 0 {
		return node
	}
	for _, child := range node.children {
		if child.visits == 0 {
			return child
		}
	}
	score := float64(0)
	result := node
	for _, child := range node.children {
		newscore := selectfn(child)
		if newscore > score {
			score = newscore
			result = child
		}
	}
	return selectNode(result)
}

func expand(node MCNode) {
	actions := findLegals(node.state)
	for _, action := range actions {
		newState := findNextState(node.state, action)
		newNode := MCNode{newState, 0, 0, &node, []MCNode{}}
		node.children = append(node.children, newNode)
	}
}

var roles = []int{}

func simulate(state State) int {
	if findTerminal(state) {
		return findReward(state)
	}
	newState := state
	for _, role := range roles {
		_ = role
		options := findLegals(newState)
		best := r.Intn(len(options))
		newState = findNextState(newState, options[best])
	}
	return simulate(newState)
}

func backPropagate(node MCNode, score int) {
	node.visits++
	node.utility += score
	if node.parent != nil {
		backPropagate(*node.parent, score)
	}
}

func findReward(state State) int {
	return 0
}

func findLegals(state State) []Action {
	return []Action{}
}

func findNextState(state State, action Action) State {
	return State{}
}

func selectfn(node MCNode) float64 {
	return float64(node.utility) + 2.0*math.Sqrt(math.Log(float64(node.parent.visits))/float64(node.visits))
}

func findTerminal(state State) bool {
	return false
}
