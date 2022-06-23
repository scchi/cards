package data

import "time"

type Deck struct {
	ID        int64     `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	Cards     []Card    `json:"cards"`
	CreatedAt time.Time `json:"-"`
}
