#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE zenports;

  CREATE USER zenports_user WITH ENCRYPTED PASSWORD 'zenports_pass';

  GRANT CONNECT ON DATABASE zenports TO zenports_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "zenports" <<-EOSQL
  -- Apply to keep modifications to the created_at column from being made
  CREATE OR REPLACE FUNCTION created_at_trigger()
  RETURNS TRIGGER AS \$\$
  BEGIN
    NEW.created_at := OLD.created_at;
    RETURN NEW;
  END;
  \$\$ language 'plpgsql';

  -- Apply to a table to automatically update update_at columns
  CREATE OR REPLACE FUNCTION updated_at_trigger()
  RETURNS TRIGGER AS \$\$
  BEGIN
     IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
        NEW.updated_at = NOW();
        RETURN NEW;
     ELSE
        RETURN OLD;
     END IF;
  END;
  \$\$ language 'plpgsql';
EOSQL


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