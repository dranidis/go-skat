package mcts

import (
	"math"
	"math/rand"
	"fmt"
	// "log"
	"time"
)

func debugLog(format string, a ...interface{}) {
	if true {
		fmt.Printf(format, a...)
	}
}



type Node struct {
	state State
	visits int
	utility float64
	action Action
	parent *Node
	children []*Node
}

type State interface {
	// players []Player
	FindLegals() []Action
	FindNextState(Action) State
	IsTerminal() bool
	FindReward() float64
	StateId() uint64
}

// type Player interface {

// }

type Action interface {
	Equals(Action) bool
}

var	rootNode *Node
var visitedStates = make(map[uint64]*Node)

// func InitSubtree(state State) {
// 	rootNode = &Node{state, 0, 0, nil, nil, []*Node{}}
// }

// func ChangeSubtree(action Action) {
// 	for _, child := range rootNode.children {
// 		if action.Equals(child.action) {
// 			rootNode = child
// 			return
// 		}
// 	}
// 	log.Fatal("Action not found")
// }

func Uct(state State, seconds int) Action {
	start := time.Now()

	var root *Node
	stateId := state.StateId()
	visitedNode, ok := visitedStates[stateId]
	if ok {
		root = visitedNode
	} else {
		rootNode = &Node{state, 0, 0, nil, nil, []*Node{}}
		root = rootNode
		visitedStates[stateId] = rootNode
	}	

	var elapsed time.Duration
	// rootNode := Node{state, 0, 0, nil, nil, []*Node{}}
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
	// return mostVisited(root.children).action
	return mostUtilityAction(root)
}

var printed = make(map[uint64]*Node)

func printTree(node *Node, depth int) {
	if depth == 0 {
		printed = make(map[uint64]*Node)
	}
	indent := ""
	for i := 0; i < depth ; i++ {
		indent += "\t"
	}

	// action := ""
	debugLog("%sState %d,Visits %v, Utility %.2f, Action %s\n", indent, node.state, node.visits, node.utility / float64(node.visits), node.action)
	_, ok := printed[node.state.StateId()]
	printed[node.state.StateId()] = node
	if !ok {
		for _, child := range node.children {
			printTree(child, depth + 1)
		}
	}

	printed[node.state.StateId()] = node
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

func mostUtilityAction(node *Node) Action {
	var best Action
	most := float64(math.MinInt64)
	actions := node.state.FindLegals()
	for _, action := range actions {
		s := node.state.FindNextState(action)
		n, ok := visitedStates[s.StateId()]
		if ok {
			u := n.utility / float64(n.visits)
			fmt.Printf(".. %s: Util: %f %d\n", action, u, n.visits)
			if u > most {
				most = u
				best = action
			}		
		} else {
			fmt.Printf(".. %s: Not visited\n", action)
		}	

	}
	return best
}


var selected = make(map[uint64]*Node)

func selectNode(node *Node, depth int) *Node {
	debugLog("selectNode at: %v\n", node.state)
	if depth == 0 {
		selected = make(map[uint64]*Node)
	}

	selected[node.state.StateId()] = node
	debugLog("SELECTED map : %v\n", selected)

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

	debugLog(".. ALL visited. Selecting Highest score\n")

	score := float64(math.MinInt64)
	result := node

	parentVisits := 0
	for _, child := range node.children {
		parentVisits += child.visits
	}

	// score := selectfn(node.children[0], parentVisits)
	// result := node.children[0]
	// debugLog(".. Child:", 0, " ", node.children[0], " score", score)

	// only choide
	if len(node.children) == 1 {
		debugLog("ONLY CHOICE: %v\n", node.children[0].state)
		return selectNode(node.children[0], depth + 1)
	}

	for i := 0; i < len(node.children); i++ {
		debugLog(".. Child: %d %v\n", i, node.children[i].state)

		_, ok := selected[node.children[i].state.StateId()]
		if !ok {
			newscore := selectfn(node.children[i], parentVisits)
			debugLog(".. Child: %d, %v score %v\n", i, node.children[i].state, newscore)
			if newscore > score {
				score = newscore
				result = node.children[i]
			}
		} else {
			debugLog(".. already Visited. Selecting not visited child..\n" )
		}

	}

	return selectNode(result, depth+1)
}

func expand(node *Node) { //*Node {
	// actions := node.state.FindLegals(node.state.player)
	actions := node.state.FindLegals()
	for _, action := range actions {
		debugLog("..Expanding ACTION %s\n", action)
		newState := node.state.FindNextState(action)
		debugLog("..state %v\n", newState)
		// check if state is visited
		stateId := newState.StateId()
		visitedNode, ok := visitedStates[stateId]
		var nNode *Node
		if ok {
			debugLog("..already visited %v\n", newState)
			// visitedNode.action = action
			node.children = append(node.children, visitedNode)
			nNode = visitedNode
			debugLog("..Propagation\n")
			reward := simulate(visitedNode.state)
			backPropagate(visitedNode, reward)	
		} else {
			debugLog("..new state %v\n", newState)
			newNode := &Node{newState, 0, 0, action, node, []*Node{}}
			visitedStates[stateId] = newNode
			node.children = append(node.children, newNode)
			nNode = newNode

			debugLog("..Propagation\n")
			reward := simulate(nNode.state)
			backPropagate(nNode, reward)	
		}
		
	}

	// return node.children[0]
}

func simulate(state State) float64 {
	if state.IsTerminal() {
		return state.FindReward()
	}
	newState := state
	// for _, player := range newState.players {
		// options := newState.FindLegals(player)
		options := newState.FindLegals()
		best := rand.Intn(len(options))
		newState = newState.FindNextState(options[best])
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

func selectfn(node *Node, parentVisits int) float64 {
	u := node.utility / float64(node.visits)
	// pv := float64(node.parent.visits)
	pv := float64(parentVisits)
	nv := float64(node.visits)
	return u + 2.0 * math.Sqrt(math.Log(pv) / nv)
}

