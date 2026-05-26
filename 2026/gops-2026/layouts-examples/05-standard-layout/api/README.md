# Office Library API

Base URL: `http://localhost:8080`

## Users

| Method | Pattern          | Description      |
|--------|------------------|------------------|
| GET    | /users           | List all users   |
| GET    | /users/{id}      | Get user by ID   |
| POST   | /users           | Create a user    |

## Books

| Method | Pattern                | Description        |
|--------|------------------------|--------------------|
| GET    | /books                 | List all books     |
| GET    | /books/{id}            | Get book by ID     |
| POST   | /books/{id}/borrow     | Borrow a book      |
| POST   | /books/{id}/return     | Return a book      |

## Reviews

| Method | Pattern                 | Description          |
|--------|-------------------------|----------------------|
| GET    | /books/{id}/reviews     | List book reviews    |
| POST   | /reviews                | Create a review      |

### POST /books/{id}/borrow

```json
{"user_id": 1}
```

### POST /users

```json
{"name": "Charlie", "email": "charlie@example.com"}
```

### POST /reviews

```json
{"book_id": 1, "user_id": 1, "text": "Great book!", "rating": 5}
```
