package gist

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatal("missing auth header")
		}

		var req gistRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.Public {
			t.Fatal("expected secret gist")
		}
		if req.Files["bookmarks.json"].Content != `{"test":[]}` {
			t.Fatalf("unexpected content: %s", req.Files["bookmarks.json"].Content)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(gistResponse{ID: "abc123"})
	}))
	defer srv.Close()

	c := &Client{Token: "test-token", BaseURL: srv.URL}
	id, err := c.Create([]byte(`{"test":[]}`))
	if err != nil {
		t.Fatal(err)
	}
	if id != "abc123" {
		t.Fatalf("expected abc123, got %s", id)
	}
}

func TestUpdate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Fatalf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/gists/abc123" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req gistRequest
		json.Unmarshal(body, &req)
		if req.Files["bookmarks.json"].Content != `{"updated":[]}` {
			t.Fatalf("unexpected content: %s", req.Files["bookmarks.json"].Content)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := &Client{Token: "test-token", BaseURL: srv.URL}
	if err := c.Update("abc123", []byte(`{"updated":[]}`)); err != nil {
		t.Fatal(err)
	}
}

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/gists/abc123" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(gistResponse{
			ID: "abc123",
			Files: map[string]gistResponseFile{
				"bookmarks.json": {Content: `{"servers":[]}`},
			},
		})
	}))
	defer srv.Close()

	c := &Client{Token: "test-token", BaseURL: srv.URL}
	data, err := c.Fetch("abc123")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"servers":[]}` {
		t.Fatalf("unexpected content: %s", data)
	}
}

func TestFetchMissingFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(gistResponse{
			ID:    "abc123",
			Files: map[string]gistResponseFile{},
		})
	}))
	defer srv.Close()

	c := &Client{Token: "test-token", BaseURL: srv.URL}
	_, err := c.Fetch("abc123")
	if err == nil {
		t.Fatal("expected error for missing bookmarks.json")
	}
}

func TestUnauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := &Client{Token: "bad-token", BaseURL: srv.URL}
	_, err := c.Create([]byte(`{}`))
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := &Client{Token: "test-token", BaseURL: srv.URL}
	_, err := c.Fetch("nonexistent")
	if err == nil {
		t.Fatal("expected not found error")
	}
}
