package main

import (
	"context"
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
			mcp.Description("The date to be passed on to the Wordle API")),
	)

	s.AddTool(tool, getWordleData)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func getWordleData(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	date, err := request.RequireString("date")

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var url string = fmt.Sprintf("https://www.nytimes.com/svc/wordle/v2/%s.json", date)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("There has been an issue with your GET request")
		fmt.Println(err)
		return mcp.NewToolResultErrorFromErr("unable to execute request", err), nil
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return mcp.NewToolResultErrorFromErr("unable to execute request", err), nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	// var data WordleAPIData
	// err = json.Unmarshal(body, &data)
	// if err != nil {
	// 	fmt.Println("Error parsing JSON response")
	// 	return mcp.NewToolResultError( err.Error() ), nil
	// }

	return mcp.NewToolResultText(string(body)), nil
}
