<!-- omit from toc -->
# Wordle MCP - Go

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/cr2007/mcp-wordle-go)

> [!NOTE]
> To check out the Python version of this MCP Server, [click here](https://github.com/cr2007/mcp-wordle-python)

A MCP Server implemented in Go for fetching the Wordle solutions via the Wordle API.

> [!IMPORTANT]
> Wordle solutions are only available from 2021-05-19, to 23 days in the future.<br>
> Any other attempts at calling other dates will return an error from the API

<!-- omit from toc -->
# Index
- [Get Started](#get-started)
  - [Using Docker (recommended)](#using-docker-recommended)
  - [Local Setup](#local-setup)
    - [Setup](#setup)
- [Examples on Claude Desktop](#examples-on-claude-desktop)
- [Contributing](#contributing)

---

# Get Started

## Using Docker (recommended)

The quickest and easiest method to get started. Ensure that you have [Docker](https://www.docker.com) installed.<br>
Add this to your MCP Server configuration:

```json
{
  "mcpServers": {
    "Wordle-MCP-Go": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "--init",
        "-e",
        "DOCKER_CONTAINER=true",
        "ghcr.io/cr2007/mcp-wordle-go:latest"
      ]
    }
  }
}
```

> [!IMPORTANT]
> If you get an error on Claude Desktop for the first time, just make sure to pull the image before running.<br>
> `docker pull ghcr.io/cr2007/mcp-wordle-go:latest`

## Local Setup

For this setup, you need to make sure that you have the [Go programming language](https://go.dev/) installed on your machine.

### Setup

Before adding this to your MCP Server, you need to do the following:

```bash
# Clone the Git repository
git clone https://github.com/cr2007/mcp-wordle-go
cd mcp-wordle-go

# Install the dependencies
go mod tidy

# Build the executable
go build main.go
```

Add this to your MCP server configuration:

```json
{
  "mcpServers": {
    "Wordle-MCP-Go":{
        "command": "ABSOLUTE//PATH//TO//main.go",
      }
  }
}
```

# Examples on Claude Desktop

<div align="center">
    <img width=75%, src="./images/Claude_Chat-Example.png">
</div>

# Contributing

Contributions are welcome! You may [fork](https://github.com/cr2007/mcp-wordle-go/fork) the repo, create your changes in a branch, and then create a [Pull Request](https://github.com/cr2007/mcp-wordle-go/compare)
