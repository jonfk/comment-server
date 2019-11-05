
CREATE TABLE IF NOT EXISTS events (
       eventId UUID PRIMARY KEY,
       event_type TEXT NOT NULL,
       timestamp TIMESTAMP WITH TIME ZONE,
       data JSONB
);


CREATE TABLE IF NOT EXISTS accounts (
       account_id UUID PRIMARY KEY,
       username TEXT UNIQUE NOT NULL,
       email TEXT UNIQUE,
       hashed_password BYTEA NOT NULL,
       hash_salt BYTEA NOT NULL,
       created_on TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS comment_threads (
       comment_thread_id UUID PRIMARY KEY,
       created_on TIMESTAMP WITH TIME ZONE,
       --website TEXT REFERENCES websites(url) NOT NULL,
       page_url TEXT NOT NULL,
       title TEXT
);

CREATE TABLE IF NOT EXISTS comments (
       comment_id UUID PRIMARY KEY,
       timestamp TIMESTAMP WITH TIME ZONE,
       data TEXT,
       parent_id UUID REFERENCES comments(comment_id),
       comment_thread_id UUID REFERENCES comment_threads(comment_thread_id),
       account_id UUID REFERENCES accounts(account_id)
);

CREATE TABLE IF NOT EXISTS websites (
       url TEXT PRIMARY KEY
);
