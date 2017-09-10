package main

import (
	"github.com/dranidis/go-skat/minimax"
	"math/rand"
)

var rminmax = rand.New(rand.NewSource(1))


type MinMaxPlayer struct {
	Player
	p1Hand []Card
	p2Hand []Card
	maxHandSize int
	// schneiderGoal bool
}

func makeMinMaxPlayer(hand []Card) MinMaxPlayer {
	return MinMaxPlayer {
		Player:     makePlayer(hand),
		p1Hand: []Card{},
		p2Hand: []Card{},
		maxHandSize: 4,
		// schneiderGoal: false,
	}
}

func (p *MinMaxPlayer) playerTactic(s *SuitState, c []Card) Card {

	minimax.DEBUG = false
	
	if len(c) == 1 {
		debugTacticsLog("..FORCED MOVE.. ")
		return c[0]
	}

	worlds := p.dealCards(s)
	debugTacticsLog("MINMAX: %d Worlds\n", len(worlds))

	if len(p.hand) <= p.maxHandSize || len(worlds) < 10 {

		// worlds := p.dealCards(s)
		// if p.getName() == s.declarer.getName() {
		// 	if p.getScore() > 60 {
		// 		p.schneiderGoal = true
		// 	} else {
		// 		p.schneiderGoal = false
		// 	}
		// } else {
		// 		partnerFunc := func(p PlayerI) PlayerI {
		// 			if s.opp1.getName() == p.getName() {
		// 				return s.opp2
		// 			}
		// 			return s.opp1
		// 		}
		// 	if p.getScore() + partnerFunc(p).getScore() > 59 {				
		// 		p.schneiderGoal = true
		// 	} else {
		// 		p.schneiderGoal = false
		// 	}
		// }

		// minimax.DEBUG = true
		cardsFreq := make(map[string]int)
		cards := make(map[string]Card)


		for i := 0; i < len(worlds); i++ {
			// SET world
			p.p1Hand = worlds[i][0]
			p.p2Hand = worlds[i][1]

			debugTacticsLog("MinMaxPlayer\n") 
			card := p.minmaxSkat(s, c)
			// card := p.minMaxTactics(s, c)
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

// func (p MinMaxPlayer) minMaxTactics(s *SuitState, c []Card) Card {
// 	// p.dealCards(s)
// 	card := p.minmaxSkat(s)

// 	return card
// }

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

func (p *MinMaxPlayer) dealCards(s *SuitState) [][][]Card {

	worlds := [][][]Card{}

	cards := makeDeck()
	cards = remove(cards, s.cardsPlayed...)
	cards = remove(cards, s.trick...)
	cards = remove(cards, p.hand...)

	if p.getName() == s.declarer.getName() {
		cards = remove(cards, s.skat...)
	}

	debugTacticsLog("ALL CARDS: %v\n", cards)

	p.p1Hand = []Card{}
	p.p2Hand = []Card{}

	if p.getName() == s.declarer.getName() {
		for _, suit := range suits {
			cards = checkVoidOpp1(s, p, cards, suit)	
			cards = checkVoidOpp2(s, p, cards, suit)	
		}
	}
	debugTacticsLog("REMAINING after void: %d cards: %v %v %v\n", len(cards), cards, p.p1Hand, p.p2Hand)

	max1 := len(p.hand)
	max2 := len(p.hand)
	if len(s.trick) == 1 {
		max2--
	}
	if len(s.trick) == 2 {
		max1--
		max2--
	}	
	// cards => ways to distribute
	// 0 (0,0) => 1
	// 1 (0,1) => 1 
	// 2 (1,1) => 2
	// 3 (1,2) => 3
	// 4 (2,2) => 6
	// 5 (2,3) => 10
	// 6 (3,3) => 

	if p.getName() == s.declarer.getName() {
		if len(cards) < 2 || len(p.p1Hand) == max1 || len(p.p2Hand) == max1 {
			p1H, p2H := p.distributeCards(s, cards)
			world := [][]Card{p1H, p2H}
			worlds = append(worlds, world)

			return worlds
		}	
		if len(cards) < 4  || len(p.p1Hand) == max1 -1 || len(p.p2Hand) == max1 -1  {
			for i := 0; i < len(cards); i++ {
				p1H, p2H := p.distributeCards(s, cards)
				world := [][]Card{p1H, p2H}
				worlds = append(worlds, world)	
				card := cards[0]	
				cards = remove(cards, card)
				cards = append(cards, card)				
			}
			return worlds
		}

	}

	copycards := make([]Card, len(cards))

	for i := 0; i < 10; i++ {
		// shuffle the rest
		copycards = ShuffleR(rminmax, cards)

		if p.getName() != s.declarer.getName() { // remove two random cards for the skat
			card1 := copycards[0]
			card2 := copycards[1]
			copycards = remove(copycards, card1, card2)
			debugTacticsLog("REMAINING after SKAT REMOVE: %d cards: %v %v %v\n", len(copycards), copycards, p.p1Hand, p.p2Hand)
		}

		p1H, p2H := p.distributeCards(s, copycards)
		world := [][]Card{p1H, p2H}
		worlds = append(worlds, world)
	}
	return worlds
}

func (p* MinMaxPlayer) distributeCards(s *SuitState, cards []Card) ([]Card, []Card) {
	handSize := len(p.hand)

	hand1 := make([]Card, len(p.p1Hand))
	copy(hand1, p.p1Hand)
	hand2 := make([]Card, len(p.p2Hand))
	copy(hand2, p.p2Hand)

	leader := 0
	middle := 0
	if len(s.trick) == 1 {
		leader = 1
	}
	if len(s.trick) == 2 {
		middle = 1
		leader = 1
	}

	nextCard := 0
	// debugTacticsLog("cards to distribute: %v\n", cards)

	for i := len(hand1); i < handSize - middle; i++ {
		debugTacticsLog("hand1: %v, i: %d, handSize: %d, middle: %d, nextCard: %d\n", hand1, i, handSize, middle, nextCard)
		hand1 = append(hand1, cards[nextCard])
		nextCard++
	}
	debugTacticsLog("completed hand1: %v\n", hand1)

	for i := len(hand2); i < handSize - leader; i++ {
		debugTacticsLog("hand2: %v, i: %d, handSize: %d, leader: %d, nextCard: %d\n", hand2, i, handSize, leader, nextCard)
		hand2 = append(hand2, cards[nextCard])
		nextCard++
	}	
	debugTacticsLog("completed hand2: %v\n", hand2)
	return hand1, hand2
}

func (p *MinMaxPlayer) minmaxSkat(s *SuitState, c []Card) Card {
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
	schneiderGoal := true

	// TODO:
	// from the point of view of the defender the skat is unknown
	// How can I add it the the evaluation function???
	if s.declarer.getScore() + sum(s.trick) > 60 || s.opp1.getScore() + s.opp2.getScore() > 59 {
		debugTacticsLog("MIN_MAX: schneiderGoal\n")
		schneiderGoal = true
		// minimax.DEBUG = true
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
		s.declarer.getScore() + sum(s.skat), 
		s.opp1.getScore() + s.opp2.getScore(),
		schneiderGoal,
	}
	debugTacticsLog("Skatstate Cards: %v\n",skatState.playerHand)

	if skatState.IsTerminal() {
		debugTacticsLog("##### !!!!!!! ~~~~~~ TERMINAL state.. back to player tactics!")
		return p.Player.playerTactic(s, c)
	}
	// for _, action := range skatState.FindLegals() {
	// 	ma := action.(SkatAction)
	// 	debugTacticsLog("LEGAL action: %v\n", ma)
	// }

	a := minimax.Minimax(&skatState)
	ma := a.(SkatAction)

	debugTacticsLog("In hand %v, MINIMAX suggesting card: %v\n", p.hand, ma.card)

	return ma.card
}


// action {{CARO K}} 
// {CARO [[{CARO K} {CARO 9} {CARO 7}] [{HEART 10} {CLUBS D} {SPADE K} {CLUBS 8}] [{CARO A} {CLUBS A} {CLUBS 9}]] [] 0 0 62 13 true}

type SkatState struct {
	trump string
	playerHand [][]Card  // YOU, 1, 2
	trick []Card
	declarer int
	turn int // 0 you, 1 player 1, 2 player2
	declScore int 
	oppScore int
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

func (m SkatState) IsOpponentTurn() bool {
	if m.declarer == 0 && m.turn == 0 {
		return false
	}
	if m.declarer == 0 && m.turn != 0 {
		return true
	}
	if m.declarer == m.turn {
		return true
	}
	return false
}

func (m *SkatState) IsTerminal() bool {
	if ! m.schneiderGoal {
		if m.declScore > 60 {
			return true
		}
		if m.oppScore > 59 {
			return true
		}
	} 
	return len(m.playerHand[0]) + len(m.playerHand[1]) + len(m.playerHand[2]) == 0
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
	if m.declScore > winsScore - 1 {
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
		m.schneiderGoal,
	}
}

