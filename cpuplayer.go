package main

type Player struct {
	PlayerData
}

func makePlayer(hand []Card) Player {
	return Player{
		PlayerData: makePlayerData(hand)}
}

func winnerCards(s *SuitState, c []Card) []Card {
	return filter(c, func(c Card) bool {
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
	cards := sortValue(c)
	return cards[len(cards)-1]
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
		if getSuite(s.trump, c) != s.trump {
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

func (p Player) declarerTactic(s *SuitState, c []Card) Card {
	//gameLog("SOLIST\n")
	if len(s.trick) == 0 {

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
		debugTacticsLog(" -TRUMPS exhausted: Hand: %v ", p.getHand())
		return HighestLong(s.trump, c)
	}
	// TODO:
	// in middlehand, if leader leads with an exhausted suit don't take
	// it with the most valuable trump. It might be taken....

	return highestValueWinnerORlowestValueLoser(s, c)
}

func HighestLong(trump string, c []Card) Card {
	s := LongestNonTrumpSuit(trump, c)
	debugTacticsLog("LongestNonTrumpSuit %v\n", s)
	cards := filter(c, func(c Card) bool {
		return c.suit == s && c.rank != "J"
	})
	debugTacticsLog("%v", cards)
	if len(cards) > 0 {
		return cards[0]
	}
	// last card?
	debugTacticsLog("... DEBUG ... VALID: %v no highest long. Returning: %v\n", c, c[0])
	return c[0]
}

func (p *Player) opponentTactic(s *SuitState, c []Card) Card {
	// OPPONENTS TACTIC

	// TODO:
	// if opp plays first high-value card and
	// partner has highest trump then he should play it

	// TODO:
	// if trick has no value, e.g. 7 8
	// don't take it with a trump if you can save it.

	if len(s.trick) == 0 {
		// if you have a card with suit played in a previous trick started from you or your partner continue with it
		// else
		prevSuit := ""
		partnerFunc := func(p PlayerI) PlayerI {
			if s.opp1 == p {
				return s.opp2
			}
			return s.opp1
		}
		if p.getPreviousSuit() != "" {
			prevSuit = p.getPreviousSuit()
		} else if partnerFunc(p).getPreviousSuit() != "" {
			prevSuit = partnerFunc(p).getPreviousSuit()
		}
		debugTacticsLog("Previous suit: %v\n", prevSuit)
		if prevSuit == "" {
			card := HighestLong(s.trump, c)
			p.setPreviousSuit(getSuite(s.trump, card))
			return card
		}
		suitCards := filter(c, func(c Card) bool {
			return c.suit == prevSuit && c.rank != "J"
		})
		if len(suitCards) > 0 {
			debugTacticsLog("Following previous suit...")
			return suitCards[0]
		}

		card := HighestLong(s.trump, c)
		p.setPreviousSuit(card.suit)
		return card
	}

	if len(s.trick) == 1 {
		debugTacticsLog("MIDDLEHAND\n")
		if s.leader == s.declarer {
			debugTacticsLog(" -- Declarer leads, ")
			// MIDDLEHAND
			// if declarer leads a low trump, and there are still HIGHER trumps
			// smear it with a high value
			if getSuite(s.trump, s.trick[0]) == s.trump && len(winnerCards(s, c)) == 0 {
				other := filter(p.otherPlayersTrumps(s), func (c Card) bool {
					return s.greater(c, s.trick[0])
					})
				if len(other) > 0 {
					debugTacticsLog("SMEAR ")
					return sortValue(c)[0]
				}
			}
			return highestValueWinnerORlowestValueLoser(s, c)
		} else {
			// TODO:TODO:TODO:TODO:TODO:TODO:
			// if high chances that teammate's card wins
			return sortValue(c)[0]
			// else
			// return sortValue(c)[len(c)-1]
		}
	}

	if len(s.trick) == 2 {
		if s.leader == s.declarer {
			if s.greater(s.trick[0], s.trick[1]) {
				return highestValueWinnerORlowestValueLoser(s, c)
			}
			return sortValue(c)[0]
		}
		if s.greater(s.trick[0], s.trick[1]) {
			return sortValue(c)[0]
		}
		return highestValueWinnerORlowestValueLoser(s, c)
	}
	return highestValueWinnerORlowestValueLoser(s, c)
}

func (p *Player) playerTactic(s *SuitState, c []Card) Card {
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
	cc := countCardsSuit(suit, cards)
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

func (p *Player) calculateHighestBid() int {
	assOtherThan := func(suit string) int {
		asses := 0
		c := in(p.getHand(), Card{CLUBS, "A"})
		s := in(p.getHand(), Card{SPADE, "A"})
		h := in(p.getHand(), Card{HEART, "A"})
		k := in(p.getHand(), Card{CARO, "A"})
		t := in(p.getHand(), Card{suit, "A"})
		if c {
			asses++
		}
		if s {
			asses++
		}
		if h {
			asses++
		}
		if k {
			asses++
		}
		if t {
			asses--
		}
		return asses
	}

	p.highestBid = 0

	suit := mostCardsSuit(p.getHand())

	largest := countTrumpsSuit(suit, p.getHand())

	prob := 0
	if largest > 4 && assOtherThan(suit) > 1 {
		prob = 80
	}
	if largest > 5 {
		prob = 85
	}
	if largest > 6 {
		prob = 99
	}

	est := p.handEstimation()
	debugTacticsLog("(%s) Hand: %v, Estimation: %d\n", p.name, p.hand, est)
	if prob < 80 && est < 45 {
		//	fmt.Printf("LOW %d %v\n", p.handEstimation(), sortSuit(p.getHand()))
		return p.highestBid
	}
	// fmt.Printf("HIGH %d %v\n", p.handEstimation(), sortSuit(p.getHand()))

	trump := p.declareTrump()
	mat := matadors(trump, p.getHand())
	if mat < 0 {
		mat *= -1
	}
	p.highestBid = (mat + 1) * trumpBaseValue(trump)
	return p.highestBid
}

func (p *Player) declareTrump() string {
	return mostCardsSuit(p.getHand())
}

func (p *Player) discardInSkat(skat []Card) {
	debugTacticsLog("FULL HAND %v\n", p.getHand())

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

	// DESPARATE???
	// DISCARD LOW CARDS
	p.setHand(sortSuit(p.declareTrump(), p.getHand()))
	if removed == 1 {
		card := p.getHand()[len(p.getHand())-1]
		p.setHand(remove(p.getHand(), card))
		skat[1] = card
		return
	}
	c1 := p.getHand()[len(p.getHand())-1]
	c2 := p.getHand()[len(p.getHand())-2]
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
