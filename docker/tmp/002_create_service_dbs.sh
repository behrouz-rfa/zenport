#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE commondb TEMPLATE template0;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "commondb" <<-EOSQL
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



psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE ntps TEMPLATE commondb;

  CREATE USER ntps_user WITH ENCRYPTED PASSWORD 'ntps_pass';
  GRANT USAGE ON SCHEMA public TO ntps_user;
  GRANT CREATE, CONNECT ON DATABASE ntps TO ntps_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "ntps" <<-EOSQL
  CREATE SCHEMA ntps;
  GRANT CREATE, USAGE ON SCHEMA ntps TO ntpst_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE notifications TEMPLATE commondb;

  CREATE USER notifications_user WITH ENCRYPTED PASSWORD 'notifications_pass';
  GRANT USAGE ON SCHEMA public TO notifications_user;
  GRANT CREATE, CONNECT ON DATABASE notifications TO notifications_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "notifications" <<-EOSQL
  CREATE SCHEMA notifications;

  GRANT CREATE, USAGE ON SCHEMA notifications TO notifications_user;
EOSQL

