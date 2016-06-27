package game_module

import "strings"

type Color int
type Suit int
type Value int

const (
	Red   Color = iota
	Black Color = iota
)

const (
	Hearts   Suit = iota
	Diamonds Suit = iota
	Spaces   Suit = iota
	Clubs    Suit = iota
)

const (
	Joker Value = iota
	Two   Value = iota
	Three Value = iota
	Four  Value = iota
	Five  Value = iota
	Six   Value = iota
	Seven Value = iota
	Eight Value = iota
	Nine  Value = iota
	Ten   Value = iota
	Jack  Value = iota
	Queen Value = iota
	King  Value = iota
)

type Card struct {
	Col Color
	Su  Suit
	V   Value
}

func NewCard(c Color, s Suit, v Value) *Card {
	return &Card{
		Col: c,
		Su:  s,
		V:   v,
	}
}

func (c *Card) GetString() string {
	getStr := []string{string(c.Col), string(c.V), "of", string(c.Su)}
	return strings.Join(getStr, " ")
}
