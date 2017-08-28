package main

import (
	_ "fmt"
)

type APlayer struct {
	Player
	moves []int
	play []int
}

func makeAPlayer(hand []Card) APlayer {
	return APlayer{
		Player:     makePlayer(hand),
		moves: []int{0,0,0,0,0,0,0,0,0,0},
		play: []int{0,0,0,0,0,0,0,0,0,0},
	}
}

var analysisEnded = false
var	tIndex = -1

func (p *APlayer) playerTactic(s *SuitState, c []Card) Card {
	gameLog("Valid cards: %v\n", c)
	tIndex++
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
	if p.play[t] == p.moves[t] - 1 { // last move
		p.next(t-1)
	} else {
		p.play[t] += 1
		for i := t + 1; i < 10; i++ {
			p.play[i] = 0
		}
	}
}

func nextGameForAnalysis() {
	tIndex = -1
}