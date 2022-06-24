CREATE TABLE IF NOT EXISTS decks (
  id uuid DEFAULT uuid_generate_v4 (),
  shuffled boolean,
  cards varchar(2)[],
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1,
  PRIMARY KEY (id)
);