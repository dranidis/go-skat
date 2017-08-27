package main

import (
	_ "fmt"
)

type Player struct {
	PlayerData
	firstCardPlay  bool
	risky          bool
	trumpToDeclare string
	afterSkat bool
}

func makePlayer(hand []Card) Player {
	return Player{
		PlayerData:     makePlayerData(hand),
		firstCardPlay:  false,
		risky:          false,
		trumpToDeclare: "",
		afterSkat: false,
	}
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

var TRYING = false

func (p Player) canWin() string {
	cs := p.getHand()
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
			} else if in(cs, Card{s, "10"}) && tenIsSupported(cs, Card{s, "10"}) {
				asses++
			}
		}
		return asses
	}

	fullOnes := assOtherThan("")
	losers := len(grandLosers(cs)) + jackLosers(cs)
	debugTacticsLog("\nLosers: %v, %d jacks\n", grandLosers(cs), jackLosers(cs))
	debugTacticsLog("\nConsidering GRAND in Hand: %v, Full ones: %v, Losers: %v\n", cs, fullOnes, losers)

	if p.afterSkat {
		debugTacticsLog("\nAFTER SKAT subtracting 2 losers\n")
		losers -= 2
	}
	p.afterSkat = false

	if fullOnes >= losers  || (p.risky && fullOnes + 1 >= losers) {
		if fullOnes < losers {
			TRYING = true
		}
	// if fullOnes >= losers  {
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
			TRYING = false
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

func (p Player) strongestLowestNotAor10(s *SuitState, cs []Card) Card {
	strongest := []Card{}
	n := firstCardTactic(cs)
	strongest = append(strongest, n)
	next := nextCard(s.trump, n)
	for in(s.cardsPlayed, next) || in(s.skat, next) || in(p.hand, next) {
		if in(p.hand, next) {
			strongest = append(strongest, next)
		}
		next = nextCard(s.trump, next)
	}
	debugTacticsLog("..Strongest lowest: %v ..", strongest)
	l := len(strongest) - 1
	for strongest[l].equals(Card{s.trump, "A"}) || strongest[l].equals(Card{s.trump, "10"}) {
		debugTacticsLog("..not playing a %v..", strongest[l])
		if l > 0 {
			l--
		} else {
			break
		}
	}
	return strongest[l]
}

func (p Player) declarerTactic(s *SuitState, c []Card) Card {
	debugTacticsLog("DECLARER ")
	if len(c) == 1 {
		debugTacticsLog("..FORCED MOVE.. ")
		return c[0]
	}	

	follows :=  followCards(s, s.follow, c)
	ownTrumps := sortRank(filter(p.hand, func(card Card) bool {
		return card.Rank == "J" || card.Suit == s.trump
	}))
	otherTrumps := p.otherPlayersTrumps(s)
	
	if len(s.trick) == 0 {
		debugTacticsLog("FOREHAND ")
		// count your own trumps and other players trump
		// if you have less you should not play trumps immediately
		if len(p.otherPlayersTrumps(s)) > 0 {

			debugTacticsLog("..other TRUMPS in game: %v", otherTrumps)
			debugTacticsLog("..own TRUMPS: %v", ownTrumps)


			// calculating sure winners and 2nd losers
			sureWinners := []Card{}
			for _,t := range ownTrumps {
				if noHigherCard(s, true, p.hand, t) {
					sureWinners = append(sureWinners, t)
				}
			}
			sureWinners = sortSuit(s.trump, sureWinners)
			debugTacticsLog("..sure winners: %v", sureWinners)

			highLosers := []Card{}
			for _,t := range ownTrumps {
				if in(sureWinners, t) {
					continue
				}
				if noSecondHigherCard(s, true, p.hand, t) {
					highLosers = append(highLosers, t)
				}
			}
			highLosers = sortValue(highLosers)
			debugTacticsLog("..high losers: %v", highLosers)

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
			if len(ownTrumps) * 2 < len(otherTrumps) + 2 {
				debugTacticsLog("Not enough trumps.  Playing suits")

				// Checking if card will not be trumped did not improve results !!???!!!
				ifMoreThanOne := func (cs []Card) bool {
					if len(cs) > 0 {
						// inPlay := suitCardsInPlay(s, true, p.hand, cs[0].Suit)
						// debugTacticsLog("..Cards in play of suit %s, %v..", cs[0].Suit, inPlay)						
						// if len(inPlay) > 1 {
							return true
						// }
					}
					return false
				}	

				suits := filter(c, func(card Card) bool {
					return card.Suit != s.trump && card.Rank != "J"
				})
				debugTacticsLog("..SUITS : %v...", suits)
				asses := filter(suits, func(card Card) bool {
					return card.Rank == "A"
				})
				if ifMoreThanOne(asses){
					return asses[0]
				}
				tens := filter(suits, func(card Card) bool {
					cardsPlayed := append(s.cardsPlayed, s.skat...)
					return card.Rank == "10" && in(cardsPlayed, Card{card.Suit, "A"})
				})
				if ifMoreThanOne(tens) {
					return tens[0]
				}					
				Ks := filter(suits, func(card Card) bool {
					cardsPlayed := append(s.cardsPlayed, s.skat...)
					return card.Rank == "K" && in(cardsPlayed, Card{card.Suit, "A"}, Card{card.Suit, "10"})
				})
				if ifMoreThanOne(Ks) {
					return Ks[0]
				}					
				Ds := filter(suits, func(card Card) bool {
					cardsPlayed := append(s.cardsPlayed, s.skat...)
					return card.Rank == "D" && in(cardsPlayed, Card{card.Suit, "A"}, Card{card.Suit, "10"}, Card{card.Suit, "K"})
				})								
				if ifMoreThanOne(Ds) {
					return Ds[0]
				}
				sortedValue := sortValue(c)
				// PLay a card with value
				// for i := len(sortedValue)-1; i >=0 ; i-- {
				// 	if cardValue(sortedValue[i]) > 0 {
				// 		return sortedValue[i]
				// 	}
				// }
				return sortedValue[len(sortedValue)-1]
			}

			if len(sureWinners) > 1 {
				debugTacticsLog("..more than one sure winner..")
				return sureWinners[len(sureWinners)-1]
			} else if len(sureWinners) == 1 {
				if len(otherTrumps) == 1 {
					return sureWinners[0]
				}
				//
				// does not improve score to HOLD the high winner for later
				//
				// if len(otherTrumps) > 2 {
				// 	debugTacticsLog("Playing low in first tricks (more than 2 trumps)")
				// 	return ownTrumps[len(ownTrumps) - 1]
				// }
			} else if len(sureWinners) == 0 {
				if len(highLosers) > 1 {
					debugTacticsLog("..Playing a high loser..")
					return highLosers[len(highLosers)-1]
				}
				debugTacticsLog("..No winner, playing low trump")
				return ownTrumps[len(ownTrumps) - 1]
			} else if len(otherTrumps) > 1 {
				debugTacticsLog("..one or less sure winner and more than one Trump in opponents. Playing a low value high loser..")
				if len(highLosers) > 0 {
					return highLosers[len(highLosers)-1]
				}
				return ownTrumps[len(ownTrumps)-1]
			}

			if p.otherPlayersHaveJs(s) {
				debugTacticsLog("... other players have Js..")
				validCards := make([]Card, len(c))
				copy(validCards, c)
				first := firstCardTactic(validCards)
				for len(validCards) > 1 && (first.equals(Card{s.trump, "A"}) || first.equals(Card{s.trump, "10"})  || first.equals(Card{s.trump, "K"})  || first.equals(Card{s.trump, "D"})) {
					debugTacticsLog("... not playing A or 10 or figures ..")
					validCards = remove(validCards, first)
					first = firstCardTactic(validCards)
				}
				return first
			} else {
				debugTacticsLog("Playing strongest lowest\n")
				return p.strongestLowestNotAor10(s, c)
				// return firstCardTactic(c)
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
				if noHigherCard(s, true, p.hand, cs[0]) {
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
	if len(s.trick) == 1 {
		debugTacticsLog("MIDDLEHAND ")
		if len(follows) == 0 && len(ownTrumps) > 0 {
			debugTacticsLog("..low strengrt/value trump.. ")

			rank := []string{"A", "10", "J", "K", "D", "9", "8", "7"}
			sortedTrumps := sortRankSpecial(ownTrumps, rank)
			return sortedTrumps[len(sortedTrumps) - 1]
		}
	}
	if len(s.trick) == 2 {
		debugTacticsLog("BACKHAND ")

		if len(follows) > 0 {
			debugTacticsLog("Following normal suit...")
			winners := sortRank(winnerCards(s, follows))
			debugTacticsLog("winners %v...", winners)
			for _, w := range winners {
				if w.Rank == "D" || w.Rank == "K" {
					debugTacticsLog("Returning  %v...", w)
					return w
				}
				if nextLowestCardsStillInPlay(s, w, follows) {
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
	debugTacticsLog("OPPONENT..")
	if len(c) == 1 {
		debugTacticsLog("FORCED MOVE\n")
		return c[0]
	}		

	if s.trump == NULL {
		return p.opponentTacticNull(s, c)
	}

	sortedValue := sortValue(c)
	noTrumps := filter(sortedValue, func(card Card) bool {
		return card.Suit != s.trump && card.Rank != "J"
	})
	ownTrumps := sortRank(filter(p.hand, func(card Card) bool {
		return card.Rank == "J" || card.Suit == s.trump
	}))
	sureWinners := []Card{}
	for _,t := range ownTrumps {
		if noHigherCard(s, true, p.hand, t) {
			sureWinners = append(sureWinners, t)
		}
	}	
	otherTrumps := p.otherPlayersTrumps(s)

	sortedValueNoTrumps := filter(sortValue(c), func(card Card) bool {
		return card.Suit != s.trump && card.Rank != "J"
	})

	if len(s.trick) == 0 {
		debugTacticsLog("FOREHAND..")

		if len(ownTrumps) >= len(otherTrumps) {
			debugTacticsLog("Many trumps..")
			if len(sureWinners) > 0 {
				debugTacticsLog("sure winners %v..", sureWinners)
				return sureWinners[0]
			}
			if len(ownTrumps) > len(otherTrumps) {
				debugTacticsLog("more than others...highest trump %v..", ownTrumps)
				return ownTrumps[0]
			}
		}
		
		// if you have a card with suit played in a previous trick
		// started from you or your partner continue with it
		prevSuit := p.FindPreviousSuit(s)
		prevSuitCards := followCards(s, prevSuit, c)
		if len(prevSuitCards) > 0 {
			debugTacticsLog("Following previous suit..")
			// TODO:
			// should I play the highest even if in the previous trick
			// declarer has taken with trump?
			// He will take it again.

			if (prevSuitCards[0].Rank == "10") && in(s.cardsPlayed, Card{prevSuit, "A"}) {
				debugTacticsLog("A and 10 seen resetting previous suit..")
				s.opp1.setPreviousSuit("")
				s.opp2.setPreviousSuit("")
				return p.opponentTactic(s, c)
			}
			return prevSuitCards[0]
		} else {
			debugTacticsLog("No cards in previous suit '%v'..", prevSuit)
		}

		var card Card
		if s.opp2 == p {
			debugTacticsLog("Declarer at MIDDLEHAND, playing LONG..")
			debugTacticsLog("not full ones if trumps in play..")

			if len(p.otherPlayersTrumps(s)) == 0 {
				debugTacticsLog("No trumps in play..")
				card = HighestLong(s.trump, c)
			} else {
				debugTacticsLog("trumps still in play..")
				card = HighestLongNoFull(s.trump, c)
			}
		} else {
			debugTacticsLog("Declarer at BACKHAND, playing SHORT..")
			debugTacticsLog("Avoiding 2 numbers and D-number suits..")
			nonTrumps := make([]Card, len(c))
			copy(nonTrumps, c)
			nonTrumps = filter(nonTrumps, func (card Card) bool {
				return getSuit(s.trump, card) != s.trump
				})
			candidates := []Card{}
			candidate := HighestShort(s.trump, c)
			candidates = append(candidates, candidate)
			for len(nonTrumps) > 0 && (twoNumbersSuit(nonTrumps, candidate.Suit) || DNumberSuit(nonTrumps, candidate.Suit)) {
				nonTrumps = filter(nonTrumps, func (crd Card) bool {
					return crd.Suit != candidate.Suit
					})
				candidate = HighestShort(s.trump, nonTrumps)
				if candidate.Suit != "" && candidate.Rank != "" {
					candidates = append(candidates, candidate)				
				}
			}
			debugTacticsLog("Candidates %v, returning last..", candidates)
			if len(candidates) > 0 {
				card = candidates[len(candidates) - 1]
			}
		}
		suit := getSuit(s.trump, card)
		if suit != s.trump {
			p.setPreviousSuit(suit)
		}
		return card
	}

	if len(s.trick) == 1 {
		debugTacticsLog("MIDDLEHAND..")
		if s.leader == s.declarer {
			debugTacticsLog("Declarer leads %v..", s.trick[0])
			// if declarer leads a low trump, and there are still HIGHER trumps
			// smear the trick with a high value
			if getSuit(s.trump, s.trick[0]) == s.trump && len(winnerCards(s, c)) == 0 {
				if len(filter(p.otherPlayersTrumps(s), func(c Card) bool {
					return s.greater(c, s.trick[0])
				})) > 0 {
					debugTacticsLog("TRUMP. There are higher trumps, SMEAR..")
					return sortValue(c)[0]
				}
			}
			if len(winnerCards(s, c)) == 0 && !noHigherCard(s, false, p.hand, s.trick[0]){
				debugTacticsLog("higher cards in play, SMEAR..")
				return sortValue(c)[0]
			}
			return highestValueWinnerORlowestValueLoser(s, c)
		} else {
			debugTacticsLog("Teammate leads %v..", s.trick[0])

			// if void at card played, and there are still cards in play
			// trump it to smear a trump and to put declarer at middlehand.
			if len(followCards(s, s.follow, c)) == 0 {
				debugTacticsLog("VOID on %s..", s.follow)
				inPlay := suitCardsInPlay(s, false, p.hand, s.follow)
				if len(inPlay) > 0 {
					debugTacticsLog("%s still in play %v..", s.follow, inPlay)

					if s.greater(s.trick[0], inPlay...) && len(sortedValueNoTrumps) > 0 {
						debugTacticsLog("partner wins the trick, smear..")
						return(sortedValueNoTrumps[0])	
					}

					sortedValueTrumps := sortValue(followCards(s, s.trump, c))
					if len(sortedValueTrumps) > 0 {
						debugTacticsLog("Playing a trump from %v..", sortedValueTrumps)
						i := 0
						for i < len(sortedValueTrumps) && cardValue(sortedValueTrumps[i]) < 3  && cardValue(sortedValueTrumps[i]) > 0 {
							i++
						}
						if i < len(sortedValueTrumps) {
							debugTacticsLog("Playing %v..", sortedValueTrumps[i])
							return sortedValueTrumps[i]
						}
						i--
						debugTacticsLog("Playing %v..", sortedValueTrumps[i])
						return sortedValueTrumps[i]
					}
					debugTacticsLog("No trump to play..")
				}
				debugTacticsLog("No %s in play..", s.follow)
			}

			debugTacticsLog("SORT-VALUE no trumps %v\n", sortedValueNoTrumps)

			cardsSuit := sortValue(followCards(s, s.trick[0].Suit, c))

			if len(cardsSuit) > 0 {
				debugTacticsLog("..FOLLOW %v..", cardsSuit)
				if noHigherCard(s, false, p.hand, cardsSuit[0]) || noHigherCard(s, false, p.hand, s.trick[0]) {
					return cardsSuit[0]
				}
				debugTacticsLog("HIGH CARD still out\n")
				if cardValue(s.trick[0]) > 0 {
					return cardsSuit[len(cardsSuit)-1]
				}
				debugTacticsLog(" increase zero value trick\n")
				for i := len(cardsSuit) - 1; i >= 0; i-- {
					if cardValue(cardsSuit[i]) > 0 {
						return cardsSuit[i]
					}
				}
				return cardsSuit[len(cardsSuit)-1]
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
			// debugTacticsLog(" -- declarer leads --\n")
			if s.greater(s.trick[0], s.trick[1]) {
				candidate := highestValueWinnerORlowestValueLoser(s, c)
				lowValue := sortedValue[0]
				if getSuit(s.trump, candidate) == s.trump && s.follow != s.trump {
					debugTacticsLog("..see if it is worth Taking with trump: %v..", candidate)
					if sum(s.trick) > 0 {
						debugTacticsLog("..taking with the trump..")
						return candidate
					}
					debugTacticsLog("..keeping the trump..")
					return lowValue
				}
				return candidate
			}
			debugTacticsLog(" largest non-trump")
			candidates := []Card{}
			candidates = sortedValue
			for len(candidates) > 0 {
				card := candidates[0]
				if in(sureWinners, card) {
					candidates = remove(candidates, card)
					continue
				}
				return card 
			}
			// if len(noTrumps) > 0 {
			// 	return noTrumps[0]
			// }
			return sortedValue[0]
		}
		debugTacticsLog(" -- teammate leads --\n")
		if s.greater(s.trick[0], s.trick[1]) {
			debugTacticsLog(" largest non-trump")

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

func (p *Player) getGamevalue(suit string) int {
	mat := matadors(suit, p.getHand())
	if mat < 0 {
		mat *= -1
	}
	return (mat + 1) * trumpBaseValue(suit)
}


func (p *Player) calculateHighestBid() int {
	p.highestBid = 0

	switch p.canWin() {
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

func (p *Player) discardInSkat(skat []Card) {
	debugTacticsLog("FULL HAND %v\n", sortSuit("", p.getHand()))

	newHbid := p.calculateHighestBid()
	debugTacticsLog("New high bid %d\n", newHbid)

	most := mostCardsSuit(p.getHand())

	removed := 0

	if p.trumpToDeclare == GRAND {
		debugTacticsLog("..GRAND..")



		// 	You	Bob	Ana
		// EURO 124.48	29.77	-154.25
		// WON   2664	 2614	 2501
		// LOST   569	  564	  562	
		// bidp    34	   34	   32	
		// pcw     82	   82	   82	
		// pcwd    18	   18	   18	
		// AVG  21.4, passed 526, won 7779, lost 1695 / 10000 games
		// Grand games 996, perc:  9.96

		// 	You	Bob	Ana
		// EURO 98.80	47.29	-146.09
		// WON   2662	 2613	 2501
		// LOST   571	  565	  562	
		// bidp    34	   34	   32	
		// pcw     82	   82	   82	
		// pcwd    18	   18	   18	
		// AVG  21.2, passed 526, won 7776, lost 1698 / 10000 games
		// Grand games 996, perc:  9.96

		tenSuits := []string{}
		counts := []int{}
		for _, s := range suits {
			if in(p.hand, Card{s, "A"}) {
				continue
			}
			if in(p.hand, Card{s, "10"}) {
				tenSuits = append(tenSuits, s)
				count := len(nonTrumpCards(s, p.hand))
				counts = append(counts, count)
			}
		}
		debugTacticsLog("..10s: %v, counts %v..", tenSuits, counts)
		nTenSuits := []string{}
		for i := 1; i < 8; i++ {
			for j, s := range tenSuits {
				if counts[j] == i {
					nTenSuits = append(nTenSuits, s)
				}
			}
		}
		debugTacticsLog("..10s: %v..", nTenSuits)

		for i := 0; removed < 2 && i < len(nTenSuits); removed++ {
			s := nTenSuits[i]
			card := Card{s, "10"}
			debugTacticsLog("REMOVING %v..", card)
			skat[removed] = card
			p.hand = remove(p.hand, card)
			i++
		}	

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

	p.afterSkat = true

	p.discardInSkat(skat)
	debugTacticsLog("SKAT AFT: %v\n", skat)
	debugTacticsLog("(%s) Hand %v\n", p.name, sortSuit("", p.hand))

	return true
}
