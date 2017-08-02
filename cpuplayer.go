package main

import (
	_ "fmt"
)

type Player struct {
	PlayerData
	firstCardPlay bool
	risky         bool
	grand bool
}

func makePlayer(hand []Card) Player {
	return Player{
		PlayerData:    makePlayerData(hand),
		firstCardPlay: false,
		risky:         false,
		grand: false,
	}
}

func winnerCards(s *SuitState, cs []Card) []Card {
	return filter(cs, func(c Card) bool {
		wins := true
		for _, t := range s.trick {
			if s.greater(t, c) {
				wins = false
			}
		}
		return wins
	})
}

func highestValueWinnerORlowestValueLoser(s *SuitState, c []Card) Card {
	winners := winnerCards(s, c)
	debugTacticsLog("Winners: %v\n", winners)
	if s.trump == s.follow {
		trumpWinnerRanks := []string{"A", "10", "K", "D", "9", "8", "7", "J"}
		winners = sortRankSpecial(winners, trumpWinnerRanks)
	} else {
		winners = sortValue(winners)
	}
	if len(winners) > 0 {
		return winners[0]
	}

	if s.trump == s.follow {
		trumpRanks := []string{"A", "10", "J", "K", "D", "9", "8", "7"}
		cards := sortRankSpecial(c, trumpRanks)
		return cards[len(cards)-1]
	}

	sortedValue := sortValue(c)
	cards := filter(sortedValue, func(card Card) bool {
		return card.suit != s.trump
	})
	debugTacticsLog("LOSING Cards (last)%v\n", cards)
	if len(cards) > 0 {
		return cards[len(cards)-1]
	}
	// LAST?
	return sortedValue[len(sortedValue)-1]
}

func (p Player) otherPlayersHaveJs(s *SuitState) bool {
	for _, suit := range suits {
		card := Card{suit, "J"}
		if in(s.trumpsInGame, card) && !in(p.getHand(), card) {
			return true
		}
	}
	return false
}

func (p Player) otherPlayersTrumps(s *SuitState) []Card {
	return filter(makeDeck(), func(c Card) bool {
		if getSuit(s.trump, c) != s.trump {
			return false
		}
		if !in(s.trumpsInGame, c) {
			return false
		}
		if in(p.getHand(), c) {
			return false
		}
		return true
	})
}

func firstCardTactic(c []Card) Card {
	return c[0]
}

func noHigherCard(s *SuitState, viewSkat bool, c Card) bool {
	allCards := makeSuitDeck(c.suit)
	allCardsPlayed := []Card{}
	if viewSkat {
		allCardsPlayed = append(allCardsPlayed, s.skat...)
	}
	allCardsPlayed = append(allCardsPlayed, s.cardsPlayed...)
	for _, cardPlayed := range allCardsPlayed {
		allCards = remove(allCards, cardPlayed)
	}
	allCards = filter(allCards, func(card Card) bool {
		return card.rank != "J"
	})
	debugTacticsLog("Cards of suit %s still in play: %v", c.suit, allCards)
	for _, card := range allCards {
		if s.greater(card, c) {
			return false
		}
	}
	return true
}

func nextLowestCardsStillInPlay(s *SuitState, w Card, followCards []Card) bool {
	next := nextCard(w)
	debugTacticsLog("Next of %v is %v...", w, next)
	// ONLY the declarer knows that. Use a flag if opp uses it.
	if in(s.skat, next) {
		return false
	}
	if in(followCards, next) {
		return false
	}
	if in(s.cardsPlayed, next) {
		return false
	}
	return true
}


func (p Player) declarerTactic(s *SuitState, c []Card) Card {
	debugTacticsLog("DECLARER ")
	if len(s.trick) == 0 {
		debugTacticsLog("FOREHAND ")
		// count your own trumps and other players trump
		// if you have less you should not play trumps immediately
		if len(p.otherPlayersTrumps(s)) > 0 {
			debugTacticsLog("other TRUMPS in game: %v", p.otherPlayersTrumps(s))

			if p.otherPlayersHaveJs(s) {
				validCards := make([]Card, len(c))
				copy(validCards, c)
				first := firstCardTactic(validCards)
				for len(validCards) > 1 && (first.equals(Card{s.trump, "A"}) || first.equals(Card{s.trump, "10"})) {
					validCards = remove(validCards, first)
					first = firstCardTactic(validCards)
				}
				return first
			} else {
				return firstCardTactic(c)
			}
		}
		// Declarer still has trumps
		if len(s.trumpsInGame) > 0 {
			// Check for A-K-x suits
			AKXsuits := filterSuit(suits, func(suit string) bool {
				return s.trump != suit && isAKX(suit, c)
			})
			// should also check that 10 is still in play
			for _, AKXsuit := range AKXsuits {
				debugTacticsLog("A-K-X: play X in hand %v\n", p.hand)
				if in(s.cardsPlayed, Card{AKXsuit, "10"}) {
					debugTacticsLog("%s 10 already played\n", AKXsuit)
					continue
				}
				cards := filter(c, func(c Card) bool {
					return c.suit == AKXsuit
				})
				return cards[len(cards)-1]
			}

			// check K-x or 10-x cases where higher cards stil in play
			for _, suit := range suits {
				if suit == s.trump {
					continue
				}
				cs := sortRank(nonTrumpCards(suit, c))
				if len(cs) == 0 {
					continue
				}
				if noHigherCard(s, true, cs[0]) {
					debugTacticsLog(" Sure winner card %v", cs[0])
					return cs[0]
				}
				if len(cs) > 1 {
					if cardValue(cs[len(cs)-1]) == 0 {
						debugTacticsLog(" Play loser card %v", cs[len(cs)-1])
						return cs[len(cs)-1]
					}
				}
			}
			//REPETITION: see below
			sortedValue := filter(sortValue(c), func(card Card) bool {
				return card.suit != s.trump && card.rank != "J"
			})
			if len(sortedValue) > 0 {
				return sortedValue[len(sortedValue)-1]
			}
		}

		debugTacticsLog(" -TRUMPS exhausted: Hand: %v ", p.getHand())
		return HighestLong(s.trump, c)
	}
	if len(s.trick) == 2 {
		debugTacticsLog("BACKHAND ")
		followCards := filter(c, func(card Card) bool {
				return card.suit != s.trump && card.rank != "J"
			})
		if len(followCards) > 0 {
			debugTacticsLog("Following normal suit...")
			winners := sortRank(winnerCards(s, followCards))
			debugTacticsLog("winners %v...", winners)
			for _, w := range winners {
				if w.rank == "D" || w.rank == "K" {
					return w
				}
				if nextLowestCardsStillInPlay(s, w, followCards) {
					debugTacticsLog("Next lower still in play...")
					continue
				}
				debugTacticsLog("Returning  %v...", w)
				return w
			}
			// if len(winners) > 0 {
			// 	debugTacticsLog("Returning last of winners %v...", winners[len(winners)-1])
			// 	return winners[len(winners)-1]
			// }
		} else {
			debugTacticsLog("TRUMP OR No cards of suit played...")
		}
		if sum(s.trick) == 0 {
			debugTacticsLog("ZERO valued trick. DO not trump!...")
			// do not trump
			nonTrumps := filter(sortValue(c), func(card Card) bool {
				return card.suit != s.trump && card.rank != "J"
			})
			debugTacticsLog("Non-trumps: %v...", nonTrumps)

			sortedValue := filter(sortValue(c), func(card Card) bool {
				return card.suit != s.trump && card.rank != "J"
			})
			if len(sortedValue) > 0 {
				return highestValueWinnerORlowestValueLoser(s, sortedValue)
			}
		}
// EURO 259.46     -256.24 -3.22
// WON  27001      18710   18886
// LOST  7706       3723    3645
// bidp    44         28      28
// pcw     78         83      84
// pcwd    16         20      20
// AVG  16.6, passed 20329, won 64597, lost 15074 / 100000 game

		// don't throw your A if not 10 in the trick and still in game
		debugTacticsLog("CHECKING the A... in trick %v, valid %v, played %v..", s.trick, c, s.cardsPlayed)
		if s.follow != s.trump && in(c, Card{s.follow, "A"}) && !in(s.trick, Card{s.follow, "10"}) && !in(append(s.cardsPlayed, s.skat...),  Card{s.follow, "10"}) {
			debugTacticsLog("Keeping the A... in trick %v..", s.trick)
			sortedRank := filter(sortRank(c), func(card Card) bool {
				return card.rank != "A" && card.rank != "J" 
			})
			if len(sortedRank) > 0 {
				debugTacticsLog("Valid: %v, Non A-cards of suit %v\n", c, sortedRank)
				return highestValueWinnerORlowestValueLoser(s, sortedRank)
			}			
		}
	}
	// TODO:exhausted
	// in middlehand, if leader leads with an  suit don't take
	// it with the most valuable trump. It might be taken....

	return highestValueWinnerORlowestValueLoser(s, c)
}

func isAKX(trump string, cs []Card) bool {
	cards := filter(cs, func(c Card) bool {
		return c.suit == trump
	})

	if !in(cards, Card{trump, "A"}, Card{trump, "K"}) {
		//	fmt.Printf("Not A-K:  %v\n", cs)
		return false
	}
	if in(cards, Card{trump, "10"}) {
		// fmt.Printf("Yes 10\n")
		return false
	}
	if len(cards) < 3 {
		// fmt.Printf("less than 3\n")
		return false
	}
	return true
}

func HighestShort(trump string, c []Card) Card {
	s := ShortestNonTrumpSuit(trump, c)
	debugTacticsLog("ShortestNonTrumpSuit %v\n", s)
	cards := sortRank(nonTrumpCards(s, c))
	debugTacticsLog("%v", cards)
	if len(cards) > 0 {
		return cards[0]
	}
	// last card?
	debugTacticsLog("... DEBUG ... VALID: %v no HighestShort. Returning: %v\n", c, c[0])
	return c[0]
}

func HighestLong(trump string, c []Card) Card {
	s := LongestNonTrumpSuit(trump, c)
	debugTacticsLog("LongestNonTrumpSuit %v\n", s)
	cards := sortRank(nonTrumpCards(s, c))
	debugTacticsLog("%v", cards)
	if len(cards) > 0 {
		return cards[0]
	}
	// last card?
	debugTacticsLog("... DEBUG ... VALID: %v no highest long. Returning: %v\n", c, c[0])
	return c[0]
}

func (p * Player) FindPreviousSuit(s *SuitState) string {
	partnerFunc := func(p PlayerI) PlayerI {
		if s.opp1 == p {
			return s.opp2
		}
		return s.opp1
	}
	if p.getPreviousSuit() != "" {
		return p.getPreviousSuit()
	} else if partnerFunc(p).getPreviousSuit() != "" {
		return partnerFunc(p).getPreviousSuit()
	}
	return ""
}

func (p *Player) opponentTactic(s *SuitState, c []Card) Card {
	// OPPONENTS TACTIC
	if len(s.trick) == 0 {
		debugTacticsLog("OPP FOREHAND\n")
		// if you have a card with suit played in a previous trick
		// started from you or your partner continue with it
		prevSuit := p.FindPreviousSuit(s)
		var prevSuitCards []Card
		if prevSuit != "" {
			prevSuitCards = filter(c, func(c Card) bool {
				return c.suit == prevSuit && c.rank != "J"
			})		
		} 
		if len(prevSuitCards) > 0 {
			debugTacticsLog("Following previous suit...")
			// TODO:
			// should I play the highest even if in the previous trick
			// declarer has taken with trump?
			// He will take it again.
			return prevSuitCards[0]
		}
		debugTacticsLog("No cards in previous suit '%v' ...", prevSuit)
		var card Card
		if s.opp2 == p {
			debugTacticsLog("Declarer at MIDDLEHAND, playing LONG...")
			card = HighestLong(s.trump, c)
		} else {
			debugTacticsLog("Declarer at BACKHAND, playing SHORT...")
			card = HighestShort(s.trump, c)
		}
		suit := getSuit(s.trump, card)
		if suit != s.trump {
			p.setPreviousSuit(suit)
		}
		return card
	}

	if len(s.trick) == 1 {
		debugTacticsLog("OPP MIDDLEHAND\n")
		if s.leader == s.declarer {
			debugTacticsLog("Declarer leads...")
			// if declarer leads a low trump, and there are still HIGHER trumps
			// smear the trick with a high value
			if getSuit(s.trump, s.trick[0]) == s.trump && len(winnerCards(s, c)) == 0 {
				if len(filter(p.otherPlayersTrumps(s), func(c Card) bool {
					return s.greater(c, s.trick[0])
				})) > 0 {
					debugTacticsLog("TRUMP...There are higher trumps...SMEAR...")
					return sortValue(c)[0]
				}
			}
			return highestValueWinnerORlowestValueLoser(s, c)
		} else {
			debugTacticsLog("Teammate leads...")
			sortedValueNoTrumps := filter(sortValue(c), func(card Card) bool {
				return card.suit != s.trump && card.rank != "J"
			})
			debugTacticsLog("SORT-VALUE no trumps %v\n", sortedValueNoTrumps)

			cardsSuit := filter(c, func(card Card) bool {
				return s.trick[0].suit == card.suit && card.rank != "J"
			})
			if len(cardsSuit) > 0 {
				cards := sortValue(cardsSuit)
				debugTacticsLog("FOLLOW ")
				if noHigherCard(s, false, cards[0]) {
					return cards[0]
				}
				debugTacticsLog("HIGH CARD still out\n")
				if cardValue(s.trick[0]) > 0 {
					return cards[len(cards)-1]
				}
				debugTacticsLog(" increase zero value trick\n")
				for i := len(cards) - 1; i >= 0; i-- {
					if cardValue(cards[i]) > 0 {
						return cards[i]
					}
				}
				return cards[len(cards)-1]
			}
			trickFollowCards := filter(makeDeck(), func(card Card) bool {
				if card.rank == "J" {
					return false
				}
				if card.suit != s.trick[0].suit {
					return false
				}
				if card.rank == s.trick[0].rank {
					return false
				}
				return true
			})
			// else
			debugTacticsLog("PLAYED from suite %v", filter(s.cardsPlayed, func(card Card) bool {
				return card.suit == s.trick[0].suit && card.rank != "J"
			}))
			if in(s.cardsPlayed, trickFollowCards...) && sum(s.trick) == 0 {
				debugTacticsLog("All suit %s played and zero trick. Increase trick value\n", s.trick[0].suit)
				for i := len(sortedValueNoTrumps) - 1; i >= 0; i-- {
					if cardValue(sortedValueNoTrumps[i]) > 0 {
						return sortedValueNoTrumps[i]
					}
				}
			}
			debugTacticsLog("PLAY lowest value card\n")
			if len(sortedValueNoTrumps) > 0 {
				return sortedValueNoTrumps[len(sortedValueNoTrumps)-1]
			}
			// LASt card?
			return sortValue(c)[0]
		}
	}

	if len(s.trick) == 2 {
		debugTacticsLog("OPP BACKHAND\n")
		if s.leader == s.declarer {
			debugTacticsLog(" -- declarer leads --\n")
			if s.greater(s.trick[0], s.trick[1]) {
				return highestValueWinnerORlowestValueLoser(s, c)
			}
			return sortValue(c)[0]
		}
		debugTacticsLog(" -- teammate leads --\n")
		if s.greater(s.trick[0], s.trick[1]) {
			debugTacticsLog(" largest non-trump")
			sortedValue := sortValue(c)
			noTrumps := filter(sortedValue, func(card Card) bool {
				return card.suit != s.trump && card.rank != "J"
			})
			if len(noTrumps) > 0 {
				return noTrumps[0]
			}
			return sortedValue[0]
		}
		return highestValueWinnerORlowestValueLoser(s, c)
	}
	return highestValueWinnerORlowestValueLoser(s, c)
}

func (p *Player) playerTactic(s *SuitState, c []Card) Card {
	if p.firstCardPlay {
		debugTacticsLog("(%s) FIRST CARD PLAY\n", p.name)
		return c[0]
	}
	if s.declarer == p {
		return p.declarerTactic(s, c)
	}
	return p.opponentTactic(s, c)
}

/*
avg: -75 -- -178 with random play
*/
// func (p *Player) accepts(bidIndex int) bool {
// 	if r.Intn(10) > 1 {
// 		return true
// 	}
// 	return false
// }

/*
avg: -25 with random play, 25 with random play and good discard
*/
func (p *Player) accepts(bidIndex int) bool {
	return bids[bidIndex] <= p.highestBid
}

//
// Der US-Amerikaner J.P. Wergin hat in seinem Buch "Wergin on Skat and Sheepshead"
// (McFarland, Wisconsin, 1975) versucht, dazu einen objektiven Berechnungsmodus zu
// finden.
func (p *Player) handEstimation() int {
	kreuzB := in(p.getHand(), Card{CLUBS, "J"})
	pikB := in(p.getHand(), Card{SPADE, "J"})
	herzB := in(p.getHand(), Card{HEART, "J"})
	karoB := in(p.getHand(), Card{CARO, "J"})

	wert := 0
	// Kreuz-B allein
	if kreuzB && !pikB && !herzB && !karoB {
		wert += 10
	}
	// Jeder andere einzelne Bube
	if !kreuzB && pikB && !herzB && !karoB {
		wert += 5
	}
	if !kreuzB && !pikB && herzB && !karoB {
		wert += 5
	}
	if !kreuzB && !pikB && !herzB && karoB {
		wert += 5
	}
	// Kreuz-B und Pik-B
	if kreuzB && pikB && !herzB && !karoB {
		wert += 25
	}
	// Jede andere 2-Buben-Kombi
	if kreuzB && !pikB && herzB && !karoB {
		wert += 20
	}
	if kreuzB && !pikB && !herzB && karoB {
		wert += 20
	}
	if !kreuzB && pikB && herzB && !karoB {
		wert += 20
	}
	if !kreuzB && pikB && !herzB && karoB {
		wert += 20
	}
	if !kreuzB && !pikB && herzB && karoB {
		wert += 20
	}
	// Kreuz-B, Pik-B, Herz-B
	if kreuzB && pikB && herzB && !karoB {
		wert += 40
	}
	// Kreuz-B, Pik-B, Karo-B
	if kreuzB && pikB && !herzB && karoB {
		wert += 37
	}
	// Kreuz-B, Herz-B, Karo-B 3
	if kreuzB && !pikB && herzB && karoB {
		wert += 35
	}
	// Pik-B, Herz-B, Karo-B
	if !kreuzB && pikB && herzB && karoB {
		wert += 35
	}
	// 4 Buben
	if kreuzB && pikB && herzB && karoB {
		wert += 50
	}

	otherCardsEstimation := func(suit string) int {
		a := in(p.getHand(), Card{suit, "A"})
		t := in(p.getHand(), Card{suit, "10"})
		k := in(p.getHand(), Card{suit, "K"})
		d := in(p.getHand(), Card{suit, "D"})
		n := in(p.getHand(), Card{suit, "9"})

		if a && t && k {
			return 25
		}
		if a && t {
			return 20
		}
		if a && k && d {
			return 15
		}
		if a && k && n {
			return 12
		}
		if a {
			return 10
		}
		return 0
	}
	wert += otherCardsEstimation(CLUBS)
	wert += otherCardsEstimation(SPADE)
	wert += otherCardsEstimation(HEART)
	wert += otherCardsEstimation(CARO)

	return wert
}

func findBlank(cards []Card, suit string) Card {
	cc := len(nonTrumpCards(suit, cards))
	if cc == 1 {
		var card Card
		for _, c := range cards {
			if c.rank == "J" {
				continue
			}
			if c.suit == suit {
				card = c
				break
			}
		}
		if card.rank != "A" {
			return card
		}
	}
	return Card{"", ""}
}

// returns blank cards in order of rank value
func findBlankCards(cards []Card) []Card {
	blankCards := []Card{}
	for _, s := range suits {
		card := findBlank(cards, s)
		if card.rank != "" {
			blankCards = append(blankCards, card)
		}
	}
	blankCards = sortRank(blankCards)
	return blankCards
}

func nonA10cards(cs []Card) []Card {
	suitf := func(suit string, cs []Card) []Card {
		cards := filter(cs, func(c Card) bool {
			return c.suit == suit && c.rank != "J"
		})
		if in(cards, Card{suit, "A"}) && !in(cards, Card{suit, "10"}) {
			cards := filter(cards, func(c Card) bool {
				return c.rank != "A"
			})
			return cards
		}
		return []Card{}
	}

	cards := []Card{}
	cards = append(cards, suitf(CLUBS, cs)...)
	cards = append(cards, suitf(SPADE, cs)...)
	cards = append(cards, suitf(HEART, cs)...)
	cards = append(cards, suitf(CARO, cs)...)
	return cards
}

func (p *Player) canPlayGrand() bool {
	return false
}

func (p *Player) canWin() bool {
	assOtherThan := func(suit string) int {
		asses := 0
		for _, s := range suits {
			if s == suit {
				continue
			}
			if in(p.getHand(), Card{s, "A"}) {
				debugTacticsLog("(A %s) ", s)
				asses++
				if in(p.getHand(), Card{s, "10"}) {
					debugTacticsLog("(10 %s) ", s)
					asses++
					if in(p.getHand(), Card{s, "K"}) {
						debugTacticsLog("(K %s) ", s)
						asses++
					}
				}
			}
		}
		return asses
	}

	p.highestBid = 0

	if p.canPlayGrand() {
		p.grand = true
	}

	suit := mostCardsSuit(p.getHand())
	largest := len(trumpCards(suit, p.getHand()))
	debugTacticsLog("Longest suit %s, %d cards\n", suit, largest)
	asses := assOtherThan(suit)
	debugTacticsLog("Extra suits: %d\n", asses)
	prob := 0

// (You) 2571.27     (Bob) -784.38     (Ana) -1786.89      EURO

// (You) 61951     (Bob) 48905     (Ana) 48242     WON
// (You) 15028     (Bob)  7516     (Ana)  7570     LOST
// AVG 62.8, passed 20451, won 64492, lost 15057 / 100000 games	
// 76979	
	// You plays RISKY and wins
	if p.risky {
		if largest > 4 && asses > 0 {
			prob = 80
		}
// (You) 1487.00     (Bob) -506.20     (Ana) -980.80       EURO

// (You) 71919     (Bob) 48536     (Ana) 48303     WON
// (You) 22176     (Bob)  6802     (Ana)  6844     LOST
// AVG 59.4, passed 15621, won 66468, lost 17911 / 100000 games		
// 94095 - 76979 = 17116 more games: 9968 won- 7148 lost
		// if largest > 3 && asses > 1 {
		// 	prob = 80
		// }
	} else {
		if largest > 4 && asses > 1 {
			prob = 80
		}
	}

	if largest > 5 {
		prob = 85
	}
	if largest > 6 {
		prob = 99
	}

	est := p.handEstimation()
	debugTacticsLog("(%s) Hand: %v, Estimation: %d\n", p.name, p.hand, est)
	if prob < 80 {
		if est < 50 {
			return false
		}
	}
	//	fmt.Printf("LOW %d %v\n", p.handEstimation(), sortSuit(p.getHand()))

	// fmt.Printf("HIGH %d %v\n", p.handEstimation(), sortSuit(p.getHand()))
	return true
}

func (p *Player) calculateHighestBid() int {
	if !p.canWin() {
		return 0
	}

	p.grand = false
	trump := mostCardsSuit(p.getHand())
	if p.grand {
		trump = GRAND
	}
	mat := matadors(trump, p.getHand())
	if mat < 0 {
		mat *= -1
	}
	p.highestBid = (mat + 1) * trumpBaseValue(trump)
	return p.highestBid
}

func (p *Player) declareTrump() string {
	if p.grand {
		return GRAND
	}
	// TODO:
	// if after SKAT pick up bid less than score use the next suit
	trump := mostCardsSuit(p.getHand())
	// gs := gameScore(trump, p.hand, 61, 18,false, false, false)
	// if gs .....
	return trump
}

func (p *Player) discardInSkat(skat []Card) {
	debugTacticsLog("FULL HAND %v\n", sortSuit("", p.getHand()))

	// discard BLANKS

	bcards := findBlankCards(p.getHand())
	debugTacticsLog("BLANK %v\n", bcards)
	removed := 0
	if len(bcards) > 0 {
		p.setHand(remove(p.getHand(), bcards[0]))
		skat[0] = bcards[0]
		//	fmt.Printf("1st %v\n", skat)
		removed++
	}
	if len(bcards) > 1 {
		p.setHand(remove(p.getHand(), bcards[1]))
		skat[1] = bcards[1]
		//	fmt.Printf("2nd %v\n", skat)
		return
	}
	// Discard high cards in non-A suits with few colors
	sranks := []string{"J", "A", "10", "K", "D", "7", "8", "9"}

	lsuit := lessCardsSuit(p.getHand())
	if lsuit != "" {
		lcards := sortRankSpecial(filter(p.getHand(), func(c Card) bool {
			return c.suit == lsuit && c.rank != "A" && c.rank != "J"
		}), sranks)
		if len(lcards) < 4 { // do not throw long fleets
			debugTacticsLog("SUIT %v LESS %v\n", lsuit, lcards)

			if len(lcards) > 1 {
				i := 0
				for removed < 2 {
					p.setHand(remove(p.getHand(), lcards[i]))
					skat[removed] = lcards[i]
					i++
					removed++
				}
				return
			}
		}
	}

	// Discard non-A-10 suit cards
	ncards := nonA10cards(p.getHand())
	ncards = findBlankCards(ncards)
	// fmt.Printf("nonA10cards %v\n", ncards)

	if len(ncards) > 1 {
		i := 0
		for removed < 2 {
			p.setHand(remove(p.getHand(), ncards[i]))
			skat[removed] = ncards[i]
			i++
			removed++
		}
		return
	}

	if len(ncards) == 1 {
		p.setHand(remove(p.getHand(), ncards[0]))
		skat[removed] = ncards[0]
		removed++

		if removed == 2 {
			return
		}
	}

	trumpToDeclare := mostCardsSuit(p.getHand())
	cardsTodiscard := filter(sortRank(p.hand), func(c Card) bool {
		return c.suit != trumpToDeclare && c.rank != "J"
	})
	if len(cardsTodiscard) < 2 {
		debugTacticsLog("ALL TRUMPS (no 2 cards to discard)? %v", p.hand)
		cardsTodiscard = sortRank(p.hand)
	}
	debugTacticsLog("HAND %v\n", cardsTodiscard)
	if removed == 1 {
		card := cardsTodiscard[len(cardsTodiscard)-1]
		p.setHand(remove(p.getHand(), card))
		skat[1] = card
		return
	}
	c1 := cardsTodiscard[len(cardsTodiscard)-1]
	c2 := cardsTodiscard[len(cardsTodiscard)-2]
	p.setHand(remove(p.getHand(), c1))
	p.setHand(remove(p.getHand(), c2))
	skat[0] = c1
	skat[1] = c2
}

func (p *Player) pickUpSkat(skat []Card) bool {
	// TODO:
	// current implementation always picks up skat
	debugTacticsLog("SKAT BEF: %v\n", skat)
	hand := make([]Card, 10)
	copy(hand, p.getHand())
	hand = append(hand, skat...)
	p.setHand(hand)

	p.discardInSkat(skat)
	debugTacticsLog("SKAT AFT: %v\n", skat)
	return true
}
