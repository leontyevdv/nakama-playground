# Nakama Playground

## Description

This playground project contains a custom RPC function written with Go
for [Nakama](https://github.com/heroiclabs/nakama).

The RPC function implements the following requirements:

- Accepts payload with type, version, hash (all parameters are optional, defaults: type=core, version=1.0.0, hash=null).
- Reads a file from the disk (path=\<type\>/\<version\>.json)
- Saves data to the Posgresql DB
- Calculates file content hash
- Responds with the following fields: type, version, hash, content
- Returns null-content if a requested and calculated hashes are not equal
- Returns an error if file doesn't exist
- Uses defaults if they are not present in the payload

The custom RPC functions is covered by unit tests with the usage of [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
lib.

## Requirements

- Docker/docker-compose
- jq
- curl

## How to run

```bash
docker-compose up -d --build nakama
```

```bash
leontyevdv_nakama_backend | {"level":"info","ts":"2024-05-04T14:39:38.826Z","caller":"main.go:204","msg":"Startup done"}
```

https://github.com/golang-migrate/migrate

```bash
docker run -v ./migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:localdb@localhost/nakama?sslmode=disable" up
```

```bash
20240504095304/u create_core_table (6.7065ms)
20240504095407/u create_score_table (11.224ms)
```

## How to test

The test_request.sh bash script authenticates the user and calls the RPC function. Use one of the following commands to
test different possible scenarios.

```bash
./test_request.sh callScore
./test_request.sh callScoreDefaultVersion
./test_request.sh callScoreDefaultHash
./test_request.sh callScoreMissingFile
./test_request.sh callWithDefaultType
./test_request.sh callWithUnknownType
./test_request.sh callWithEmptyPayload
```

## How to stop and remove everything

```bash
docker compose down -v --remove-orphans
```

## Possible improvements and production considerations

- Down for rollback
- Externalize configuration
- Use https://github.com/ascii8/nktest
- Separate business logic from the Nakama implementation details
- Create layers (such as Controller/Service/Repository) or a hexagonal architecture where the business logic is a core
  whereas DB and Nakama are outgoing and incoming ports 
- Externalize configuration
- Store and request all credentials from a secret vault 