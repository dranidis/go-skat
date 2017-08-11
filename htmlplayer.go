package main

import (
	"bufio"
	"fmt"
	"os"
)

type HtmlPlayer struct {
	PlayerData

	handGame bool
}

func makeHtmlPlayer(hand []Card) HtmlPlayer {
	return HtmlPlayer{
		PlayerData: makePlayerData(hand)}
}

func (p *HtmlPlayer) setPartner(partner PlayerI) {

}

func (p *HtmlPlayer) accepts(bidIndex int) bool {

	fmt.Printf("HAND: %v", p.getHand())

	if !getYes(" Do you accept %d? (y/n/q)", bids[bidIndex]) {
		return false
	}
	return true
}

func (p *HtmlPlayer) declareTrump() string {
	fmt.Printf("HAND: %v\n", p.getHand())
	for {
		fmt.Printf("TRUMP? (1 for CLUBS, 2 for SPADE, 3 for HEART, 4 for CARO, g for GRAND, n for NULL) ")
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
		case 'n':
			if p.declaredBid > 23 {
				if p.handGame {
					if p.declaredBid > 35 {
						gameLog("Your bid %s is higher than Null Hand 35\n", p.declaredBid)
						continue
					}
				}
				gameLog("Your bid %s is higher than Null 23\n", p.declaredBid)
				continue
			}
			return NULL
		default:
			continue
		}
	}

	return mostCardsSuit(p.getHand())
}

func (p *HtmlPlayer) calculateHighestBid() int {
	return 0
}

func (p *HtmlPlayer) discardInSkat(skat []Card) {
	p.setHand(sortSuit("", p.getHand()))
	sorting := true
	for {
		gameLog("Full Hand : %v\n", p.getHand())
		gameLog("DISCARD CARDS? (1 to %d) [0 to toggle SUIT/NULL sorting]", len(p.getHand()))

		var i1, i2 int
		_, err := fmt.Scanf("%d", &i1)
		if err != nil {
			gameLog("%v", err)
			continue
		}
		if i1 == 0 {
			if sorting {
				gameLog("Change to Null sorting\n")
				p.setHand(sortSuit(NULL, p.getHand()))
				sorting = false
				continue
			} else {
				gameLog("Change to Suit sorting\n")
				p.setHand(sortSuit("", p.getHand()))
				sorting = true
				continue
			}

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

func (p *HtmlPlayer) pickUpSkat(skat []Card) bool {
	gameLog("HAND: %v", p.getHand())

	if !getYes(" Pick up SKAT? (y/n/q) ") {
		p.handGame = true
		return false
	}

	gameLog("SKAT: %v\n", skat)
	hand := make([]Card, 10)
	copy(hand, p.getHand())
	hand = append(hand, skat...)
	p.setHand(hand)

	p.discardInSkat(skat)

	p.handGame = false
	return true
}

func (p *HtmlPlayer) playerTactic(s *SuitState, c []Card) Card {

	gameLog("Your Hand : %v\n", p.getHand())
	gameLog("Valid: %v\n", c)

	gameLog("Waiting card at trickChannel\n")
	card := <- trickChannel
	return card
}
