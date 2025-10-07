package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Joke represents a joke item.
type Joke struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// JokeStore is a threadsafe in-memory store.
type JokeStore struct {
	sync.RWMutex
	jokes map[int64]*Joke
	next  int64
	r     *rand.Rand
}

func NewJokeStore() *JokeStore {
	src := rand.NewSource(time.Now().UnixNano())
	return &JokeStore{
		jokes: make(map[int64]*Joke),
		next:  1,
		r:     rand.New(src),
	}
}

func (s *JokeStore) Create(content, author string) *Joke {
	s.Lock()
	defer s.Unlock()
	j := &Joke{
		ID:        s.next,
		Content:   strings.TrimSpace(content),
		Author:    strings.TrimSpace(author),
		CreatedAt: time.Now().UTC(),
	}
	s.jokes[j.ID] = j
	s.next++
	return j
}

func (s *JokeStore) GetAll() []*Joke {
	s.RLock()
	defer s.RUnlock()
	out := make([]*Joke, 0, len(s.jokes))
	for _, j := range s.jokes {
		out = append(out, j)
	}
	return out
}

func (s *JokeStore) Get(id int64) (*Joke, bool) {
	s.RLock()
	defer s.RUnlock()
	j, ok := s.jokes[id]
	return j, ok
}

func (s *JokeStore) Delete(id int64) bool {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.jokes[id]; !ok {
		return false
	}
	delete(s.jokes, id)
	return true
}

func (s *JokeStore) Random() (*Joke, error) {
	s.RLock()
	defer s.RUnlock()
	n := len(s.jokes)
	if n == 0 {
		return nil, errors.New("no jokes available")
	}
	// collect keys
	keys := make([]int64, 0, n)
	for k := range s.jokes {
		keys = append(keys, k)
	}
	// pick random index
	idx := s.r.Intn(n)
	return s.jokes[keys[idx]], nil
}

// JSON helpers
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	store := NewJokeStore()
	seed(store)

	mux := http.NewServeMux()

	// GET /jokes -> list
	// POST /jokes -> create { "content": "...", "author": "..." }
	mux.HandleFunc("/jokes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			all := store.GetAll()
			writeJSON(w, http.StatusOK, all)
		case http.MethodPost:
			var req struct {
				Content string `json:"content"`
				Author  string `json:"author,omitempty"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			req.Content = strings.TrimSpace(req.Content)
			if req.Content == "" {
				http.Error(w, "content is required", http.StatusBadRequest)
				return
			}
			j := store.Create(req.Content, req.Author)
			writeJSON(w, http.StatusCreated, j)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /jokes/random -> random joke
	mux.HandleFunc("/jokes/random", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		j, err := store.Random()
		if err != nil {
			http.Error(w, "no jokes available", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, j)
	})

	// GET /jokes/{id}  DELETE /jokes/{id}
	mux.HandleFunc("/jokes/", func(w http.ResponseWriter, r *http.Request) {
		// trim trailing slash and split
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 2 {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			j, ok := store.Get(id)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, j)
		case http.MethodDelete:
			ok := store.Delete(id)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := ":8081"
	fmt.Printf("Joke API running at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, loggingMiddleware(mux)))
}

// simple logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// seed with a few jokes so /jokes/random works immediately
func seed(s *JokeStore) {
	s.Create("I told my computer I needed a break, and it said: 'No problem — I'll go to sleep.'", "unknown")
	s.Create("Why do programmers prefer dark mode? Because light attracts bugs.", "classic")
	s.Create("There's no place like 127.0.0.1", "nerd")
}
