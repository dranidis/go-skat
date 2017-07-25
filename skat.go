package main

import (
	"bufio"
	"fmt"
	"os"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const CLUBS = "CLUBS"
const SPADE = "SPADE"
const HEART = "HEART"
const CARO = "CARO"

var _ = rand.New(rand.NewSource(1))
var r = rand.New(rand.NewSource(time.Now().Unix()))

type Card struct {
	suit string
	rank string
}

func (c Card) equals(o Card) bool {
	return c.suit == o.suit && c.rank == o.rank
}

func (c Card) String() string {
	black := color.New(color.Bold, color.FgWhite).SprintFunc()
	green := color.New(color.Bold, color.FgGreen).SprintFunc()
	red := color.New(color.Bold, color.FgRed).SprintFunc()
	yellow := color.New(color.Bold, color.FgYellow).SprintFunc()
	switch c.suit {
	case CLUBS:
		return black(c.rank)
	case SPADE:
		return green(c.rank)
	case HEART:
		return red(c.rank)
	case CARO:
		return yellow(c.rank)
	}
	return ""
}

func cardValue(c Card) int {
	switch c.rank {
	case "J":
		return 2
	case "A":
		return 11
	case "10":
		return 10
	case "K":
		return 4
	case "D":
		return 3
	}
	return 0
}

func trumpBaseValue(s string) int {
	switch s {
	case CLUBS:
		return 12
	case SPADE:
		return 11
	case HEART:
		return 10
	case CARO:
		return 9
	}
	return 0
}

func getSuite(trump string, card Card) string {
	if card.rank == "J" {
		return trump
	}
	return card.suit
}

var suits = []string{CLUBS, SPADE, HEART, CARO}
var ranks = []string{"J", "A", "10", "K", "D", "9", "8", "7"}

func sortRankSpecial(cs []Card, ranks []string) []Card {
	cards := []Card{}

	for _, r := range ranks {
		for _, s := range suits {
			if in(cs, Card{s, r}) {
				cards = append(cards, Card{s, r})
			}
		}

	}
	return cards
}

func sortRank(cs []Card) []Card {
	return sortRankSpecial(cs, ranks)
}

func sortValue(cs []Card) []Card {
	valueRanks := []string{"A", "10", "K", "D", "J", "7", "8", "9"}
	return sortRankSpecial(cs, valueRanks)
}

func sortSuit(trump string, cs []Card) []Card {
	cards := []Card{}

	cardJs := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
	}
	cardsSuit := func(suit string) []Card {
		return []Card{
			Card{suit, "A"},
			Card{suit, "10"},
			Card{suit, "K"},
			Card{suit, "D"},
			Card{suit, "9"},
			Card{suit, "8"},
			Card{suit, "7"},
		}
	}
	for _, c := range cardJs {
		if in(cs, c) {
			cards = append(cards, c)
		}
	}
	if trump != "" {
		switch trump {
		case CLUBS:
		case SPADE:
			suits = []string{SPADE, CLUBS, HEART, CARO}
		case HEART:
			suits = []string{HEART, CLUBS, SPADE, CARO}
		case CARO:
			suits = []string{CARO, CLUBS, SPADE, HEART}
		}
	}
	for _, s := range suits {
		for _, c := range cardsSuit(s) {
			if in(cs, c) {
				cards = append(cards, c)
			}
		}
	}
	return cards
}

type SuitState struct {
	declarer *Player
	opp1 *Player
	opp2 *Player
	trump  string
	leader *Player
	follow string
	trick  []Card
}

func makeSuitState() SuitState {
	return SuitState{nil, nil, nil, "", nil, "", []Card{}}
}


type tactic func([]Card) Card

type Player struct {
	isHuman	bool
//	isDeclarer bool
	hand       []Card
	highestBid int
	score      int
	schwarz    bool
	totalScore int
	previousSuit string
}

func Shuffle(cards []Card) []Card {
	//r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]Card, len(cards))
	perm := r.Perm(len(cards))
	for i, randIndex := range perm {
		ret[i] = cards[randIndex]
	}
	return ret
}

func sum(trick []Card) int {
	s := 0
	for _, c := range trick {
		s += cardValue(c)
	}
	return s
}

func makeDeck() []Card {
	makeSuitDeck := func(suit string) []Card {
		return []Card{
			Card{suit, "J"},
			Card{suit, "A"},
			Card{suit, "10"},
			Card{suit, "K"},
			Card{suit, "D"},
			Card{suit, "9"},
			Card{suit, "8"},
			Card{suit, "7"},
		}
	}
	cards := []Card{}
	cards = append(cards, makeSuitDeck(CLUBS)...)
	cards = append(cards, makeSuitDeck(SPADE)...)
	cards = append(cards, makeSuitDeck(HEART)...)
	cards = append(cards, makeSuitDeck(CARO)...)
	return cards
}

// CARD MANIPULATION FUNCTIONS
func in(cs []Card, c Card) bool {
	for _, card := range cs {
		if card.equals(c) {
			return true
		}
	}
	return false
}

func filter(cards []Card, f func(Card) bool) []Card {
	cs := []Card{}
	for _, c := range cards {
		if f(c) {
			cs = append(cs, c)
		}
	}
	return cs
}

func remove(cs []Card, c Card) []Card {
	return filter(cs, func(cc Card) bool {
		return !(cc.equals(c))
	})
}

func matadors(trump string, cs []Card) int {
	cards := []Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{trump, "A"},
		Card{trump, "10"},
		Card{trump, "K"},
		Card{trump, "D"},
		Card{trump, "9"},
		Card{trump, "8"},
	}
	m := 0
	if in(cs, Card{CLUBS, "J"}) {
		m++
		for _, card := range cards {
			if !in(cs, card) {
				break
			}
			m++
		}
		return m
	}
	m--
	for _, card := range cards {
		if in(cs, card) {
			return m
		}
		m--
	}
	return m
}

func setNextTrickOrder(s *SuitState, players []*Player) []*Player {
	var newPlayers []*Player
	var winner *Player
	if s.greater(s.trick[0], s.trick[1]) && s.greater(s.trick[0], s.trick[2]) {
		winner = players[0]
		newPlayers = players
	} else if s.greater(s.trick[1], s.trick[2]) {
		winner = players[1]
		newPlayers = []*Player{players[1], players[2], players[0]}
	} else {
		winner = players[2]
		newPlayers = []*Player{players[2], players[0], players[1]}
	}

	winner.score += sum(s.trick)

	if s.declarer != nil && s.opp1 != nil &&  s.opp2 != nil {
		gameLog("TRICK %v : %d points: %d - %d\n", s.trick, sum(s.trick), s.declarer.score, s.opp1.score + s.opp2.score)
	}

	winner.schwarz = false
	s.trick = []Card{}
	s.leader = newPlayers[0]

	return newPlayers
}

func round(s *SuitState, players []*Player) []*Player {
	//fmt.Printf("TRICK: %v\n", s.trick)
	//var trick [3]Card
	players[0].play(s)
	s.follow = getSuite(s.trump, s.trick[0])
	players[1].play(s)
	players[2].play(s)


	//fmt.Println(players)
	players = setNextTrickOrder(s, players)
	//fmt.Println(players)


	s.follow = ""
	return players
}

func highestValueWinnerORlowestValueLoser(s *SuitState, c []Card) Card {
	winners := filter(c, func(c Card) bool {
		wins := true
		for _, t := range s.trick {
			if s.greater(t, c) {
				wins = false
			}
		}
		return wins
	})
	

	if s.trump == s.follow {
		trumpWinnerRanks := []string{"A","10","K","D","9","8","7","J"}
		winners = sortRankSpecial(winners, trumpWinnerRanks)
	} else {
		winners = sortValue(winners)
	}
	//gameLog("Trick: %v Winners: %v\n", s.trick, winners)
	if len(winners) > 0 {
		// return winners[len(winners)-1]
		return winners[0]
	}


	if s.trump == s.follow {
		trumpRanks := []string{"A","10","J","K","D","9","8","7"}
		cards := sortRankSpecial(c, trumpRanks)
		return cards[len(cards)-1]		
	}
	cards := sortValue(c)
	return cards[len(cards)-1]
}

func HighestLong(trump string, c []Card) Card {
	//gameLog("HIGHESTLONG\n")
	s := LongestNonTrumpSuit(trump, c)
	fmt.Println("LongestNonTrumpSuit", s)
	cards := filter(c, func (c Card) bool {
		return c.suit == s && c.rank != "J"
		})
	fmt.Println(cards)
	//fmt.Println(c)
	if len(cards) > 0 {
		return cards[0]
	}
	// last card?
	gameLog("... DEBUG ... VALID: %v no highest long. Returning: %v\n", c, c[0])
	return c[0]
}

func (p Player) declarerTactic(s *SuitState, c []Card) Card {
	//gameLog("SOLIST\n")
	if len(s.trick) == 0 {
		return firstCardTactic(c)
	}
	return highestValueWinnerORlowestValueLoser(s, c)
}

func (p *Player) playerTactic(s *SuitState, c []Card) Card {
	c = sortSuit(s.trump, c)
	//gameLog("Trick: %v\n", s.trick)
	gameLog("Trick: %v,   valid: %v\n", s.trick, c)
	if s.declarer == p {
		gameLog("SOLIST ")
	} 

	if p.isHuman {
		gameLog("Hand : %v\n", p.hand)
		gameLog("Valid: %v\n", c)
		for {
			fmt.Printf("CARD? ")
			// reader := bufio.NewReader(os.Stdin)
			// char, _, err := reader.ReadRune()

			// if err != nil {
			// 	fmt.Println(err)
			// }
	    	var i int
	    	_, err := fmt.Scanf("%d", &i)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if i > len(c) {
				continue
			}    	
	    	return c[i-1]
			//return c[int(char - '0')-1]			
		}

	}

	if s.declarer == p {
		return p.declarerTactic(s, c)
	}
	// OPPONENTS TACTIC

	// if opp plays first high-value card and
	// partner has highest trump then he should play it

	// TODO:
	// if trick has no value, e.g. 7 8
	// don't take it with a trump if you can save it.

	if len(s.trick) == 0 {
		// if you have a card with suit played in a previous trick started from you or your partner continue with it
		// else
		prevSuit := ""
		partnerFunc := func(p *Player) *Player {
			if s.opp1 == p {
				return s.opp2
			}
			return s.opp1
		}
		if p.previousSuit != "" {
			prevSuit = p.previousSuit
		} else if partnerFunc(p).previousSuit != "" {
			prevSuit = partnerFunc(p).previousSuit
		}
		if prevSuit == "" {
			card := HighestLong(s.trump, c)
			p.previousSuit = getSuite(s.trump, card)
			return card
		}
		suitCards := filter(c, func (c Card) bool {
			return c.suit == prevSuit
			})
		if len(suitCards) > 0 {
			fmt.Println("Following previous suit...")
			return suitCards[0]
		}

		card := HighestLong(s.trump, c)
		p.previousSuit = card.suit
		return card
	}

	if len(s.trick) == 1 {
		if s.leader == s.declarer {
			return highestValueWinnerORlowestValueLoser(s, c)
		} else {
			// TODO:
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

func firstCardTactic(c []Card) Card {
	return c[0]
}

func randomCardTactic(c []Card) Card {
	cardIndex := r.Intn(len(c))
	return c[cardIndex]
}

func (p *Player) play(s *SuitState) Card {
	// fmt.Println("\nPLAYER PLAYS")
	valid := validCards(*s, p.hand)
	card := p.playerTactic(s, valid)
	// fmt.Println(s, card, p.hand)
	p.hand = remove(p.hand, card)
	s.trick = append(s.trick, card)
	// fmt.Println(p.hand)
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
	return Player{false, 
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
	if p.isHuman {
		fmt.Printf("HAND: %v", p.hand)
		for {
			fmt.Printf("BID %d? (y/n/q)", bids[bidIndex])
			reader := bufio.NewReader(os.Stdin)
			char, _, err := reader.ReadRune()

			if err != nil {
				fmt.Println(err)
				continue
			}

			switch char {
			case 'y':
				return true
			case 'n':
				return false			
			case 'q':
				os.Exit(0)
			default:
				fmt.Printf("... don't understand! ")
				continue
			}
		}

	}
	return bids[bidIndex] <= p.highestBid
}

//
// Der US-Amerikaner J.P. Wergin hat in seinem Buch "Wergin on Skat and Sheepshead"
// (McFarland, Wisconsin, 1975) versucht, dazu einen objektiven Berechnungsmodus zu
// finden.
func (p *Player) handEstimation() int {
	kreuzB := in(p.hand, Card{CLUBS, "J"})
	pikB := in(p.hand, Card{SPADE, "J"})
	herzB := in(p.hand, Card{HEART, "J"})
	karoB := in(p.hand, Card{CARO, "J"})

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
		a := in(p.hand, Card{suit, "A"})
		t := in(p.hand, Card{suit, "10"})
		k := in(p.hand, Card{suit, "K"})
		d := in(p.hand, Card{suit, "D"})
		n := in(p.hand, Card{suit, "9"})

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
		c := in(p.hand, Card{CLUBS, "A"})
		s := in(p.hand, Card{SPADE, "A"})
		h := in(p.hand, Card{HEART, "A"})
		k := in(p.hand, Card{CARO, "A"})
		t := in(p.hand, Card{suit, "A"})
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

	suit := mostCardsSuit(p.hand)

	largest := countTrumpsSuit(suit, p.hand)

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
		//	fmt.Printf("LOW %d %v\n", p.handEstimation(), sortSuit(p.hand))
		return
	}
	// fmt.Printf("HIGH %d %v\n", p.handEstimation(), sortSuit(p.hand))

	trump := p.declareTrump()
	mat := matadors(trump, p.hand)
	if mat < 0 {
		mat *= -1
	}
	p.highestBid = (mat + 1) * trumpBaseValue(trump)
}

func bidLog(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
func gameLog(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func bidding(listener, speaker *Player, bidIndex int) (int, *Player) {
	for speaker.accepts(bidIndex) {
		bidLog("Bid %d\n", bids[bidIndex])
		if listener.accepts(bidIndex) {
			bidLog("Yes %d\n", bids[bidIndex])
			bidIndex++
		} else {
			bidLog("Listener Pass %d\n", bids[bidIndex])
			return bidIndex, speaker
		}
	}
	bidLog("Speaker Pass %d\n", bids[bidIndex])
	bidIndex--
	return bidIndex, listener
}

func bid(players []*Player) (int, *Player) {
	bidLog("FOREHAND vs MIDDLEHAND\n")
	bidIndex, p := bidding(players[0], players[1], 0)
	bidIndex++
	bidLog("WINNER vs BACKHAND\n")
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

func countTrumpsSuit(suit string, cards []Card) int {
	count := 0
	for _, c := range cards {
		if c.suit == suit || c.rank == "J" {
			count++
		}
	}
	return count
}

func countCardsSuit(suit string, cards []Card) int {
	count := 0
	for _, c := range cards {
		if c.suit == suit && c.rank != "J" {
			count++
		}
	}
	return count
}

func LongestNonTrumpSuit(trump string, cards []Card) string {
	maxI, maxCount := -1, -1
	for i, s := range suits {
		if s == trump {
			continue
		} 
		c := countCardsSuit(s, cards)
		if countCardsSuit(s, cards) > maxCount {
			maxI = i
			maxCount = c
		} 
	}
	return suits[maxI]
}

func mostCardsSuit(cards []Card) string {
	c := countTrumpsSuit(CLUBS, cards)
	s := countTrumpsSuit(SPADE, cards)
	h := countTrumpsSuit(HEART, cards)
	k := countTrumpsSuit(CARO, cards)
	if c > s && c > h && c > k {
		return CLUBS
	}
	if s > h && s > k {
		return SPADE
	}
	if h > k {
		return HEART
	}
	return CARO
}

func lessCardsSuit(cards []Card) string {
	c := countCardsSuit(CLUBS, cards)
	s := countCardsSuit(SPADE, cards)
	h := countCardsSuit(HEART, cards)
	k := countCardsSuit(CARO, cards)

	// we don't want to discard a suit having an A
	if in(cards, Card{CLUBS, "A"}) {
		c = 100
	}
	if in(cards, Card{SPADE, "A"}) {
		s = 100
	}
	if in(cards, Card{HEART, "A"}) {
		h = 100
	}
	if in(cards, Card{CARO, "A"}) {
		k = 100
	}
	if c == 0 {
		c = 9999 // no cards for comparison below
	}
	if s == 0 {
		s = 9999
	}
	if h == 0 {
		h = 9999
	}
	if k == 0 {
		k = 9999
	}
	// fmt.Println("lessCardsSuit", c, s, h, k)
	if c != 0 && c < s && c < h && c < k {
		return CLUBS
	}
	if s != 0 && s < h && s < k {
		return SPADE
	}
	if h != 0 && h < k {
		return HEART
	}
	if k != 0 {
		return CARO
	}
	return ""
}

func (p *Player) declareTrump() string {
	p.hand = sortSuit("", p.hand)
	if p.isHuman {
		fmt.Printf("HAND: %v\n", p.hand)
		for {
			fmt.Printf("TRUMP? (1 for CLUBS, 2 for SPADE, 3 for HEART, 4 for CARO)")
			reader := bufio.NewReader(os.Stdin)
			char, _, err := reader.ReadRune()

			if err != nil {
				fmt.Println(err)
				continue
			}

			switch char {
			case '1':
				return CLUBS
			case '2':
				return SPADE
			case '3':
				return HEART
			case '4':
				return CARO
			default:
				continue
			}			
		}

	}

	return mostCardsSuit(p.hand)
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
	if p.isHuman {
		p.hand = sortSuit("", p.hand)
		gameLog("Full Hand : %v\n", p.hand)
		for {
			fmt.Printf("DISCARD CARDS?")

	    	var i1, i2 int
	    	_, err := fmt.Scanf("%d", &i1)
			if err != nil {
				fmt.Println(err)
				continue
			}  
			 _, err = fmt.Scanf("%d", &i2)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//fmt.Println(i1, i2)
			if i1 > len(p.hand) || i2 > len(p.hand) || i1 == i2 {
				continue
			}
			card1 :=  p.hand[i1-1]   	
			card2 :=  p.hand[i2-1]   	
			p.hand = remove(p.hand, card1)
			p.hand = remove(p.hand, card2)
			skat[0] = card1
			skat[1] = card2
			return			
		}

	}
	// fmt.Printf("FULL HAND %v\n", p.hand)

	// discard BLANKS

	bcards := findBlankCards(p.hand)
	// fmt.Printf("BLANK %v\n", bcards)
	removed := 0
	if len(bcards) > 0 {
		p.hand = remove(p.hand, bcards[0])
		skat[0] = bcards[0]
		//	fmt.Printf("1st %v\n", skat)
		removed++
	}
	if len(bcards) > 1 {
		p.hand = remove(p.hand, bcards[1])
		skat[1] = bcards[1]
		//	fmt.Printf("2nd %v\n", skat)
		return
	}
	// Discard high cards in non-A suits with few colors
	sranks := []string{"J", "A", "10", "K", "D", "7", "8", "9"}

	lsuit := lessCardsSuit(p.hand)
	if lsuit != "" {
		lcards := sortRankSpecial(filter(p.hand, func(c Card) bool {
			return c.suit == lsuit && c.rank != "A" && c.rank != "J"
		}), sranks)
		if len(lcards) < 4 { // do not throw long fleets
			// fmt.Printf("SUIT %v LESS %v\n", lsuit, lcards)

			if len(lcards) > 1 {
				i := 0
				for removed < 2 {
					p.hand = remove(p.hand, lcards[i])
					skat[removed] = lcards[i]
					i++
					removed++
				}
				return
			}
		}
	}

	// Discard non-A-10 suit cards
	ncards := nonA10cards(p.hand)
	ncards = findBlankCards(ncards)
	// fmt.Printf("nonA10cards %v\n", ncards)

	if len(ncards) > 1 {
		i := 0
		for removed < 2 {
			p.hand = remove(p.hand, ncards[i])
			skat[removed] = ncards[i]
			i++
			removed++
		}
		return
	}

	if len(ncards) == 1 {
		p.hand = remove(p.hand, ncards[0])
		skat[removed] = ncards[0]
		removed++

		if removed == 2 {
			return
		}
	}

	// DESPARATE???
	// DISCARD LOW CARDS
	p.hand = sortSuit(p.declareTrump(), p.hand)
	if removed == 1 {
		card := p.hand[len(p.hand)-1]
		p.hand = remove(p.hand, card)
		skat[1] = card
		return
	}
	c1 := p.hand[len(p.hand)-1]
	c2 := p.hand[len(p.hand)-2]
	p.hand = remove(p.hand, c1)
	p.hand = remove(p.hand, c2)
	skat[0] = c1
	skat[1] = c2
}

func (p *Player) pickUpSkat(skat []Card) bool {
	if p.isHuman {
		fmt.Printf("HAND: %v", p.hand)
		yes := false
		for !yes {
			fmt.Printf("Pick up SKAT? (y/n/q) ")
			reader := bufio.NewReader(os.Stdin)
			char, _, err := reader.ReadRune()

			if err != nil {
				fmt.Println(err)
				continue
			}

			switch char {
			case 'y':
				yes = true
			case 'n':
				return false			
			case 'q':
				os.Exit(0)
			default:
				fmt.Printf("... don't understand! ")
				continue
			}
		}
	}
	hand := make([]Card, 10)
	copy(hand, p.hand)
	hand = append(hand, skat...)
	p.hand = hand

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

func game(players []*Player) (int, int) {
	//fmt.Println("------------NEW GAME----------")
	// DEALING
	cards := Shuffle(makeDeck())
	//fmt.Printf("CARDS %d %v\n", -1, cards)
	players[0].hand = sortSuit("", cards[:10])
	players[1].hand = sortSuit("", cards[10:20])
	players[2].hand = sortSuit("", cards[20:30])

	if players[0].isHuman {
		fmt.Printf("%v\n", players[0].hand)
	}

	skat := make([]Card, 2)
	copy(skat, cards[30:32])

	for _, p := range players {
		if !p.isHuman {
			p.calculateHighestBid()
		}
	}

	// BIDDING
	bidIndex, declarer := bid(players)
	if bidIndex == -1 {
		// fmt.Println("ALL PASSED")
		return 0, 0
	}
	var opp1, opp2 *Player
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
	// fmt.Printf("\nHAND bef: %v\n", sortSuit(declarer.hand))
	// fmt.Printf("SKAT bef: %v\n", skat)

	if declarer.pickUpSkat(skat) {
		// fmt.Printf("HAND aft: %v\n", sortSuit(declarer.hand))
		handGame = false
		// fmt.Printf("SKAT aft: %v\n", skat)
	}

	// DECLARE
	state := SuitState{declarer, opp1, opp2, declarer.declareTrump(), players[0], "", []Card{}}
	players[0].hand = sortSuit(state.trump, players[0].hand)
	players[1].hand = sortSuit(state.trump, players[1].hand)
	players[2].hand = sortSuit(state.trump, players[2].hand)

	gameLog("TRUMP: %v\n", state.trump)
	declarerCards := make([]Card, len(declarer.hand))
	copy(declarerCards, declarer.hand)
	declarerCards = append(declarerCards, skat...)

	// fmt.Println("DECLARER Hand after SKAT: %v" , declarer.hand)

	// gameLog("BID: %d, SUIT: %d %s",
	// 	bids[bidIndex], countTrumpsSuit(state.trump, declarer.hand), state.trump)

	// PLAY
	for i := 0; i < 10; i++ {
		gameLog("\n")
		players = round(&state, players)
	}
	// gameLog("SKAT: %v, %d\n", skat, sum(skat))
	declarer.score += sum(skat)

	gs := gameScore(state, declarerCards, declarer.score, bids[bidIndex],
		declarer.schwarz, opp1.schwarz && opp2.schwarz, handGame)

	declarer.totalScore += gs

	if declarer.score > 60 {
		fmt.Printf(" VICTORY: %d - %d, SCORE: %d\n",
		declarer.score, opp1.score + opp2.score, gs)
	} else {
		fmt.Printf(" LOSS: %d - %d, SCORE: %d\n",
		declarer.score, opp1.score + opp2.score, gs)
	}

	return declarer.score, opp1.score + opp2.score

}

func rotatePlayers(players []*Player) []*Player {
	newPlayers := []*Player{}
	newPlayers = append(newPlayers, players[2])
	newPlayers = append(newPlayers, players[0])
	newPlayers = append(newPlayers, players[1])
	return newPlayers
}

func main() {
	player1 := makePlayer([]Card{})
	player2 := makePlayer([]Card{})
	player3 := makePlayer([]Card{})

	players := []*Player{&player1, &player2, &player3}
	players[0].isHuman = true

	passed := 0
	won := 0
	lost := 0
	totalGames := 9
	for times := totalGames; times > 0; times-- {
		for _, p := range players {
			p.score = 0
			p.schwarz = true
			p.previousSuit = ""
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
		fmt.Println(player1.totalScore, player2.totalScore, player3.totalScore)
		//time.Sleep(1000 * time.Millisecond)
		players = rotatePlayers(players)
	}
	avg := float64(player1.totalScore+player2.totalScore+player3.totalScore) / float64(totalGames-passed)
	fmt.Printf("AVG %3.1f, passed %d, won %d, lost %d\n", avg, passed, won, lost)
}
