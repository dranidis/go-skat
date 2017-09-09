package main

import (
	 // "fmt"
	"github.com/dranidis/go-skat/minimax"
)

type MinMaxPlayer struct {
	Player
	p1Hand []Card
	p2Hand []Card
	maxHandSize int
}

func makeMinMaxPlayer(hand []Card) MinMaxPlayer {
	return MinMaxPlayer {
		Player:     makePlayer(hand),
		p1Hand: []Card{},
		p2Hand: []Card{},
		maxHandSize: 4,
	}
}

func (p *MinMaxPlayer) playerTactic(s *SuitState, c []Card) Card {

	if len(c) == 1 {
		debugTacticsLog("..FORCED MOVE.. ")
		return c[0]
	}

	if len(p.hand) <= p.maxHandSize {
		// minimax.DEBUG = true
		cardsFreq := make(map[string]int)
		cards := make(map[string]Card)
		for i := 0; i < 10; i++ {
			debugTacticsLog("MinMaxPlayer\n") 
			card := p.minMaxTactics(s, c)
			v, ok := cardsFreq[card.String()]	
			if ok {
				cardsFreq[card.String()] = v + 1
			} else {
				cardsFreq[card.String()] = 1
			}
			cards[card.String()] = card	
		}
		most := 0
		var card Card
		for k, v := range cardsFreq { 
			if v > most {
				most, card = v, cards[k] 
			}
		}
		return card
	}

	player := p.Player
	debugTacticsLog("MINMAX player %v", player)
	return player.playerTactic(s, c)
}

func (p MinMaxPlayer) minMaxTactics(s *SuitState, c []Card) Card {
	p.dealCards(s)
	card := p.minmaxSkat(s)

	return card
}

func checkVoidOpp1(s *SuitState, p *MinMaxPlayer, cards []Card, suit string) []Card {
	if s.opp1VoidSuit[suit] {
		suitCards := filter(cards, func(c Card) bool {
			return getSuit(s.trump, c) == suit
			})
		p.p2Hand = append(p.p2Hand, suitCards...)
		cards = remove(cards, suitCards...)
	}
	return cards		
}

func checkVoidOpp2(s *SuitState, p *MinMaxPlayer, cards []Card, suit string) []Card {
	if s.opp2VoidSuit[suit] {
		suitCards := filter(cards, func(c Card) bool {
			return getSuit(s.trump, c) == suit
			})
		p.p1Hand = append(p.p1Hand, suitCards...)
		cards = remove(cards, suitCards...)
	}
	return cards	
}

func (p *MinMaxPlayer) dealCards(s *SuitState) {
	cards := makeDeck()
	cards = remove(cards, s.cardsPlayed...)
	cards = remove(cards, s.trick...)
	cards = remove(cards, p.hand...)

	if p.getName() == s.declarer.getName() {
		cards = remove(cards, s.skat...)
	}

	debugTacticsLog("ALL CARDS: %v\n", cards)

	if p.getName() == s.declarer.getName() {
		for _, suit := range suits {
			cards = checkVoidOpp1(s, p, cards, suit)	
			cards = checkVoidOpp2(s, p, cards, suit)	
		}
	}
	debugTacticsLog("REMAINING after void: %d cards: %v %v %v\n", len(cards), cards, p.p1Hand, p.p2Hand)

	// TODO
	// depending on the number of remaining cards
	// limit the number of possible worlds.


	// shuffle the rest
	cards = Shuffle(cards)

	if p != s.declarer { // remove two random cards for the skat
		card1 := cards[0]
		card2 := cards[1]
		cards = remove(cards, card1, card2)
	}

	handSize := len(p.hand)
	nextCard := 0

	leader := 0
	middle := 0
	if len(s.trick) == 1 {
		leader = 1
	}
	if len(s.trick) == 2 {
		middle = 1
		leader = 1
	}

	for i := len(p.p1Hand); i < handSize - middle; i++ {
		p.p1Hand = append(p.p1Hand, cards[nextCard])
		nextCard++
	}

	for i := len(p.p2Hand); i < handSize - leader; i++ {
		p.p2Hand = append(p.p2Hand, cards[nextCard])
		nextCard++
	}
}

func (p *MinMaxPlayer) minmaxSkat(s *SuitState) Card {
	var player1 PlayerI
	var player2 PlayerI
	if len(s.trick) == 0 {
		player1 = players[1]
		player2 = players[2]
	}
	if len(s.trick) == 1 {
		player1 = players[2]
		player2 = players[0]
	}
	if len(s.trick) == 2 {
		player1 = players[0]
		player2 = players[1]
	}
	var decl int // 0 is you, 1 is next player1, 2 is next player 2
	if s.declarer.getName() == p.getName() {
		decl = 0
	} else if s.declarer.getName() == player1.getName() {
		decl = 1
	} else {
		decl = 2
	}

	debugTacticsLog("MINMAX: cards %s: %v, %s: %v\n", player1.getName(), p.p1Hand,  player2.getName(), p.p2Hand)
	debugTacticsLog("Decl: %d\n", decl)

	strick := make([]Card, len(s.trick))
	copy(strick, s.trick)

	dist := make([][]Card, 3)
	dist[0] = p.hand
	dist[1] = p.p1Hand
	dist[2] = p.p2Hand

	debugTacticsLog("CARDS: %v", dist)
	debugTacticsLog("CARDS[0]: %v", dist[0])
	debugTacticsLog("CARDS[1]: %v", dist[1])
	debugTacticsLog("CARDS[2]: %v", dist[2])

	skatState := SkatState{
		s.trump,
		dist, 			
		strick, 
		decl, 
		0, 
		s.declarer.getScore(), 
		s.opp1.getScore() + s.opp2.getScore(),
	}
	debugTacticsLog("Skatstate Cards: %v\n",skatState.playerHand)

	a := minimax.Minimax(&skatState)
	ma := a.(SkatAction)

	debugTacticsLog("In hand %v, MINIMAX suggesting card: %v\n", p.hand, ma.card)

	return ma.card
}


type SkatState struct {
	trump string
	playerHand [][]Card  // YOU, 1, 2
	trick []Card
	declarer int
	turn int // 0 you, 1 player 1, 2 player2
	declScore int 
	oppScore int
}

type SkatAction struct {
	card Card 
}

func (m SkatState) Heuristic() float64 {
	if m.IsTerminal() {
		return m.FindReward()
	} else {
		return 0 /// ????????????????
	}
}

func (m SkatState) IsOpponentTurn() bool {
	if m.turn == 0 {
		return false
	}
	if m.declarer == m.turn {
		return true
	}
	return false
}

func (m *SkatState) IsTerminal() bool {
	return len(m.playerHand[0]) + len(m.playerHand[1]) + len(m.playerHand[2]) == 0
}

func (m *SkatState) FindReward() float64 {
	if m.declarer == 0 { //YOU
		if m.declScore > 60 {
			return float64(1.0)
		} else {
			return float64(0.0)
		}
	}
	if m.declScore > 60 {
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

func (m *SkatState) FindLegals() []minimax.Action {
	actions := []minimax.Action{}
	for _, card := range m.validCards(m.playerHand[m.turn]) {
		actions = append(actions, SkatAction{card})
	}
	return actions
}

func (m *SkatState) FindNextState(a minimax.Action) minimax.State {
	ma := a.(SkatAction)
	
	// deep copy before you make any changes
	newState := copySkatState(*m)
	// remove the card from the player and add it to the trick
	newState.playerHand[m.turn] = remove(newState.playerHand[m.turn], ma.card)
	newState.trick = append(newState.trick, ma.card)

	strump := newState.trump
	sfollow := getSuit(newState.trump, newState.trick[0])

	if len(newState.trick) == 3 {
		winnerCard := -1
		// find winner
		if greater(strump, sfollow, newState.trick[0], newState.trick[1], newState.trick[2]) {
			winnerCard = 0 
		}
		if greater(strump, sfollow, newState.trick[1], newState.trick[0], newState.trick[2]) {
			winnerCard = 1 
		} else {
			winnerCard = 2
		}

		// set the scores, depending on who played when
		if newState.turn == 0 { // YOU PLAYED THE 3RD CARD
			newState.setScores(2, winnerCard) 	
		}
		if newState.turn == 1 { // YOU PLAYED THE 2nd CARD
			newState.setScores(1, winnerCard) 	
		}		
		if newState.turn == 2 { // YOU PLAYED THE 1st CARD
			newState.setScores(0, winnerCard) 	
		}		
	}

	// set next player turn
	newState.turn++
	if newState.turn > 2 {
		newState.turn = 0
	}

	var state minimax.State
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
	}
}

