CREATE TABLE IF NOT EXISTS users (
    id varchar(256) PRIMARY KEY NOT NULL,
    name varchar(256) NOT NULL,
    email varchar(256) NOT NULL,
    ts_created DATETIME NOT NULL,
    ts_updated DATETIME NOT NULL
);