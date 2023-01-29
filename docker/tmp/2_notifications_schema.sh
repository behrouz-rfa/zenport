#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "zenports" <<-EOSQL
  CREATE SCHEMA notifications;

  CREATE TABLE notifications.ntp_cache
  (
      id         text NOT NULL,
      time       text NOT NULL,

      created_at timestamptz NOT NULL DEFAULT NOW(),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  GRANT USAGE ON SCHEMA notifications TO zenports_user;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA notifications TO zenports_user;
EOSQL
