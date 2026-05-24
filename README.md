<p align="center">
  <img src="./gitclaw-logo.png" alt="GitClaw Logo" width="200" />
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22-brightgreen?style=flat-square" alt="go version" />
  <img src="https://img.shields.io/github/license/open-gitagent/gitclaw?style=flat-square" alt="license" />
  <img src="https://img.shields.io/badge/Runtime-Go-blue?style=flat-square&logo=go&logoColor=white" alt="go runtime" />
</p>

<h1 align="center">Gitclaw</h1>

<p align="center">
  <strong>A universal git-native multimodal always learning AI Agent (TinyHuman)</strong><br/>
  Your agent lives inside a git repo — identity, rules, memory, tools, and skills are all version-controlled files.
</p>

<p align="center">
  <a href="#one-command-install">Install</a> &bull;
  <a href="#quick-start">Quick Start</a> &bull;
  <a href="#sdk">SDK</a> &bull;
  <a href="#architecture">Architecture</a> &bull;
  <a href="#tools">Tools</a> &bull;
  <a href="#hooks">Hooks</a> &bull;
  <a href="#skills">Skills</a> &bull;
  <a href="#plugins">Plugins</a>
</p>

---

## Why Gitclaw?

Most agent frameworks treat configuration as code scattered across your application. Gitclaw flips this — **your agent IS a git repository**:

- **`agent.yaml`** — model, tools, runtime config
- **`SOUL.md`** — personality and identity
- **`RULES.md`** — behavioral constraints
- **`memory/`** — git-committed memory with full history
- **`tools/`** — declarative YAML tool definitions
- **`skills/`** — composable skill modules
- **`hooks/`** — lifecycle hooks (script or programmatic)

Fork an agent. Branch a personality. `git log` your agent's memory. Diff its rules. This is **agents as repos**.

## One-Command Install

Copy, paste, run. That's it — no cloning, no manual setup. The installer handles everything:

```bash
bash <(curl -fsSL "https://raw.githubusercontent.com/open-gitagent/gitagent/main/install.sh?$(date +%s)")
```

This will:
- Build the Go CLI binary and place it in `~/.local/bin/gitclaw`

> **Requirements:** Go 1.22+, git

### Or install manually:

```bash
go build -o ~/.local/bin/gitclaw ./cmd/gitclaw
```

## Quick Start

**Run your first agent in one line:**

```bash
export OPENAI_API_KEY="sk-..."
gitclaw run --dir ~/my-project --prompt "Explain this project and suggest improvements"
```

That's it. Gitclaw auto-scaffolds everything on first run — `agent.yaml`, `SOUL.md`, `memory/` — and drops you into the agent.

### Guard Policy

Each agent can ship its own allowlist policy at `agents/<name>/guard.json`. The Go runtime enforces this policy before any tool call.

### CLI Options

| Command | Description |
|---|---|
| `gitclaw run --dir <path> --prompt <text> [--model provider:model]` | Run an agent prompt |
| `gitclaw diff --dir <path> [--json]` | Semantic diff of agent changes |
| `gitclaw bench --file bench.yaml --a <dir> [--b <dir>] [--json]` | Benchmark agent behavior |

### SDK (Go)

```bash
go get github.com/open-gitagent/gitagent/sdk
```

```go
package main

import (
  "fmt"
  "github.com/open-gitagent/gitagent/sdk"
)

func main() {
  out, err := sdk.Run(sdk.RunOptions{Dir: ".", Prompt: "Summarize this repo", MaxTurns: 10})
  if err != nil {
    panic(err)
  }
  fmt.Println(out)
}
```

## SDK

The Go SDK exposes `Run`, `SemanticDiff`, and `RunBench` for embedding the runtime.

```go
package main

import (
  "fmt"
  "github.com/open-gitagent/gitagent/sdk"
)

func main() {
  result, err := sdk.SemanticDiff(".")
  if err != nil {
    panic(err)
  }
  fmt.Println(result.Human())
}
```
| `maxTurns` | `number` | Max agent turns |
| `abortController` | `AbortController` | Cancellation signal |
| `constraints` | `object` | `temperature`, `maxTokens`, `topP`, `topK` |

### Message Types

| Type | Description | Key Fields |
|---|---|---|
| `delta` | Streaming text/thinking chunk | `deltaType`, `content` |
| `assistant` | Complete LLM response | `content`, `model`, `usage`, `stopReason` |
| `tool_use` | Tool invocation | `toolName`, `args`, `toolCallId` |
| `tool_result` | Tool output | `content`, `isError`, `toolCallId` |
| `system` | Lifecycle events | `subtype`, `content`, `metadata` |
| `user` | User message (multi-turn) | `content` |

## Architecture

```
my-agent/
├── agent.yaml          # Model, tools, runtime config
├── SOUL.md             # Agent identity & personality
├── RULES.md            # Behavioral rules & constraints
├── DUTIES.md           # Role-specific responsibilities
├── memory/
│   └── MEMORY.md       # Git-committed agent memory
├── tools/
│   └── *.yaml          # Declarative tool definitions
├── skills/
│   └── <name>/
│       ├── SKILL.md    # Skill instructions (YAML frontmatter)
│       └── scripts/    # Skill scripts
├── workflows/
│   └── *.yaml|*.md     # Multi-step workflow definitions
├── agents/
│   └── <name>/         # Sub-agent definitions
├── plugins/
│   └── <name>/         # Local plugins (plugin.yaml + tools/hooks/skills)
├── hooks/
│   └── hooks.yaml      # Lifecycle hook scripts
├── knowledge/
│   └── index.yaml      # Knowledge base entries
├── config/
│   ├── default.yaml    # Default environment config
│   └── <env>.yaml      # Environment overrides
├── examples/
│   └── *.md            # Few-shot examples
└── compliance/
    └── *.yaml          # Compliance & audit config
```

### Agent Manifest (`agent.yaml`)

```yaml
spec_version: "0.1.0"
name: my-agent
version: 1.0.0
description: An agent that does things

model:
  preferred: "anthropic:claude-sonnet-4-5-20250929"
  fallback: ["openai:gpt-4o"]
  constraints:
    temperature: 0.7
    max_tokens: 4096

tools: [cli, read, write, memory]

runtime:
  max_turns: 50
  timeout: 120

# Optional
extends: "https://github.com/org/base-agent.git"
skills: [code-review, deploy]
delegation:
  mode: auto
compliance:
  risk_level: medium
  human_in_the_loop: true
```

## Tools

### Built-in Tools

| Tool | Description |
|---|---|
| `cli` | Execute shell commands |
| `read` | Read files with pagination |
| `write` | Write/create files |
| `memory` | Load/save git-committed memory |

### Declarative Tools

Define tools as YAML in `tools/`:

```yaml
# tools/search.yaml
name: search
description: Search the codebase
input_schema:
  properties:
    query:
      type: string
      description: Search query
    path:
      type: string
      description: Directory to search
  required: [query]
implementation:
  script: search.sh
  runtime: sh
```

The script receives args as JSON on stdin and returns output on stdout.

## Hooks

Script-based hooks in `hooks/hooks.yaml`:

```yaml
hooks:
  on_session_start:
    - script: validate-env.sh
      description: Check environment is ready
  pre_tool_use:
    - script: audit-tools.sh
      description: Log and gate tool usage
  post_response:
    - script: notify.sh
  on_error:
    - script: alert.sh
```

Hook scripts receive context as JSON on stdin and return:

```json
{ "action": "allow" }
{ "action": "block", "reason": "Not permitted" }
{ "action": "modify", "args": { "modified": "args" } }
```

## Skills

Skills are composable instruction modules in `skills/<name>/`:

```
skills/
  code-review/
    SKILL.md
    scripts/
      lint.sh
```

```markdown
---
name: code-review
description: Review code for quality and security
---

# Code Review

When reviewing code:
1. Check for security vulnerabilities
2. Verify error handling
3. Run the lint script for style checks
```

Invoke via CLI: `/skill:code-review Review the auth module`

## Plugins

Plugins are reusable extensions that can provide tools, hooks, skills, prompts, and memory layers. They follow the same git-native philosophy — a plugin is a directory with a `plugin.yaml` manifest.

### CLI Commands

```bash
# Install from git URL
gitclaw plugin install https://github.com/org/my-plugin.git

# Install from local path
gitclaw plugin install ./path/to/plugin

# Install with options
gitclaw plugin install <source> --name custom-name --force --no-enable

# List all discovered plugins
gitclaw plugin list

# Enable / disable
gitclaw plugin enable my-plugin
gitclaw plugin disable my-plugin

# Remove
gitclaw plugin remove my-plugin

# Scaffold a new plugin
gitclaw plugin init my-plugin
```

| Flag | Description |
|---|---|
| `--name <name>` | Custom plugin name (default: derived from source) |
| `--force` | Reinstall even if already present |
| `--no-enable` | Install without auto-enabling |

### Plugin Manifest (`plugin.yaml`)

```yaml
id: my-plugin                    # Required, kebab-case
name: My Plugin
version: 0.1.0
description: What this plugin does
author: Your Name
license: MIT
engine: ">=0.3.0"               # Min gitclaw version

provides:
  tools: true                    # Load tools from tools/*.yaml
  skills: true                   # Load skills from skills/
  prompt: prompt.md              # Inject into system prompt
  hooks:
    pre_tool_use:
      - script: hooks/audit.sh
        description: Audit tool calls

config:
  properties:
    api_key:
      type: string
      description: API key
      env: MY_API_KEY            # Env var fallback
    timeout:
      type: number
      default: 30
  required: [api_key]

entry: index.ts                  # Optional programmatic entry point
```

### Plugin Config in `agent.yaml`

```yaml
plugins:
  my-plugin:
    enabled: true
    source: https://github.com/org/my-plugin.git  # Auto-install on load
    version: main                                   # Git branch/tag
    config:
      api_key: "${MY_API_KEY}"                      # Supports env interpolation
      timeout: 60
```

Config resolution priority: `agent.yaml config` > `env var` > `manifest default`.

### Discovery Order

Plugins are discovered in this order (first match wins):

1. **Local** — `<agent-dir>/plugins/<name>/`
2. **Global** — `~/.gitclaw/plugins/<name>/`
3. **Installed** — `<agent-dir>/.gitagent/plugins/<name>/`

### Programmatic Plugins

Plugins with an `entry` field in their manifest get a full API:

```typescript
// index.ts
import type { GitclawPluginApi } from "gitclaw";

export async function register(api: GitclawPluginApi) {
  // Register a tool
  api.registerTool({
    name: "search_docs",
    description: "Search documentation",
    inputSchema: {
      properties: { query: { type: "string" } },
      required: ["query"],
    },
    handler: async (args) => {
      const results = await search(args.query);
      return { text: JSON.stringify(results) };
    },
  });

  // Register a lifecycle hook
  api.registerHook("pre_tool_use", async (ctx) => {
    api.logger.info(`Tool called: ${ctx.tool}`);
    return { action: "allow" };
  });

  // Add to system prompt
  api.addPrompt("Always check docs before answering questions.");

  // Register a memory layer
  api.registerMemoryLayer({
    name: "docs-cache",
    path: "memory/docs-cache.md",
    description: "Cached documentation lookups",
  });
}
```

**Available API methods:**

| Method | Description |
|---|---|
| `registerTool(def)` | Register a tool the agent can call |
| `registerHook(event, handler)` | Register a lifecycle hook (`on_session_start`, `pre_tool_use`, `post_response`, `on_error`) |
| `addPrompt(text)` | Append text to the system prompt |
| `registerMemoryLayer(layer)` | Register a memory layer |
| `logger.info/warn/error(msg)` | Prefixed logging (`[plugin:id]`) |
| `pluginId` | Plugin identifier |
| `pluginDir` | Absolute path to plugin directory |
| `config` | Resolved config values |

### Plugin Structure

```
my-plugin/
├── plugin.yaml          # Manifest (required)
├── tools/               # Declarative tool definitions
│   └── *.yaml
├── hooks/               # Hook scripts
├── skills/              # Skill modules
├── prompt.md            # System prompt addition
└── index.ts             # Programmatic entry point
```

## Multi-Model Support

Gitclaw works with any LLM provider supported by [pi-ai](https://github.com/badlogic/pi-mono/tree/main/packages/ai):

```yaml
# agent.yaml
model:
  preferred: "anthropic:claude-sonnet-4-5-20250929"
  fallback:
    - "openai:gpt-4o"
    - "google:gemini-2.0-flash"
```

Supported providers: `anthropic`, `openai`, `google`, `xai`, `groq`, `mistral`, and more.

## Inheritance & Composition

Agents can extend base agents:

```yaml
# agent.yaml
extends: "https://github.com/org/base-agent.git"

# Dependencies
dependencies:
  - name: shared-tools
    source: "https://github.com/org/shared-tools.git"
    version: main
    mount: tools

# Sub-agents
delegation:
  mode: auto
```

## Compliance & Audit

Built-in compliance validation and audit logging:

```yaml
# agent.yaml
compliance:
  risk_level: high
  human_in_the_loop: true
  data_classification: confidential
  regulatory_frameworks: [SOC2, GDPR]
  recordkeeping:
    audit_logging: true
    retention_days: 90
```

Audit logs are written to `.gitagent/audit.jsonl` with full tool invocation traces.

## Telemetry

Gitclaw ships with built-in OpenTelemetry instrumentation. Set `OTEL_EXPORTER_OTLP_ENDPOINT` and telemetry is on; leave it unset and runtime cost is zero.

Three layers of signals:

1. **HTTP-level** — `@opentelemetry/instrumentation-undici` auto-patches `fetch`/`undici`, so every LLM provider call (Anthropic, OpenAI, Google, …) gets a client span with URL, status code, and timing.
2. **`gen_ai.chat` spans** — emitted on every assistant `message_end`. Carry `gen_ai.system`, `gen_ai.request.model`, `gen_ai.usage.input_tokens`, `gen_ai.usage.output_tokens`, `gen_ai.response.finish_reasons`, and `gitclaw.cost_usd`. Span/metric content never contains the prompt or completion text.
3. **`gitclaw.tool.execute` spans** — wrap every tool call with `tool.name`, `tool.call_id`, `tool.status` (`ok`/`error`), and `tool.error_message` on failure.

A root `gitclaw.agent.session` span opens at agent construction and closes on every exit path (success, hook-block, SIGINT, error).

### CLI usage

Just set the endpoint — no `--import` flag, no extra install steps:

```bash
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318 gitclaw -p "your prompt"
```

Telemetry is enabled automatically when the endpoint is set and disabled when it is not. To force-disable even when the endpoint is set, pass `GITCLAW_OTEL_ENABLED=false`.

### Environment variables

| Variable | Description | Default |
|----------|-------------|---------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP/HTTP collector base URL (e.g. `http://localhost:4318`). When set, telemetry is auto-enabled. | (unset → telemetry off) |
| `GITCLAW_OTEL_ENABLED` | Set to `false` to disable telemetry even when the endpoint is set | (unset = auto) |
| `OTEL_SERVICE_NAME` | Resource `service.name` | `gitclaw` |
| `OTEL_SERVICE_VERSION` | Resource `service.version` | (unset) |
| `OTEL_EXPORTER_OTLP_HEADERS` | Comma-separated key=value pairs, no quotes (e.g. `Authorization=Bearer xyz,x-tenant=abc`) | (unset) |
| `OTEL_TRACES_EXPORTER` | Set to `console` to print spans to stdout — no collector needed | (unset) |

### SDK usage

For programmatic embedders, call `initTelemetry` explicitly — you control when initialisation happens:

```ts
import { initTelemetry, shutdownTelemetry, query } from "gitclaw";

await initTelemetry({ serviceName: "my-app" });

for await (const msg of query({ prompt: "hello", model: "anthropic:claude-4-6-sonnet-latest" })) {
  // …
}

await shutdownTelemetry();
```

`OTEL_EXPORTER_OTLP_ENDPOINT` and `OTEL_EXPORTER_OTLP_HEADERS` are read automatically by the OTLP exporter when not supplied programmatically. Pass `exporterEndpoint` / `headers` only when you need to override env-based config in code.

### Emitted spans

| Name | Kind | Key attributes |
|------|------|----------------|
| `gitclaw.agent.session` | INTERNAL | `gitclaw.entry` (`sdk` / `cli`), `gitclaw.cost_usd`, `gitclaw.session.duration_ms` |
| `gitclaw.tool.execute` | INTERNAL | `tool.name`, `tool.call_id`, `tool.status`, `tool.error_message` |
| `gen_ai.chat` | CLIENT | `gen_ai.system`, `gen_ai.request.model`, `gen_ai.usage.input_tokens`, `gen_ai.usage.output_tokens`, `gen_ai.response.finish_reasons`, `gitclaw.cost_usd` |
| `HTTP …` | CLIENT | URL, status code, duration (auto from `instrumentation-undici`) |

### Emitted metrics

| Name | Type | Description |
|------|------|-------------|
| `gitclaw.tool.calls` | counter | Number of tool executions, labelled by `tool.name` |
| `gitclaw.tool.duration_ms` | histogram | Tool execution duration |
| `gitclaw.session.duration_ms` | histogram | Session duration |
| `gitclaw.session.cost_usd` | counter (USD) | Cumulative session cost |
| `gen_ai.client.token.usage` | counter | Token usage by `gen_ai.system`, `gen_ai.request.model`, `gen_ai.token.type` |
| `gen_ai.client.operation.duration` | histogram | LLM call duration |

### Console quickstart (no collector)

Print spans directly to stdout — useful for local debugging:

```bash
OTEL_TRACES_EXPORTER=console gitclaw -p "test"
```

### Local Jaeger quickstart

```bash
docker run --rm -p 16686:16686 -p 4318:4318 jaegertracing/all-in-one:latest

OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318 gitclaw -p "test"

# Open http://localhost:16686 → service "gitclaw"
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## ❓ FAQ

### General

**What is Gitclaw?**
Gitclaw (formerly GitAgent) is a git-native AI agent framework where the agent IS a git repository. Identity, rules, memory, tools, and skills are all version-controlled files, enabling "agents as repos" paradigm.

**How does Gitclaw differ from other agent frameworks?**
Unlike frameworks that scatter configuration across application code, Gitclaw makes the agent itself a git repo:
- Fork an agent → inherit personality, rules, tools
- Branch → create alternate personality versions
- `git log` → see agent's memory evolution
- Diff → track rule changes over time

**What is the "agents as repos" concept?**
Your agent lives in a git repository with structured files:
- `agent.yaml` — model, tools, runtime config
- `SOUL.md` — personality and identity
- `RULES.md` — behavioral constraints
- `memory/` — git-committed memory with full history
- `tools/` — declarative YAML tool definitions
- `skills/` — composable skill modules
- `hooks/` — lifecycle hooks

### The Go Migration (Recent Changes)

**We recently completed a major architectural overhaul, migrating Gitclaw from TypeScript to Go.** Here is a comparison of what changed:

| Feature/Metric | Before (TypeScript) | After (Go) | Why we changed |
|---|---|---|---|
| **Runtime Engine** | Node.js (V8) | Compiled Go Binary | Go provides a single, dependency-free binary with sub-50ms cold starts, making CLI executions near-instant. |
| **Concurrency & State** | Ad-hoc async/await handling | MVCC Write Ledger | We introduced Multi-Version Concurrency Control (MVCC) to ensure conflict-free file management when multiple agents or users write simultaneously. |
| **Security & Safety** | Basic script wrappers | Stateless Circuit-Breaker Pipeline | Go's robust standard library allowed us to build a pipeline to intercept, validate, and block unauthorized or runaway agent tasks before they execute. |
| **Testing & Stability** | ~120 mixed Jest tests | 42 Comprehensive Unit & Integration Tests | Streamlined, strict Go tests covering the new ledger and guard policies ensure maximum stability. |

### Installation & Setup

**What are the requirements?**
Go 1.22+ and git.

**How do I install Gitclaw?**
Run the installer for guided setup:
```bash
bash <(curl -fsSL "https://raw.githubusercontent.com/open-gitagent/gitagent/main/install.sh")
```
Or set manually:
```bash
export OPENAI_API_KEY="sk-..."
```

**Which LLM providers are supported?**
- OpenAI (GPT-4o, GPT-4o-mini, etc.)
- Anthropic (Claude models)
- Google (Gemini)
- Any OpenAI-compatible provider

Use `--model` flag to override: `gitclaw --model anthropic:claude-sonnet-4-5-20250929`

### Core Concepts

**What is the SDK and how do I use it?**
The SDK provides programmatic access via `sdk.Run()` that streams agent events:
```go
package main

import (
  "fmt"
  "github.com/open-gitagent/gitagent/sdk"
)

func main() {
  out, err := sdk.Run(sdk.RunOptions{Dir: ".", Prompt: "hello", MaxTurns: 10})
  if err != nil {
    panic(err)
  }
  fmt.Println(out)
}
```

**How do local repo mode sessions work?**
Clone a GitHub repo, run an agent on it, auto-commit to a session branch:
```bash
gitclaw run --dir . --prompt "Fix the bug"
```

**What hooks are available?**
Hooks are lifecycle scripts or programmatic handlers in the `hooks/` directory. They trigger on agent events like tool execution, session start/end, or memory updates.

### Development

**How do I create custom tools?**
Define tools in `tools/` directory using declarative YAML format. Each tool specifies name, description, parameters, and execution logic.

**How do I add skills?**
Create skill modules in `skills/` directory. Skills are composable and can be imported from installed packages or defined locally.

**What telemetry options are available?**
OpenTelemetry integration for observability:
- Set `OTEL_EXPORTER_OTLP_ENDPOINT` for auto-enable
- Use `OTEL_TRACES_EXPORTER=console` for local debugging
- Jaeger quickstart with Docker

### Troubleshooting

**Why is my agent not responding?**
- Check API key is set (`OPENAI_API_KEY` or equivalent)
- Verify network connectivity to LLM provider
- Use `--verbose` flag for detailed logs
- Check `agent.yaml` model configuration

**How do I debug agent behavior?**
- Use console exporter: `OTEL_TRACES_EXPORTER=console gitclaw run -p "test"`
- Check spans in Jaeger: `docker run -p 16686:16686 -p 4318:4318 jaegertracing/all-in-one`
- Inspect `memory/` directory for agent state

**Where can I get help?**
- GitHub Issues: https://github.com/open-gitagent/gitagent/issues
- Examples: See README SDK section and CLI options
- Contributing: See CONTRIBUTING.md for guidelines

## License

This project is licensed under the [MIT License](./LICENSE).
