package gist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://api.github.com"

// Client talks to the GitHub Gist API.
type Client struct {
	Token   string
	BaseURL string // override for testing; empty uses default
}

func (c *Client) baseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return defaultBaseURL
}

func (c *Client) do(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return http.DefaultClient.Do(req)
}

type gistFile struct {
	Content string `json:"content"`
}

type gistRequest struct {
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]gistFile `json:"files"`
}

type gistResponse struct {
	ID    string                     `json:"id"`
	Files map[string]gistResponseFile `json:"files"`
}

type gistResponseFile struct {
	Content string `json:"content"`
}

// Create creates a new secret gist containing bookmarks.json.
// Returns the gist ID.
func (c *Client) Create(content []byte) (string, error) {
	payload := gistRequest{
		Description: "TuiBookie bookmarks",
		Public:      false,
		Files: map[string]gistFile{
			"bookmarks.json": {Content: string(content)},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := c.do("POST", c.baseURL()+"/gists", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if err := checkStatus(resp, http.StatusCreated); err != nil {
		return "", err
	}

	var result gistResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("invalid response: %w", err)
	}
	return result.ID, nil
}

// Update replaces the bookmarks.json content in an existing gist.
func (c *Client) Update(gistID string, content []byte) error {
	payload := gistRequest{
		Files: map[string]gistFile{
			"bookmarks.json": {Content: string(content)},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.do("PATCH", c.baseURL()+"/gists/"+gistID, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	return checkStatus(resp, http.StatusOK)
}

// Fetch downloads the bookmarks.json content from a gist.
func (c *Client) Fetch(gistID string) ([]byte, error) {
	resp, err := c.do("GET", c.baseURL()+"/gists/"+gistID, nil)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if err := checkStatus(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var result gistResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("invalid response: %w", err)
	}

	f, ok := result.Files["bookmarks.json"]
	if !ok {
		return nil, fmt.Errorf("gist does not contain bookmarks.json")
	}
	return []byte(f.Content), nil
}

func checkStatus(resp *http.Response, expected int) error {
	if resp.StatusCode == expected {
		return nil
	}
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("auth failed — check your token")
	case http.StatusNotFound:
		return fmt.Errorf("gist not found — it may have been deleted")
	default:
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
}
