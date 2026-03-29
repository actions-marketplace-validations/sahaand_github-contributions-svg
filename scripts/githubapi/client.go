package githubapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// apiReq performs an authenticated HTTP request to the GitHub API.
func apiReq(token, method, url string, body []byte) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "readme-asset-generator/1.0")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// RestGet fetches a single REST API endpoint and returns the parsed JSON object.
func RestGet(token, path string) (map[string]any, error) {
	data, err := apiReq(token, "GET", "https://api.github.com"+path, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	return result, json.Unmarshal(data, &result)
}

// RestAll fetches all pages of a paginated REST API endpoint (100 items/page).
func RestAll(token, path string) ([]map[string]any, error) {
	var results []map[string]any
	page := 1
	for {
		sep := "&"
		if !strings.Contains(path, "?") {
			sep = "?"
		}
		url := fmt.Sprintf("https://api.github.com%s%sper_page=100&page=%d", path, sep, page)
		data, err := apiReq(token, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		var chunk []map[string]any
		if err := json.Unmarshal(data, &chunk); err != nil {
			return nil, fmt.Errorf("restAll %s page %d: %w", path, page, err)
		}
		results = append(results, chunk...)
		if len(chunk) < 100 {
			break
		}
		page++
	}
	return results, nil
}

// GqlQuery executes a GraphQL query against the GitHub GraphQL API.
func GqlQuery(token, query string, variables map[string]any) (map[string]any, error) {
	body, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, err
	}
	data, err := apiReq(token, "POST", "https://api.github.com/graphql", body)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	return result, json.Unmarshal(data, &result)
}

// AsInt safely extracts an int from a JSON-decoded map value (stored as float64 by Go).
func AsInt(m map[string]any, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}
