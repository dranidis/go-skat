package main

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
	followsSuite := func(s SuitState, card Card) bool {
		if card.rank == "J" {
			return s.trump == s.follow
		}
		return s.follow == card.suit
	}
	for _, c := range playerHand {
		if followsSuite(s, c) {
			return followsSuite(s, card)
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
