package main

import (
	"testing"
)


func TestGrandLosers1(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "K"},

		Card{CLUBS, "9"}, //LOSER

		Card{SPADE, "A"},
		Card{SPADE, "K"},  //LOSER
		Card{HEART, "10"}, //LOSER
	}

	losers := grandLosers(cards)
	if !in(losers,Card{SPADE, "K"}) {
		// , Card{HEART, "10"}) {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosers2(t *testing.T) {
	cards := []Card{
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},

		Card{SPADE, "8"},  
		Card{SPADE, "7"},  

		Card{CLUBS, "K"},

		Card{HEART, "10"},  // will be discarded!
		
		Card{CARO, "9"},
		Card{CARO, "8"},
	}

	losers := grandLosers(cards)
	if len(losers) != 3 {
		t.Errorf("Wrong losers: %v", losers)
	}
}

func TestGrandLosers3(t *testing.T) {
	cards := []Card{
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "9"},  
		Card{SPADE, "8"},  
		Card{SPADE, "7"},  
	}

	losers := grandLosers(cards)
	if len(losers) != 0 {
		t.Errorf("Wrong losers: %v", losers)
	}
}

func TestGrandLosers4(t *testing.T) {
	cards := []Card{
		Card{SPADE, "A"},
		Card{SPADE, "K"},
		Card{SPADE, "Q"},
		Card{SPADE, "9"},  
		Card{SPADE, "8"},  
		Card{SPADE, "7"},  
	}

	losers := grandLosers(cards)
	if len(losers) != 0 {
		t.Errorf("Wrong losers: %v", losers)
	}
}

func TestGrandLosersCount(t *testing.T) {
	cs := []Card{
		Card{SPADE, "J"}, //LOSER
		Card{CLUBS, "A"},
		Card{CLUBS, "8"}, //LOSER
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "9"}, //LOSER
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "K"},
		Card{CARO, "9"}, //LOSER
	}

	losers := len(grandLosers(cs)) + jackLosers(cs)
	if losers != 4 {
		t.Errorf("Losers: %v, Expected 4", losers)
	}
}

func TestGrandLosersJ0(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJNo(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
	}

	losers := jackLosers(cards)
	if losers != 0 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ1(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
	}

	losers := jackLosers(cards)
	if losers != 0 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ2(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ3(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ4(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"},
		Card{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 2 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ5(t *testing.T) {
	cards := []Card{
		Card{HEART, "J"},
		Card{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 2 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ6(t *testing.T) {
	cards := []Card{
		Card{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ7(t *testing.T) {
	cards := []Card{
		Card{HEART, "J"},
		Card{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 2 {
		t.Errorf("Losers: %v", losers)
	}
}
