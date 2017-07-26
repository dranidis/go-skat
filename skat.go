package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
)

var logFile io.Writer = nil
var debugTacticsLogFlag = false
var gameLogFlag = true

func logToFile(format string, a ...interface{}) {
	if logFile != nil {
		fmt.Fprintf(logFile, format, a...)
	}
}
func bidLog(format string, a ...interface{}) {
	if gameLogFlag {
		fmt.Printf(format, a...)
	}
	logToFile(format, a...)
}

func gameLog(format string, a ...interface{}) {
	if gameLogFlag {
		fmt.Printf(format, a...)
	}
	logToFile(format, a...)
}

func debugTacticsLog(format string, a ...interface{}) {
	if debugTacticsLogFlag {
		fmt.Printf(format, a...)
	}
	logToFile(format, a...)
}

type SuitState struct {
	declarer PlayerI
	opp1     PlayerI
	opp2     PlayerI
	trump    string
	leader   PlayerI
	follow   string
	trick    []Card
	// not necessary for game but for tactics
	trumpsInGame []Card
}

func makeSuitState() SuitState {
	return SuitState{nil, nil, nil, "", nil, "", []Card{}, []Card{}}
}

type Player struct {
	name	string
	hand         []Card
	highestBid   int
	score        int
	schwarz      bool
	totalScore   int
	previousSuit string
}

func (p * Player) incTotalScore(s int) {
	p.totalScore += s
}

func (p * Player) setHand(cs []Card) {
	p.hand = cs
}

func (p * Player) setScore(s int) {
	p.score = s
}

func (p * Player) setHuman(b bool) {
}

func (p * Player) setSchwarz(b bool) {
	p.schwarz = b
}
func (p * Player) setPreviousSuit(s string) {
	p.previousSuit = s
}

func (p * Player) getScore() int {
	return p.score
}

func (p * Player) getPreviousSuit() string {
	return p.previousSuit
}

func (p * Player) getTotalScore() int {
	return p.totalScore
}

func (p * Player) setName(n string)  {
	p.name = n
}

func (p * Player) getName() string {
	return p.name
}

func (p * Player) getHand() []Card {
	return p.hand
}

func (p * Player) isHuman() bool {
	return false
}

func (p * Player) isSchwarz() bool {
	return p.schwarz
}

func setNextTrickOrder(s *SuitState, players []PlayerI) []PlayerI {
	var newPlayers []PlayerI
	var winner PlayerI
	if s.greater(s.trick[0], s.trick[1]) && s.greater(s.trick[0], s.trick[2]) {
		winner = players[0]
		newPlayers = players
	} else if s.greater(s.trick[1], s.trick[2]) {
		winner = players[1]
		newPlayers = []PlayerI{players[1], players[2], players[0]}
	} else {
		winner = players[2]
		newPlayers = []PlayerI{players[2], players[0], players[1]}
	}

	winner.setScore(winner.getScore() + sum(s.trick))

	if s.declarer != nil && s.opp1 != nil && s.opp2 != nil {
		gameLog("TRICK %v\n", s.trick)
		debugTacticsLog("%d points: %d - %d\n", sum(s.trick), s.declarer.getScore(), s.opp1.getScore()+s.opp2.getScore())
	}

	winner.setSchwarz(false)
	s.trick = []Card{}
	s.leader = newPlayers[0]

	return newPlayers
}

func round(s *SuitState, players []PlayerI) []PlayerI {
	play(s, players[0])
	s.follow = getSuite(s.trump, s.trick[0])
	play(s, players[1])
	play(s, players[2])

	players = setNextTrickOrder(s, players)
	s.follow = ""
	return players
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

func (p Player) otherPlayersHaveJs(s *SuitState) bool {
	for _, suit := range suits {
		card := Card{suit, "J"}
		if in(s.trumpsInGame, card) && !in(p.getHand(), card) {
			return true
		}
	}
	return false
}

func (p Player) otherPlayersTrumps(s * SuitState) []Card {
	return filter(makeDeck(), func (c Card) bool {
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
			debugTacticsLog(" -- Declarer leads ")
		// MIDDLEHAND
			// if declarer leads a low trump, and there are still higher trumps
			// smear it with a high value
			if getSuite(s.trump, s.trick[0]) == s.trump && len(winnerCards(s, c)) == 0 {
				other := p.otherPlayersTrumps(s)
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

func firstCardTactic(c []Card) Card {
	return c[0]
}

func randomCardTactic(c []Card) Card {
	cardIndex := r.Intn(len(c))
	return c[cardIndex]
}

func play(s *SuitState, p PlayerI) Card {
	valid := sortSuit(s.trump, validCards(*s, p.getHand()))

	p.setHand(sortSuit(s.trump, p.getHand()))
	gameLog("Trick: %v\n", s.trick)
	debugTacticsLog("valid: %v\n", valid)
	if s.opp1 != nil && s.opp2 != nil {
	debugTacticsLog("Previous suit: %s:%v, %s:%v\n", 
			s.opp1.getName(), s.opp1.getPreviousSuit(),
			s.opp2.getName(), s.opp2.getPreviousSuit())		
	}
	if s.declarer == p {
		red := color.New(color.Bold, color.FgRed).SprintFunc()
		gameLog("(%s) ", red(p.getName()))
	} else {
		gameLog("(%v) ", p.getName())
	}
	card := p.playerTactic(s, valid)
	p.setHand(remove(p.getHand(), card))
	s.trick = append(s.trick, card)
	if getSuite(s.trump, card) == s.trump {
		s.trumpsInGame = remove(s.trumpsInGame, card)
	}
	return card
}

// Returns a list of all cards that are playeable from the player's hand.
func validCards(s SuitState, playerHand []Card) []Card {
	return filter(playerHand, func(c Card) bool {
		return s.valid(playerHand, c)
	})
}

func (s SuitState) valid(playerHand []Card, card Card) bool {
	for _, c := range playerHand {
		// if there is at least one card in your hand matching the followed suit
		// your played card should follow
		if s.follow == getSuite(s.trump, c) {
			return s.follow == getSuite(s.trump, card)
		}
	}
	// otherwise any card is playable
	return true
}

func (s SuitState) greater(card1, card2 Card) bool {
	rank := map[string]int{
		"A":  13,
		"10": 12,
		"K":  11,
		"D":  10,
		"9":  9,
		"8":  8,
		"7":  7,
	}
	JRank := map[string]int{
		CLUBS: 4,
		SPADE: 3,
		HEART: 2,
		CARO:  1,
	}

	if card1.rank == "J" {
		if card2.rank == "J" {
			return JRank[card1.suit] > JRank[card2.suit]
		}
		return true
	}

	if card2.rank == "J" {
		return false
	}

	if card1.suit == s.trump {
		if card2.suit == s.trump {
			return rank[card1.rank] > rank[card2.rank]
		}
		return true
	}

	if card2.suit == s.trump {
		return false
	}

	if card1.suit == s.follow {
		if card2.suit == s.follow {
			return rank[card1.rank] > rank[card2.rank]
		}
		return true
	}

	if card2.suit == s.follow {
		return false
	}

	return rank[card1.rank] > rank[card2.rank]
}

func makePlayer(hand []Card) Player {
	return Player{"dummy", 
	//false,
		// false,
		hand, 0, 0, true, 0, ""}
}

var bids = []int{
	18, 20, 22, 23, 24,
	27, 30, 33, 35, 36,
	40, 44, 45, 46, 48, 50,
	54, 55, 59, 60,
	63, 66, 70, 72, 77,
	80, 81, 84, 88, 90, 96, 99, 100, 108, 110, 117,
	121, 126, 130, 132, 135, 140, 143, 144,
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
	if !kreuzB && pikB && herzB && !karoB {
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

func (p *Player) calculateHighestBid() {
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

	if prob < 80 && p.handEstimation() < 45 {
		//	fmt.Printf("LOW %d %v\n", p.handEstimation(), sortSuit(p.getHand()))
		return
	}
	// fmt.Printf("HIGH %d %v\n", p.handEstimation(), sortSuit(p.getHand()))

	trump := p.declareTrump()
	mat := matadors(trump, p.getHand())
	if mat < 0 {
		mat *= -1
	}
	p.highestBid = (mat + 1) * trumpBaseValue(trump)
}

func bidding(listener, speaker PlayerI, bidIndex int) (int, PlayerI) {
	for speaker.accepts(bidIndex) {
		bidLog("Bid %d\n", bids[bidIndex])
		if listener.accepts(bidIndex) {
			bidLog("Yes %d\n", bids[bidIndex])
			bidIndex++
		} else {
			bidLog("Listener (%v) Pass %d\n", listener.getName(), bids[bidIndex])
			return bidIndex, speaker
		}
	}
	bidLog("Speaker (%v) Pass %d\n", speaker.getName(), bids[bidIndex])
	bidIndex--
	return bidIndex, listener
}

func bid(players []PlayerI) (int, PlayerI) {
	bidLog("FOREHAND (%v) vs MIDDLEHAND (%v)\n", players[0].getName(), players[1].getName())
	bidIndex, p := bidding(players[0], players[1], 0)
	bidIndex++
	bidLog("WINNER (%v) vs BACKHAND (%v)\n", p.getName(), players[2].getName())
	bidIndex, p = bidding(p, players[2], bidIndex)
	if bidIndex == -1 {
		if players[0].accepts(0) {
			bidLog("Yes %d\n", bids[0])
			return 0, players[0]
		} else {
			bidLog("Listener Pass %d\n", bids[0])
			return -1, nil
		}
	}
	//	p.isDeclarer = true
	return bidIndex, p
}


func (p *Player) declareTrump() string {
	return mostCardsSuit(p.getHand())
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

func (p *Player) discardInSkat(skat []Card) {
	// fmt.Printf("FULL HAND %v\n", p.getHand())

	// discard BLANKS

	bcards := findBlankCards(p.getHand())
	// fmt.Printf("BLANK %v\n", bcards)
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
			// fmt.Printf("SUIT %v LESS %v\n", lsuit, lcards)

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
	hand := make([]Card, 10)
	copy(hand, p.getHand())
	hand = append(hand, skat...)
	p.setHand(hand)

	p.discardInSkat(skat)

	return true
}

func gameScore(state SuitState, cs []Card, score, bid int,
	decSchwarz, oppSchwarz, handGame bool) int {
	mat := matadors(state.trump, cs)
	if mat < 0 {
		mat = mat * -1
	}
	multiplier := mat + 1

	gameLog("With %d ", mat)

	base := trumpBaseValue(state.trump)

	if handGame {
		multiplier++
		gameLog("Hand ")
	}
	// Schneider?
	if score > 89 || score < 31 {
		multiplier++
		gameLog("Schneider ")
	}

	if decSchwarz || oppSchwarz {
		multiplier++
		gameLog("Schwarz ")
	}
	gameLog("\n")
	gs := multiplier * base

	// OVERBID?
	if gs < bid {
		gameLog("OVERBID!!! Game Value: %d < Bid: %d", gs, bid)
		leastMult := 0
		for leastMult*base < bid {
			leastMult++
		}
		return -2 * leastMult * base
	}

	if score > 60 {
		return gs
	} else {
		return -2 * gs
	}
}

func game(players []PlayerI) (int, int) {
	//fmt.Println("------------NEW GAME----------")
	// DEALING
	cards := Shuffle(makeDeck())
	players[0].setHand(sortSuit("", cards[:10]))
	players[1].setHand(sortSuit("", cards[10:20]))
	players[2].setHand(sortSuit("", cards[20:30]))

	for _, p := range players {
		debugTacticsLog("(%v) hand: %v\n", p.getName(), p.getHand())
	}
	if players[0].isHuman() {
		gameLog("%v\n", players[0].getHand())
	}

	skat := make([]Card, 2)
	copy(skat, cards[30:32])

	for _, p := range players {
		if !p.isHuman() {
			p.calculateHighestBid()
		}
	}

	// BIDDING
	bidIndex, declarer := bid(players)
	if bidIndex == -1 {
		gameLog("ALL PASSED\n")
		return 0, 0
	}
	var opp1, opp2 PlayerI
	if declarer == players[0] {
		opp1, opp2 = players[1], players[2]
	}
	if declarer == players[1] {
		opp1, opp2 = players[0], players[2]
	}
	if declarer == players[2] {
		opp1, opp2 = players[0], players[1]
	}

	// HAND GAME?
	handGame := true
	// fmt.Printf("\nHAND bef: %v\n", sortSuit(declarer.getHand()))
	// fmt.Printf("SKAT bef: %v\n", skat)

	if declarer.pickUpSkat(skat) {
		// fmt.Printf("HAND aft: %v\n", sortSuit(declarer.getHand()))
		handGame = false
		// fmt.Printf("SKAT aft: %v\n", skat)
	}

	trump := declarer.declareTrump()
	allTrumps := filter(makeDeck(), func(c Card) bool {
		return getSuite(trump, c) == trump
	})
	// DECLARE
	state := SuitState{
		declarer, opp1, opp2,
		trump,
		players[0],
		"",
		[]Card{},
		allTrumps,
	}
	players[0].setHand(sortSuit(state.trump, players[0].getHand()))
	players[1].setHand(sortSuit(state.trump, players[1].getHand()))
	players[2].setHand(sortSuit(state.trump, players[2].getHand()))

	gameLog("TRUMP: %v\n", state.trump)
	declarerCards := make([]Card, len(declarer.getHand()))
	copy(declarerCards, declarer.getHand())
	declarerCards = append(declarerCards, skat...)

	// fmt.Println("DECLARER Hand after SKAT: %v" , declarer.getHand())

	// gameLog("BID: %d, SUIT: %d %s",
	// 	bids[bidIndex], countTrumpsSuit(state.trump, declarer.getHand()), state.trump)

	// PLAY
	for i := 0; i < 10; i++ {
		debugTacticsLog("TRUMPS IN PLAY %v\n", state.trumpsInGame)
		gameLog("\n")
		players = round(&state, players)
	}
	// gameLog("SKAT: %v, %d\n", skat, sum(skat))
	declarer.setScore(declarer.getScore() + sum(skat))

	gs := gameScore(state, declarerCards, declarer.getScore(), bids[bidIndex],
		declarer.isSchwarz(), opp1.isSchwarz() && opp2.isSchwarz(), handGame)

	declarer.incTotalScore(gs)

	if declarer.getScore() > 60 {
		gameLog(" VICTORY: %d - %d, SCORE: %d\n",
			declarer.getScore(), opp1.getScore()+opp2.getScore(), gs)
	} else {
		gameLog(" LOSS: %d - %d, SCORE: %d\n",
			declarer.getScore(), opp1.getScore()+opp2.getScore(), gs)
	}

	return declarer.getScore(), opp1.getScore() + opp2.getScore()

}

func rotatePlayers(players []PlayerI) []PlayerI {
	newPlayers := []PlayerI{}
	newPlayers = append(newPlayers, players[2])
	newPlayers = append(newPlayers, players[0])
	newPlayers = append(newPlayers, players[1])
	return newPlayers
}

func main() {
	file, err := os.Create("gameLog.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	logFile = file
	defer file.Close()

	player1 := makeHumanPlayer([]Card{})
	player2 := makePlayer([]Card{})
	player3 := makePlayer([]Card{})

	players := []PlayerI{&player1, &player2, &player3}
	players[0].setHuman(true)
	players[0].setName("You")
	players[1].setName("Bob")
	players[2].setName("Ana")

	passed := 0
	won := 0
	lost := 0
	totalGames := 9
	for times := totalGames; times > 0; times-- {
		for _, p := range players {
			p.setScore(0)
			p.setSchwarz(true)
			p.setPreviousSuit("")
		}
		score, oppScore := game(players)
		if score == 0 && oppScore == 0 {
			passed++
		}
		if score > 60 {
			won++
		} else {
			lost++
		}
		fmt.Println(player1.getTotalScore(), player2.getTotalScore(), player3.getTotalScore())
		//time.Sleep(1000 * time.Millisecond)
		players = rotatePlayers(players)
	}
	avg := float64(player1.getTotalScore()+player2.getTotalScore()+player3.getTotalScore()) / float64(totalGames-passed)
	fmt.Printf("AVG %3.1f, passed %d, won %d, lost %d\n", avg, passed, won, lost)
}
