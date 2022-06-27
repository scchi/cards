package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/scchi/cards/internal/validator"
)

type Deck struct {
	ID          string    `json:"deck_id,omitempty"`
	Shuffled    bool      `json:"shuffled"`
	Remaining   int       `json:"remaining,omitempty"`
	Cards       []Card    `json:"cards,omitempty"`
	StringCards []string  `json:"-"`
	CreatedAt   time.Time `json:"-"`
	Version     int       `json:"-"`
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

	fmt.Printf("%+v", result)
	return result
}

func ValidateDeck(v *validator.Validator, deck *Deck) {
	v.Check(validator.Unique(deck.Cards), "cards", "must not contain duplicated values")
	v.Check(len(deck.Cards) <= 52, "cards", "must not contain more than 52 cards")

	permittedCards := GenerateCards()
	v.Check(validator.PermittedValues(deck.Cards, permittedCards), "cards", "contains invalid card")
}

type DeckModel struct {
	DB *sql.DB
}

func (d DeckModel) Insert(deck *Deck) error {
	query := `
		INSERT INTO decks (shuffled, cards)
		VALUES ($1, $2)
		RETURNING id, created_at`

	args := []any{deck.Shuffled, pq.Array(deck.Cards)}
	return d.DB.QueryRow(query, args...).Scan(&deck.ID, &deck.CreatedAt)
}

func (d DeckModel) Get(id string) (*Deck, error) {
	query := `
		SELECT id, shuffled, cards
		FROM decks
		WHERE id::text = $1`

	var deck Deck

	err := d.DB.QueryRow(query, id).Scan(
		&deck.ID,
		&deck.Shuffled,
		pq.Array(&deck.StringCards),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &deck, nil
}

func (d DeckModel) Update(deck *Deck) error {
	query := `
		UPDATE decks
		SET cards = $1, version = version + 1
		WHERE id::text = $2
		RETURNING version`

	args := []any{
		pq.Array(deck.Cards),
		deck.ID,
	}

	return d.DB.QueryRow(query, args...).Scan(&deck.Version)
}

func (d DeckModel) Draw(deck *Deck) error {
	return nil
}

// -------------------------------------------------

type MockDeckModel struct{}

func (m MockDeckModel) Insert(deck *Deck) error {
	deck.ID = "a23d446a-f01a-4d6e-bec3-f928a3457ac7"
	return nil
}

func (m MockDeckModel) Get(id string) (*Deck, error) {
	if id == "wrongid" {
		return nil, ErrRecordNotFound
	}

	deck := Deck{
		ID:       id,
		Shuffled: true,
		StringCards: []string{
			"AS",
			"9D",
		},
	}

	return &deck, nil
}

func (m MockDeckModel) Update(deck *Deck) error {
	return nil
}

func (m MockDeckModel) Draw(deck *Deck) error {

	return nil
}
