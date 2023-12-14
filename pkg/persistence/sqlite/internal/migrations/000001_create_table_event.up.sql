CREATE TABLE IF NOT EXISTS event
(
    row_id         INTEGER PRIMARY KEY AUTOINCREMENT,
    stream_name    TEXT             NOT NULL,
    stream_version UNSIGNED BIG INT NOT NULL,
    event_id       TEXT             NOT NULL,
    event_name     TEXT             NOT NULL,
    event_data     BLOB             NOT NULL,
    happened_on    TIMESTAMP        NOT NULL,
    content_type   TEXT             NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS stream_name_version_unique_idx ON event (stream_name, stream_version);
CREATE UNIQUE INDEX IF NOT EXISTS event_id_unique_idx ON event (event_id);
