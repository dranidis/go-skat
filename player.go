package main

type PlayerI interface {
	playerTactic(s *SuitState, c []Card) Card
	accepts(bidIndex int) bool
	declareTrump() string
	discardInSkat(skat []Card)
	pickUpSkat(skat []Card) bool 
	calculateHighestBid()
	//
	getName() string
	setHand(cs []Card)
	getHand() []Card
	isHuman() bool
}