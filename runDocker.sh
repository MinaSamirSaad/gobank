#!/bin/bash


# Check if the PostgreSQL image exists locally
if [[ "$(docker images -q postgres:latest 2> /dev/null)" == "" ]]; then
    echo "PostgreSQL image not found locally. Pulling the image..."
    docker pull postgres:latest
else
    echo "PostgreSQL image already exists locally."
fi

# Check if PostgreSQL container is running
if [[ $(docker ps -q -f name=some-postgres) ]]; then
    echo "Container 'some-postgres' is already running."
else
    if [[ $(docker ps -a -q -f name=some-postgres) ]]; then
        echo "Container 'some-postgres' exists but is stopped. Restarting it..."
        docker start some-postgres
    else
        echo "Container 'some-postgres' does not exist. Creating and running a new one..."
        docker run --name some-postgres -e POSTGRES_PASSWORD=gobank -p 5432:5432 -d postgres
    fi
fi

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until docker exec some-postgres pg_isready -U postgres > /dev/null 2>&1; do
    sleep 1
done

echo "PostgreSQL is ready!"

