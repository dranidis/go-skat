package main 

import (
	// "github.com/dranidis/go-skat/minimax"
	"github.com/dranidis/go-skat/mcts"
)

type SkatState struct {
	trump         string
	playerHand    [][]Card // YOU, 1, 2
	trick         []Card
	declarer      int
	turn          int // 0 you, 1 player 1, 2 player2
	declScore     int
	oppScore      int
	schneiderGoal bool
}

type SkatAction struct {
	card Card
}

func (m SkatState) Heuristic() float64 {
	if m.IsTerminal() {
		return m.FindRewardNum()
	} else {
		return 0 /// ????????????????
	}
}

// func (m SkatState) IsOpponentTurn() bool {
// 	if m.declarer == 0 && m.turn == 0 {
// 		return false
// 	}
// 	if m.declarer == 0 && m.turn != 0 {
// 		return true
// 	}
// 	if m.declarer == m.turn {
// 		return true
// 	}
// 	return false
// }


func (m SkatState) IsOpponentTurn() bool {
	return m.declarer != m.turn
}

func (m *SkatState) IsTerminal() bool {
	if !m.schneiderGoal {
		if m.declScore > 60 {
			return true
		}
		if m.oppScore > 59 {
			return true
		}
	}
	return len(m.playerHand[0])+len(m.playerHand[1])+len(m.playerHand[2]) == 0
}

func (m *SkatState) FindRewardNum() float64 {
	// winsScore := 61
	// if m.schneiderGoal {
	// 	winsScore = 90
	// }

	// if m.declarer == 0 { //YOU
	// 	if m.declScore >= winsScore {
	// 		return float64(1.0)
	// 	} else {
	// 		return float64(0.0)
	// 	}
	// }
	// if m.declScore > winsScore - 1 {
	// 	return float64(0.0)
	// }
	return float64(m.declScore) // TODO
}

func (m *SkatState) FindReward() float64 {
	winsScore := 61
	if m.schneiderGoal {
		winsScore = 90
	}

	if m.declarer == 0 { //YOU
		if m.declScore >= winsScore {
			return float64(1.0)
		} else {
			return float64(0.0)
		}
	}
	if m.declScore > winsScore-1 {
		return float64(0.0)
	}
	return float64(1.0) // TODO
}

func (m SkatState) validCards(cards []Card) []Card {
	if len(m.trick) == 0 {
		return cards
	}
	return filter(cards, func(c Card) bool {
		return valid(getSuit(m.trump, m.trick[0]), m.trump, cards, c)
	})
}

func (m *SkatState) FindLegals() []mcts.Action {
	actions := []mcts.Action{}
	for _, card := range m.validCards(m.playerHand[m.turn]) {
		actions = append(actions, SkatAction{card})
	}
	return actions
}

func (m *SkatState) FindNextState(a mcts.Action) mcts.State {
	ma := a.(SkatAction)

	// deep copy before you make any changes
	newState := copySkatState(*m)

	if len(newState.trick) == 3 {
		newState.trick = []Card{}
	}
	// remove the card from the player and add it to the trick
	newState.playerHand[m.turn] = remove(newState.playerHand[m.turn], ma.card)
	newState.trick = append(newState.trick, ma.card)

	strump := newState.trump
	sfollow := getSuit(strump, newState.trick[0])

	if len(newState.trick) == 3 {
		winnerCard := -1
		// find winner
		if greater(strump, sfollow, newState.trick[0], newState.trick[1], newState.trick[2]) {
			winnerCard = 0
			newState.turn = m.turn + 1
		} else if greater(strump, sfollow, newState.trick[1], newState.trick[0], newState.trick[2]) {
			winnerCard = 1
			newState.turn = m.turn + 2
		} else {
			winnerCard = 2
			newState.turn = m.turn
		}
		if newState.turn > 2 {
			newState.turn -= 3
		}

		// set the scores, depending on who played when
		if m.turn == 0 { // YOU PLAYED THE 3RD CARD
			newState.setScores(2, winnerCard)
		}
		if m.turn == 1 { // YOU PLAYED THE 2nd CARD
			newState.setScores(1, winnerCard)
		}
		if m.turn == 2 { // YOU PLAYED THE 1st CARD
			newState.setScores(0, winnerCard)
		}

	} else {
		// set next player turn
		newState.turn++
		if newState.turn > 2 {
			newState.turn = 0
		}
	}

	var state mcts.State
	state = &newState
	return state
}

func (m *SkatState) setScores(w int, winner int) {
	if winner == w {
		// you won
		if m.declarer == 0 {
			m.declScore += sum(m.trick)
		} else {
			m.oppScore += sum(m.trick)
		}
	} else {
		// you lost
		if m.declarer == 0 {
			m.oppScore += sum(m.trick)
		} else {
			m.declScore += sum(m.trick)
		}
	}
}

func copySkatState(m SkatState) SkatState {
	strick := make([]Card, len(m.trick))
	copy(strick, m.trick)

	p0hand := make([]Card, len(m.playerHand[0]))
	copy(p0hand, m.playerHand[0])
	p1hand := make([]Card, len(m.playerHand[1]))
	copy(p1hand, m.playerHand[1])
	p2hand := make([]Card, len(m.playerHand[2]))
	copy(p2hand, m.playerHand[2])

	return SkatState{
		m.trump,
		[][]Card{
			p0hand,
			p1hand,
			p2hand,
		},
		strick,
		m.declarer,
		m.turn,
		m.declScore,
		m.oppScore,
		m.schneiderGoal,
	}
}
