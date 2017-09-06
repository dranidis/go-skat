package mcts

import (
	"math"
	"math/rand"
	"fmt"
	"log"
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
}

// type Player interface {

// }

type Action interface {
	Equals(Action) bool
}

var	rootNode *Node

func InitSubtree(state State) {
	rootNode = &Node{state, 0, 0, nil, nil, []*Node{}}
}

func ChangeSubtree(action Action) {
	for _, child := range rootNode.children {
		if action.Equals(child.action) {
			rootNode = child
			return
		}
	}
	log.Fatal("Action not found")
}

func Uct(state State, iters int) Action {
	// rootNode := Node{state, 0, 0, nil, nil, []*Node{}}
	root := rootNode
	for i := 0; i < iters; i++ {
		// fmt.Println("Iteration: ", i)
		// fmt.Println("..Selection")
		node := selectNode(root)
		// fmt.Println("..Expansion")
		node = expand(node)
		// printTree(root, 0)
		// fmt.Println("..Simulation")
		reward := simulate(node.state)
		// fmt.Println("..Propagation")
		backPropagate(node, reward)
	}
	printTree(root, 0)
	// return mostVisited(root.children).action
	return mostUtility(root.children).action
}

func printTree(node *Node, depth int) {
	indent := ""
	for i := 0; i < depth ; i++ {
		indent += "\t"
	}

	// action := ""
	fmt.Printf("%sState %d,Visits %v, Utility %v, Action %s\n", indent, node.state, node.visits, node.utility, node.action)
	for _, child := range node.children {
		printTree(child, depth + 1)
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

func mostUtility(nodes []*Node) *Node {
	// fmt.Println(len(nodes))
	most := float64(math.MinInt64)
	imost := -1
	for i, node := range nodes {
		if node.utility > most {
			most = node.utility
			imost = i
		}
	}
	if imost >= len(nodes) {
		imost = 0
	}
	return nodes[imost]
}

func selectNode(node *Node) *Node {
		// fmt.Println("selectNode:", node)
	if node.visits == 0 || node.children == nil || len(node.children) == 0{
		return node
	}
	for _, child := range node.children {
		// fmt.Println(".. Child:", i)
		if child.visits == 0 { 
			// fmt.Println(".. Returning..")
			return child
		}
	}
	score := float64(math.MinInt64)
	result := node

	for _, child := range node.children {
		newscore := selectfn(child)
		// fmt.Println(".. Child:", i, " ", child, " score", newscore)
		if newscore > score {
			score = newscore
			result = child
		}
	}
	return selectNode(result)
}

func expand(node *Node) *Node {
	// actions := node.state.FindLegals(node.state.player)
	actions := node.state.FindLegals()
	for _, action := range actions {
		newState := node.state.FindNextState(action)
		newNode := &Node{newState, 0, 0, action, node, []*Node{}}
		node.children = append(node.children, newNode)
	}
	return node.children[0]
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
	node.utility /= float64(node.visits)
	if node.parent != nil {
		backPropagate(node.parent, score)
	}
}

func selectfn(node *Node) float64 {
	return float64(node.utility) + 2.0 * math.Sqrt(math.Log(float64(node.parent.visits)) / float64(node.visits))
}

