package main

import (
	"fmt"
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
	return <-declareChannel
}

func (p *HtmlPlayer) calculateHighestBid(b bool) int {
	return 0
}

func (p *HtmlPlayer) discardInSkat(skat []Card) {
	card1 := <- discardChannel	
	card2 := <- discardChannel	

	p.setHand(remove(p.getHand(), card1))
	p.setHand(remove(p.getHand(), card2))
	skat[0] = card1
	skat[1] = card2
}

func (p *HtmlPlayer) pickUpSkat(skat []Card) bool {
	gameLog("HAND: %v", p.getHand())

	pickUp := <-pickUpChannel
	if pickUp == "hand" {
		gameLog("HAND game\n")
		p.handGame = true
		return false
	}

	gameLog("SKAT: %v\n", skat)
	hand := make([]Card, 10)
	copy(hand, p.getHand())
	hand = append(hand, skat...)
	p.setHand(sortSuit("", hand))

	// find indices of skat cards in hand
	for i, c := range(p.hand) {
		if c.equals(skat[0]) {
			skatPositionChannel <- i
		} 
		if c.equals(skat[1]) {
			skatPositionChannel <- i
		} 		
	}
	// SEND HAND TO SERVER

	p.discardInSkat(skat)

	p.handGame = false
	return true
}

func (p *HtmlPlayer) playerTactic(s *SuitState, c []Card) Card {

	gameLog("Your Hand : %v\n", p.getHand())
	gameLog("Valid: %v\n", c)
	printCollectedInfo(s)

	htmlLog("Reading card at trickChannel...\n")
	card := <-trickChannel
	htmlLog("read %v\n", card)
	return card
}
