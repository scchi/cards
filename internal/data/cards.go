package data

import (
	"encoding/json"
	"strconv"
)

type Card string

var suits = map[string]string{
	"S": "SPADES",
	"D": "DIAMONDS",
	"C": "CLUBS",
	"H": "HEARTS",
}

var values = map[string]string{
	"A": "ACE",
	"J": "JACK",
	"Q": "QUEEN",
	"K": "KING",
}

func (c Card) MarshalJSON() ([]byte, error) {
	value := c.GetValue()
	suit := c.GetSuit()

	jsonValue := map[string]string{
		"value": value,
		"suit":  suit,
		"code":  string(c),
	}

	js, _ := json.Marshal(jsonValue)

	return js, nil
}

func (c Card) GetSuit() string {
	return suits[string(c[1])]
}

func (c Card) GetValue() string {
	value := string(c[0])
	if number, _ := strconv.Atoi(value); number == 0 {
		value = values[value]
	}

	return value
}
