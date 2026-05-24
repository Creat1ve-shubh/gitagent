# How to Use Gitclaw's New Features

Welcome to the new era of Gitclaw! Following our major architectural migration to Go, we have introduced powerful enterprise-grade features that make agent execution faster, safer, and completely reliable. 

Here is a guide on how to take advantage of these new capabilities.

---

## 1. The MVCC Write Ledger (Conflict-Free File Management)

### What is it?
Multi-Version Concurrency Control (MVCC) is a backend technology (commonly used in databases like PostgreSQL) that we have integrated into our agent file system. It ensures that when multiple agents or parallel processes try to edit the same file simultaneously, there are **no race conditions or corrupted files**.

### How to use it:
The beauty of the MVCC ledger is that it works **automatically in the background**. You do not need to change your scripts or agent prompts.

However, if you want to trace file operations or rollback changes, the MVCC ledger provides an immutable audit trail.

**Usage:**
- Simply let your agents use the `write`, `edit`, and `sandbox-edit` tools as normal. 
- The runtime automatically serializes concurrent writes. If Agent A and Agent B modify the same file at the exact same millisecond, the ledger ensures one succeeds safely and the other is either merged cleanly or rejected with an actionable error.
- Check the `.gitagent/ledger.log` (if enabled in your config) to see the transaction history of all file mutations.

---

## 2. The Stateless Circuit-Breaker Guard Pipeline

### What is it?
Security is paramount when giving AI access to your file system and terminal. The stateless circuit-breaker pipeline is a zero-latency policy engine that intercepts **every single tool call** an agent attempts to make. If a command or file access violates your policy, it is instantly blocked.

### How to use it:
You can enforce strict boundaries by defining a guard policy. 

1. Create a `guard.json` file in your agent's directory (e.g., `agents/my-agent/guard.json`).
2. Define your allowlists and blocklists.

**Example `guard.json`:**
```json
{
  "version": "1.0",
  "filesystem": {
    "allowed_paths": ["./src", "./docs", "./tests"],
    "blocked_paths": ["./secrets", ".env", "node_modules"],
    "read_only": ["./docs"]
  },
  "commands": {
    "allowed_binaries": ["npm", "git", "go", "ls", "cat", "grep"],
    "blocked_args": ["rm -rf", "sudo"]
  },
  "network": {
    "allow_egress": false
  }
}
```

**What happens at runtime?**
If the agent tries to run `cli` with `rm -rf ./src`, the circuit-breaker will instantly return a `tool_result` error to the LLM: 
*`Error: Command blocked by guard policy: 'rm -rf' is not allowed.`* 
The LLM can then apologize and try an allowed command instead.

---

## 3. Sub-50ms Cold Starts (The Go Runtime)

### What is it?
By compiling Gitclaw to a single, dependency-free Go binary, we've eliminated the V8 engine startup time. The CLI now boots in under 50 milliseconds.

### How to use it:
Just run it! The speed is built-in. This makes Gitclaw perfect for:
- **CI/CD Pipelines:** Run agent tasks in GitHub Actions without waiting for `npm install` or node module resolution.
- **Pre-commit hooks:** You can now configure Gitclaw to run seamlessly before a git commit because the cold start is imperceptible.
- **Shell aliases:** Alias your common terminal commands to Gitclaw without feeling any lag.

**Example: Super-fast ad-hoc run:**
```bash
# This will execute almost instantly
gitclaw run --dir . --prompt "Summarize the recent commits"
```

---

## 4. Comprehensive OpenTelemetry Observability

### What is it?
We rebuilt our telemetry from the ground up in Go. Every LLM call, tool execution, and session is instrumented with zero overhead.

### How to use it:
To see exactly what your agents are doing in real-time, just set your endpoint:

```bash
# Export the OpenTelemetry endpoint
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318

# Run your agent
gitclaw run -p "Refactor the authentication module"
```

If you just want to see spans in your terminal locally:
```bash
OTEL_TRACES_EXPORTER=console gitclaw run -p "Fix the typo in README"
```

Enjoy the power, speed, and safety of the new Go-based Gitclaw!
