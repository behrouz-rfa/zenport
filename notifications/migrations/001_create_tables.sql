CREATE TABLE notifications.ntp_cache
(
    id         text NOT NULL,
    time       text NOT NULL,

    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);
