package main

import (
	"github.com/dranidis/go-skat/mcts"
	"fmt"
	"log"
	"math/rand"
	"time"
)


type TTTState struct {
	board []int // -1 for opp, 1 for player, 0 empty
	player int // whose turn it is, // -1 opp, 1 player
}

type TTTAction struct {
	square int // 0..8
}

func (m TTTState) IsOpponentTurn() bool {
	return m.player == -1
}


func (m TTTState) String() string {
	s := "\t"
	for i :=0 ; i < 9; i ++ {
		if i % 3 == 0 {
			s += "\n\t"
		}
		switch m.board[i] {
		case 0:
			s += "."
		case 1:
			s += "X"
		case -1:
			s += "O"
		}
	}
	return s
}

func (ma TTTAction) String() string {
	return fmt.Sprintf("%d", ma.square)
}

func (ma TTTAction) Equals(a mcts.Action) bool {
	mma := a.(TTTAction)
	return mma.square == ma.square
}

func (m *TTTState) StateId() uint64 {
	// convert to a decimal number using 3-base
	n := uint64(0)
	for i :=0 ; i < 9; i ++ {
		n *= uint64(3)
		n += (uint64(m.board[i]) + 1)
	}
	return n
}

func (m *TTTState) FindNextState(a mcts.Action) mcts.State {
	ma := a.(TTTAction)
	
	if m.board[ma.square] != 0 {
		log.Fatal("illegal move ", ma.square, " on board ", m.String())
	}
	newBoard := make([]int, 9)
	copy(newBoard, m.board)

	newBoard[ma.square] = m.player

	var state mcts.State
	state = &TTTState{newBoard, -1 * m.player}
	return state
}

func (m *TTTState) sumRow(row int) int {
	sum := 0
	for c:= 0; c < 3 ; c++ {
		sum += m.board[3 * row + c]
	}
	return sum
}

func (m *TTTState) sumCol(col int) int {
	sum := 0
	for r:= 0; r < 3 ; r++ {
		sum += m.board[col + r *3]
	}
	return sum
}

func (m *TTTState) winnerAux(s int) bool {
	r0 := m.sumRow(0)
	r1 := m.sumRow(1)
	r2 := m.sumRow(2)
	if r0 == s || r1 == s ||  r2 == s {
		return true
	}
	c0 := m.sumCol(0)
	c1 := m.sumCol(1)
	c2 := m.sumCol(2)
	if c0 == s || c1 == s ||  c2 == s {
		return true
	}
	if m.board[0] + m.board[4] + m.board[8] == s {
		return true
	}
	if m.board[2] + m.board[4] + m.board[6] == s {
		return true
	}
	return false
}

func (m *TTTState) winner() int {
	if m.winnerAux(3) {
		return 1
	}
	if m.winnerAux(-3) {
		return -1
	} 
	return 0
}

func (m *TTTState) IsTerminal() bool {
	w := m.winner()
	if w == 1 || w == -1 {
		return true
	}
	for _,c := range m.board {
		if c == 0 {return false}
	}
	return true
}

func (m *TTTState) FindReward() float64 {
	return float64(m.winner())
}

func (m *TTTState) FindLegals() []mcts.Action {
	var actions []mcts.Action

	for i := 0; i < len(m.board); i++ {
		// fmt.Printf("Board: %d\n", m.board[i])
		if m.board[i] == 0 {
			actions = append(actions, TTTAction{i})
		}
	}
	return actions
}

func makeTTT() TTTState {
	board := []int{
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,	
	}
	player := 1
	return TTTState{board, player}
}

func main() {

	r := rand.New(rand.NewSource(time.Now().Unix()))
	randSeed := int64(r.Intn(9999))
	rand.Seed(randSeed)

	mcts.SimulationRuns = 1000
	mcts.ExplorationParameter = 2.0
	mcts.DEBUG = false // You have to increase delay to 2000 if you are dedugging to give time for runs
	runMilliseconds := 500

	gameStats := []int{0,0,0}
	for i := 0; i < 1; i++ {
		fmt.Print(".")
		var state mcts.State
		initial := makeTTT()
		state = &initial
		for !state.IsTerminal() {
			a, _ := mcts.Uct(state, runMilliseconds)
			fmt.Println("PERFORMING ACTION: ", a)		
			state = state.FindNextState(a)	
			fmt.Printf("%v\n\n", state)	
		}
		game := state.(*TTTState)
		w := game.winner()
		// fmt.Printf("Winner: %d\n", w)	
		gameStats[w + 1] ++
	}

	fmt.Println(gameStats)
	
}