package main

import (
	_ "fmt"
)

type Player struct {
	PlayerData
	firstCardPlay  bool
	risky          bool
	trumpToDeclare string
}

func makePlayer(hand []Card) Player {
	return Player{
		PlayerData:     makePlayerData(hand),
		firstCardPlay:  false,
		risky:          false,
		trumpToDeclare: "",
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
		return card.Suit != s.trump
	})
	debugTacticsLog("LOSING Cards (last)%v\n", cards)
	if len(cards) > 0 {
		return cards[len(cards)-1]
	}
	// LAST?
	return sortedValue[len(sortedValue)-1]
}

func (p *Player) setPartner(partner PlayerI) {

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
	allCards := makeSuitDeck(c.Suit)
	allCardsPlayed := []Card{}
	if viewSkat {
		allCardsPlayed = append(allCardsPlayed, s.skat...)
	}
	allCardsPlayed = append(allCardsPlayed, s.cardsPlayed...)
	for _, cardPlayed := range allCardsPlayed {
		allCards = remove(allCards, cardPlayed)
	}
	allCards = filter(allCards, func(card Card) bool {
		return card.Rank != "J"
	})
	debugTacticsLog("Cards of suit %s still in play: %v", c.Suit, allCards)
	for _, card := range allCards {
		if s.greater(card, c) {
			return false
		}
	}
	return true
}

func nextLowestCardsStillInPlay(s *SuitState, w Card, followCards []Card) bool {
	next := nextCard(w)
	// debugTacticsLog("Next of %v is %v...", w, next)
	// ONLY the declarer knows that. Use a flag if opp uses it.
	if in(s.skat, next) || in(followCards, next) || in(s.cardsPlayed, next) {
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
			ownTrumps := sortRank(filter(p.hand, func(card Card) bool {
				return card.Rank == "J" || card.Suit == s.trump
			}))
			otherTrumps := p.otherPlayersTrumps(s)
			debugTacticsLog("..other TRUMPS in game: %v", otherTrumps)
			debugTacticsLog("..own TRUMPS: %v", ownTrumps)

			// 5, 2 => 10 < 4
			// 2, 4 => 4 < 6
			// 2, 3 => 4 < 5
			// 2, 2 => 4 < 4
			// 1, 0 => 2 < 2
			// 1, 1 => 2 < 3
			// 1, 2 => 2 < 4
			if len(ownTrumps) > 0 && len(otherTrumps) == 1 {
				if s.greater(ownTrumps[0], otherTrumps[0]) {
					return ownTrumps[0]
				}
			}
			if len(ownTrumps)*2 < len(otherTrumps)+2 {
				debugTacticsLog("Not enough trumps.  Playing suits")
				suits := filter(c, func(card Card) bool {
					return card.Suit != s.trump && card.Rank != "J"
				})
				debugTacticsLog("..SUITS : %v...", suits)
				asses := filter(suits, func(card Card) bool {
					return card.Rank == "A"
				})
				if len(asses) > 0 {
					return asses[0]
				}
				tens := filter(suits, func(card Card) bool {
					cardsPlayed := append(s.cardsPlayed, s.skat...)
					return card.Rank == "10" && in(cardsPlayed, Card{card.Suit, "A"})
				})
				if len(tens) > 0 {
					return tens[0]
				}
			}

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
		// TRUMP MONOPOLY
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
					return c.Suit == AKXsuit
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
				return card.Suit != s.trump && card.Rank != "J"
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
			return card.Suit == s.follow && card.Rank != "J"
		})
		if len(followCards) > 0 {
			debugTacticsLog("Following normal suit...")
			winners := sortRank(winnerCards(s, followCards))
			debugTacticsLog("winners %v...", winners)
			for _, w := range winners {
				if w.Rank == "D" || w.Rank == "K" {
					debugTacticsLog("Returning  %v...", w)
					return w
				}
				if nextLowestCardsStillInPlay(s, w, followCards) {
					debugTacticsLog("Next lower still in play...")
					continue
				}
				debugTacticsLog("Returning  %v...", w)
				return w
			}
		} else {
			debugTacticsLog("TRUMP OR No cards of suit played...")
		}
		if sum(s.trick) == 0 {
			debugTacticsLog("ZERO valued trick. DO not trump!...")
			// do not trump
			nonTrumps := filter(sortValue(c), func(card Card) bool {
				return card.Suit != s.trump && card.Rank != "J"
			})
			debugTacticsLog("Non-trumps: %v...", nonTrumps)

			sortedValue := filter(sortValue(c), func(card Card) bool {
				return card.Suit != s.trump && card.Rank != "J"
			})
			if len(sortedValue) > 0 {
				return highestValueWinnerORlowestValueLoser(s, sortedValue)
			}
		}
		// don't throw your A if not 10 in the trick and still in game
		debugTacticsLog("CHECKING the A... in trick %v, valid %v, played %v..", s.trick, c, s.cardsPlayed)
		if s.follow != s.trump && in(c, Card{s.follow, "A"}) && !in(s.trick, Card{s.follow, "10"}) && !in(append(s.cardsPlayed, s.skat...), Card{s.follow, "10"}) {
			debugTacticsLog("Keeping the A... in trick %v..", s.trick)
			sortedRank := filter(sortRank(c), func(card Card) bool {
				return card.Rank != "A" && card.Rank != "J"
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

func (p *Player) FindPreviousSuit(s *SuitState) string {
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

func (p *Player) opponentTacticNull(s *SuitState, c []Card) Card {
	revValue := sortValueNull(c)
	debugTacticsLog("NULL: rev value %v\n", revValue)

	notExhausted := sortValueNull(filter(c, func(card Card) bool {
		return !isVoid(s, c, card.Suit)
	}))
	debugTacticsLog("NULL: not exhausted %v\n", notExhausted)

	if len(s.trick) == 0 {
		debugTacticsLog("NULL FOREHAND..")
		prevSuit := p.FindPreviousSuit(s)
		debugTacticsLog("Prev suit: %v..", prevSuit)
		if prevSuit != "" && !s.declarerVoidSuit[prevSuit] {
			debugTacticsLog("VOID suits: %v", s.declarerVoidSuit)
			var prevSuitCards []Card
			if prevSuit != "" {
				prevSuitCards = filter(notExhausted, func(c Card) bool {
					return c.Suit == prevSuit
				})
			}
			if len(prevSuitCards) > 0 {
				debugTacticsLog("Following previous suit: cards %v %v..", prevSuitCards, prevSuitCards[0])
				return prevSuitCards[0]
			}
		}

		singles := singletons(notExhausted)
		debugTacticsLog("Singles %v..", singles)
		if len(singles) > 0 {
			s := singles[0]
			debugTacticsLog("PLAYING singleton %v..", s)
			p.previousSuit = s.Suit
			return s
		}

		if len(notExhausted) > 0 {
			p.previousSuit = notExhausted[0].Suit
			return notExhausted[0]
		}
		p.previousSuit = revValue[0].Suit
		return revValue[0]

	}

	if len(s.trick) > 0 {
		if len(filter(c, func(card Card) bool {
			return card.Suit == s.trick[0].Suit
		})) == 0 {
			debugTacticsLog("THROWING OFF..")
			return revValue[len(revValue)-1]
		}
	}

	if len(s.trick) == 1 {
		debugTacticsLog("NULL MIDHAND..")

		if s.leader == s.declarer {
			debugTacticsLog("Declarer opened..")
			// if s.greater(revValue[0], s.trick[0]) {
			// 	return revValue[len(revValue)-1]
			// }
			return revValue[0]
		}
		debugTacticsLog("Declarer at Back..")
		//NOT SURE ABOUT THIS::
		if smallerCardsInPlay(s, s.trick[0], c) {
			debugTacticsLog("Smaller still in play, throwing off ...")
			pr := previousNull(s.trick[0])
			debugTacticsLog("Previous %v..", pr)
			for in(s.cardsPlayed, pr) {
				pr = previousNull(pr)
				debugTacticsLog("Previous %v..", pr)
			}
			if in(c, pr) {
				debugTacticsLog("Returninhg Previous %v..", pr)
				return pr
			}
			//	return revValue[len(revValue)-1]
		}
		return revValue[0]
	}
	if len(s.trick) == 2 {
		debugTacticsLog("NULL BACKHAND..")
		debugTacticsLog("RevVal %v..", revValue)
		if s.leader == s.declarer {
			debugTacticsLog("Declarer leads..")
			if s.greater(revValue[0], s.trick[0]) || s.greater(s.trick[1], s.trick[0]) {
				debugTacticsLog("Returning last %v...", revValue[len(revValue)-1])
				return revValue[len(revValue)-1]
			}
			return revValue[0]
		}
		debugTacticsLog("Teammate leads..")
		if s.greater(revValue[0], s.trick[1]) || s.greater(s.trick[0], s.trick[1]) {
			debugTacticsLog("Returning last %v...", revValue[len(revValue)-1])
			return revValue[len(revValue)-1]
		}
		return revValue[0]
	}

	return revValue[0]
}

func (p *Player) opponentTactic(s *SuitState, c []Card) Card {
	if s.trump == NULL {
		return p.opponentTacticNull(s, c)
	}
	// OPPONENTS TACTIC
	if len(s.trick) == 0 {
		debugTacticsLog("OPP FOREHAND\n")
		// if you have a card with suit played in a previous trick
		// started from you or your partner continue with it
		prevSuit := p.FindPreviousSuit(s)
		var prevSuitCards []Card
		if prevSuit != "" {
			prevSuitCards = filter(c, func(c Card) bool {
				return c.Suit == prevSuit && c.Rank != "J"
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
				return card.Suit != s.trump && card.Rank != "J"
			})
			debugTacticsLog("SORT-VALUE no trumps %v\n", sortedValueNoTrumps)

			cardsSuit := filter(c, func(card Card) bool {
				return s.trick[0].Suit == card.Suit && card.Rank != "J"
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
				if card.Rank == "J" {
					return false
				}
				if card.Suit != s.trick[0].Suit {
					return false
				}
				if card.Rank == s.trick[0].Rank {
					return false
				}
				return true
			})
			// else
			debugTacticsLog("PLAYED from suite %v", filter(s.cardsPlayed, func(card Card) bool {
				return card.Suit == s.trick[0].Suit && card.Rank != "J"
			}))
			if in(s.cardsPlayed, trickFollowCards...) && sum(s.trick) == 0 {
				debugTacticsLog("All suit %s played and zero trick. Increase trick value\n", s.trick[0].Suit)
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
				return card.Suit != s.trump && card.Rank != "J"
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
	debugTacticsLog("(%s) Bid: %d Highest: %d\n", p.name, bids[bidIndex], p.highestBid)
	return bids[bidIndex] <= p.highestBid
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

func canWin(cs []Card) string {
	assOtherThan := func(suit string) int {
		asses := 0
		for _, s := range suits {
			if s == suit {
				continue
			}
			if in(cs, Card{s, "A"}) {
				//	debugTacticsLog("(A %s) ", s)
				asses++
				if in(cs, Card{s, "10"}) {
					//		debugTacticsLog("(10 %s) ", s)
					asses++
					// if in(p.getHand(), Card{s, "K"}) {
					// 	debugTacticsLog("(K %s) ", s)
					// 	asses++
					// }
				}
			}
		}
		return asses
	}

	fullOnes := assOtherThan("")
	losers := len(grandLosers(cs)) + jackLosers(cs)
	debugTacticsLog("\nLosers: %v, %d jacks\n", grandLosers(cs), jackLosers(cs))
	debugTacticsLog("\nConsidering GRAND in Hand: %v, Full ones: %v, Losers: %v\n", cs, fullOnes, losers)
	if fullOnes > losers {
		asuits := 0
		for _, s := range suits {
			if in(cs, Card{s, "A"}) {
				asuits++
			}
		}
		Js := filter(cs, func(c Card) bool {
			return c.Rank == "J"
		})
		debugTacticsLog("Js %v, Asuits %d\n", Js, asuits)
		if len(Js) > 1 {
			debugTacticsLog("WILL PLAY GRAND with Jacks: %v\n", Js)
			return "GRAND"
		}
		if len(Js) == 1 {
			if asuits >= 4 {
				debugTacticsLog("WILL PLAY GRAND with 1 Jack and 4 suits covered with A: %v\n", Js)
				return "GRAND"
			}
		}
		//return "GRAND"

	}
	// if len(filter(p.hand, func (c Card) bool {
	// 	return c.Rank == "J"
	// })) > 1 {
	// 	asses := assOtherThan("")
	// 	if asses > 3 {
	// 		debugTacticsLog(" - Can win GRAND - ")
	// 		return "GRAND"
	// 	}
	// }

	suit := mostCardsSuit(cs)
	largest := len(trumpCards(suit, cs))
	debugTacticsLog("Longest suit %s, %d cards\n", suit, largest)
	asses := assOtherThan(suit)
	debugTacticsLog("Extra suits: %d\n", asses)
	prob := 0

	if largest > 4 && asses > 0 {
		prob = 80
	}

	if largest > 5 {
		prob = 85
	}
	if largest > 6 {
		prob = 99
	}

	est := handEstimation(cs)
	debugTacticsLog("Hand: %v, Estimation: %d\n", cs, est)
	if prob < 80 {
		if est < 50 {
			return ""
		}
	}
	//	fmt.Printf("LOW %d %v\n", p.handEstimation(), sortSuit(p.getHand()))

	// fmt.Printf("HIGH %d %v\n", p.handEstimation(), sortSuit(p.getHand()))
	return "SUIT"
}

func (p *Player) getGamevalue(suit string) int {
	mat := matadors(suit, p.getHand())
	if mat < 0 {
		mat *= -1
	}
	return (mat + 1) * trumpBaseValue(suit)
}

func (p *Player) calculateHighestBid() int {
	p.highestBid = 0

	switch canWin(p.hand) {
	case "":
		return 0
	case "SUIT":
		p.trumpToDeclare = mostCardsSuit(p.getHand())
		p.highestBid = p.getGamevalue(p.trumpToDeclare)
	case "GRAND":
		p.trumpToDeclare = GRAND
		p.highestBid = p.getGamevalue(p.trumpToDeclare)
	default:
		return 0
	}
	if matadors(p.trumpToDeclare, p.hand) < -1 {
		// maybe you pick the CLUBS J from the skat 1/4
		worstCaseScore := 2 * trumpBaseValue(p.trumpToDeclare)
		if worstCaseScore < p.highestBid {
			debugTacticsLog("(%s) I will not raise more. Bid: %d, worst: %d\n", p.name, p.declaredBid, worstCaseScore)
			p.highestBid = worstCaseScore
		}
	}
	return p.highestBid
}

func (p *Player) declareTrump() string {
	if p.trumpToDeclare == GRAND {
		return GRAND
	}
	// TODO:
	// if after SKAT pick up bid less than score use the next suit
	trump := mostCardsSuit(p.getHand())

	// 	EURO 487.79     -364.63 -123.16
	// WON  27225      18802   18999
	// LOST  7482       3631    3532
	// bidp    44         28      28
	// pcw     78         84      84
	// pcwd    16         19      19
	// AVG  17.2, passed 20329, won 65026, lost 14645 / 100000 game
	if p.getGamevalue(trump) < p.declaredBid {
		debugTacticsLog("Game Value: %d. Declared bid: %d. TO AVOID OVERBID I will play first trump %s and not new %s.\n",
			p.getGamevalue(trump), p.declaredBid, p.trumpToDeclare, trump)
		trump = p.trumpToDeclare
	}
	// EURO -524.86    229.37  295.49
	// WON  26504      18436   18633
	// LOST  8205       3996    3897
	// bidp    44         28      28
	// pcw     76         82      83
	// pcwd    18         21      21
	// AVG  18.2, passed 20329, won 63573, lost 16098 / 100000 game

	// gs := gameScore(trump, p.hand, 61, 18,false, false, false)
	// if gs .....
	return trump
}

//         You     Bob     Ana
// EURO 342.82     -315.20 -27.62
// WON  24407      24244   24411
// LOST  5289       5485    5404
// bidp    33         33      33
// pcw     82         82      82
// pcwd    18         18      18
// AVG  17.9, passed 10760, won 73062, lost 16178 / 100000 game

// WITH GRAND
//         You     Bob     Ana
// EURO 345.80     -324.31 -21.49
// WON  24410      24249   24437
// LOST  5283       5476    5385
// bidp    33         33      33
// pcw     82         82      82
// pcwd    18         18      18
// AVG  18.6, passed 10760, won 73096, lost 16144 / 100000 games
// Grand games 1351, perc:  1.35

func grandSuitLosers(cs []Card) []Card {
	if len(cs) == 0 {
		return cs
	}
	s := cs[0].Suit
	if in(cs, Card{s, "A"}) {
		cs = remove(cs, Card{s, "A"})
		if in(cs, Card{s, "10"}) {
			cs = remove(cs, Card{s, "10"})
			if in(cs, Card{s, "K"}) {
				cs = remove(cs, Card{s, "K"})
				if in(cs, Card{s, "D"}) {
					cs = remove(cs, Card{s, "D"})
				}
				return cs
			}
			return cs
		}
		return cs
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
		losers = append(losers, grandSuitLosers(cards)...)
	}
	return losers
}

func (p *Player) discardInSkat(skat []Card) {
	debugTacticsLog("FULL HAND %v\n", sortSuit("", p.getHand()))
	most := mostCardsSuit(p.getHand())

	removed := 0

	if p.trumpToDeclare == GRAND {
		debugTacticsLog("..GRAND..")
		losers := sortValue(grandLosers(p.hand))
		debugTacticsLog("..GRAND..losers %v..", losers)
		for ; removed < 2 && len(losers) > 0 && cardValue(losers[0]) > 0; removed++ {
			card := losers[0]
			debugTacticsLog("REMOVING %v..", card)
			skat[removed] = card
			p.hand = remove(p.hand, card)
			losers = remove(losers, card)
		}
		bcards := findBlankCards(p.getHand())
		for ; removed < 2 && len(bcards) > 0; removed++ {
			card := bcards[0]
			skat[removed] = card
			p.hand = remove(p.hand, card)
			bcards = remove(bcards, card)
		}
		if removed == 2 {
			return
		}
	}

	// discard BLANKS

	bcards := findBlankCards(p.getHand())
	debugTacticsLog("BLANK %v\n", bcards)
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
	debugTacticsLog("..Less cards suit %v..", lsuit)
	if lsuit != "" {
		lcards := sortRankSpecial(filter(p.getHand(), func(c Card) bool {
			return c.Suit == lsuit && c.Rank != "A" && c.Rank != "J"
		}), sranks)
		debugTacticsLog(".. TRUMP to DECLARE [%s]..", p.trumpToDeclare)
		if lsuit != p.trumpToDeclare { //len(lcards) < 4 { // do not throw long fleets
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

	cardsTodiscard := filter(sortRank(p.hand), func(c Card) bool {
		return c.Suit != most && c.Rank != "J"
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
	debugTacticsLog("SKAT BEF: %v\n", skat)
	hand := make([]Card, 10)
	copy(hand, p.getHand())
	hand = append(hand, skat...)
	p.setHand(hand)

	p.discardInSkat(skat)
	debugTacticsLog("SKAT AFT: %v\n", skat)
	debugTacticsLog("(%s) Hand %v\n", p.name, sortSuit("", p.hand))
	return true
}
