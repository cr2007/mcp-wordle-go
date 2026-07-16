package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func callRequest(args map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "get_wordle_solution",
			Arguments: args,
		},
	}
}

// withTestServer points wordleAPIBaseURL at a test server for the duration of the test.
func withTestServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	original := wordleAPIBaseURL
	wordleAPIBaseURL = server.URL
	t.Cleanup(func() { wordleAPIBaseURL = original })
}

func TestGetWordleData_Success(t *testing.T) {
	want := WordleAPIData{
		ID:              2191,
		Solution:        "whine",
		PrintDate:       "2023-01-01",
		DaysSinceLaunch: 561,
		Editor:          "Tracy Bennett",
	}

	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if wantPath := "/" + want.PrintDate + ".json"; r.URL.Path != wantPath {
			t.Errorf("unexpected request path: got %s, want %s", r.URL.Path, wantPath)
		}
		json.NewEncoder(w).Encode(want)
	})

	result, err := getWordleData(context.Background(), callRequest(map[string]any{"date": want.PrintDate}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got error result: %+v", result)
	}

	got, ok := result.StructuredContent.(*WordleAPIData)
	if !ok {
		t.Fatalf("structured content is not *WordleAPIData: %#v", result.StructuredContent)
	}
	if *got != want {
		t.Errorf("got %+v, want %+v", *got, want)
	}
}

func TestGetWordleData_MissingDate(t *testing.T) {
	result, err := getWordleData(context.Background(), callRequest(map[string]any{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected an error result when date argument is missing")
	}
}

func TestGetWordleData_UpstreamError(t *testing.T) {
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	result, err := getWordleData(context.Background(), callRequest(map[string]any{"date": "1999-01-01"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected an error result for a non-200 upstream response")
	}
}

func TestGetWordleData_MalformedResponse(t *testing.T) {
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	})

	result, err := getWordleData(context.Background(), callRequest(map[string]any{"date": "2023-01-01"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected an error result for a malformed upstream response")
	}
}
