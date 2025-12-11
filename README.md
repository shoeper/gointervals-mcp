
# Build

## go

```
go build -o server cmd/server/main.go
```

## docker

```
docker build -t mcp .
```

# Run

Configure the following environment variables in .env or via the environment
```
INTERVALS_API_BASE_URL=https://intervals.icu
# get from developer settings at https://intervals.icu/settings
INTERVALS_ATHLETE_ID=
INTERVALS_API_KEY=
MCP_AUTH_TOKEN=
PORT=8000
```

```
./server
```

# Project structure

```
cmd/server/main.go: entry point to run intervals mcp server
internal/config: config
internal/intervalsclient: generated client for intervals
internal/tools: mcp tools offered by this server, one tool per file
internal/tools/tools.go: registration for all tools
```

# Generate intervals api client

Intervals client is generated from the OpenAPI specification.

Get latest spec
```
curl -ointernal/intervalsclient/openapi-spec.json https://intervals.icu/api/v1/docs
```

Currently, the following is exempt from generation as it generates invalid code.
```
DeleteEventsResponse
```

And the spec has to be edited, changing the type of skyline_chart_bytes to string
```
#"skyline_chart_bytes": {
#                        "type": "string"
#                    },
```

Run go generate to update the generated client
```
cd internal/intervalsclient
go generate
```

oapi-codegen
```
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

# Debugging

## Debugging with curl

### Get mcp-session-id

```
# docker exec -it intervals-mcp 
curl -i -s http://localhost:8000/mcp \
    -H "Authorization: Bearer <MCP_AUTH_TOKEN>" \
    -H "Accept: text/event-stream,application/json" \
    -H "Content-Type: application/json" \
    -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "id": 1,
        "params": {
            "protocolVersion": "2024-11-05",
            "capabilities": {
            "roots": { "listChanged": true },
            "sampling": {}
            },
            "clientInfo": {
                "name": "curl-client",
                "version": "1.0.0"
            }
        }
    }' | grep -i 'Mcp-Session-Id:' | awk '{print $2}'
```

### List tools
```
# docker exec -it intervals-mcp 
curl \
    -H "Authorization: Bearer <MCP_AUTH_TOKEN>" \
    -H "Accept: text/event-stream,application/json" \
    -H "Mcp-Session-Id: mcp-session-d7b5663f-c68d-46d3-88e3-eeef2700f489" \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc": "2.0", "method": "tools/list", "id": 1}' -X POST http://127.0.0.1:8000/mcp
```

### Invoke tool
```
# docker exec -it intervals-mcp 
curl \
    -H "Authorization: Bearer <MCP_AUTH_TOKEN>" \
    -H "Accept: text/event-stream,application/json" \
    -H "Mcp-Session-Id: mcp-session-d7b5663f-c68d-46d3-88e3-eeef2700f489" \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "tools/call",
        "params": {
            "name": "get_activities"
        },
        "id": 1
    }' http://127.0.0.1:8000/mcp
```