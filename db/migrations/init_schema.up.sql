CREATE TABLE "Videos" (
        id uuid PRIMARY KEY,
        youtube_id varchar UNIQUE,
        title varchar,
        description varchar,
        created_at timestamptz NOT NULL DEFAULT (now()),
        updated_at timestamptz NOT NULL DEFAULT (now()),
        deleted_at timestamptz
);

-- index will optimise paginated query
CREATE INDEX id_created_at_idx ON Videos (id, created_at);

-- index will optimise search query
CREATE INDEX text_search_idx ON Videos USING GIN (to_tsvector('english', title || ' ' || description));
