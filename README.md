# Movies API
A REST API built with Go and SQLite for managing a movie database.

## Project Overview
The API supports CRUD (Create, Read, Update, Delete) operations on three independently managed entities: `movies`, `genres`, and `actors`.

Movies can have many-to-many relationships with genres and actors, and functionality exists to create these relationships.

Extra features include filtering, search, and pagination.

 The [`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) driver is used to  interface with SQLite.
 Validation libary, [`go-playground/validator`](https://pkg.go.dev/github.com/go-playground/validator), is used for input validation of struct fields.
## Setup Instructions
Clone the respository:
```bash
git clone https://gitea.kood.tech/georgiisenotrusov/movies-api.git
```
## Usage Guide
Start the server, using default database `movies.db`:
```bash
cd movies-api
go run .
```
By default, the server starts on http://localhost:8080.

### CLI Flags
| Flag | Use | Default value |
| ---- | ---- | ------------ |
| `db` | database file name | `movies.db` |
| `reset` | resets database and seeds it with dummy data for testing | `false` |

### Example usage
Use default `movies.db` database, resetting it with dummy data.
```bash
go run . -reset
```

---------------------


Use `my-database.db`. If the database file does not already exist, it will be created and initialized with all necessary tables. No dummy data added.
```bash
go run . -db="my-database.db"
```
---------------------

Use `my-database.db` *and* reset it with dummy data.
```bash
go run . -db="my-database.db" -reset
```

## API Usage

### Basic CRUD operations
The following endpoints can be used with all entities, where `{entity}` is `movies`, `genres`, or `actors`. `{id}` must be a positive integer.
| Method | Endpoint | Description |
|--------|----------|-------------|
| **POST**    | `/api/{entity}` | Create a new entity |
| **GET**    | `/api/{entity}` | Retrieve all entities (within pagination and filter parameters)|
| **GET** | `/api/{entity}/{id}` | Retrieve a specific entity by ID |
| **PATCH**    | `/api/{entity}/{id}` | Partially update an existing entity |
| **DELETE** | `/api/{entity}/{id}` | Delete an entity |

### Filtering Movies
When retrieving movies using endpoint `GET /api/movies`, the following query parameters can be used as filters:
| Query parameter |  Description |
|-----------------|-------------|
| year     | Filter movies by release year |
| genre     | Filter movies by genre |
| actor     | Filter movies by actor (featuring the actor with the given ID) |

Query parameters can be combined, including both filter and pagination parameters.

#### Examples
Get movies released in 1999:
`GET /api/movies?year=1999`

Get movies released in 1999 and featuring actor with ID 4: `GET /api/movies?year=1999&actor=16`


### Other filters
**Retrieve all actors in a movie:** `GET /api/movies/{movieId}/actors`

**Retrieve all actors filtered by name:** `GET /api/actors?name={name}`

## Pagination
Pagination is implemented for `GET` requests returning multiple entities.

| Query parameter |  Description |
|-----------------|-------------|
| page     | Specifies which page of results to return. Starts at 0. |
| size     | Results per page. Max size = 100. |

**Example:** `GET /api/movies?page=0&size=10`

If pagination parameters are not specified, the result is automatically paginated with default parameters of `page=0` and `size=10`.

## Required Format for POST & PATCH Requests
This section shows the required formats for different entities when **adding** or **updating** an entity with a `POST` or `PATCH` request. The request body must be in JSON.

### Movie submission
Example JSON request body:
```json
{
  "title": "The Great Dawn",
  "release_year": 2026,
  "duration": 120,
  "genre_ids": [1, 4],
  "actor_ids": [5, 7, 19, 22]
}
```
**Required fields for adding a movie:** `title`, `release_year`, `duration`

### Genre submissions
Example JSON request body:
```json
{
  "name": "Romantic Comedy"
}
```

### Actor submissions
Example request body:
```json
{
  "name": "Bob Smith",
  "birth_date": "1991-04-23",
  "movie_ids": [12, 19]
}
```
**Required fields for adding an actor:** `name`, `birth_date`

## Bonus Features
- **Pagination** for GET requests returning multiple entities, e.g. `GET /api/movies?page=0&size=10`
- **Search** movies by title (case-insensitive, partial match search), e.g. `GET /api/movies/search?title=last`
- **Prevention of SQL injection attacks** by using placeholder parameters

- Use of **transactions** to execute multiple SQL statements in one atomic action to prevent partial database updates during execution failures

- Use of **context** to end functions early if client disconnects

## Team Members
[Ellin Park](https://github.com/Ellin), [PLACEHOLDER](), [PLACEHOLDER]()