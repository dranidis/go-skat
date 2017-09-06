package mcts

import (
	"math"
	"math/rand"
	"fmt"
	// "log"
	"time"
)

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
		// fmt.Println("Iteration: ", i)
		// fmt.Println("..Selection")
		node := selectNode(root)
		// fmt.Println("..Expansion")
		// node = expand(node)
		expand(node)
		// printTree(root, 0)
		// fmt.Println("..Simulation")
		// fmt.Printf(".. Simulating form state: %v\n", node.state)
		printTree(root, 0)
		fmt.Println("")
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
	fmt.Printf("%sState %d,Visits %v, Utility %v, Action %s\n", indent, node.state, node.visits, node.utility, node.action)
	printed[node.state.StateId()] = node

	for _, child := range node.children {
		_, ok := printed[child.state.StateId()]
		if !ok {
			printTree(child, depth + 1)
		}
	}
}

func mostVisited(nodes []*Node) *Node {
	// fmt.Println(len(nodes))
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

func selectNode(node *Node) *Node {
	// fmt.Println("selectNode at:", node)
	if node.visits == 0 || node.children == nil || len(node.children) == 0{
		return node
	}
	// fmt.Println(".. not visited")

	for _, child := range node.children {
		// fmt.Println(".. Child: %s", child.action)
		if child.visits == 0 { 
			// fmt.Println(".. zero visits..")
			return child
		}
	}

	// fmt.Println(".. ALL visited. Selecting Highest score")

	score := float64(math.MinInt64)
	result := node

	parentVisits := 0
	for _, child := range node.children {
		parentVisits += child.visits
	}
	for _, child := range node.children {
		newscore := selectfn(child, parentVisits)
		// fmt.Println(".. Child:", i, " ", child, " score", newscore)
		if newscore > score {
			score = newscore
			result = child
		}
	}
	return selectNode(result)
}

func expand(node *Node) { //*Node {
	// actions := node.state.FindLegals(node.state.player)
	actions := node.state.FindLegals()
	for _, action := range actions {
		// fmt.Printf("..Expanding ACTION %s\n", action)
		newState := node.state.FindNextState(action)
		// fmt.Printf("..state %v\n", newState)
		// check if state is visited
		stateId := newState.StateId()
		visitedNode, ok := visitedStates[stateId]
		if ok {
			// fmt.Printf("..already visited %v\n", newState)
			visitedNode.action = action
			node.children = append(node.children, visitedNode)
		} else {
			// fmt.Printf("..new state %v\n", newState)
			newNode := &Node{newState, 0, 0, action, node, []*Node{}}
			visitedStates[stateId] = newNode
			node.children = append(node.children, newNode)

			reward := simulate(newNode.state)
		fmt.Println("..Propagation")
			backPropagate(newNode, reward)



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

