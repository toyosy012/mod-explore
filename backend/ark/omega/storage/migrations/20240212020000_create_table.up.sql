CREATE TABLE "tests"
(
    id          SERIAL  PRIMARY KEY,
    name        VARCHAR NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

INSERT INTO tests(name) VALUES ('test');
