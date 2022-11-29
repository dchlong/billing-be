### How to run the code locally
Make sure you have docker installed.

Clone the project:
```txt
git clone https://github.com/dchlong/billing-be.git
```

Update deployments/base.env file

```txt
HTTP_ADDR=:8088
DATABASE_CONFIG_DATA_SOURCE=root:change_me@tcp(localhost:3306)/billing?timeout=10s&parseTime=true
NUMBER_OF_SECONDS_IN_A_BLOCK=30
```

Run command:

```bash
./scripts/bin.sh run
```

##### or by Docker

Update docker-compose.local.yml file (update port)

```txt
...
    environment:
      - HTTP_ADDR=:80
      - DATABASE_CONFIG_DATA_SOURCE=root:change_me@tcp(billing_db:3306)/billing?timeout=10s&parseTime=true
    ports:
      - 8088:80
...
```

Run command:

```bash
./scripts/bin.sh run docker
```

### Then open postman:

#### Create call history

```bash
curl --location --request PUT 'http://localhost:8088/mobile/dchlong/call' \
--header 'Content-Type: application/json' \
--data-raw '{
    "call_duration": 60000
}'
```

Responses
```json
{
    "id": 8
}
```

```json
{
    "error": {
        "code": "invalid_input",
        "status_code": 400,
        "message": "call_duration is a required field"
    }
}
```

#### Get bill

```bash
curl --location --request GET 'http://localhost:8088/mobile/dchlong/billing'
```

Responses
```json
{
    "call_count": 5,
    "block_count": 3
}
```

```json
{
    "error": {
        "code": "invalid_user_name",
        "status_code": 400,
        "message": "username must be from 1-32 characters"
    }
}
```

### Tests
Command:

```bash
 ./scripts/bin.sh test
```

#### Unit test

Command:

```bash
 ./scripts/bin.sh unit_test
```

#### Integration test

Command:

```bash
./scripts/bin.sh integration_test
```


#### Generate mock and initialize function
Command:

```bash
 ./scripts/bin.sh generate
```

#### Linting
Linting is the automated checking of your source code for programmatic and stylistic errors

Command:

```bash
 ./scripts/bin.sh lint
```