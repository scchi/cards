package data

import (
	"encoding/json"
	"fmt"
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

func GenerateCards() []Card {
	var result []Card

	values := []string{
		"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
	}
	suits := []string{
		"S", "D", "C", "H",
	}

	for _, suit := range suits {
		for _, value := range values {
			code := fmt.Sprintf("%s%s", value, suit)
			card := Card(code)
			result = append(result, card)
		}
	}

	return result
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
