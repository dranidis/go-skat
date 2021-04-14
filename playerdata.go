package main

import (
	"fmt"
)

type PlayerData struct {
	name         string
	hand         []Card
	highestBid   int
	score        int
	schwarz      bool
	totalScore   int
	previousSuit string
	won          int
	lost         int
	defWon       int
	declaredBid  int
	Inference
}

var randomNameNumber = 0

func makePlayerData(hand []Card) PlayerData {
	randomNameNumber++
	inf := makeInference()
	return PlayerData{"dummy" + fmt.Sprintf("%d", randomNameNumber),
		//false,
		// false,
		hand, 0, 0, true, 0, "", 0, 0, 0, 0, inf}
}

func (p *PlayerData) clone() PlayerData {
	newPlayerData := makePlayerData([]Card{})

	newPlayerData.name = p.name
	newHand := make([]Card, len(p.hand))
	copy(newHand, p.hand)

	newPlayerData.hand = newHand
	newPlayerData.highestBid = p.highestBid
	newPlayerData.score = p.score
	newPlayerData.schwarz = p.schwarz
	newPlayerData.totalScore = p.totalScore
	newPlayerData.previousSuit = p.previousSuit
	newPlayerData.won = p.won
	newPlayerData.lost = p.lost
	newPlayerData.defWon = p.defWon
	newPlayerData.declaredBid = p.declaredBid

	newPlayerData.Inference = p.cloneInference()

	return newPlayerData
}

func (p PlayerData) String() string {
	return fmt.Sprintf("(%s):%v [%d]", p.name, p.hand, p.score)
}

func (p *PlayerData) ResetPlayer() {
	p.setScore(0)
	p.setSchwarz(true)
	p.setPreviousSuit("")
	p.setDeclaredBid(0)

	p.Inference = makeInference()
}

func (p *PlayerData) getInference() Inference {
	return p.Inference
}

func (p *PlayerData) setDeclaredBid(b int) {
	p.declaredBid = b
}
func (p *PlayerData) getDeclaredBid() int {
	return p.declaredBid
}

func (p *PlayerData) wonAsDefenders() {
	p.defWon++
}

func (p *PlayerData) getWonAsDefenders() int {
	return p.defWon
}

func (p *PlayerData) getWon() int {
	return p.won
}

func (p *PlayerData) getLost() int {
	return p.lost
}

func (p *PlayerData) incTotalScore(s int) {
	if s > 0 {
		p.won++
	} else {
		p.lost++
	}
	p.totalScore += s
}

func (p *PlayerData) setHand(cs []Card) {
	p.hand = cs
}

func (p *PlayerData) setScore(s int) {
	p.score = s
}

func (p *PlayerData) setSchwarz(b bool) {
	p.schwarz = b
}
func (p *PlayerData) setPreviousSuit(s string) {
	p.previousSuit = s
}

func (p *PlayerData) getScore() int {
	return p.score
}

func (p *PlayerData) getPreviousSuit() string {
	return p.previousSuit
}

func (p *PlayerData) getTotalScore() int {
	return p.totalScore
}

func (p *PlayerData) setName(n string) {
	p.name = n
}

func (p *PlayerData) getName() string {
	return p.name
}

func (p *PlayerData) getHand() []Card {
	return p.hand
}

func (p *PlayerData) isSchwarz() bool {
	return p.schwarz
}
