package main

import (
	_ "fmt"
)

type APlayer struct {
	Player
	moves     []int
	play      []int
	prevPlay  []int
	forcedBid int
}

func makeAPlayer(hand []Card) APlayer {
	return APlayer{
		Player:    makePlayer(hand),
		moves:     []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		play:      []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		prevPlay:  []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		forcedBid: 0,
	}
}

var analysisEnded = false
var tIndex = -1
var previousGameAnalysis = false

func (p *APlayer) playerTactic(s *SuitState, c []Card) Card {
	gameLog("Valid cards: %v\n", c)
	tIndex++

	if previousGameAnalysis {
		gameLog("Index: %v\n", p.prevPlay[tIndex])
		return c[p.prevPlay[tIndex]]
	}

	p.prevPlay[tIndex] = p.play[tIndex]
	gameLog("Index: %v\n", p.play[tIndex])
	p.moves[tIndex] = len(c)
	current := p.play[tIndex]
	if tIndex == 9 {
		p.next(tIndex - 1)
	}
	return c[current]
}

func (p *APlayer) next(t int) {
	if t == -1 {
		analysisEnded = true
		return
	}
	if p.play[t] == p.moves[t]-1 { // last move
		p.next(t - 1)
	} else {
		p.play[t] += 1
		for i := t + 1; i < 10; i++ {
			p.play[i] = 0
		}
	}
}

func (p *APlayer) calculateHighestBid(afterSkat bool) int {
	if p.forcedBid == 0 {
		return p.Player.calculateHighestBid(afterSkat)
	}
	p.highestBid = p.forcedBid
	return p.forcedBid
}

func nextGameForAnalysis() {
	tIndex = -1
}
