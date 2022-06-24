package data

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/scchi/cards/internal/validator"
)

type Card string

type Deck struct {
	ID        int64     `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	Cards     []Card    `json:"cards"`
	CreatedAt time.Time `json:"-"`
}

func ValidateDeck(v *validator.Validator, deck *Deck) {
	// validate shuffled exists. at the moment, missing shuffled means false

	// missing permitted values for cards
	v.Check(deck.Cards != nil, "cards", "must be provided")
	v.Check(validator.Unique(deck.Cards), "cards", "must not contain duplicated values")
	v.Check(len(deck.Cards) <= 52, "cards", "must not contain more than 52 cards")
}

type DeckModel struct {
	DB *sql.DB
}

func (d DeckModel) Insert(deck *Deck) error {
	query := `
		INSERT INTO decks (shuffled, remaining, cards)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	args := []any{deck.Shuffled, deck.Remaining, pq.Array(deck.Cards)}
	return d.DB.QueryRow(query, args...).Scan(&deck.ID, &deck.CreatedAt)
}

func (d DeckModel) Get(id int64) (*Deck, error) {
	return nil, nil
}

func (d DeckModel) Update(deck *Deck) error {
	return nil
}

func (d DeckModel) Draw(deck *Deck) error {
	return nil
}

type MockDeckModel struct{}

func (m MockDeckModel) Insert(deck *Deck) error {
	return nil
}

func (m MockDeckModel) Get(id int64) (*Deck, error) {
	return nil, nil
}

func (m MockDeckModel) Update(deck *Deck) error {
	return nil
}

func (m MockDeckModel) Draw(deck *Deck) error {
	return nil
}
