package main

import (
	"math"
)
// 01 function minimax(node, depth, maximizingPlayer)
// 02     if depth = 0 or node is a terminal node
// 03         return the heuristic value of node

// 04     if maximizingPlayer
// 05         bestValue := −∞
// 06         for each child of node
// 07             v := minimax(child, depth − 1, FALSE)
// 08             bestValue := max(bestValue, v)
// 09         return bestValue

// 10     else    (* minimizing player *)
// 11         bestValue := +∞
// 12         for each child of node
// 13             v := minimax(child, depth − 1, TRUE)
// 14             bestValue := min(bestValue, v)
// 15         return bestValue
type Node struct {

}
func terminal(node Node) bool {
	return true
}

func heuristic(node Node) int {
	return 0
}

func children(node Node) []Node {
	return nil
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