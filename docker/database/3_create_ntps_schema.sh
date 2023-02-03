#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "zenports" <<-EOSQL
  CREATE SCHEMA ntps;

  CREATE TABLE ntps.ntps
  (
      id            text NOT NULL,
      time          text NOT NULL,
      created_at    timestamptz NOT NULL DEFAULT NOW(),
      updated_at    timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );



  CREATE TABLE ntps.events
  (
      stream_id      text        NOT NULL,
      stream_name    text        NOT NULL,
      stream_version int         NOT NULL,
      event_id       text        NOT NULL,
      event_name     text        NOT NULL,
      event_data     bytea       NOT NULL,
      occurred_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (stream_id, stream_name, stream_version)
  );


  GRANT USAGE ON SCHEMA ntps TO zenports_user;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA ntps TO zenports_user;
EOSQL
