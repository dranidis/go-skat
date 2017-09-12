package main

type PlayerI interface {
	playerTactic(s *SuitState, c []Card) Card
	accepts(bidIndex int, listens bool) bool
	declareTrump() string
	discardInSkat(skat []Card)
	pickUpSkat(skat []Card) bool
	calculateHighestBid(bool) int
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
	getWon() int
	getLost() int
	wonAsDefenders()
	getWonAsDefenders() int
	setDeclaredBid(int)
	getDeclaredBid() int
	setPartner(p PlayerI)

	ResetPlayer()
}

