package main

import (
	_ "fmt"
	"log"
)

func winnerCards(s *SuitState, cs []Card) []Card {
	if len(s.trick) > 0 {
		s.follow = getSuit(s.trump, s.trick[0])
	}
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

	trumpSuits := []string{CARO, HEART, SPADE, CLUBS}
	trumpWinnerRanks := []string{"A", "10", "K", "D", "9", "8", "7", "J"}
	if s.trump == s.follow {
		winners = sortSuitRankSpecial(winners, trumpSuits, trumpWinnerRanks)
	} else {
		// winners = sortValue(winners)
		winners = sortSuitRankSpecial(winners, trumpSuits, trumpWinnerRanks)
	}

	if len(winners) > 0 {
		debugTacticsLog("Winners: %v\n", winners)
		return winners[0]
	}

	if s.trump == s.follow {
		trumpRanks := []string{"A", "10", "J", "K", "D", "9", "8", "7"}
		cards := sortRankSpecial(c, trumpRanks)
		if len(cards) > 0 {
			return cards[len(cards)-1]
		}
	}

	sortedValue := sortValue(c)
	cards := filter(sortedValue, func(card Card) bool {
		return card.Suit != s.trump
	})
	debugTacticsLog("LOSING Cards (last)%v\n", cards)
	if len(cards) > 0 {
		return cards[len(cards)-1]
	}
	// LAST?
	return sortedValue[len(sortedValue)-1]
}

func firstCardTactic(c []Card) Card {
	return c[0]
}

// func noHigherCard(s *SuitState, viewSkat bool, c Card) bool {
// 	allCards := makeSuitDeck(c.Suit)
// 	allCardsPlayed := []Card{}
// 	if viewSkat {
// 		allCardsPlayed = append(allCardsPlayed, s.skat...)
// 	}
// 	allCardsPlayed = append(allCardsPlayed, s.cardsPlayed...)
// 	allCardsPlayed = append(allCardsPlayed, s.trick...)
// 	for _, cardPlayed := range allCardsPlayed {
// 		allCards = remove(allCards, cardPlayed)
// 	}
// 	allCards = filter(allCards, func(card Card) bool {
// 		return card.Rank != "J"
// 	})
// 	debugTacticsLog("Cards of suit %s still in play: %v", c.Suit, allCards)
// 	for _, card := range allCards {
// 		if s.greater(card, c) {
// 			return false
// 		}
// 	}
// 	return true
// }
func noHigherCard(s *SuitState, viewSkat bool, hand []Card, c Card) bool {
	suit := getSuit(s.trump, c)
	scip := suitCardsInPlay(s, viewSkat, hand, suit)
	debugTacticsLog("Cards of suit %s still in play: %v", suit, scip)
	for _, card := range scip {
		if s.greater(card, c) {
			return false
		}
	}
	return true
}

func HigherCards(s *SuitState, viewSkat bool, hand []Card, c Card) []Card {
	suit := getSuit(s.trump, c)
	scip := suitCardsInPlay(s, viewSkat, hand, suit)
	debugTacticsLog("Cards of suit %s still in play: %v", suit, scip)
	higherCards := []Card{}
	for _, card := range scip {
		if s.greater(card, c) {
			higherCards = append(higherCards, card)
		}
	}
	debugTacticsLog("Cards of suit %s higher than %v, %v", suit, c, higherCards)
	return higherCards
}

func noSecondHigherCard(s *SuitState, viewSkat bool, hand []Card, c Card) bool {
	firstFound := false
	scip := suitCardsInPlay(s, viewSkat, hand, getSuit(s.trump, c))
	// debugTacticsLog("Cards of suit %s still in play: %v", c.Suit, scip)
	for _, card := range scip {
		if !firstFound && s.greater(card, c) {
			firstFound = true
			continue
		}
		if s.greater(card, c) {
			return false
		}
	}
	return true
}

func suitCardsInPlay(s *SuitState, viewSkat bool, hand []Card, suit string) []Card {
	allCards := []Card{}
	if s.trump == suit {
		allCards = makeTrumpDeck(suit)
	} else {
		allCards = makeNoTrumpDeck(suit)
	}
	allCardsPlayed := []Card{}
	if viewSkat {
		allCardsPlayed = append(allCardsPlayed, s.skat...)
	}
	allCardsPlayed = append(allCardsPlayed, s.cardsPlayed...)
	allCardsPlayed = append(allCardsPlayed, s.trick...)
	allCardsPlayed = append(allCardsPlayed, hand...)
	// for _, cardPlayed := range allCardsPlayed {
	// allCards = remove(allCards, cardPlayed...)
	// }
	allCards = remove(allCards, allCardsPlayed...)
	return allCards
}

func nextLowestCardsStillInPlay(s *SuitState, c Card, followCards []Card) bool {
	next := nextCard(s.trump, c)
	// debugTacticsLog("Next of %v is %v...", c, next)
	// ONLY the declarer knows that. Use a flag if opp uses it.
	if in(s.skat, next) || in(followCards, next) || in(s.cardsPlayed, next) {
		return false
	}
	return true
}

func isAKX(trump string, cs []Card) bool {
	cards := filter(cs, func(c Card) bool {
		return c.Suit == trump
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

func is10X(suit string, cs []Card) bool {
	cards := filter(cs, func(c Card) bool {
		return c.Suit == suit
	})

	if in(cards, Card{suit, "A"}) {
		return false
	}
	if !in(cards, Card{suit, "10"}) {
		return false
	}
	if len(cards) > 2 {
		// fmt.Printf("less than 3\n")
		return false
	}
	return true
}

func HighestShort(trump string, c []Card) Card {
	s := ShortestNonTrumpSuit(trump, c)
	debugTacticsLog("ShortestNonTrumpSuit %v\n", s)
	cards := sortRank(nonTrumpCards(s, c))
	debugTacticsLog("HighestShort %v..", cards)
	if len(cards) > 0 {
		return cards[0]
	}
	// last card?
	if len(c) > 0 {
		debugTacticsLog("... DEBUG ... VALID: %v no HighestShort. Returning: %v\n", c, c[0])
		return c[0]
	}
	debugTacticsLog(".. NO CARDS..")
	return Card{"", ""}
}

func HighestShortNotFull(trump string, c []Card) Card {
	s := ShortestNonTrumpSuit(trump, c)
	debugTacticsLog("ShortestNonTrumpSuit %v\n", s)
	cards := sortRank(nonTrumpCards(s, c))
	debugTacticsLog("HighestShort not full %v..", cards)
	for i := 0; i < len(cards); i++ {
		if cardValue(cards[i]) > 4 {
			continue
		}
		return cards[i]
	}
	debugTacticsLog(".. only fulls..")
	if len(cards) > 0 {
		return cards[len(cards)-1]
	}
	// last card?
	if len(c) > 0 {
		debugTacticsLog("... DEBUG ... VALID: %v no HighestShort. Returning: %v\n", c, c[len(c)-1])
		return c[len(c)-1]
	}
	debugTacticsLog(".. NO CARDS..")
	return Card{"", ""}
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

func HighestLongNoFull(trump string, c []Card) Card {
	s := LongestNonTrumpSuit(trump, c)
	debugTacticsLog("LongestNonTrumpSuit %v\n", s)
	cards := sortRank(nonTrumpCards(s, c))
	debugTacticsLog("%v", cards)
	for i := 0; i < len(cards); i++ {
		if cardValue(cards[i]) > 4 {
			continue
		}
		return cards[i]
	}
	if len(cards) > 0 {
		return cards[len(cards)-1]
	}
	// last card?
	debugTacticsLog("... DEBUG ... VALID: %v no highest long. Returning: %v\n", c, c[0])
	return c[0]
}

// {8, 10} J  {K,9} ===> 7 in play
func smallerCardsInPlay(s *SuitState, trick Card, c []Card) bool {
	suit := trick.Suit
	for i := len(nullRanks) - 1; i >= 0; i-- {
		if trick.Rank == nullRanks[i] {
			break
		}
		cardsNotInPlay := append(s.cardsPlayed, c...)
		if in(cardsNotInPlay, Card{suit, nullRanks[i]}) {
			continue
		}
		return true
	}
	return false
}

func isVoid(s *SuitState, cs []Card, suit string) bool {
	if s.declarerVoidSuit[suit] {
		return true
	}
	cardsPlayed := filter(append(s.cardsPlayed, cs...), func(c Card) bool {
		return c.Suit == suit
	})

	if len(cardsPlayed) == 8 {
		return true
	}

	return false
}

//
// Der US-Amerikaner J.P. Wergin hat in seinem Buch "Wergin on Skat and Sheepshead"
// (McFarland, Wisconsin, 1975) versucht, dazu einen objektiven Berechnungsmodus zu
// finden.
func handEstimation(cs []Card) int {
	kreuzB := in(cs, Card{CLUBS, "J"})
	pikB := in(cs, Card{SPADE, "J"})
	herzB := in(cs, Card{HEART, "J"})
	karoB := in(cs, Card{CARO, "J"})

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
		a := in(cs, Card{suit, "A"})
		t := in(cs, Card{suit, "10"})
		k := in(cs, Card{suit, "K"})
		d := in(cs, Card{suit, "D"})
		n := in(cs, Card{suit, "9"})

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
			if c.Rank == "J" {
				continue
			}
			if c.Suit == suit {
				card = c
				break
			}
		}
		if card.Rank != "A" {
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
		if card.Rank != "" {
			blankCards = append(blankCards, card)
		}
	}
	blankCards = sortRank(blankCards)
	return blankCards
}

func nonA10cards(cs []Card) []Card {
	suitf := func(suit string, cs []Card) []Card {
		cards := filter(cs, func(c Card) bool {
			return c.Suit == suit && c.Rank != "J"
		})
		if in(cards, Card{suit, "A"}) && !in(cards, Card{suit, "10"}) {
			cards := filter(cards, func(c Card) bool {
				return c.Rank != "A"
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

func tenIsSupported(cs []Card, ten Card) bool {
	if ten.Rank != "10" {
		log.Fatal("ten argument is not a 10 card")
	}
	s := ten.Suit
	suitCards := filter(cs, func(c Card) bool {
		return c.Suit == s && c.Rank != "J"
	})
	if len(suitCards) > 2 {
		debugTacticsLog("..Supported 10 in %v..", suitCards)
		return true
	}
	return false
}

func grandSuitLosers(cs []Card) []Card {
	// debugTacticsLog("Card: %v\n", cs)
	l := len(cs)
	if l == 0 {
		return cs
	}
	s := cs[0].Suit
	if in(cs, Card{s, "A"}) {
		if l >= 6 {
			return []Card{}
		}
		cs = remove(cs, Card{s, "A"})
		if in(cs, Card{s, "10"}) {
			if l >= 5 {
				return []Card{}
			}
			cs = remove(cs, Card{s, "10"})
			if in(cs, Card{s, "K"}) {
				return []Card{}
				// if len >= 4 {
				// 	return []Card{}
				// }
				// cs = remove(cs, Card{s, "K"})
				// if in(cs, Card{s, "D"}) {
				// 	cs = []Card{}
				// }
				// return cs
			}
			return cs
		}
		return cs
	} else if in(cs, Card{s, "10"}) {
		// && tenIsSupported(cs, Card{s, "10"}) {
		debugTacticsLog(".. removing 10s, they will be discarded..")
		cs = remove(cs, Card{s, "10"}) // will be discarded
	}
	return cs
}

func jackLosers(cs []Card) int {
	c := in(cs, Card{CLUBS, "J"})
	s := in(cs, Card{SPADE, "J"})
	h := in(cs, Card{HEART, "J"})
	k := in(cs, Card{CARO, "J"})
	if c {
		if s {
			if h {
				return 0
			}
			if k {
				return 1
			}
			// return 0
		} else if h || k {
			return 1
		}
		return 0
	} else {
		if s {
			if h {
				return 1
			}
			if k {
				return 2
			}
			return 1
		} else if h {
			if k {
				return 2
			}
			return 1
		} else if k {
			return 1
		} else if h && k {
			return 2
		}
		return 0
	}
}

func grandLosers(cs []Card) []Card {
	losers := []Card{}
	for _, s := range suits {
		cards := filter(cs, func(c Card) bool {
			return c.Rank != "J" && c.Suit == s
		})
		gsl := grandSuitLosers(cards)
		debugTacticsLog("Suit: %s, losers: %v\n", s, gsl)
		losers = append(losers, gsl...)
	}
	return losers
}

func twoNumbersSuit(cs []Card, suit string) bool {
	debugTacticsLog("..Checking cards %v for 2 number suit %s..", cs, suit)
	suitCards := filter(cs, func(c Card) bool {
		return c.Suit == suit && c.Rank != "J"
	})
	if len(suitCards) == 2 && sum(suitCards) == 0 {
		debugTacticsLog("FOUND..")
		return true
	}
	return false
}

func DNumberSuit(cs []Card, suit string) bool {
	debugTacticsLog("..Checking cards %v for D number suit %s..", cs, suit)
	suitCards := filter(cs, func(c Card) bool {
		return c.Suit == suit && c.Rank != "J"
	})
	if len(suitCards) == 2 && sum(suitCards) == 3 {
		debugTacticsLog("FOUND..")
		return true
	}
	return false
}

// Returns the cards that follow the suit argument
func followCards(s *SuitState, suit string, cards []Card) []Card {
	if suit == s.trump {
		return trumpCards(suit, cards)
	}
	return nonTrumpCards(suit, cards)
}

// declarer is going to play the card now
func declarerCardIsLosingTrick(s *SuitState, p PlayerI, card Card) bool {
	for _, c := range s.trick {
		if s.greater(c, card) {
			return true
		}
	}
	if len(s.trick) == 2 {
		return false
	}
	if noHigherCard(s, true, p.getHand(), card) {
		return false
	}
	return true
}
