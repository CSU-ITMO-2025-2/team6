#!/bin/bash

source .env

echo "Waiting for database to be ready..."

# Формируем DSN для миграций
export MIGRATION_DSN="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_ADDR/$POSTGRES_DB?sslmode=disable"

echo "Running migrations with DSN: postgres://$POSTGRES_USER:****@db:5432/$POSTGRES_DB?sslmode=disable"

migrate -path "${MIGRATION_DIR}" -database "${MIGRATION_DSN}" up

if [ $? -eq 0 ]; then
    echo "Migrations completed successfully"
else
    echo "Migrations failed"
    exit 1
fi