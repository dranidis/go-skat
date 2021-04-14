package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ISSPlayer struct {
	PlayerData
	shortcut []Card
}

func makeISSPlayer(hand []Card) ISSPlayer {
	return ISSPlayer{
		PlayerData: makePlayerData(hand),
		shortcut:   []Card{},
	}
}

func (p *ISSPlayer) clone() PlayerI {
	newPlayer := makeISSPlayer([]Card{})

	newPlayer.PlayerData = p.PlayerData.clone()
	newPlayer.shortcut = p.shortcut
	return &newPlayer
}

func (p *ISSPlayer) setPartner(partner PlayerI) {

}

func (p *ISSPlayer) accepts(bidIndex int, listens bool) bool {
	fmt.Printf("Waiting bid from %s..", p.getName())
	bid := <-bidChannel
	fmt.Printf(" .. bid received %s\n", bid)

	if bidNr, err := strconv.ParseInt(bid, 10, 64); err == nil {
		if int64(bids[bidIndex]) == bidNr {
			return true
		}
	}
	if bid == "p" {
		return false
	}
	if bid == "y" {
		return true
	}
	log.Fatal("Unrecognised bid from ISS player")
	return false
}

func (p *ISSPlayer) declareTrump() string {
	fmt.Printf("Waiting declaration from %s..", p.getName())
	declare := <-declareChannel
	fmt.Printf(".. received %s", declare)

	return declare
}

func (p *ISSPlayer) calculateHighestBid(b bool) int {
	return 0
}

func (p *ISSPlayer) discardInSkat(skat []Card) {
	// fmt.Printf("Waiting 2 discard cards from %s..", p.getName())
	// card1 := <- discardChannel
	// card2 := <- discardChannel
	// fmt.Printf(" .. received %v and %v\n", card1, card2)

	// p.setHand(remove(p.getHand(), card1))
	// p.setHand(remove(p.getHand(), card2))
	// skat[0] = card1
	// skat[1] = card2
}

func (p *ISSPlayer) pickUpSkat(skat []Card) bool {

	// WHAT HAPPENS IN A HAND GAME???

	// fmt.Printf("Waiting pick up from %s..", p.getName())
	// pickUp := <-pickUpChannel
	// fmt.Printf(" .. received %s\n", pickUp)

	// if pickUp == "hand" {
	// 	fmt.Printf("Player %s plays a hand game\n", p.getName())
	// 	return false
	// } else {
	// 	fmt.Printf("Player %s picks up the skat\n", p.getName())
	// }
	// p.discardInSkat(skat)
	return true
}

func (p *ISSPlayer) playerTactic(s *SuitState, c []Card) Card {
	var card Card
	if len(p.shortcut) > 0 {
		card = p.shortcut[0]
		p.shortcut = remove(p.shortcut, card)
		fmt.Printf(" .. playing SC card %v\n", card)
		return card
	}

	fmt.Printf("Waiting card up from %s..", p.getName())
	card = <-trickChannel

	if card.Suit == "SC" { // many cards packed in the rank
		p.shortcut = parseCards(strings.Split(card.Rank, "."))
		fmt.Printf(" .. received shortcut %v\n", p.shortcut)
		card = p.shortcut[0]
		p.shortcut = remove(p.shortcut, card)
		fmt.Printf(" .. playing SC card %v\n", card)
		return card
	}
	p.shortcut = []Card{} // empty shortcut if normal card was played
	fmt.Printf(" .. received %v\n", card)

	suit := getSuit(s.trump, card)
	if suit != s.trump {
		p.setPreviousSuit(suit)
	}

	return card
}
