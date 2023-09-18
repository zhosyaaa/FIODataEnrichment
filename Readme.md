# FIO Enrichment Service

This repository contains an example implementation of a Go service for enriching information about individuals based on their full names (FIO)
and storing this information in a PostgreSQL database. The service also provides REST and GraphQL APIs for accessing the data and supports data caching in Redis.

---

## Installation and Setup
Requirements
Before getting started, ensure you have the following components installed:

- Go (version 1.16 or higher)
- PostgreSQL
- Redis
- Docker (optional, for convenient dependency setup)

Clone the Repository

    git clone https://github.com/zhosyaaa/FIODataEnrichment.git 
    cd FIODataEnrichment

Install Dependencies

    go mod tidy

## Configuration

Create a '.env' file in the project's root directory and set the following environment variables:

      POSTGRES_URL=postgres://postgres:1079@localhost:5432/TestCase
      DB_CONNECTION_STRING=postgresql://postgres:1079@localhost:5432/TestCase?sslmode=disable
      KAFKA_BROKER_URL=localhost:9021
      REDIS_ADDR=localhost:6379
      API_KEY=your_api_key_here

## Usage
### REST API
The service provides the following REST API endpoints:

1. Adding New People
'POST /people'
Example request body:

        {
        "name": "Dmitriy",
        "surname": "Ushakov",
        "patronymic": "Vasilevich"
        }
2. Fetching Data with Various Filters and Pagination
'GET /people'
You can use query parameters to filter and paginate the results. 
3. Deleting by Identifier
'DELETE /people/{id}' 
4. Updating an Entity
'PUT /people/{id}' 
5. Retrieving Error Information
'GET /errors'

## GraphQL API
The service also provides a GraphQL API accessible at 'http://localhost:8080/graphql'. You can use the GraphiQL interface to execute queries.

Example GraphQL query for adding a new person:

    mutation {
        createPerson(input: {
            name: "Dmitriy",
            surname: "Ushakov",
            patronymic: "Vasilevich"
        }) {
          id
        }
    }

## Caching
Data is cached in Redis to improve performance. Caching is configured using environment variables in the .env file.

## Logging
Logging is done using the standard Go log package. Logs are written to the app.log file in the project's root directory.

## Testing
You can run unit tests using the following command:

    go test ./...
