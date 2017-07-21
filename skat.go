package main

import (
	_ "fmt"
	"math/rand"
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
	score int
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

func setNextTrickOrder(s *SuitState, players []*Player, trick []Card) []*Player {
	if s.greater(trick[0], trick[1]) && s.greater(trick[0], trick[2]) {
		players[0].score += sum(trick)
		return players
	}
	if s.greater(trick[1], trick[2]) {
		players[1].score += sum(trick)
		return []*Player{players[1], players[2], players[0]}
	}
	players[2].score += sum(trick)
	return []*Player{players[2], players[0], players[1]}
}

func round(s *SuitState, players []*Player) {
	var trick [3]Card
	trick[0] = players[0].play(s)
	s.follow = getSuite(*s, trick[0])
	trick[1] = players[1].play(s)
	trick[2] = players[2].play(s)

	players = setNextTrickOrder(s, players, trick[:])
	s.follow = ""
}

var r = rand.New(rand.NewSource(99))

func firstCardTactic(c []Card) Card {
	return c[0]
}

func randomCardTactic(c []Card) Card {
	cardIndex := r.Intn(len(c))
	return c[cardIndex]
}

func (p *Player) play(s *SuitState) Card{
	valid := validCards(*s, p.hand)
	card := p.playerTactic(valid)
	//fmt.Println(s, card, p.hand)
	p.hand = remove(p.hand, card)
	//fmt.Println(p.hand)
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

// func (s SuitState) valid(playerHand []Card, card Card) bool {
// 	followsSuite := func(s SuitState, card Card) bool {
// 		if card.rank == "J" {
// 			return s.trump == s.follow
// 		}
// 		return s.follow == card.suit
// 	}
// 	for _, c := range playerHand {
// 		if followsSuite(s, c) {
// 			return followsSuite(s, card)
// 		}
// 	}
// 	return true
// }
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

func main() {

}
