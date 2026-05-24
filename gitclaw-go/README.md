# GitClaw Go Runtime

A high-performance Go rewrite of the GitClaw agent runtime with built-in safety guardrails.

## Quick Start

```bash
# Build
go build -o gitclaw ./cmd/gitclaw/

# Run tests (no API key needed)
go test ./... -v

# Run with a model
export OPENAI_API_KEY="sk-..."
./gitclaw --model "openai:gpt-4o-mini"
```

## What's New vs TypeScript

| Feature | TypeScript | Go |
|---------|-----------|-----|
| Binary | 150MB (node + deps) | **9.2 MB** |
| Startup | ~1.5s | **<50ms** |
| Concurrency | Single-threaded | **Goroutines** |
| Write safety | None | **MVCC ledger** |
| Tool safety | Shell hooks only | **Circuit breaker pipeline** |
| Dependencies | 15 npm packages | **1** |

See [COMPARISON.md](COMPARISON.md) for the full breakdown.

## Project Structure

```
gitclaw-go/
├── cmd/gitclaw/          # CLI entrypoint + REPL
├── internal/
│   ├── agent/            # Agent loop, LLM client, event system
│   ├── config/           # agent.yaml parser
│   ├── guard/            # Circuit breaker middleware (rate, policy, breaker, cost)
│   ├── hooks/            # Shell script hook runner
│   ├── state/            # MVCC write ledger + git commit serializer
│   └── tools/            # Built-in tools (cli, read, write, memory)
└── pkg/sdk/              # Public programmatic API
```

## Configuration

Uses the same `agent.yaml` format as the TypeScript version:

```yaml
name: my-agent
model:
  preferred: "openai:gpt-4o-mini"
tools: [cli, read, write, memory]
runtime:
  max_turns: 50
  guard:
    rate_limit:
      max_per_window: 100
    circuit_breaker:
      failure_threshold: 5
    cost_ceiling:
      max_usd: 10.0
    policies:
      - tool: cli
        deny: "args.command matches 'rm -rf.*'"
      - tool: "*"
        deny: "args contains 'sudo'"
```

## Tests

```bash
go test ./... -v -count=1    # 42 tests, all packages
go test ./internal/guard/ -v # 20 circuit breaker tests
go test ./internal/state/ -v # 7 write ledger tests
go test ./internal/tools/ -v # 10 tool tests
```

## License

Same as the parent GitClaw project.
