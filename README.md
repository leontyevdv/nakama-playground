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

### Run containers

Run the docker-compose command to build and run a container with Nakama that includes a custom RPC function.
This command will run a PostgreSQL container.

```bash
docker-compose up -d --build nakama
```

To check if everything works as expected use the following command:

```bash
docker ps
```

You should see two containers in the output: nakama-playground-nakama and postgres:12.2-alpine.

Run the following command to check logs:

```bash
docker logs leontyevdv_nakama_backend
```

You should see the following line in there:

```bash
{"level":"info","ts":"2024-05-05T20:01:34.313Z","caller":"main.go:204","msg":"Startup done"}
```

### Run a migration

To create DB tables I use a [golang-migrate](https://github.com/golang-migrate/migrate) tool. This is supposed to be a
separate step in the CD pipeline.

Migration files reside in the ./migrations folder. Run the following command to apply them:

```bash
docker run -v ./migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:localdb@localhost/nakama?sslmode=disable" up
```

You should see something like this in Terminal:

```bash
20240504095304/u create_core_table (6.7065ms)
20240504095407/u create_score_table (11.224ms)
```

## How to test

The test_request.sh bash script (resides in the project folder) authenticates the user and calls the RPC function. Use
one of the following commands to test different possible scenarios.

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

- When creating migrations, there should be rollback migrations defined.
- Externalize configuration. All the credentials and configurations should be provided to the app through the
  environment variables or a config-server. Store all the credentials from a secret vault.
- Separate business logic from the Nakama implementation details. Create layers (such as Controller/Service/Repository)
  or a hexagonal architecture where the business logic is a core whereas DB and Nakama are outgoing and incoming ports
- Use https://github.com/ascii8/nktest to write Nakama-specific tests