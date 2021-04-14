package main

import (
	"testing"
)

func TestGrandLosers1(t *testing.T) {
	cards := []Card{
		{CLUBS, "A"},
		{CLUBS, "10"},
		{CLUBS, "K"},

		{CLUBS, "9"}, //LOSER

		{SPADE, "A"},
		{SPADE, "K"},  //LOSER
		{HEART, "10"}, //LOSER
	}

	losers := grandLosers(cards)
	if !in(losers, Card{SPADE, "K"}) {
		// , Card{HEART, "10"}) {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosers2(t *testing.T) {
	cards := []Card{
		{SPADE, "A"},
		{SPADE, "10"},
		{SPADE, "K"},

		{SPADE, "8"},
		{SPADE, "7"},

		{CLUBS, "K"},

		{HEART, "10"}, // will be discarded!

		{CARO, "9"},
		{CARO, "8"},
	}

	losers := grandLosers(cards)
	if len(losers) != 3 {
		t.Errorf("Wrong losers: %v", losers)
	}
}

func TestGrandLosers3(t *testing.T) {
	cards := []Card{
		{SPADE, "A"},
		{SPADE, "10"},
		{SPADE, "9"},
		{SPADE, "8"},
		{SPADE, "7"},
	}

	losers := grandLosers(cards)
	if len(losers) != 0 {
		t.Errorf("Wrong losers: %v", losers)
	}
}

func TestGrandLosers4(t *testing.T) {
	cards := []Card{
		{SPADE, "A"},
		{SPADE, "K"},
		{SPADE, "Q"},
		{SPADE, "9"},
		{SPADE, "8"},
		{SPADE, "7"},
	}

	losers := grandLosers(cards)
	if len(losers) != 0 {
		t.Errorf("Wrong losers: %v", losers)
	}
}

func TestGrandLosersCount(t *testing.T) {
	cs := []Card{
		{SPADE, "J"}, //LOSER
		{CLUBS, "A"},
		{CLUBS, "8"}, //LOSER
		{SPADE, "A"},
		{SPADE, "10"},
		{SPADE, "9"}, //LOSER
		{HEART, "A"},
		{HEART, "10"},
		{HEART, "K"},
		{CARO, "9"}, //LOSER
	}

	losers := len(grandLosers(cs)) + jackLosers(cs)
	if losers != 4 {
		t.Errorf("Losers: %v, Expected 4", losers)
	}
}

func TestGrandLosersJ0(t *testing.T) {
	cards := []Card{
		{SPADE, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJNo(t *testing.T) {
	cards := []Card{
		{CLUBS, "J"},
	}

	losers := jackLosers(cards)
	if losers != 0 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ1(t *testing.T) {
	cards := []Card{
		{CLUBS, "J"},
		{SPADE, "J"},
	}

	losers := jackLosers(cards)
	if losers != 0 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ2(t *testing.T) {
	cards := []Card{
		{CLUBS, "J"},
		{HEART, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ3(t *testing.T) {
	cards := []Card{
		{SPADE, "J"},
		{HEART, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ4(t *testing.T) {
	cards := []Card{
		{SPADE, "J"},
		{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 2 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ5(t *testing.T) {
	cards := []Card{
		{HEART, "J"},
		{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 2 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ6(t *testing.T) {
	cards := []Card{
		{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ7(t *testing.T) {
	cards := []Card{
		{HEART, "J"},
		{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 2 {
		t.Errorf("Losers: %v", losers)
	}
}
