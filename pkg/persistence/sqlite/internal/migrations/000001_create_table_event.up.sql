CREATE TABLE IF NOT EXISTS event
(
    stream_name    TEXT      NOT NULL,
    stream_version INTEGER   NOT NULL,
    event_name     TEXT      NOT NULL,
    event_data     BLOB      NOT NULL,
    happened_on    TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS stream_name_version_unique_idx ON event (stream_name, stream_version);
