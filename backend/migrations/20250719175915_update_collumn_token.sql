ALTER TABLE user_tokens DROP COLUMN token;
ALTER TABLE user_tokens ADD COLUMN token text NOT NULL DEFAULT '';