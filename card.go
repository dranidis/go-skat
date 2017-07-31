package main

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const CLUBS = "CLUBS"
const SPADE = "SPADE"
const HEART = "HEART"
const CARO = "CARO"
const GRAND = "Grand"

var r = rand.New(rand.NewSource(1))
var _ = rand.New(rand.NewSource(time.Now().Unix()))

var black = color.New(color.Bold, color.FgWhite).SprintFunc()
var green = color.New(color.Bold, color.FgGreen).SprintFunc()
var red = color.New(color.Bold, color.FgRed).SprintFunc()
var yellow = color.New(color.Bold, color.FgYellow).SprintFunc()

type Card struct {
	suit string
	rank string
}

func (c Card) equals(o Card) bool {
	return c.suit == o.suit && c.rank == o.rank
}

func (c Card) String() string {

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
	case GRAND:
		return 24
	}
	return 0
}

func getSuit(trump string, card Card) string {
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
			suits = []string{CLUBS, SPADE, HEART, CARO}
		case SPADE:
			suits = []string{SPADE, CLUBS, HEART, CARO}
		case HEART:
			suits = []string{HEART, CLUBS, SPADE, CARO}
		case CARO:
			suits = []string{CARO, CLUBS, SPADE, HEART}
		default:
			suits = []string{CLUBS, SPADE, HEART, CARO}
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

// CARD MANIPULATION FUNCTIONS

func in(cs []Card, card ...Card) bool {
	inOne := func(cs []Card, c Card) bool {
		for _, card := range cs {
			if card.equals(c) {
				return true
			}
		}
		return false
	}
	for _, c := range card {
		if !inOne(cs, c) {
			return false
		}
	}
	return true
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

func filterSuit(suits []string, f func(string) bool) []string {
	cs := []string{}
	for _, c := range suits {
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

func trumpCards(trump string, cards []Card) []Card {
	return filter(cards, func(c Card) bool {
		return c.suit == trump || c.rank == "J"
	})
}

func nonTrumpCards(suit string, cards []Card) []Card {
	return filter(cards, func(c Card) bool {
		return c.suit == suit && c.rank != "J"
	})
}

// TACTICS aux functions

func ShortestNonTrumpSuit(trump string, cards []Card) string {
	minI, minCount := -1, 99
	for i, s := range suits {
		if s == trump {
			continue
		}
		c := len(nonTrumpCards(s, cards))
		if c < minCount && c > 0 {
			minI = i
			minCount = c
		}
	}
	return suits[minI]
}

func LongestNonTrumpSuit(trump string, cards []Card) string {
	maxI, maxCount := -1, -1
	for i, s := range suits {
		if s == trump {
			continue
		}
		c := len(nonTrumpCards(s, cards))
		if c > maxCount {
			maxI = i
			maxCount = c
		}
	}
	return suits[maxI]
}

// With a preference to non-A suits
// and a preference to stronger cards (between A-suits)
func mostCardsSuit(cards []Card) string {
	maxCount := 0
	maxI := -1
	for i, s := range suits {
		cs := trumpCards(s, cards)
		count := 200 * len(cs)
		if !in(cards, Card{s, "A"}) {
			count += 100
		}
		count += sum(cs)
		if count > maxCount {
			maxCount = count
			maxI = i
		}
	}
	return suits[maxI]
}

func lessCardsSuit(cards []Card) string {
	c := len(nonTrumpCards(CLUBS, cards))
	s := len(nonTrumpCards(SPADE, cards))
	h := len(nonTrumpCards(HEART, cards))
	k := len(nonTrumpCards(CARO, cards))

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
