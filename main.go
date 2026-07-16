package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Wordle MCP",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	tool := mcp.NewTool("get_wordle_solution",
		mcp.WithDescription("Fetches the Wordle of a particular date provided between 2021-05-19 to 23 days future"),
		mcp.WithString("date",
			mcp.Required(),
			mcp.Description("The date to be passed on to the Wordle API (YYYY-MM-DD format)")),
	)

	s.AddTool(tool, getWordleData)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// wordleAPIBaseURL is a var (not a const) so tests can point it at an
// httptest.Server instead of the real NYT API.
var wordleAPIBaseURL = "https://www.nytimes.com/svc/wordle/v2"

// WordleAPIData is the shape of a single day's response from the Wordle API.
type WordleAPIData struct {
	ID              int    `json:"id"`
	Solution        string `json:"solution"`
	PrintDate       string `json:"print_date"`
	DaysSinceLaunch int    `json:"days_since_launch"`
	Editor          string `json:"editor"`
}

// fetchWordleData retrieves and decodes the Wordle solution for the given date.
func fetchWordleData(ctx context.Context, date string) (*WordleAPIData, error) {
	url := fmt.Sprintf("%s/%s.json", wordleAPIBaseURL, date)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wordle api returned status %d: %s", res.StatusCode, body)
	}

	var data WordleAPIData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unable to parse response: %w", err)
	}

	return &data, nil
}

func getWordleData(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	date, err := request.RequireString("date")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, err := fetchWordleData(ctx, date)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to fetch wordle solution", err), nil
	}

	result, err := mcp.NewToolResultJSON(data)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to encode wordle solution", err), nil
	}

	return result, nil
}
