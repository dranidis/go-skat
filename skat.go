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
	accepts      func(int) bool
	score        int
	schwarz	bool
	totalScore int
}

func Shuffle(cards []Card) []Card {
	r := rand.New(rand.NewSource(time.Now().Unix()))
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
	return val(trick[0]) + val(trick[1]) + val(trick[2])
}

func inHand(cs []Card, c Card) bool {
	for _, card := range cs {
		if card.suit == c.suit && card.rank == c.rank {
			return true
		}
	}
	return false
}
func matadors(s SuitState, cs []Card) int {
	cards := []Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{s.trump, "A"},
		Card{s.trump, "10"},
		Card{s.trump, "K"},
		Card{s.trump, "D"},
		Card{s.trump, "9"},
		Card{s.trump, "8"},
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

	s.follow = ""
	return players
}

//var r = rand.New(rand.NewSource(99))
var	r = rand.New(rand.NewSource(time.Now().Unix()))

func firstCardTactic(c []Card) Card {
	return c[0]
}

func randomCardTactic(c []Card) Card {
	cardIndex := r.Intn(len(c))
	return c[cardIndex]
}

func (p *Player) play(s *SuitState) Card {
	valid := validCards(*s, p.hand)
	card := p.playerTactic(valid)
	//fmt.Println(s, card, p.hand)
	p.hand = remove(p.hand, card)
	//fmt.Println(p)
	return card
}

func remove(cs []Card, c Card) []Card {
	for i, cc := range cs {
		if cc.suit == c.suit && cc.rank == c.rank {
			front := cs[:i]
			back := cs[i+1:]
			return append(front, back...)
		}
	}
	return cs
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
	return Player{hand, firstCardTactic, accepts, 0, true, 0}
}

var bids = []int{
	18, 20, 22, 23, 24,
	27, 30, 33, 35, 36,
	40, 44, 45, 46, 48, 50,
	54, 55, 59, 60,
	63, 66, 70, 72, 77,
	80, 81, 84, 88, 90, 96, 99, 100,
}

func accepts(bidIndex int) bool {
	if r.Intn(10) > 1 {
		return true
	}
	return false
}
func bidLog(format string, a ...interface{}) {

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
	bidLog("FOREHAND vs MIDDLEHAND")
	bidIndex, p := bidding(players[0], players[1], 0)
	bidIndex++
	bidLog("WINNER vs BACKHAND")
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

func (p *Player) declareTrump() SuitState {
	return SuitState{mostCardsSuit(p.hand), ""}
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
	mat := matadors(state, cs)
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
		for leastMult * base < bid {
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

func game(players []*Player) {
	fmt.Println("------------NEW GAME----------")
	// DEALING
	cards := Shuffle(makeDeck())
	players[0].hand = cards[:10]
	players[1].hand = cards[10:20]
	players[2].hand = cards[20:30]

	skat := cards[30:32]

	// BIDDING
	bidIndex, declarer := bid(players)
	if bidIndex == -1 {
		fmt.Println("ALL PASSED")
		return
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
	state := declarer.declareTrump()
	declarerCards := append(declarer.hand, skat...)

	fmt.Printf("BID: %d, SUIT: %d %s", 
		bids[bidIndex], countCardsSuit(state.trump, declarer.hand), state.trump)

	// PLAY
	for i := 0; i < 10; i++ {
		players = round(&state, players)
	}

	gs := gameScore(state, declarerCards, declarer.score, bids[bidIndex], 
		declarer.schwarz, opp1.schwarz && opp2.schwarz, handGame)

	declarer.totalScore += gs

	if declarer.score > 60 {
		fmt.Printf(" VICTORY: %d, SCORE: %d\n", declarer.score, gs)
	} else {
		fmt.Printf(" LOSS: %d, SCORE: %d\n", declarer.score, gs)
	}

}

func main() {
	player1 := makePlayer([]Card{})
	player2 := makePlayer([]Card{})
	player3 := makePlayer([]Card{})	

	players := []*Player{&player1, &player2, &player3}

	for times := 12; times > 0; times-- {
		for _, p := range players {
			p.score = 0
			p.schwarz = true
		}
		game(players)
		fmt.Println(player1.totalScore, player2.totalScore, player3.totalScore)
		time.Sleep(1000 * time.Millisecond)
	}

}
