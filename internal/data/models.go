package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Decks interface {
		Insert(deck *Deck) error
		Get(id string) (*Deck, error)
		Update(deck *Deck) error
	}
	Data map[string]string
}

func NewModels(db *sql.DB) Models {
	return Models{
		Decks: DeckModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Decks: MockDeckModel{},
	}
}
