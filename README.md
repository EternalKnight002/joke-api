## 🃏 Joke API

![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go\&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-yellow.svg)
![Build](https://img.shields.io/badge/build-passing-brightgreen.svg)

A simple and lightweight **REST API built with Go** that serves random programming jokes.
Perfect for learning backend development, exploring Go’s standard library, or integrating into fun frontend projects.

---

### 🚀 Features

* RESTful endpoints (`GET`, `POST`, `DELETE`)
* Random joke endpoint → `/jokes/random`
* Thread-safe in-memory store (no database needed)
* Middleware for clean request logging
* `/favicon.ico` silencer (no unwanted logs)
* Configurable port (`PORT` env var) — defaults to **8081**
* Built entirely with Go’s standard library (no external dependencies)

---

### 📁 Folder Structure

```
joke-api/
├── go.mod
└── main.go
```

---

### 🧰 Requirements

* Go **1.20+**
* Internet (only needed once for `go mod tidy`)

---

### ⚙️ Setup & Run

```bash
# 1. Clone the repo
git clone https://github.com/EternalKnight002/joke-api.git
cd joke-api

# 2. Install dependencies
go mod tidy

# 3. Run the API
go run .
```

You’ll see:

```
Joke API running at :8081
```

Then open in your browser:

* ➜ [http://localhost:8081/jokes](http://localhost:8081/jokes)
* ➜ [http://localhost:8081/jokes/random](http://localhost:8081/jokes/random)

---

### 🔌 Environment Variables

| Variable | Description                | Default |
| -------- | -------------------------- | ------- |
| `PORT`   | Port number for the server | `8081`  |

Example:

```bash
PORT=8082 go run .
```

---

### 📡 API Endpoints

| Method   | Endpoint        | Description                   |
| -------- | --------------- | ----------------------------- |
| `GET`    | `/jokes`        | Returns all jokes             |
| `GET`    | `/jokes/random` | Returns one random joke       |
| `GET`    | `/jokes/{id}`   | Returns a specific joke by ID |
| `POST`   | `/jokes`        | Adds a new joke               |
| `DELETE` | `/jokes/{id}`   | Deletes a joke by ID          |

---

### 🧪 Example Usage

#### Get all jokes

```bash
curl -s http://localhost:8081/jokes | jq
```

#### Get a random joke

```bash
curl -s http://localhost:8081/jokes/random | jq
```

#### Add a new joke

```bash
curl -X POST http://localhost:8081/jokes \
  -H "Content-Type: application/json" \
  -d '{"content":"Go programmers never panic; they recover.", "author":"you"}' | jq
```

#### Delete a joke

```bash
curl -X DELETE http://localhost:8081/jokes/3 -v
```

---

### 🧠 Example JSON Output

```json
{
  "id": 2,
  "content": "Why do programmers prefer dark mode? Because light attracts bugs.",
  "author": "classic",
  "created_at": "2025-10-07T11:41:27Z"
}
```

---

### 🧱 How It Works

* **Thread-safe store:** Uses Go’s `sync.RWMutex` to safely handle concurrent requests.
* **Random jokes:** Stored and retrieved from memory using Go’s `math/rand`.
* **Clean architecture:** Each endpoint is minimal and focused.
* **Logging middleware:** Logs method, path, and response time for every request.

---

### 🎯 Future Enhancements

* [ ] Add persistence (SQLite or JSON file)
* [ ] Add `/health` endpoint
* [ ] Add simple web frontend to show random jokes
* [ ] Add Dockerfile for containerization
* [ ] Add automated tests

---

### 🏁 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
