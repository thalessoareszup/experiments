# User Story: MCP Integration for mini-swe-agent

## Problem

mini-swe-agent currently has a fixed set of tools (file read/write, bash, grep, etc.). Users want to extend the agent's capabilities without modifying the core codebase. For example, connecting to databases, fetching from APIs, or integrating with internal company tools.

## Who is it for

Developers and teams who use mini-swe-agent and want to:
- Add custom tools specific to their workflow
- Connect the agent to external services (databases, APIs, etc.)
- Use existing MCP servers from the ecosystem

## What it should do

1. **Configure MCP servers** - Users should be able to specify MCP servers in a config file. Each server entry includes the command to run it and any environment variables needed.

2. **Discover tools** - On startup, mini-swe-agent connects to configured MCP servers and discovers available tools (name, description, input schema).

3. **Use MCP tools** - The agent can call MCP tools just like built-in tools. When the agent decides to use an MCP tool, mini-swe-agent sends the request to the appropriate server and returns the result.

4. **Handle failures gracefully** - If an MCP server is unavailable or a tool call fails, the agent should get a clear error message and continue working.

## Constraints

- Use stdio transport only (no HTTP/SSE for now)
- Config file should be simple (YAML or JSON)
- MCP servers are started on-demand when mini-swe-agent launches
- Keep the implementation minimal - just tool discovery and execution, no resources or prompts yet
