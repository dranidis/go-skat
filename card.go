package main

import (
	"github.com/fatih/color"
	"log"
	"math/rand"
	"runtime/debug"
	"sort"
)

const CLUBS = "CLUBS"
const SPADE = "SPADE"
const HEART = "HEART"
const CARO = "CARO"
const GRAND = "Grand"
const NULL = "Null"

var black = color.New(color.Bold, color.FgWhite).SprintFunc()
var green = color.New(color.Bold, color.FgGreen).SprintFunc()
var red = color.New(color.Bold, color.FgRed).SprintFunc()
var yellow = color.New(color.Bold, color.FgYellow).SprintFunc()

type Card struct {
	Suit string
	Rank string
}

func (c Card) equals(o Card) bool {
	return c.Suit == o.Suit && c.Rank == o.Rank
}

func (c Card) String() string {

	switch c.Suit {
	case CLUBS:
		return black(c.Rank)
	case SPADE:
		return green(c.Rank)
	case HEART:
		return red(c.Rank)
	case CARO:
		return yellow(c.Rank)
	}
	return ""
}

func cardValue(c Card) int {
	switch c.Rank {
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
	case NULL:
		debug.PrintStack()

		log.Fatal("Error! No base value for NULL")
	}
	return 0
}

func sevens(cs []Card) []Card {
	return filter(cs, func(c Card) bool {
		return c.Rank == "7"
	})
}

func getSuit(trump string, card Card) string {
	if trump == NULL {
		return card.Suit
	}
	if card.Rank == "J" {
		return trump
	}
	return card.Suit
}

var suits = []string{CLUBS, SPADE, HEART, CARO}
var ranks = []string{"J", "A", "10", "K", "D", "9", "8", "7"}
var nullRanks = []string{"A", "K", "D", "J", "10", "9", "8", "7"}
var nullRanksRev = []string{"7", "8", "9", "10", "J", "D", "K", "A"}

var ranksNum = map[string]int{
	"J":  14,
	"A":  13,
	"10": 12,
	"K":  11,
	"D":  10,
	"9":  9,
	"8":  8,
	"7":  7,
}

var suitNum = map[string]int{
	CLUBS: 4,
	SPADE: 3,
	HEART: 2,
	CARO:  1,
}

type ByRank []Card

func (a ByRank) Len() int      { return len(a) }
func (a ByRank) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByRank) Less(i, j int) bool {
	if ranksNum[a[i].Rank] < ranksNum[a[j].Rank] {
		return false
	} else if ranksNum[a[i].Rank] > ranksNum[a[j].Rank] {
		return true
	}
	return suitNum[a[i].Suit] > suitNum[a[j].Suit]
}

// func sortRank(cs []Card) []Card {
// 	return sortRankSpecial(cs, ranks)
// }

func sortRank(cs []Card) []Card {
	sort.Sort(ByRank(cs))
	return cs
}

type ByRankSpecial []Card

func (a ByRankSpecial) Len() int      { return len(a) }
func (a ByRankSpecial) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByRankSpecial) Less(i, j int) bool {
	if ranksSpecialNum[a[i].Rank] < ranksSpecialNum[a[j].Rank] {
		return false
	} else if ranksSpecialNum[a[i].Rank] > ranksSpecialNum[a[j].Rank] {
		return true
	}
	return suitNum[a[i].Suit] > suitNum[a[j].Suit]
}

var ranksSpecialNum = map[string]int{
	"J":  14,
	"A":  13,
	"10": 12,
	"K":  11,
	"D":  10,
	"9":  9,
	"8":  8,
	"7":  7,
}

func sortRankSpecial(cs []Card, ranks []string) []Card {
	for i, r := range ranks {
		ranksSpecialNum[r] = 12 - i
	}
	sort.Sort(ByRankSpecial(cs))
	return cs
}

// func sortRankSpecial(cs []Card, ranks []string) []Card {
// 	cards := []Card{}

// 	for _, r := range ranks {
// 		for _, s := range suits {
// 			if in(cs, Card{s, r}) {
// 				cards = append(cards, Card{s, r})
// 			}
// 		}
// 	}
// 	return cards
// }

func sortSuitRankSpecial(cs []Card, suits, ranks []string) []Card {
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

func previousNull(c Card) Card {
	for i := len(nullRanks) - 1; i >= 0; i-- {
		if c.Rank != nullRanks[i] {
			continue
		}
		if i > 0 {
			return Card{c.Suit, nullRanks[i-1]}
		}
	}
	return Card{"", ""}
}

func singletons(cs []Card) []Card {
	singles := []Card{}
	for _, s := range suits {
		cards := filter(cs, func(c Card) bool {
			return c.Suit == s
		})
		if len(cards) == 1 {
			singles = append(singles, cards[0])
		}
	}
	return singles
}

func similar(s *SuitState, cards []Card) []Card {
	sim := []Card{}
	cards = sortSuit(s.trump, cards)
	if len(cards) == 0 {
		return sim
	}
	card := cards[0]
	next := nextCard(s.trump, card)
	for in(s.cardsPlayed, next) && next.Rank != "A" && next.Rank != "K" && next.Rank != "9" {
		next = nextCard(s.trump, next)
	}
	sim = append(sim, card)
	for i := 1; i < len(cards); i++ {
		// debugTacticsLog("next: %v cards[%d]: %v\n", next, i, cards[i])
		if cards[i].Rank == "A" || cards[i].Rank == "K" || cards[i].Rank == "9" {
			card = cards[i]
			sim = append(sim, card)
			next = nextCard(s.trump, card)
			for in(s.cardsPlayed, next) && next.Rank != "A" && next.Rank != "K" && next.Rank != "9" {
				next = nextCard(s.trump, next)
			}
		} else if cards[i].equals(next) {
			next = nextCard(s.trump, next)
			for in(s.cardsPlayed, next) && next.Rank != "A" && next.Rank != "K" && next.Rank != "9" {
				next = nextCard(s.trump, next)
			}
			continue
		} else {
			card = cards[i]
			sim = append(sim, card)
			next = nextCard(s.trump, card)
			for in(s.cardsPlayed, next) && next.Rank != "A" && next.Rank != "K" && next.Rank != "9" {
				next = nextCard(s.trump, next)
			}
		}
	}
	return sim
}

// func equivalent(s *SuitState, cards []Card) []Card {
// 	eqCards := []Card{}
// 	for _, c := range cards {
// 		if cardValue(c) > 2 {
// 			eqCards = append(eqCards, c)
// 		}
// 	}
// 	for _, s := range suits {
// 		if in(cards, Card{s, "9"}) {
// 			eqCards = append(eqCards, Card{s, "9"})
// 			if !in(cards, Card{s, "8"}) && in(cards, Card{s, "7"}) {
// 				eqCards = append(eqCards, Card{s, "7"})
// 			}
// 		} else if in(cards, Card{s, "8"}) {
// 			eqCards = append(eqCards, Card{s, "8"})
// 		} else if in(cards, Card{s, "7"}) {
// 			eqCards = append(eqCards, Card{s, "8"})
// 		}
// 	}

// 	if in(cards, Card{CLUBS, "J"}) {
// 		eqCards = append(eqCards, Card{CLUBS, "J"})
// 		if !in(cards, Card{SPADE, "J"}) && in(cards, Card{HEART, "J"}) {
// 			eqCards = append(eqCards, Card{HEART, "J"})
// 		}
// 	} else if in(cards, Card{SPADE, "J"}) {
// 		eqCards = append(eqCards, Card{SPADE, "J"})
// 		if !in(cards, Card{HEART, "J"}) && in(cards, Card{CARO, "J"}) {
// 			eqCards = append(eqCards, Card{CARO, "J"})
// 		}
// 	} else if in(cards, Card{HEART, "J"}) {
// 		eqCards = append(eqCards, Card{HEART, "J"})
// 	} else if in(cards, Card{CARO, "J"}) {
// 		eqCards = append(eqCards, Card{CARO, "J"})
// 	}
// 	return eqCards
// }

func equivalent(s *SuitState, cards []Card) []Card {
	eqCards := []Card{}
	for _, c := range cards {
		if cardValue(c) > 2 {
			eqCards = append(eqCards, c)
		}
	}

	suits := []string{"CARO", "HEART", "SPADE", "CLUBS"}
	ins := true
	for _, suit := range suits {
		card := Card{suit, "J"}
		if in(cards, card) && ins {
			eqCards = append(eqCards, card)
			ins = false
		} else if !in(cards, card) {
			if !in(s.cardsPlayed, card) {
				ins = true
			} else if in(s.trick, card) {
				ins = true
			}
		}
	}
	ranks := []string{"9", "8", "7"}

	for _, suit := range suits {
		ins = true
		for _, r := range ranks {
			card := Card{suit, r}
			if in(cards, card) && ins {
				eqCards = append(eqCards, card)
				ins = false
			} else if !in(cards, card) {
				if !in(s.cardsPlayed, card) {
					ins = true
				} else if in(s.trick, card) {
					ins = true
				}
			}
		}
	}
	return eqCards
}

func nextCard(trump string, c Card) Card {
	if c.equals(Card{CLUBS, "J"}) {
		return Card{SPADE, "J"}
	}
	if c.equals(Card{SPADE, "J"}) {
		return Card{HEART, "J"}
	}
	if c.equals(Card{HEART, "J"}) {
		return Card{CARO, "J"}
	}
	if c.equals(Card{CARO, "J"}) {
		return Card{trump, "A"} // returnin {"", A} in Grand
	}
	i, r := -1, ""
	for i, r = range ranks {
		if c.Rank == r {
			break
		}
	}
	if i+1 < len(ranks) {
		return Card{c.Suit, ranks[i+1]}
	}
	return Card{"", ""}
}

func sortValue(cs []Card) []Card {
	valueRanks := []string{"A", "10", "K", "D", "J", "7", "8", "9"}
	return sortRankSpecial(cs, valueRanks)
}

func realSortValue(cs []Card) []Card {
	valueRanks := []string{"A", "10", "K", "D", "J", "9", "8", "7"}
	return sortRankSpecial(cs, valueRanks)
}

func sortValueNull(cs []Card) []Card {
	valueRanks := []string{"7", "8", "9", "10", "J", "D", "K", "A"}
	return sortRankSpecial(cs, valueRanks)
}

func sortSuit(trump string, cs []Card) []Card {
	cards := []Card{}

	cardJs := []Card{
		{CLUBS, "J"},
		{SPADE, "J"},
		{HEART, "J"},
		{CARO, "J"},
	}
	cardsSuit := func(suit string) []Card {
		return []Card{
			{suit, "A"},
			{suit, "10"},
			{suit, "K"},
			{suit, "D"},
			{suit, "9"},
			{suit, "8"},
			{suit, "7"},
		}
	}
	if trump == NULL {
		cardsSuit = func(suit string) []Card {
			return []Card{
				{suit, "A"},
				{suit, "K"},
				{suit, "D"},
				{suit, "J"},
				{suit, "10"},
				{suit, "9"},
				{suit, "8"},
				{suit, "7"},
			}
		}
	}
	if trump != NULL {
		for _, c := range cardJs {
			if in(cs, c) {
				cards = append(cards, c)
			}
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

func ShuffleR(r *rand.Rand, cards []Card) []Card {
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
		{suit, "J"},
		{suit, "A"},
		{suit, "10"},
		{suit, "K"},
		{suit, "D"},
		{suit, "9"},
		{suit, "8"},
		{suit, "7"},
	}
}

func makeTrumpDeck(suit string) []Card {
	return []Card{
		{CLUBS, "J"},
		{SPADE, "J"},
		{HEART, "J"},
		{CARO, "J"},
		{suit, "A"},
		{suit, "10"},
		{suit, "K"},
		{suit, "D"},
		{suit, "9"},
		{suit, "8"},
		{suit, "7"},
	}
}

func makeNoTrumpDeck(suit string) []Card {
	return []Card{
		{suit, "A"},
		{suit, "10"},
		{suit, "K"},
		{suit, "D"},
		{suit, "9"},
		{suit, "8"},
		{suit, "7"},
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

// func remove(cs []Card, c Card) []Card {
// 	return filter(cs, func(cc Card) bool {
// 		return !(cc.equals(c))
// 	})
// }

// a = append(a[:i], a[i+1:]...)

func remove(cs []Card, card ...Card) []Card {
	for _, c := range card {
		cs = removeOne(cs, c)
	}
	return cs
}

func removeOne(cs []Card, card Card) []Card {
	for i, nc := range cs {
		if nc.equals(card) {
			cs = append(cs[:i], cs[i+1:]...)
			break
		}
	}
	return cs

}

// func remove(cs []Card, card ...Card) []Card {
// 	ncs := []Card{}
// 	for _, nc := range cs {
// 		found := false
// 		for _, c := range card {
// 			if nc.equals(c) {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			ncs = append(ncs, nc)
// 		}
// 	}
// 	return ncs
// }

func matadors(trump string, cs []Card) int {
	cards := []Card{
		{SPADE, "J"},
		{HEART, "J"},
		{CARO, "J"},
		{trump, "A"},
		{trump, "10"},
		{trump, "K"},
		{trump, "D"},
		{trump, "9"},
		{trump, "8"},
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
		return c.Suit == trump || c.Rank == "J"
	})
}

func nonTrumpCards(suit string, cards []Card) []Card {
	return filter(cards, func(c Card) bool {
		return c.Suit == suit && c.Rank != "J"
	})
}

// TACTICS aux functions

func strength(cs []Card) int {
	s := sum(cs)
	if in(cs, Card{CLUBS, "J"}) {
		s += 25
	}
	if in(cs, Card{SPADE, "J"}) {
		s += 20
	}
	if in(cs, Card{HEART, "J"}) {
		s += 15
	}
	if in(cs, Card{CARO, "J"}) {
		s += 12
	}
	return s
}

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
	if minI >= 0 {
		return suits[minI]
	}
	return ""
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
//   unless there are at least 3 suit
// and a preference to weaker cards (between A-suits)
func mostCardsSuit(cards []Card) string {
	spCardValuefunc := func(c Card) int {
		switch c.Rank {
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
		case "9":
			return 2
		case "8":
			return 1
		}
		return 0
	}
	spSum := func(trick []Card) int {
		s := 0
		for _, c := range trick {
			s += spCardValuefunc(c)
		}
		return s
	}

	asuits := 0
	for _, s := range suits {
		if in(cards, Card{s, "A"}) {
			asuits++
		}
	}
	acePenalty := 100
	if asuits > 2 {
		acePenalty = -100
	}
	debugTacticsLog("..mostCardsSuit")
	maxCount := 0
	maxI := -1
	for i, s := range suits {
		cs := trumpCards(s, cards)
		count := 200 * len(cs)
		count -= spSum(cs)
		if in(cards, Card{s, "A"}) {
			count -= acePenalty
		}
		if count > maxCount {
			maxCount = count
			maxI = i
		}
		debugTacticsLog(".%s %v: %d  ", s, cs, count)
	}
	debugTacticsLog("MOST: %s\n", suits[maxI])
	return suits[maxI]
}

func lessCardsSuitExcept(suitsToExclude []string, cards []Card) string {
	copyCards := filter(cards, func(c Card) bool {
		for _, s := range suitsToExclude {
			if c.Suit == s {
				return false
			}
		}
		return true
	})
	return lessCardsSuit(copyCards)
}

func lessCardsSuit(cards []Card) string {
	nonA := func(cs []Card) []Card {
		return filter(cs, func(c Card) bool {
			return c.Rank != "A"
		})
	}
	clubs := nonTrumpCards(CLUBS, cards)
	spades := nonTrumpCards(SPADE, cards)
	hearts := nonTrumpCards(HEART, cards)
	caro := nonTrumpCards(CARO, cards)
	c := 100*len(clubs) - sum(nonA(clubs)) // we want to discard higher value cards
	s := 100*len(spades) - sum(nonA(spades))
	h := 100*len(hearts) - sum(nonA(hearts))
	k := 100*len(caro) - sum(nonA(caro))

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
	// we don't want to discard a suit having an A
	if in(cards, Card{CLUBS, "A"}) {
		c += 100
	}
	if in(cards, Card{SPADE, "A"}) {
		s += 100
	}
	if in(cards, Card{HEART, "A"}) {
		h += 100
	}
	if in(cards, Card{CARO, "A"}) {
		k += 100
	}
	debugTacticsLog("...in cards %v,%v,%v,%v lessCardsSuit CLUBS: %v, SPADES: %v, HEARTS: %v, CARO: %v...", clubs, spades, hearts, caro, c, s, h, k)
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

func (s SuitState) nullGreater(card1, card2 Card) bool {
	return nullGreater(s.follow, card1, card2)
}

func nullGreater(sfollow string, card1, card2 Card) bool {
	rank := map[string]int{
		"A":  13,
		"K":  12,
		"D":  11,
		"J":  10,
		"10": 9,
		"9":  8,
		"8":  7,
		"7":  6,
	}
	if card1.Suit == sfollow {
		if card2.Suit == sfollow {
			return rank[card1.Rank] > rank[card2.Rank]
		}
		return true
	}

	if card2.Suit == sfollow {
		return false
	}

	return rank[card1.Rank] > rank[card2.Rank]
}

func (s SuitState) greater(card1 Card, cards ...Card) bool {
	for _, card2 := range cards {
		if s.greaterOne(card1, card2) {
			continue
		}
		return false
	}
	return true
}

func greater(strump, sfollow string, card1 Card, cards ...Card) bool {
	for _, card2 := range cards {
		if greaterOne(strump, sfollow, card1, card2) {
			continue
		}
		return false
	}
	return true
}

func (s SuitState) greaterOne(card1, card2 Card) bool {
	return greaterOne(s.trump, s.follow, card1, card2)
}

func greaterOne(strump, sfollow string, card1, card2 Card) bool {
	if strump == NULL {
		return nullGreater(sfollow, card1, card2)
	}

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

	if card1.Rank == "J" {
		if card2.Rank == "J" {
			return JRank[card1.Suit] > JRank[card2.Suit]
		}
		return true
	}

	if card2.Rank == "J" {
		return false
	}

	if card1.Suit == strump {
		if card2.Suit == strump {
			return rank[card1.Rank] > rank[card2.Rank]
		}
		return true
	}

	if card2.Suit == strump {
		return false
	}

	if card1.Suit == sfollow {
		if card2.Suit == sfollow {
			return rank[card1.Rank] > rank[card2.Rank]
		}
		return true
	}

	if card2.Suit == sfollow {
		return false
	}

	return rank[card1.Rank] > rank[card2.Rank]
}
