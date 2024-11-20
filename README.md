# CLI-Golang-Quiz

## Backend

Run the backend server with the following command:

```bash
cd backend
go run main.go
```

Retrieve questions with curl:

```bash
curl http://localhost:8080/questions
```

Submit answers with the following command:

```bash
curl -X POST http://localhost:8080/submission \
    -H "Content-Type: application/json" \
    -d '[
    {"question_id": 1, "answer": 0},
    {"question_id": 2, "answer": 1},
    {"question_id": 3, "answer": 2},
    {"question_id": 4, "answer": 0},
    {"question_id": 5, "answer": 0},
    {"question_id": 6, "answer": 1}
    ]'
```

## Play the Quiz using the CLI

Make sure you have the backend running. Run the following command from the root directory:

```bash
go run main.go start-quiz
```

and answer the questions as prompted.
