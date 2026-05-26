package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/08-monorepo/libs/library"
)

func setup() (*Handler, *http.ServeMux) {
	store := library.NewStore()
	h := &Handler{Store: store}
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return h, mux
}

func TestListUsers(t *testing.T) {
	_, mux := setup()
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var users []library.User
	json.NewDecoder(w.Body).Decode(&users)
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestGetUser(t *testing.T) {
	_, mux := setup()
	req := httptest.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var u library.User
	json.NewDecoder(w.Body).Decode(&u)
	if u.Name != "Alice" {
		t.Fatalf("expected Alice, got %s", u.Name)
	}
}

func TestGetUserNotFound(t *testing.T) {
	_, mux := setup()
	req := httptest.NewRequest("GET", "/users/999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCreateUser(t *testing.T) {
	_, mux := setup()
	body := `{"Name":"Charlie","Email":"charlie@example.com"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var u library.User
	json.NewDecoder(w.Body).Decode(&u)
	if u.ID != 3 {
		t.Fatalf("expected ID 3, got %d", u.ID)
	}
}

func TestListBooks(t *testing.T) {
	_, mux := setup()
	req := httptest.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var books []library.Book
	json.NewDecoder(w.Body).Decode(&books)
	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}
}

func TestGetBook(t *testing.T) {
	_, mux := setup()
	req := httptest.NewRequest("GET", "/books/1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var b library.Book
	json.NewDecoder(w.Body).Decode(&b)
	if b.Title != "The Go Programming Language" {
		t.Fatalf("unexpected title: %s", b.Title)
	}
}

func TestBorrowAndReturn(t *testing.T) {
	_, mux := setup()

	// Borrow
	body := `{"UserID":1}`
	req := httptest.NewRequest("POST", "/books/1/borrow", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}

	// Verify borrowed
	req = httptest.NewRequest("GET", "/books/1", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	var b library.Book
	json.NewDecoder(w.Body).Decode(&b)
	if b.Borrowed == nil || b.Borrowed.Name != "Alice" {
		t.Fatal("expected book to be borrowed by Alice")
	}

	// Borrow again -> fail
	body = `{"UserID":2}`
	req = httptest.NewRequest("POST", "/books/1/borrow", strings.NewReader(body))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for double borrow, got %d", w.Code)
	}

	// Return
	req = httptest.NewRequest("POST", "/books/1/return", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	// Verify returned
	req = httptest.NewRequest("GET", "/books/1", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	json.NewDecoder(w.Body).Decode(&b)
	if b.Borrowed != nil {
		t.Fatal("expected book to be returned")
	}
}

func TestReviews(t *testing.T) {
	_, mux := setup()

	// No reviews initially
	req := httptest.NewRequest("GET", "/books/1/reviews", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	var reviews []library.Review
	json.NewDecoder(w.Body).Decode(&reviews)
	if len(reviews) != 0 {
		t.Fatalf("expected 0 reviews, got %d", len(reviews))
	}

	// Create review
	body := `{"BookID":1,"UserID":1,"Text":"Great book!","Rating":5}`
	req = httptest.NewRequest("POST", "/reviews", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// Check reviews
	req = httptest.NewRequest("GET", "/books/1/reviews", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	json.NewDecoder(w.Body).Decode(&reviews)
	if len(reviews) != 1 || reviews[0].Text != "Great book!" {
		t.Fatalf("unexpected reviews: %+v", reviews)
	}
}
