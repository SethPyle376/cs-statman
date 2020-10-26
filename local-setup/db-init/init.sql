CREATE TABLE match (
    id BIGINT PRIMARY KEY,
    map TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    roundCount INT NOT NULL
);

CREATE TABLE player (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE match_player (
    id SERIAL PRIMARY KEY,
    matchID BIGINT NOT NULL,
    userID BIGINT NOT NULL,
    hltv REAL NOT NULL,
    kills INT NOT NULL,
    deaths INT NOT NULL,
    adr REAL NOT NULL
);