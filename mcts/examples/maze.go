package main

import (
	"github.com/dranidis/go-skat/mcts"
	"fmt"
	"log"
)


type MazeState struct {
	tile int
	fromTile int
}

type MazeAction struct {
	direction int // 0 up, 1 right, 2 down, 3 left
}

func (ma MazeAction) String() string {
	d := ""
	switch ma.direction {
	case 0: d = "UP"
	case 1: d = "RI"
	case 2: d = "DN"
	case 3: d = "LE"
	}
	return d
}

func (ma MazeAction) Equals(a mcts.Action) bool {
	mma := a.(MazeAction)
	if mma.direction == ma.direction {
		return true
	}
	return false
}

func (m *MazeState) StateId() uint64 {
	return uint64(m.tile * 100 + m.fromTile)
}

func (m *MazeState) FindNextState(a mcts.Action) mcts.State {
	currentTile := m.tile
	var nextTile int
	ma := a.(MazeAction)
	switch ma.direction {
	case u: 
		nextTile = currentTile - 5
	case d:
		nextTile = currentTile + 5
	case r:
		nextTile = currentTile + 1
	case l:
		nextTile = currentTile - 1
	}
	if nextTile < 0 || nextTile > 14 {
		log.Fatal("FindNextState, Out of bounds index, nextTile:", nextTile, " currentTile", currentTile, " action ", ma)
	}
	steps ++
	var state mcts.State
	state = &MazeState{nextTile, currentTile}
	// fmt.Printf("Moving from tile %d to tile %d\n", currentTile, nextTile)
	return state
}

func (m *MazeState) IsTerminal() bool {
	return m.tile == 1
}

func (m *MazeState) FindReward() float64 {
	reward := -1.0  * float64(steps)
	steps = 0
	return reward
}

func (m *MazeState) FindLegals() []mcts.Action {
	if m.tile < 0 || m.tile > 14 {
		log.Fatal("FindLegals, Out of bounds index, m.tile:", m.tile)
	}	
	var actions []mcts.Action
	for _, move := range moves[m.tile] {
		actions = append(actions, MazeAction{move})
	}
	return actions
}

const (
	u = 0
	r = 1
	d = 2
	l = 3
)

var moves [][]int
var steps = 0

func makeMaze() {
	moves = [][]int{
		[]int{d},
		[]int{u, r},
		[]int{l, r},
		[]int{l, d},
		[]int{d},

		[]int{u,d,r},
		[]int{l,r},
		[]int{l,r},
		[]int{l,u},
		[]int{d,u},

		[]int{r,u},
		[]int{l,r},
		[]int{l,r},
		[]int{l,r},
		[]int{l,u},		
	}
	// moves[0] = []int{d}
	// moves[1] = []int{u, r}
	// moves[2] = []int{l, r}
	// moves[3] = []int{l, d}
	// moves[4] = []int{d}
	// moves[5] = []int{u,d,r}
	// moves[6] = []int{l,r}
	// moves[7] = []int{l,r}
	// moves[8] = []int{l,u}
	// moves[9] = []int{d,u}
	// moves[10] = []int{r,u}
	// moves[11] = []int{l,r}
	// moves[12] = []int{l,r}
	// moves[13] = []int{l,r}
	// moves[14] = []int{l,u}
}

func main() {
	makeMaze()

	// s :=  MazeState{10}
	// actions := s.FindLegals()
	// for _, action := range actions {
	// 	fmt.Println(action)
	// 	s.FindNextState(action)
	// }


	var state mcts.State
	initial := &MazeState{0, 0}
	state = initial
	// mcts.InitSubtree(state)
	// for ! state.IsTerminal() {
	for i := 0; i < 100; i++ {
		a := mcts.Uct(state, 3000)
		fmt.Println("PERFORMING ACTION: ", a)		
		state = state.FindNextState(a)		
		// mcts.ChangeSubtree(a)
		// fmt.Println("PERFORMING ACTION: ", a)		
	}
}