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
	calculateHighestBid()
	//

	incTotalScore(s int)
	setHand(cs []Card)
	setScore(s int)
	setHuman(b bool)
	setSchwarz(b bool)
	setPreviousSuit(s string)
	getScore() int
	getPreviousSuit() string
	getTotalScore() int
	setName(n string)
	getName() string
	getHand() []Card
	isHuman() bool
	isSchwarz() bool
}

type HumanPlayer struct {
	name         string
	hand         []Card
	highestBid   int
	score        int
	schwarz      bool
	totalScore   int
	previousSuit string
}

func makeHumanPlayer(hand []Card) HumanPlayer {
	return HumanPlayer{"dummy",
		//false,
		// false,
		hand, 0, 0, true, 0, ""}
}

func (p *HumanPlayer) accepts(bidIndex int) bool {

	fmt.Printf("HAND: %v", p.getHand())
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

func (p *HumanPlayer) declareTrump() string {
	fmt.Printf("HAND: %v\n", p.getHand())
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

	return mostCardsSuit(p.getHand())
}

func (p *HumanPlayer) calculateHighestBid() {

}

func (p *HumanPlayer) discardInSkat(skat []Card) {
	p.setHand(sortSuit("", p.getHand()))
	gameLog("Full Hand : %v\n", p.getHand())
	for {
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
		p.setHand(remove(p.getHand(), card1))
		p.setHand(remove(p.getHand(), card2))
		skat[0] = card1
		skat[1] = card2
		return
	}
}

func (p *HumanPlayer) pickUpSkat(skat []Card) bool {
		gameLog("HAND: %v", p.getHand())
		yes := false
		for !yes {
			gameLog("Pick up SKAT? (y/n/q) ")
			reader := bufio.NewReader(os.Stdin)
			char, _, err := reader.ReadRune()

			if err != nil {
				gameLog("%v", err)
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
				gameLog("... don't understand! ")
				continue
			}
		}

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

func (p *HumanPlayer) incTotalScore(s int) {
	p.totalScore += s
}

func (p *HumanPlayer) setHand(cs []Card) {
	p.hand = cs
}

func (p *HumanPlayer) setScore(s int) {
	p.score = s
}

func (p *HumanPlayer) setHuman(b bool) {
}

func (p *HumanPlayer) setSchwarz(b bool) {
	p.schwarz = b
}
func (p *HumanPlayer) setPreviousSuit(s string) {
	p.previousSuit = s
}

func (p *HumanPlayer) getScore() int {
	return p.score
}

func (p *HumanPlayer) getPreviousSuit() string {
	return p.previousSuit
}

func (p *HumanPlayer) getTotalScore() int {
	return p.totalScore
}

func (p *HumanPlayer) setName(n string) {
	p.name = n
}

func (p *HumanPlayer) getName() string {
	return p.name
}

func (p *HumanPlayer) getHand() []Card {
	return p.hand
}

func (p *HumanPlayer) isHuman() bool {
	return true
}

func (p *HumanPlayer) isSchwarz() bool {
	return p.schwarz
}
