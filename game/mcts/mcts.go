package mcts

import (
	"math"
	"math/rand"
	"fmt"
	// "log"
	"time"
	"github.com/dranidis/go-skat/game"
)

var SimulationRuns = 100
var ExplorationParameter = 2.0
var DEBUG = false
var MostVisited = true

func debugLog(format string, a ...interface{}) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}

type Node struct {
	state game.State
	visits int
	utility float64
	action game.Action
	parent *Node
	children []*Node
}

var	rootNode *Node

func Uct(state game.State, seconds int) (game.Action, float64) {
	start := time.Now()

	var root *Node
	rootNode = &Node{state, 0, 0, nil, nil, []*Node{}}
	root = rootNode

	var elapsed time.Duration
	for elapsed < time.Duration(seconds) * time.Millisecond {
		// debugLog("Iteration: ", i)
		debugLog("..Selection\n")
		node := selectNode(root, 0)
		debugLog(".... selected: %v\n", node.state)
		debugLog("..Expansion\n")
		// node = expand(node)
		expand(node)
		// printTree(root, 0)
		// debugLog("..Simulation")
		// debugLog(".. Simulating form state: %v\n", node.state)
		printTree(root, 0)
		debugLog("")
		t := time.Now()
		elapsed = t.Sub(start)
	}
	// printTree(root, 0)
	if MostVisited {
		best, u := mostVisitedAction(root)
		debugLog("\nRETURNING MOST VISITED MOVE: %v\n\n", best)
		return best, u		
	} else {
		best, u := mostUtilityAction(root)
		debugLog("\nRETURNING BEST MOVE: %v\n\n", best)
		return best, u		

	}

}

func printTree(node *Node, depth int) {
	indent := ""
	for i := 0; i < depth ; i++ {
		indent += "\t"
	}
	debugLog("%sState %s,Visits %v, Utility %.2f, game.Action %s\n", indent, node.state, node.visits, node.utility / float64(node.visits), node.action)
	for _, child := range node.children {
		printTree(child, depth + 1)
	}
}

func mostVisited(nodes []*Node) *Node {
	// debugLog(len(nodes))
	most := 0
	imost := -1
	for i, node := range nodes {
		if node.visits > most {
			most = node.visits
			imost = i
		}
	}
	if imost >= len(nodes) {
		imost = 0
	}
	return nodes[imost]
}


// func mostVisitedAction(node *Node) game.Action {
// 	var best game.Action
// 	most := math.MinInt64
// 	actions := node.state.FindLegals()
// 	for _, action := range actions {
// 		s := node.state.FindNextState(action)
// 		n, ok := visitedStates[s.StateId()]
// 		if ok {
// 			u := n.visits
// 			// fmt.Printf(".. %s: Util: %f %d\n", action, u, n.visits)
// 			if u > most {
// 				most = u
// 				best = action
// 			}		
// 		} else {
// 			// fmt.Printf(".. %s: Not visited\n", action)
// 		}	
// 	}
// 	return best
// }
func mostVisitedAction(node *Node) (game.Action, float64) {
	// debugLog(len(nodes))
	var bestAction game.Action
	mostVisits := 0
	mostU := -1.0
	for _, child := range node.children {
		if child.visits > mostVisits {
			mostVisits = child.visits
			bestAction = child.action
			mostU = child.utility / float64(child.visits)
		}
	}
	return bestAction, mostU 
}

func mostUtilityAction(node *Node) (game.Action, float64) {
	var bestAction game.Action
	most := float64(math.MinInt64)
	for _, child := range node.children {
		action := child.action
		u := child.utility / float64(child.visits)
		// fmt.Printf(".. %s: Util: %f %d\n", action, u, n.visits)
		if u > most {
			most = u
			bestAction = action
		}		
	}
	return bestAction, most
}

// func mostUtilityAction(node *Node) (game.Action, float64) {
// 	var best game.Action
// 	most := float64(math.MinInt64)
// 	actions := node.state.FindLegals()
// 	for _, action := range actions {
// 		s := node.state.FindNextState(action)
// 		u := n.utility / float64(s.visits)
// 		// fmt.Printf(".. %s: Util: %f %d\n", action, u, n.visits)
// 		if u > most {
// 			most = u
// 			best = action
// 		}		
// 	}
// 	return best, most
// }

func selectNode(node *Node, depth int) *Node {
	debugLog("selectNode at: %v\n", node.state)

	if node.visits == 0 || node.children == nil || len(node.children) == 0 {
		return node
	}
	// debugLog(".. not visited")

	for _, child := range node.children {
		debugLog(".. Child: %s\n", child.action)
		if child.visits == 0 { 
			debugLog(".. zero visits..\n")
			return child
		}
	}

	debugLog(".. ALL visited. Selecting Highest/Lowest score\n")

	var score float64
	if !node.state.IsOpponentTurn() {
		score = float64(math.MinInt64)
	} else {
		score = float64(math.MaxInt64)
	}

	result := node

	// parentVisits := 0
	// for _, child := range node.children {
	// 	parentVisits += child.visits
	// }

	// score := selectfn(node.children[0], parentVisits)
	// result := node.children[0]
	// debugLog(".. Child:", 0, " ", node.children[0], " score", score)

	// only choide
	if len(node.children) == 1 {
		debugLog("ONLY CHOICE: %v\n", node.children[0].state)
		return selectNode(node.children[0], depth + 1)
	}

	for i := 0; i < len(node.children); i++ {
		// debugLog(".. Child: %d %v\n", i, node.children[i].state)

		if !node.state.IsOpponentTurn() {
			newscore := selectfn(node.children[i], 2.0)
			// debugLog(".. Child: %d, %v score %v\n", i, node.children[i].state, newscore)
			if newscore > score { // MAX
				score = newscore
				result = node.children[i]
			}
		} else {
			newscore := selectfn(node.children[i], -2.0)
			// debugLog(".. Child: %d, %v score %v\n", i, node.children[i].state, newscore)
			if newscore < score { //MIN
				score = newscore
				result = node.children[i]
			}
		}
	}
	debugLog(".. Selected:  %v\n", result.state)
	return selectNode(result, depth+1)
}

func selectfn(node *Node, factor float64) float64 {
	u := node.utility / float64(node.visits)
	pv := float64(node.parent.visits)
	nv := float64(node.visits)
	return u + ExplorationParameter * math.Sqrt(math.Log(pv) / nv)
}

func expand(node *Node) { //*Node {
	// actions := node.state.FindLegals(node.state.player)
	actions := node.state.FindLegals()
	for _, action := range actions {
		// debugLog("..Expanding ACTION %s\n", action)
		newState := node.state.FindNextState(action)
		// debugLog("..state %v\n", newState)
		// check if state is visited
		var nNode *Node
		// debugLog("..new state %v\n", newState)
		newNode := &Node{newState, 0, 0, action, node, []*Node{}}
		node.children = append(node.children, newNode)
		nNode = newNode

		// debugLog("..Simulation\n")
		for i :=0 ; i < SimulationRuns; i++ {
			reward := simulate(nNode.state)
			backPropagate(nNode, reward)	
		}		
	}
}

func simulate(state game.State) float64 {
	if state.IsTerminal() {
		return state.FindReward()
	}
	newState := state
	// for _, player := range newState.players {
		// options := newState.FindLegals(player)
		options := newState.FindLegals()
		// debugLog("Options: %v\n", options)
		best := rand.Intn(len(options))
		newState = newState.FindNextState(options[best])
		// debugLog("game.Action: %v\n", options[best])
		// debugLog("%v", newState)
	// } 
	return simulate(newState)
}

func backPropagate(node *Node, score float64) {
	node.visits++
	node.utility += score
	// node.utility /= float64(node.visits)
	if node.parent != nil {
		backPropagate(node.parent, score)
	}
}


