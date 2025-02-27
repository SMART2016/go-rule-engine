#!/bin/bash
while true; do
    sleep 5
    psql -U dbuser -d rule_engine < /docker-entrypoint-initdb.d/schema.sql
done
