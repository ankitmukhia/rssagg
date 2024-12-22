-- +goose Up
ALTER TABLE users ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
	/*-- Chatgpted --*/
	/* Note: random(), generates random num between 0 and 1 */
	/* Note: ::text, type casting, casting result of random, which from float to string  */
	/* Note: ::bytea, another type cast, convert text back into byte of array  */
	/* Note: CG hash, takes binary data, & spit out SHA-256 hash, which is also binary data */
	/* Note: convert binary data(byte array), into hex string */
	encode(sha256(random()::text::bytea), 'hex')
);

-- +goose Down
ALTER TABLE users DROP COLUMN api_key;
