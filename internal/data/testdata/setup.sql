CREATE TABLE IF NOT EXISTS decks (
  id uuid DEFAULT uuid_generate_v4 (),
  shuffled boolean,
  cards varchar(3)[],
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1,
  PRIMARY KEY (id)
);

ALTER TABLE decks ADD CONSTRAINT cards_length_check CHECK (array_length(cards, 1) BETWEEN 0 AND 52);

INSERT INTO decks (shuffled, cards) VALUES (TRUE, '{"AS", "5D", "KC", "2C", "AH"}');