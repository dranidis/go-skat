package main

import (
	"fmt"
	"math/rand"
	"time"
)

const CLUBS = "CLUBS"
const SPADE = "SPADE"
const HEART = "HEART"
const CARO = "CARO"

type Card struct {
	suit string
	rank string
}

type SuitState struct {
	trump  string
	follow string
}

type tactic func([]Card) Card

type Player struct {
	hand         []Card
	playerTactic tactic
	highestBid   int
	score        int
	schwarz      bool
	totalScore   int
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
	val := func(c Card) int {
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
	s := 0
	for _,c := range trick {
		s += val(c)
	}
	return s
}

func inHand(cs []Card, c Card) bool {
	for _, card := range cs {
		if card.suit == c.suit && card.rank == c.rank {
			return true
		}
	}
	return false
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
	if inHand(cs, Card{CLUBS, "J"}) {
		m++
		for _, card := range cards {
			if !inHand(cs, card) {
				break
			}
			m++
		}
		return m
	}
	m--
	for _, card := range cards {
		if inHand(cs, card) {
			return m
		}
		m--
	}
	return m
}

func setNextTrickOrder(s *SuitState, players []*Player, trick []Card) []*Player {
	if s.greater(trick[0], trick[1]) && s.greater(trick[0], trick[2]) {
		players[0].score += sum(trick)
		players[0].schwarz = false
		return players
	}
	if s.greater(trick[1], trick[2]) {
		players[1].score += sum(trick)
		players[1].schwarz = false
		return []*Player{players[1], players[2], players[0]}
	}
	players[2].score += sum(trick)
	players[2].schwarz = false
	return []*Player{players[2], players[0], players[1]}
}

func round(s *SuitState, players []*Player) []*Player {
	var trick [3]Card
	trick[0] = players[0].play(s)
	s.follow = getSuite(*s, trick[0])
	trick[1] = players[1].play(s)
	trick[2] = players[2].play(s)

	//fmt.Println(players)
	players = setNextTrickOrder(s, players, trick[:])
	//fmt.Println(players)

	// fmt.Printf("TRICK %v : %d\n", trick, sum(trick[:]))

	s.follow = ""
	return players
}

//var r = rand.New(rand.NewSource(99))
var r = rand.New(rand.NewSource(time.Now().Unix()))

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
	card := p.playerTactic(valid)
	// fmt.Println(s, card, p.hand)
	p.hand = remove(p.hand, card)
	// fmt.Println(p.hand)
	return card
}

// func remove(cs []Card, c Card) []Card {
// 	for i, cc := range cs {
// 		if cc.suit == c.suit && cc.rank == c.rank {
// 			front := cs[:i]
// 			back := cs[i+1:]
// 			return append(front, back...)
// 		}
// 	}
// 	return cs
// }
func remove(cs []Card, c Card) []Card {
	cards := []Card{}
	for _, cc := range cs {
		if !(cc.suit == c.suit && cc.rank == c.rank) {
			cards = append(cards, cc)
		}
	}
	return cards
}

func getSuite(s SuitState, card Card) string {
	if card.rank == "J" {
		return s.trump
	}
	return card.suit
}

// Returns a list of all cards that are playeable from the player's hand.
func validCards(s SuitState, playerHand []Card) []Card {
	cards := []Card{}
	for _, c := range playerHand {
		if s.valid(playerHand, c) {
			cards = append(cards, c)
		}
	}
	return cards
}

func (s SuitState) valid(playerHand []Card, card Card) bool {
	for _, c := range playerHand {
		if s.follow == getSuite(s, c) {
			return s.follow == getSuite(s, card)
		}
	}
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

func makeSuitDeck(suit string) []Card {
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

func makeDeck() []Card {
	cards := []Card{}
	cards = append(cards, makeSuitDeck(CLUBS)...)
	cards = append(cards, makeSuitDeck(SPADE)...)
	cards = append(cards, makeSuitDeck(HEART)...)
	cards = append(cards, makeSuitDeck(CARO)...)
	return cards
}

func makePlayer(hand []Card) Player {
	return Player{hand, firstCardTactic, 0, 0, true, 0}
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
avg: -25 with random play
*/
func (p *Player) accepts(bidIndex int) bool {
	return bids[bidIndex] <= p.highestBid
}

//
// Der US-Amerikaner J.P. Wergin hat in seinem Buch "Wergin on Skat and Sheepshead"
// (McFarland, Wisconsin, 1975) versucht, dazu einen objektiven Berechnungsmodus zu
// finden.
func (p *Player) handEstimation() int {
	kreuzB := inHand(p.hand, Card{CLUBS, "J"})
	pikB := inHand(p.hand, Card{SPADE, "J"})
	herzB := inHand(p.hand, Card{HEART, "J"})
	karoB := inHand(p.hand, Card{CARO, "J"})

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

	wert += p.otherCardsEstimation(CLUBS)
	wert += p.otherCardsEstimation(SPADE)
	wert += p.otherCardsEstimation(HEART)
	wert += p.otherCardsEstimation(CARO)

	return wert
}

func (p *Player) otherCardsEstimation(suit string) int {
	a := inHand(p.hand, Card{suit, "A"})
	t := inHand(p.hand, Card{suit, "10"})
	k := inHand(p.hand, Card{suit, "K"})
	d := inHand(p.hand, Card{suit, "D"})
	n := inHand(p.hand, Card{suit, "9"})

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

func sort(cs []Card) []Card {
	cards := []Card{}

	cardJs := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
	}
	suits := []string{CLUBS, SPADE, HEART, CARO}
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
		if inHand(cs, c) {
			cards = append(cards, c)
		}
	}

	for _, s := range suits {
		for _, c := range cardsSuit(s) {
			if inHand(cs, c) {
				cards = append(cards, c)
			}
		}
	}
	return cards
}

func (p *Player) calculateHighestBid() {
	assOtherThan := func(suit string) int {
		asses := 0
		c := inHand(p.hand, Card{CLUBS, "A"})
		s := inHand(p.hand, Card{SPADE, "A"})
		h := inHand(p.hand, Card{HEART, "A"})
		k := inHand(p.hand, Card{CARO, "A"})
		t := inHand(p.hand, Card{suit, "A"})
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

	largest := countCardsSuit(suit, p.hand)

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
	//	fmt.Printf("LOW %d %v\n", p.handEstimation(), sort(p.hand))
		return
	}
	// fmt.Printf("HIGH %d %v\n", p.handEstimation(), sort(p.hand))

	trump := p.declareTrump()
	mat := matadors(trump, p.hand)
	if mat < 0 {
		mat *= -1
	}
	p.highestBid = (mat + 1) * trumpBaseValue(trump)
}

func bidLog(format string, a ...interface{}) {
	//	fmt.Printf(format, a)
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
	return bidIndex, p
}

func countCardsSuit(suit string, cards []Card) int {
	count := 0
	for _, c := range cards {
		if c.suit == suit || c.rank == "J" {
			count++
		}
	}
	return count
}

func mostCardsSuit(cards []Card) string {
	c := countCardsSuit(CLUBS, cards)
	s := countCardsSuit(SPADE, cards)
	h := countCardsSuit(HEART, cards)
	k := countCardsSuit(CARO, cards)
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

func (p *Player) declareTrump() string {
	return mostCardsSuit(p.hand)
}

func (p *Player) pickUpSkat(skat []Card) bool {
	return false
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

func gameScore(state SuitState, cs []Card, score, bid int,
	decSchwarz, oppSchwarz, handGame bool) int {
	mat := matadors(state.trump, cs)
	if mat < 0 {
		mat = mat * -1
	}
	multiplier := mat + 1
	base := trumpBaseValue(state.trump)

	if handGame {
		multiplier++
	}
	// Schneider?
	if score > 89 || score < 31 {
		multiplier++
	}

	if decSchwarz || oppSchwarz {
		multiplier++
	}

	gs := multiplier * base

	// OVERBID?
	if gs < bid {
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

func game(players []*Player) bool {
	//fmt.Println("------------NEW GAME----------")
	// DEALING
	cards := Shuffle(makeDeck())
	//fmt.Printf("CARDS %d %v\n", -1, cards)
	players[0].hand = cards[:10]
	players[1].hand = cards[10:20]
	players[2].hand = cards[20:30]
	skat := cards[30:32]

	// fmt.Printf("HAND 1 %v\n", players[0].hand)
	// fmt.Printf("HAND 2 %v\n", players[1].hand)
	// fmt.Printf("HAND 3 %v\n", players[2].hand)
	// fmt.Printf("SKAT %v\n", skat)

	sumAll := sum(players[0].hand) + sum(players[1].hand) + sum(players[2].hand) + sum(skat)
	if sumAll != 120 {
		fmt.Printf("DEAL PROBLEM: %d", sumAll)
	}
	//fmt.Printf("HAND %d %v\n", 0, players[0].hand)
	//fmt.Printf("HAND %d %v\n", 1, players[1].hand)
	//fmt.Printf("HAND %d %v\n", 2, players[2].hand)
	for _, p := range players {
		p.calculateHighestBid()
	}


	// BIDDING
	bidIndex, declarer := bid(players)
	if bidIndex == -1 {
		// fmt.Println("ALL PASSED")
		return false
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
	if declarer.pickUpSkat(skat) {
		handGame = false
	}

	// DECLARE
	state := SuitState{declarer.declareTrump(), ""}
	declarerCards := make([]Card, len(declarer.hand))
	copy(declarerCards, declarer.hand)
	declarerCards = append(declarerCards, skat...)

	// fmt.Println("DECLARER Hand after SKAT: %v" , declarer.hand)

	// fmt.Printf("BID: %d, SUIT: %d %s",
	// 	bids[bidIndex], countCardsSuit(state.trump, declarer.hand), state.trump)


	// fmt.Printf("HAND 1 %v\n", players[0].hand)
	// fmt.Printf("HAND 2 %v\n", players[1].hand)
	// fmt.Printf("HAND 3 %v\n", players[2].hand)
	// fmt.Printf("SKAT %v\n", skat)
	// PLAY
	for i := 0; i < 10; i++ {
		players = round(&state, players)
	}
	// fmt.Printf("SKAT: %v, %d\n", skat, sum(skat))
	declarer.score += sum(skat)
 
	gs := gameScore(state, declarerCards, declarer.score, bids[bidIndex],
		declarer.schwarz, opp1.schwarz && opp2.schwarz, handGame)

	declarer.totalScore += gs

	if declarer.score > 60 {
			// fmt.Printf(" VICTORY: %d - %d, SCORE: %d\n", 
				// declarer.score, opp1.score + opp2.score, gs)
	} else {
			// fmt.Printf(" LOSS: %d - %d, SCORE: %d\n", 
				// declarer.score, opp1.score + opp2.score, gs)
	}

	return true

}

func main() {
	player1 := makePlayer([]Card{})
	player2 := makePlayer([]Card{})
	player3 := makePlayer([]Card{})

	players := []*Player{&player1, &player2, &player3}

	passed := 0
	totalGames := 3600
	for times := totalGames; times > 0; times-- {
		for _, p := range players {
			p.score = 0
			p.schwarz = true
		}
		if !game(players) {
			passed++
		}
		fmt.Println(player1.totalScore, player2.totalScore, player3.totalScore)
		//time.Sleep(1000 * time.Millisecond)
	}
	avg := float64(player1.totalScore+player2.totalScore+player3.totalScore) / float64(totalGames-passed)
	fmt.Printf("AVG %3.1f, passed %d\n", avg, passed)
}
