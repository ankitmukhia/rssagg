-- +goose Up
CREATE TABLE follow_feed (
	id uuid PRIMARY KEY,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL, 
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
	UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE follow_feed;
