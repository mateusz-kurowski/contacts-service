#!/bin/bash

atlas schema apply \
  -u "${DB_URL}" \
  --to file://sql/schema.sql \
  --dev-url "docker://postgres/18/contacts_db"