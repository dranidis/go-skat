package main

import (
	 // "fmt"
	// "github.com/dranidis/go-skat/minimax"
)

type MinMaxPlayer struct {
	Player
	p1Hand []Card
	p2Hand []Card
}

func makeMinMaxPlayer(hand []Card) MinMaxPlayer {
	return MinMaxPlayer {
		Player:     makePlayer(hand),
		p1Hand: []Card{},
		p2Hand: []Card{},
	}
}

func (p *MinMaxPlayer) playerTactic(s *SuitState, c []Card) Card {
	debugTacticsLog("MinMaxPlayer\n") 

	if len(c) == 1 {
		debugTacticsLog("..FORCED MOVE.. ")
		return c[0]
	}

	if len(p.hand) < 4 {
		return p.minMaxTactics(s, c)
	}

	return p.Player.playerTactic(s, c)
}

func (p MinMaxPlayer) minMaxTactics(s *SuitState, c []Card) Card {

	p.dealCardsToOpponents(s)


	return c[0]
}

func (p *MinMaxPlayer) dealCardsToOpponents(s *SuitState) {
	cards := makeDeck()
	cards = remove(cards, s.cardsPlayed...)
	cards = remove(cards, s.trick...)
	cards = remove(cards, p.hand...)

	if p == s.declarer {
		cards = remove(cards, s.skat...)
	}

	cards = Shuffle(cards)

	if p != s.declarer { // remove two random cards for the skat
		card1 := cards[0]
		card2 := cards[1]
		cards = remove(cards, card1, card2)
	}
	half := len(cards) / 2

	var player1 PlayerI
	var player2 PlayerI
	if len(s.trick) == 0 {
		player1 = players[1]
		player2 = players[2]
	}
	if len(s.trick) == 1 {
		player1 = players[0]
		player2 = players[2]
	}
	if len(s.trick) == 2 {
		player1 = players[0]
		player2 = players[1]
	}
	
	// } else { // len 1: opp2 you opp1
	// 	p.p1Hand = cards[:half]
	// 	p.p2Hand = cards[half:len(cards)]		
	// }
	
	p.p1Hand = cards[:half]
	p.p2Hand = cards[half:len(cards)]	

	debugTacticsLog("MINMAX: Remaining cards %v, %s: %v, %s: %v\n", cards, player1.getName(), p.p1Hand,  player2.getName(), p.p2Hand)
}


// type SkatState struct {
// 	board []int // -1 for opp, 1 for player, 0 empty
// 	player int // whose turn it is, // -1 opp, 1 player
// }

// type TTTAction struct {
// 	square int // 0..8
// }

// func (m SkatState) Heuristic() float64 {
// 	if m.IsTerminal() {
// 		return m.FindReward()
// 	} else {
// 		return 0 /// ????????????????
// 	}
// }

// func (m SkatState) IsOpponentTurn() bool {
// 	return m.player == -1
// }


// func (m SkatState) String() string {
// 	// return fmt.Sprintf("%v", m.board)
// 	s := "\t"
// 	for i :=0 ; i < 9; i ++ {
// 		if i % 3 == 0 {
// 			s += "\n\t"
// 		}
// 		switch m.board[i] {
// 		case 0:
// 			s += "."
// 		case 1:
// 			s += "X"
// 		case -1:
// 			s += "O"
// 		}
// 	}
// 	return s
// }

// func (ma TTTAction) String() string {
// 	return fmt.Sprintf("%d", ma.square)
// }

// func (ma TTTAction) Equals(a minimax.Action) bool {
// 	mma := a.(TTTAction)
// 	return mma.square == ma.square
// }

// func (m *SkatState) StateId() uint64 {
// 	// convert to a decimal number using 3-base
// 	n := uint64(0)
// 	for i :=0 ; i < 9; i ++ {
// 		n *= uint64(3)
// 		n += (uint64(m.board[i]) + 1)
// 	}
// 	return n
// }

// func (m *SkatState) FindNextState(a minimax.Action) minimax.State {
// 	ma := a.(TTTAction)
	
// 	if m.board[ma.square] != 0 {
// 		log.Fatal("illegal move ", ma.square, " on board ", m.String())
// 	}
// 	newBoard := make([]int, 9)
// 	copy(newBoard, m.board)

// 	newBoard[ma.square] = m.player

// 	var state minimax.State
// 	state = &SkatState{newBoard, -1 * m.player}
// 	return state
// }
