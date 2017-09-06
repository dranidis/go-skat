package mcts

import (
	"math"
)

type Node struct {
	score int
	children []Node
}

func terminal(node Node) bool {
	return len(node.children) == 0
}

func heuristic(node Node) int {
	if terminal(node) {
		if node.score > 60 {
			return 1
		}
		return -1 
	}
	return 0
}

func children(node Node) []Node {
	return node.children
}

func minimax(node Node, depth int, maximizingPlayer bool) int {
	if depth == 0 || terminal(node) {
		return heuristic(node)
	}

	if maximizingPlayer {
		bestValue := math.MinInt32
		for _, child := range children(node) {
			v := minimax(child, depth - 1, false)
			if v > bestValue {
				bestValue = v
			}
		}
		return bestValue
	} else {
		bestValue := math.MaxInt32
		for _, child := range children(node) {
			v := minimax(child, depth - 1, true)
			if v < bestValue {
				bestValue = v
			}
		}
		return bestValue	
	}
}