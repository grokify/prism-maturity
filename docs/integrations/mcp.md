# MCP Integration

The `mcp` package provides a Model Context Protocol (MCP) server for exposing maturity management tools to AI assistants.

## Overview

MCP enables AI assistants like Claude to interact with PRISM maturity data through a standardized protocol. The server exposes tools for:

- Goal category management
- Goal CRUD operations
- Maturity model operations

## Installation

The MCP server is included with prism-maturity. No additional installation required.

## Usage

### Starting the Server

```bash
prism mcp serve
```

The server runs on stdio transport, suitable for integration with AI assistants.

### Programmatic Usage

```go
import (
    "context"
    "github.com/grokify/prism-maturity/mcp"
)

server := mcp.NewServer().WithDataDir("./data")
if err := server.Run(context.Background()); err != nil {
    log.Fatal(err)
}
```

## Available Tools

### Category Tools

| Tool | Description |
|------|-------------|
| `category_list` | List all valid goal categories |
| `category_validate` | Validate a category string |
| `category_revenue_target` | Create a revenue target |
| `category_adoption_target` | Create an adoption target |

### Goal Tools

| Tool | Description |
|------|-------------|
| `goal_list` | List goals from a PRISM document |
| `goal_get` | Get a specific goal by ID |
| `goal_create` | Create a new goal |
| `goal_update` | Update an existing goal |

### Maturity Tools

| Tool | Description |
|------|-------------|
| `maturity_score` | Calculate maturity score |
| `maturity_gaps` | Identify maturity gaps |
| `maturity_roadmap` | Generate improvement roadmap |

## Configuration

Configure the MCP server in your AI assistant's settings. For Claude Desktop:

```json
{
  "mcpServers": {
    "prism-maturity": {
      "command": "prism",
      "args": ["mcp", "serve"],
      "env": {
        "PRISM_DATA_DIR": "/path/to/data"
      }
    }
  }
}
```

## Example Interactions

### List Categories

```
Tool: category_list
Result:
- operations: Operational excellence - SLI/SLO targets, reliability, performance
- revenue: Revenue growth - ARR/MRR targets, financial metrics
- adoption: User adoption - MAU/DAU, retention, NPS, engagement
...
```

### Create Revenue Target

```
Tool: category_revenue_target
Input: {
  "target_arr": 10000000000,
  "current_arr": 8000000000,
  "currency": "USD"
}
Result: {
  "progress_percent": 80.0,
  "gap_arr_dollars": 20000000,
  "is_met": false
}
```

## See Also

- [Goal Categories](../schema/goal-categories.md) - Category types and targets
- [Model Context Protocol](https://modelcontextprotocol.io/) - MCP specification
