package data

import (
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/lib/pq"
	"github.com/scchi/cards/internal/validator"
)

type Deck struct {
	ID          string    `json:"deck_id"`
	Shuffled    bool      `json:"shuffled"`
	Remaining   int       `json:"remaining"`
	Cards       []Card    `json:"cards,omitempty"`
	StringCards []string  `json:"-"`
	CreatedAt   time.Time `json:"-"`
	Version     int       `json:"-"`
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

func ShuffleDeck(deck *Deck) {
	rand.Seed(time.Now().Unix())

	rand.Shuffle(len(deck.Cards), func(i, j int) {
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	})
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

// -------------------------------------------------

type MockDeckModel struct{}

var MockID = "a23d446a-f01a-4d6e-bec3-f928a3457ac7"
var MockShuffled = true
var MockCards = []string{
	"AS",
	"9D",
}

func (m MockDeckModel) Insert(deck *Deck) error {
	deck.ID = MockID
	return nil
}

func (m MockDeckModel) Get(id string) (*Deck, error) {
	if len(id) != 36 {
		return nil, ErrRecordNotFound
	}

	if id != MockID {
		return nil, ErrRecordNotFound
	}

	deck := Deck{
		ID:          id,
		Shuffled:    MockShuffled,
		StringCards: MockCards,
	}

	return &deck, nil
}

func (m MockDeckModel) Update(deck *Deck) error {
	return nil
}
