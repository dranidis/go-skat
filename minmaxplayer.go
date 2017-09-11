package main

import (
	"github.com/dranidis/go-skat/minimax"
	"math/rand"
)

var rminmax = rand.New(rand.NewSource(1))


type MinMaxPlayer struct {
	Player
	p1Hand      []Card
	p2Hand      []Card
	skat        []Card
	maxHandSize int
	maxWorlds   int
	// schneiderGoal bool
}

func makeMinMaxPlayer(hand []Card) MinMaxPlayer {
	return MinMaxPlayer{
		Player:      makePlayer(hand),
		p1Hand:      []Card{},
		p2Hand:      []Card{},
		maxHandSize: 5,
		maxWorlds:   12,
		// schneiderGoal: false,
	}
}

func (p *MinMaxPlayer) playerTactic(s *SuitState, c []Card) Card {

	minimax.DEBUG = false

	if len(c) == 1 {
		debugMinmaxLog("..FORCED MOVE.. ")
		return c[0]
	}
	if s.trump == NULL {
		return p.Player.playerTactic(s, c)
	}
	// for the moment do not use minimax in defense
	// if s.declarer.getName() != p.getName(){
	// 	return p.opponentTactic(s, c)
	// }

	worlds := p.dealCards(s)
	debugMinmaxLog("(%s) %d Worlds\n", p.name, len(worlds))

	if len(p.hand) <= p.maxHandSize || len(worlds) < p.maxWorlds {
		cardsFreq := make(map[string]int)
		cardsTotal := make(map[string]float64)
		cards := make(map[string]Card)

		i := 0
		for i = 0; i < len(worlds); i++ {
			// SET world
			p.p1Hand = worlds[i][0]
			p.p2Hand = worlds[i][1]

			debugMinmaxLog("MinMaxPlayer\n")
			card, value := p.minmaxSkat(s, c)
			// card := p.minMaxTactics(s, c)
			v, ok := cardsFreq[card.String()]
			if ok {
				cardsFreq[card.String()] = v + 1
				cardsTotal[card.String()] = cardsTotal[card.String()] + value
				// if cardsFreq[card.String()] > len(worlds)/2 {
				// 	// half of the worlds suggest this card
				// 	debugMinmaxLog("..PRELIMINARY END!\n")
				// 	break
				// }
			} else {
				cardsFreq[card.String()] = 1
				cardsTotal[card.String()] = value
			}
			cards[card.String()] = card
		}

		mostFrequent := 0
		mostFrequentKey := ""
		var bestAvg float64
		if p.name == s.declarer.getName() {
			bestAvg = 0.0
		} else 	{
			bestAvg = 120.0
		}
		var mostFrequentCard Card
		var bestAvgCard Card
		for k, v := range cardsFreq {
			if v > mostFrequent {
				mostFrequent, mostFrequentCard = v, cards[k]
				mostFrequentKey = k
			}
			cardsTotal[k] = cardsTotal[k]/float64(v)
			if p.name == s.declarer.getName() && cardsTotal[k] > bestAvg {
				bestAvg, bestAvgCard = cardsTotal[k], cards[k]
			}
			if p.name != s.declarer.getName() && cardsTotal[k] < bestAvg {
				bestAvg, bestAvgCard = cardsTotal[k], cards[k]
			}

		}
		if mostFrequent > 0 {
			debugMinmaxLog("(%s) In hand %v, examined %d/%d worlds\n", p.name, p.hand, i, len(worlds))
			if p.losingScore(s, cardsTotal[mostFrequentKey]) {
				if p.losingScore(s, bestAvg) {
					debugMinmaxLog("(%s) Losing AVG (BACK TO NORMAL TACTICS)\n", p.name)
					return p.Player.playerTactic(s, c)
				} 
				debugMinmaxLog("(%s)Playing best AVG card: %v (%.0f)\n", p.name, bestAvgCard, bestAvg)
				return bestAvgCard
			} else {
				debugMinmaxLog("(%s)Playing card: %v (at least %d times)\n", p.name, mostFrequentCard, mostFrequent)
				return mostFrequentCard
			}
		}
		// NO CARDS
		debugMinmaxLog("MINMAX Failed.. back no normal tactics")
	}

	debugMinmaxLog("(%s) (NORMAL TACTICS)\n", p.name)
	return p.Player.playerTactic(s, c)
}

func (p *MinMaxPlayer) losingScore(s *SuitState, score float64) bool {
	if p.name == s.declarer.getName() {
		return score < float64(61)
	}
	return score >= float64(61)
}

func checkVoidDecl(s *SuitState, p *MinMaxPlayer, cards []Card, suit string, IsDeclarerP1 bool) []Card {
	if s.declarerVoidSuit[suit] {
		debugMinmaxLog("..Declarer VOID in %s\n", suit)
		suitCards := filter(cards, func(c Card) bool {
			return getSuit(s.trump, c) == suit
		})
		if IsDeclarerP1 {
			p.p2Hand = append(p.p2Hand, suitCards...)
		} else {
			p.p1Hand = append(p.p1Hand, suitCards...)
		}

		cards = remove(cards, suitCards...)
	}
	return cards
}

// FROM THE POINT OF VIEW of A DECLARER
func checkVoidOpp1(s *SuitState, p *MinMaxPlayer, cards []Card, suit string) []Card {
	if s.opp1VoidSuit[suit] {
		debugMinmaxLog("..Opponent 1 is VOID in %s\n", suit)
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
		debugMinmaxLog("..Opponent 2 is VOID in %s\n", suit)
		suitCards := filter(cards, func(c Card) bool {
			return getSuit(s.trump, c) == suit
		})
		p.p1Hand = append(p.p1Hand, suitCards...)
		cards = remove(cards, suitCards...)
	}
	return cards
}

// FROM THE POINT OF VIEW of A DEFENDER (OPPONENT)
// When opp1 is void cards go to p1
func partnerCheckVoidOpp1(s *SuitState, p *MinMaxPlayer, cards []Card, suit string) []Card {
	if s.opp1VoidSuit[suit] {
		debugMinmaxLog("..Opponent 1 is VOID in %s\n", suit)
		suitCards := filter(cards, func(c Card) bool {
			return getSuit(s.trump, c) == suit
		})
		p.p1Hand = append(p.p1Hand, suitCards...)
		cards = remove(cards, suitCards...)
	}
	return cards
}

// When opp2 is void cards go to p2
func partnerCheckVoidOpp2(s *SuitState, p *MinMaxPlayer, cards []Card, suit string) []Card {
	if s.opp2VoidSuit[suit] {
		debugMinmaxLog("..Opponent 2 is VOID in %s\n", suit)
		suitCards := filter(cards, func(c Card) bool {
			return getSuit(s.trump, c) == suit
		})
		p.p2Hand = append(p.p2Hand, suitCards...)
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

	// only declarer knows the skat
	if p.getName() == s.declarer.getName() {
		cards = remove(cards, s.skat...)
		p.skat = []Card{s.skat[0], s.skat[1]}
	}

	debugMinmaxLog("ALL REMAINING CARDS (%d): %v\n", len(cards), cards)

	p.p1Hand = []Card{}
	p.p2Hand = []Card{}

	if p.getName() == s.declarer.getName() {
		for _, suit := range suits {
			cards = checkVoidOpp1(s, p, cards, suit)
			cards = checkVoidOpp2(s, p, cards, suit)
		}
		debugMinmaxLog("REMAINING after void: %d cards: %v %v %v\n", len(cards), cards, p.p1Hand, p.p2Hand)
	} 

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
	// 6 (3,3) => 20

	if p.getName() == s.declarer.getName() {
	// why not for everybody?
	// because defenders DO NOT know the SKAT!
	// if true {
		if len(cards) < 2 || len(p.p1Hand) == max1 || len(p.p2Hand) == max2 {
			p1H, p2H := p.distributeCards(s, cards)
			world := [][]Card{p1H, p2H}
			worlds = append(worlds, world)

			return worlds
		}
		if len(cards) < 4 || len(p.p1Hand) == max1-1 || len(p.p2Hand) == max2-1 {
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
		if len(cards) == 4 {
			// STORE THE hands in order to restore them
			originalP1 := make([]Card, len(p.p1Hand))
			copy(originalP1, p.p1Hand)
			originalP2 := make([]Card, len(p.p2Hand))
			copy(originalP2, p.p2Hand)
			originalCards := make([]Card, len(cards))
			copy(originalCards, cards)

			for i := 0; i < len(originalCards)-1; i++ {
				for j := i + 1; j < len(originalCards); j++ {
					p.p1Hand = originalP1
					p.p2Hand = originalP2
					cards = originalCards

					card1 := cards[i]
					card2 := cards[j]
					p.p1Hand = append(p.p1Hand, card1)
					p.p1Hand = append(p.p1Hand, card2)
					cards = remove(cards, card1, card2)
					// distribute the rest
					p1H, p2H := p.distributeCards(s, cards)
					world := [][]Card{p1H, p2H}
					worlds = append(worlds, world)
				}
			}
			return worlds
		}
		if len(cards) == 5 || len(p.p1Hand) == max1-2 || len(p.p2Hand) == max2-2 {
			var restoreHands func()
			var appendCards func(card1, card2 Card)
			// STORE THE hands in order to restore them
			originalCards := make([]Card, len(cards))
			copy(originalCards, cards)

			if len(p.p1Hand) == max1-2 {
				originalP1 := make([]Card, len(p.p1Hand))
				copy(originalP1, p.p1Hand)
				originalP2 := make([]Card, len(p.p2Hand))
				copy(originalP2, p.p2Hand)
				restoreHands = func() {
					p.p1Hand = originalP1
					p.p2Hand = originalP2
					cards = originalCards
				}
				appendCards = func(card1, card2 Card) {
					p.p1Hand = append(p.p1Hand, card1)
					p.p1Hand = append(p.p1Hand, card2)
				}
			} else {
				originalP1 := make([]Card, len(p.p2Hand))
				copy(originalP1, p.p2Hand)
				originalP2 := make([]Card, len(p.p1Hand))
				copy(originalP2, p.p1Hand)
				restoreHands = func() {
					p.p2Hand = originalP1
					p.p1Hand = originalP2
					cards = originalCards
				}
				appendCards = func(card1, card2 Card) {
					p.p2Hand = append(p.p2Hand, card1)
					p.p2Hand = append(p.p2Hand, card2)
				}
			}

			for i := 0; i < len(originalCards)-1; i++ {
				for j := i + 1; j < len(originalCards); j++ {
					restoreHands()
					card1 := cards[i]
					card2 := cards[j]
					appendCards(card1, card2)
					cards = remove(cards, card1, card2)
					// distribute the rest
					p1H, p2H := p.distributeCards(s, cards)
					world := [][]Card{p1H, p2H}
					worlds = append(worlds, world)
				}
			}
			return worlds
		}

		if len(cards) == 6 || len(p.p1Hand) == max1 - 3 || len(p.p2Hand) == max2 - 3 {
			var restoreHands func()
			var appendCards func(card1, card2, card3 Card)
			// STORE THE hands in order to restore them
			originalCards := make([]Card, len(cards))
			copy(originalCards, cards)

			if len(p.p1Hand) == max1-3 {
				originalP1 := make([]Card, len(p.p1Hand))
				copy(originalP1, p.p1Hand)
				originalP2 := make([]Card, len(p.p2Hand))
				copy(originalP2, p.p2Hand)
				restoreHands = func() {
					p.p1Hand = originalP1
					p.p2Hand = originalP2
					cards = originalCards
				}
				appendCards = func(card1, card2, card3 Card) {
					p.p1Hand = append(p.p1Hand, card1)
					p.p1Hand = append(p.p1Hand, card2)
					p.p1Hand = append(p.p1Hand, card3)
				}
			} else {
				originalP1 := make([]Card, len(p.p2Hand))
				copy(originalP1, p.p2Hand)
				originalP2 := make([]Card, len(p.p1Hand))
				copy(originalP2, p.p1Hand)
				restoreHands = func() {
					p.p2Hand = originalP1
					p.p1Hand = originalP2
					cards = originalCards
				}
				appendCards = func(card1, card2, card3 Card) {
					p.p2Hand = append(p.p2Hand, card1)
					p.p2Hand = append(p.p2Hand, card2)
					p.p1Hand = append(p.p1Hand, card3)
				}
			}

			for i := 0; i < len(originalCards)-2; i++ {
				for j := i + 1; j < len(originalCards)-1; j++ {
					for k := j + 1; k < len(originalCards); k++ {
						restoreHands()
						card1 := cards[i]
						card2 := cards[j]
						card3 := cards[k]
						appendCards(card1, card2, card3)
						cards = remove(cards, card1, card2, card3)
						// distribute the rest
						p1H, p2H := p.distributeCards(s, cards)
						world := [][]Card{p1H, p2H}
						worlds = append(worlds, world)
					}
				}
			}
			return worlds
		}

	}

	copycards := make([]Card, len(cards))

	for i := 0; i < p.maxWorlds; i++ {
		// shuffle the rest
		copycards = ShuffleR(rminmax, cards)

		if p.getName() != s.declarer.getName() { // remove two random cards for the skat
			// TODO
			// what do we usually discard in skat?
			// Definitely no trumps
			// Most probably no Aces
			p.skat = make([]Card, 2)
			var card Card
			for i := 0; i < 2; i++ {
				for {
					card = copycards[rminmax.Intn(len(copycards))]
					if getSuit(s.trump, card) == s.trump {
						continue
					}
					if card.Rank == "A" {
						continue
					}
					copycards = remove(copycards, card)
					break
				}
				p.skat[i] = card
			}
			debugMinmaxLog("Removing (SKAT): %v \n", p.skat)
			debugMinmaxLog("REMAINING after SKAT REMOVE: %d cards: %v %v %v\n", len(copycards), copycards, p.p1Hand, p.p2Hand)
		}

		// check VOIDS for defenders. FOr declarers already done before
		decl, _, _ := p.getDeclarerNrAndPlayers(s) // (int, PlayerI, PlayerI)
		IsDeclarerP1 := decl == 1
		if IsDeclarerP1 {
			debugMinmaxLog("Declarer is P1\n") // you are s.opp2
		} else {
			debugMinmaxLog("Declarer is P2\n") // you are s.opp1
		} 

		if p.getName() != s.declarer.getName() {
			p.p1Hand = []Card{}
			p.p2Hand = []Card{}
			for _, suit := range suits {
				if IsDeclarerP1 {
					copycards = checkVoidDecl(s, p, copycards, suit, IsDeclarerP1)
					// debugMinmaxLog("..copycards: %v\n", copycards)
				} else {
					copycards = checkVoidDecl(s, p, copycards, suit, IsDeclarerP1)
					// debugMinmaxLog("..copycards: %v\n", copycards)
				}	
				if p.getName() == s.opp1.getName() {
					copycards = partnerCheckVoidOpp2(s, p, copycards, suit)
					// debugMinmaxLog("..copycards after partner void: %v\n", copycards)
				} else {
					copycards = partnerCheckVoidOpp1(s, p, copycards, suit)
					// debugMinmaxLog("..copycards after partner void: %v\n", copycards)
				}
			}
			debugMinmaxLog("REMAINING after void: %d cards: %v %v %v\n", len(copycards), copycards, p.p1Hand, p.p2Hand)
		}
		if len(p.p1Hand) > max1 || len(p.p2Hand) > max2 {
			debugMinmaxLog("IMPOSSIBLE!")
			continue
		}
		p1H, p2H := p.distributeCards(s, copycards)
		if len(p1H) > max1 || len(p2H) > max2 {
			debugMinmaxLog("IMPOSSIBLE!")
			continue
		}		
		debugMinmaxLog("DISTRIBUTION: %v %v\n", p1H, p2H)
		world := [][]Card{p1H, p2H}
		worlds = append(worlds, world)
	}
	return worlds
}

func (p *MinMaxPlayer) distributeCards(s *SuitState, cards []Card) ([]Card, []Card) {
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
	// debugMinmaxLog("cards to distribute: %v\n", cards)
	// debugMinmaxLog("hand1: %v\n", hand1)
	// debugMinmaxLog("hand2: %v\n", hand2)

	// p1 has more cards than p2 if you are in the middle: p2 you p1
	// [] => you p1 p2
	// [x] => p2 you p1
	// [x,x] => p1 p2 you

	for i := len(hand1); i < handSize-middle; i++ {
		// debugMinmaxLog("hand1: %v, i: %d, handSize: %d, middle: %d, nextCard: %d\n", hand1, i, handSize, middle, nextCard)
		hand1 = append(hand1, cards[nextCard])
		nextCard++
	}
	// debugMinmaxLog("completed hand1: %v\n", hand1)

	for i := len(hand2); i < handSize-leader; i++ {
		// debugMinmaxLog("hand2: %v, i: %d, handSize: %d, leader: %d, nextCard: %d\n", hand2, i, handSize, leader, nextCard)
		hand2 = append(hand2, cards[nextCard])
		nextCard++
	}
	// debugMinmaxLog("completed hand2: %v\n", hand2)
	return hand1, hand2
}

func (p *MinMaxPlayer) minmaxSkat(s *SuitState, c []Card) (Card, float64) {
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
	if s.declarer.getScore()+sum(s.trick) > 60 || s.opp1.getScore()+s.opp2.getScore() > 59 {
		debugMinmaxLog("MIN_MAX: schneiderGoal\n")
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

	debugMinmaxLog("MINMAX: cards %s: %v, %s: %v, SKAT:%v\n", player1.getName(), p.p1Hand, player2.getName(), p.p2Hand, p.skat)
	debugMinmaxLog("Decl: %d\n", decl)

	strick := make([]Card, len(s.trick))
	copy(strick, s.trick)

	dist := make([][]Card, 3)
	dist[0] = p.hand
	dist[1] = p.p1Hand
	dist[2] = p.p2Hand

	// debugMinmaxLog("CARDS: %v", dist)
	// debugMinmaxLog("CARDS[0]: %v", dist[0])
	// debugMinmaxLog("CARDS[1]: %v", dist[1])
	// debugMinmaxLog("CARDS[2]: %v", dist[2])
	declScore := s.declarer.getScore()
	if p.name == s.declarer.getName() {
		declScore += sum(s.skat)
	}

	skatState := SkatState{
		s.trump,
		dist,
		strick,
		decl,
		0,
		declScore,
		s.opp1.getScore() + s.opp2.getScore(),
		schneiderGoal,
	}

	// if skatState.IsTerminal() {
	// 	debugMinmaxLog("##### !!!!!!! ~~~~~~ TERMINAL state.. back to player tactics!")
	// 	return p.Player.playerTactic(s, c), 
	// }
	// for _, action := range skatState.FindLegals() {
	// 	ma := action.(SkatAction)
	// 	debugMinmaxLog("LEGAL action: %v\n", ma)
	// }

	a, value := minimax.Minimax(&skatState)
	ma := a.(SkatAction)

	debugMinmaxLog("Suggesting card: %v with value %.0f\n", ma.card, value)

	return ma.card, value
}

// TODO:
// refactor code above by calling this function
func (p *MinMaxPlayer) getDeclarerNrAndPlayers(s *SuitState) (int, PlayerI, PlayerI) {
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
	return decl, player1, player2
}


