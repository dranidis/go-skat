package main

import (
	"fmt"
	"github.com/dranidis/go-skat/game"
	"log"
)

type SkatState struct {
	SuitState
	players []PlayerI // YOU, 1, 2
}

type SkatAction struct {
	card Card
}

func (s SkatAction) String() string {
	return fmt.Sprintf("%s", s.card)
}

func (s SkatState) String() string {
	return fmt.Sprintf("%v, %s, %v %T, TRICK: %v, SCORE: %d - %d",
		s.SuitState, s.trump, s.players, s.players[0], s.trick,
		s.declarer.getScore(), s.opp1.getScore()+s.opp2.getScore())
}

func (m SkatState) Heuristic() float64 {
	if m.IsTerminal() {
		// fmt.Printf("R: %v\n", m.FindRewardNum())
		return m.FindRewardNum()
	} else {
		m.playToTheEndWithTactics()
		// fmt.Printf("R: %v\n", m.FindRewardNum())
		return m.FindRewardNum()
	}
}

func copyPlayer(p PlayerI) Player {
	player := makePlayer(p.getHand())
	player.name = p.getName()
	player.score = p.getScore()
	player.previousSuit = p.getPreviousSuit()
	player.schwarz = p.isSchwarz()
	// player.highestBid = p.getHighestBid()
	player.declaredBid = p.getDeclaredBid()
	return player
}

func copyPlayerMM(p PlayerI) MinMaxPlayer {
	player := makeMinMaxPlayer(p.getHand())
	player.name = p.getName()
	player.score = p.getScore()
	player.previousSuit = p.getPreviousSuit()
	player.schwarz = p.isSchwarz()
	// player.highestBid = p.getHighestBid()
	player.declaredBid = p.getDeclaredBid()
	// debugTacticsLog("copyPlayerMM: %v %v\n", p, player)
	return player
}

func (m *SkatState) playToTheEndWithTactics() {
	cpuPlayers := make([]Player, 3)
	for i := 0; i < 3; i++ {
		cpuPlayers[i] = copyPlayer(m.players[i])
		if m.declarer.getName() == m.players[i].getName() {
			m.declarer = &cpuPlayers[i]
		}
		if m.opp1.getName() == m.players[i].getName() {
			m.opp1 = &cpuPlayers[i]
		}
		if m.opp2.getName() == m.players[i].getName() {
			m.opp2 = &cpuPlayers[i]
		}
		if m.leader.getName() == m.players[i].getName() {
			m.leader = &cpuPlayers[i]
		}
		// debugTacticsLog("CPU Player %v - %v\n", cpuPlayers[i], m.players[i])
		m.players[i] = &cpuPlayers[i]
		// debugTacticsLog("CPU Player %v - %v\n", cpuPlayers[i], m.players[i])
	}

	disableLogs()

	for len(m.players[2].getHand()) > 0 {
		_ = m.moveOne()
		// fmt.Printf("%v ", c)
	}

	restoreLogs()
}

// TODO:
// Replay all played tricks so that they get analysed
// to be used by tactics.
func (m SkatState) GetTacticsMove() game.Action {
	tmpState := m.copySkatState()
	cpuPlayers := make([]Player, 3)
	for i := 0; i < 3; i++ {
		cpuPlayers[i] = makePlayer(m.players[i].getHand())
		cpuPlayers[i].name = m.players[i].getName()
		// cpuPlayers[i].name += "-TAC1"
		if m.declarer.getName() == m.players[i].getName() {
			tmpState.declarer = &cpuPlayers[i]
		}
		if m.opp1.getName() == m.players[i].getName() {
			tmpState.opp1 = &cpuPlayers[i]
		}
		if m.opp2.getName() == m.players[i].getName() {
			tmpState.opp2 = &cpuPlayers[i]
		}
		if m.leader.getName() == m.players[i].getName() {
			tmpState.leader = &cpuPlayers[i]
		}
	}

	disableLogs()

	card := tmpState.moveOne()

	restoreLogs()

	return SkatAction{card}
}

func (m *SkatState) moveOne() Card {
	// debugTacticsLog("MOVEONE %v\n", m)
	l := len(m.trick)
	var card Card
	if l == 0 {
		card = play(&m.SuitState, m.players[0])
		m.follow = getSuit(m.trump, m.trick[0])
	}
	if l == 1 { // USING else if bevause play changes the s.trick
		card = play(&m.SuitState, m.players[1])
	}
	if l == 2 {
		card = play(&m.SuitState, m.players[2])
		// var players = []PlayerI{&m.players[0], &m.players[1], &m.players[2]}
		m.players = setNextTrickOrder(&m.SuitState, m.players)
		// for i := 0; i < 3; i++ {
		// 	var p = players[i].(*MinMaxPlayer)
		// 	m.players[i] = *p
		// }
		m.follow = ""
	}

	return card
}

func (m SkatState) IsOpponentTurn() bool {
	return m.players[len(m.trick)].getName() != m.declarer.getName()
}

func (m *SkatState) IsTerminal() bool {
	for i := 0; i < 3; i++ {
		if len(m.players[i].getHand()) > 0 {
			// debugTacticsLog("Player %v\n", m.players[i])
			return false
		}
	}
	return true
}

func (m *SkatState) FindRewardNum() float64 {
	// return float64(m.declarer.getScore() - m.opp1.getScore() - m.opp2.getScore())
	return float64(m.declarer.getScore())
}

func (m *SkatState) FindReward() float64 {
	log.Fatal("Not used.")
	return float64(0.0) // TODO
}

func (m SkatState) validCards(cards []Card) []Card {
	if len(m.trick) == 0 {
		return cards
	}
	return filter(cards, func(c Card) bool {
		return valid(getSuit(m.trump, m.trick[0]), m.trump, cards, c)
	})
}

func (m *SkatState) cardIsLosingTheTrick(card Card) bool {
	// return only following cards
	if getSuit(m.trump, card) != m.follow {
		return false
	}

	if len(m.trick) == 2 {
		if m.players[2].getName() == m.opp1.getName() || m.players[2].getName() == m.opp2.getName() {
			//opponent is playing last
			dIndex := 0
			for dIndex, _ = range m.players {
				if m.players[dIndex].getName() == m.declarer.getName() {
					break
				}
			}
			other := 0
			if dIndex == 0 {
				other = 1
			}
			if m.greater(m.trick[dIndex], card) || m.greater(m.trick[dIndex], m.trick[other]) {
				return true
			}
		} else {
			// declarer is playing last
			if m.greater(m.trick[0], card) && m.greater(m.trick[1], card) {
				return true
			}
		}
	}
	// currently only for last card
	return false
}

func (m *SkatState) FindLegals() []game.Action {
	actions := []game.Action{}

	validCards := m.validCards(m.players[len(m.trick)].getHand())
	// if len(m.trick) == 2 {
	// 	losers := []Card{}
	// 	for _, card := range validCards {
	// 		if m.cardIsLosingTheTrick(card) {
	// 			losers = append(losers, card)
	// 		}
	// 	}
	// 	if len(losers) > 1 {
	// 		losers = realSortValue(losers)
	// 		validCards = remove(validCards, losers...)
	// 		validCards = append(validCards, losers[len(losers) - 1])
	// 	}
	// }
	for _, card := range validCards {
		actions = append(actions, SkatAction{card})
	}
	return actions
}

func (m *SkatState) FindNextState(a game.Action) game.State {
	disableLogs()
	ma := a.(SkatAction)

	// deep copy before you make any changes
	newState := m.copySkatState()

	currentP := newState.players[len(newState.trick)]

	analysePlay(&newState.SuitState, currentP, ma.card)

	currentP.setHand(remove(currentP.getHand(), ma.card))
	newState.players[len(newState.trick)] = currentP // STRANGE!!
	newState.trick = append(newState.trick, ma.card)
	if getSuit(newState.trump, ma.card) == newState.trump {
		newState.trumpsInGame = remove(newState.trumpsInGame, ma.card)
	}
	newState.cardsPlayed = append(newState.cardsPlayed, ma.card)

	if len(newState.trick) == 1 {
		newState.follow = getSuit(newState.trump, newState.trick[0])
	}

	if len(newState.trick) == 3 {
		newState.players = setNextTrickOrder(&newState.SuitState, newState.players)
		newState.follow = ""
	}
	var state game.State
	state = &newState
	restoreLogs()
	return state
}

func (m SkatState) copySkatState() SkatState {
	suitState := m.SuitState.cloneSuitStateNotPlayers()
	// copyPlayers := make([]MinMaxPlayer, 3)
	copyIPlayers := make([]PlayerI, 3)
	for i := 0; i < 3; i++ {
		p := m.players[i].clone()

		// var clone = p.(*MinMaxPlayer)
		// copyPlayers[i] = *clone
		copyIPlayers[i] = p

		if m.declarer.getName() == copyIPlayers[i].getName() {
			suitState.declarer = copyIPlayers[i]
		}
		if m.opp1.getName() == copyIPlayers[i].getName() {
			suitState.opp1 = copyIPlayers[i]
		}
		if m.opp2.getName() == copyIPlayers[i].getName() {
			suitState.opp2 = copyIPlayers[i]
		}
		if m.leader.getName() == copyIPlayers[i].getName() {
			suitState.leader = copyIPlayers[i]
		}
	}

	return SkatState{
		suitState,
		copyIPlayers,
	}
}
