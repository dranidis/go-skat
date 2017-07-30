package main

import (
	"bufio"
	"fmt"
	"os"
)

type PlayerI interface {
	playerTactic(s *SuitState, c []Card) Card
	accepts(bidIndex int) bool
	declareTrump() string
	discardInSkat(skat []Card)
	pickUpSkat(skat []Card) bool
	calculateHighestBid() int
	//

	incTotalScore(s int)
	setHand(cs []Card)
	setScore(s int)
	setSchwarz(b bool)
	setPreviousSuit(s string)
	getScore() int
	getPreviousSuit() string
	getTotalScore() int
	setName(n string)
	getName() string
	getHand() []Card
	isSchwarz() bool
}

type HumanPlayer struct {
	PlayerData
}

func makeHumanPlayer(hand []Card) HumanPlayer {
	return HumanPlayer{
		PlayerData: makePlayerData(hand)}
}

func (p *HumanPlayer) accepts(bidIndex int) bool {

	fmt.Printf("HAND: %v", p.getHand())

	if !getYes(" Do you accept %d? (y/n/q)", bids[bidIndex]) {
		return false
	}
	return true
}

func (p *HumanPlayer) declareTrump() string {
	fmt.Printf("HAND: %v\n", p.getHand())
	for {
		fmt.Printf("TRUMP? (1 for CLUBS, 2 for SPADE, 3 for HEART, 4 for CARO, G for GRAND) ")
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
		case 'G':
			return GRAND
		case 'g':
			return GRAND
		default:
			continue
		}
	}

	return mostCardsSuit(p.getHand())
}

func (p *HumanPlayer) calculateHighestBid() int {
	return 0
}

func (p *HumanPlayer) discardInSkat(skat []Card) {
	p.setHand(sortSuit("", p.getHand()))
	for {
		gameLog("Full Hand : %v\n", p.getHand())
		gameLog("DISCARD CARDS? (1 to %d) ", len(p.getHand()))

		var i1, i2 int
		_, err := fmt.Scanf("%d", &i1)
		if err != nil {
			gameLog("%v", err)
			continue
		}
		_, err = fmt.Scanf("%d", &i2)
		if err != nil {
			gameLog("%v", err)
			continue
		}
		//fmt.Println(i1, i2)
		if i1 > len(p.getHand()) || i2 > len(p.getHand()) || i1 == i2 {
			continue
		}

		card1 := p.getHand()[i1-1]
		card2 := p.getHand()[i2-1]

		if !getYes("Discard %v and %v (y/n/q) ", card1, card2) {
			continue
		}

		p.setHand(remove(p.getHand(), card1))
		p.setHand(remove(p.getHand(), card2))
		skat[0] = card1
		skat[1] = card2
		return
	}
}

func getYes(format string, a ...interface{}) bool {
	for {
		gameLog(format, a...)
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()

		if err != nil {
			gameLog("%v", err)
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
			gameLog("... don't understand! ")
			continue
		}
	}
	return false
}

func (p *HumanPlayer) pickUpSkat(skat []Card) bool {
	gameLog("HAND: %v", p.getHand())

	if !getYes(" Pick up SKAT? (y/n/q) ") {
		return false
	}

	gameLog("SKAT: %v\n", skat)
	hand := make([]Card, 10)
	copy(hand, p.getHand())
	hand = append(hand, skat...)
	p.setHand(hand)

	p.discardInSkat(skat)

	return true
}

func (p *HumanPlayer) playerTactic(s *SuitState, c []Card) Card {
	gameLog("Your Hand : %v\n", p.getHand())
	gameLog("Valid: %v\n", c)
	for {
		fmt.Printf("CARD? (1 to %d) ", len(c))
		var i int
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if i > len(c) {
			continue
		}
		card := c[i-1]
		if len(s.trick) == 0 {
			p.setPreviousSuit(card.suit)
		}
		return c[i-1]
	}
}