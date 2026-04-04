<p align="center">
  <img src="https://img.shields.io/badge/LogicaProxy-v1.0.0-8b5cf6?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAzMiAzMiI+PHJlY3Qgd2lkdGg9IjMyIiBoZWlnaHQ9IjMyIiByeD0iNiIgZmlsbD0iIzhiNWNmNiIvPjx0ZXh0IHg9IjE2IiB5PSIyMiIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZmlsbD0id2hpdGUiIGZvbnQtZmFtaWx5PSJzeXN0ZW0tdWkiIGZvbnQtd2VpZ2h0PSI3MDAiIGZvbnQtc2l6ZT0iMTYiPkxQPC90ZXh0Pjwvc3ZnPg==&logoColor=white" alt="LogicaProxy" height="32"/>
</p>

<h1 align="center">LogicaProxy</h1>

<p align="center">
  <strong>The intelligent CLI-to-API proxy that routes your AI models where they belong.</strong>
</p>

<p align="center">
  <a href="#features"><img src="https://img.shields.io/badge/Providers-6+-22c55e?style=flat-square" alt="Providers"/></a>
  <a href="#ai-router"><img src="https://img.shields.io/badge/AI_Router-Intelligent-8b5cf6?style=flat-square" alt="AI Router"/></a>
  <a href="#dashboard"><img src="https://img.shields.io/badge/Dashboard-Built--in-3b82f6?style=flat-square" alt="Dashboard"/></a>
  <a href="#vault"><img src="https://img.shields.io/badge/Vault-AES--256-ef4444?style=flat-square" alt="Vault"/></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-71717a?style=flat-square" alt="MIT License"/></a>
  <a href="https://github.com/Rovemark/LogicaProxy/releases"><img src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go"/></a>
  <a href="https://github.com/Rovemark/LogicaProxy/releases"><img src="https://img.shields.io/github/v/release/Rovemark/LogicaProxy?style=flat-square&color=8b5cf6" alt="Release"/></a>
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> &bull;
  <a href="#why-logicaproxy">Why LogicaProxy</a> &bull;
  <a href="#ai-router">AI Router</a> &bull;
  <a href="#dashboard">Dashboard</a> &bull;
  <a href="#vault">Vault</a> &bull;
  <a href="#api-reference">API</a> &bull;
  <a href="#sdk">SDK</a> &bull;
  <a href="#deployment">Deployment</a>
</p>

---

## The Problem

You pay for CLI subscriptions — Claude Code, Gemini CLI, OpenAI Codex. You want to use those models in your own tools, SDKs, and applications. But CLI subscriptions don't give you API access.

**LogicaProxy bridges that gap.** It authenticates via OAuth to your CLI subscriptions and exposes standard API endpoints that any compatible client can consume. No separate API keys needed.

But unlike other proxy tools that just forward requests, LogicaProxy **thinks about where your request should go**.

---

## Why LogicaProxy

Most CLI proxies are dumb pipes. They take a request and forward it. LogicaProxy is different.

### The Routing Engine

Every request that hits LogicaProxy gets analyzed. The AI Router reads your prompt, classifies the task, and routes it to the optimal model — automatically.

| Task Type | What it detects | Routes to |
|-----------|----------------|-----------|
| **Code** | `implement`, `debug`, `refactor`, language keywords | Best coding model |
| **Reasoning** | `analyze`, `compare`, `strategy`, `architecture` | Strongest reasoning model |
| **Fast** | `translate`, `summarize`, `list`, `simple` | Fastest model |
| **Creative** | `write`, `story`, `blog`, `marketing` | Creative model |
| **Analysis** | `data`, `metrics`, `forecast`, `report` | Analytical model |

Set your model to `auto` and let the router decide. Or pin a specific model — your call.

### Feature Comparison

| Feature | LogicaProxy | Other proxies |
|---------|:-----------:|:-------------:|
| Intelligent routing by task type | ✅ | ❌ |
| Built-in web dashboard | ✅ | ❌ |
| Encrypted credential storage | ✅ | ❌ |
| Prometheus metrics | ✅ | ❌ |
| Health & readiness probes | ✅ | ❌ |
| Multi-account load balancing | ✅ | ✅ |
| Multi-provider OAuth | ✅ | ✅ |
| Streaming & WebSocket | ✅ | ✅ |
| Protocol translation | ✅ | ✅ |
| Go SDK for embedding | ✅ | ✅ |

---

<h2 id="features">Features</h2>

### Multi-Provider OAuth Routing

Connect your CLI subscriptions and use them as standard API endpoints:

| Provider | Auth Method | Models |
|----------|------------|--------|
| **Claude Code** | Anthropic OAuth | claude-opus, claude-sonnet, claude-haiku |
| **Gemini CLI** | Google OAuth | gemini-pro, gemini-flash |
| **OpenAI Codex** | OpenAI OAuth | gpt-4o, gpt-4o-mini |
| **Qwen Code** | Alibaba OAuth | qwen-plus, qwen-turbo |
| **iFlow** | iFlow OAuth | iflow models |
| **Antigravity** | AG OAuth | ag models |

All providers support **multi-account round-robin load balancing** — add multiple accounts per provider to distribute load and avoid rate limits.

### Streaming & Tools

- Server-Sent Events (SSE) streaming
- WebSocket streaming for real-time responses
- Function calling / tool use
- Multimodal input (text + images)
- Thinking/reasoning mode support

### Protocol Translation

One proxy speaks every format. Send requests in any format, receive responses in any format:

```
OpenAI format  →  ┌──────────────┐  →  Claude API
Claude format  →  │  LogicaProxy  │  →  Gemini API
Gemini format  →  └──────────────┘  →  OpenAI API
```

| Endpoint | Format | Use case |
|----------|--------|----------|
| `/v1/chat/completions` | OpenAI | Most SDKs and tools |
| `/v1/messages` | Claude | Anthropic SDK |
| `/v1/responses` | OpenAI Responses | Streaming + WebSocket |
| `/v1beta/models/:model` | Gemini | Google AI SDK |

---

<h2 id="ai-router">AI Router</h2>

The AI Router is a middleware that classifies your prompt and selects the optimal model. It reads keywords and patterns in your request to determine the task type.

### Configuration

```yaml
# config.yaml
ai-router:
  enabled: true
  default-model: claude-sonnet-4-6
  routes:
    code: claude-sonnet-4-6
    reasoning: claude-opus-4-6
    fast: claude-haiku-4-5-20251001
    creative: claude-sonnet-4-6
    analysis: claude-opus-4-6
```

### Usage

Send a request with `"model": "auto"` and the router picks the best model:

```bash
curl http://localhost:8317/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "x-api-key: your-key" \
  -d '{
    "model": "auto",
    "messages": [{"role": "user", "content": "refactor this function to use async/await"}]
  }'
# → Routes to code model (detected: refactor, function, async/await)
```

### Response Headers

Every response includes routing metadata:

```
X-AI-Category: code
X-AI-Route: claude-sonnet-4-6
```

---

<h2 id="dashboard">Dashboard</h2>

LogicaProxy ships with a built-in web dashboard. No external dependencies, no separate installation, no configuration.

```
http://localhost:8317/dashboard/
```

### What's on the dashboard

- **Status & Uptime** — live health with sparkline charts
- **Request metrics** — total requests, error rate, avg latency, req/s
- **Request volume chart** — bar chart with hover tooltips showing timestamp and count
- **Available models** — auto-discovered from connected providers
- **Per-endpoint breakdown** — requests, errors, latency per route
- **Recent requests log** — last 10 requests in real-time
- **Connection test** — ping button with latency history
- **Supported providers** — status of all 6 providers
- **API reference** — complete endpoint documentation
- **Dark / Light theme** — toggle with localStorage persistence

The dashboard is embedded in the binary via Go's `embed` package — zero runtime dependencies.

---

<h2 id="vault">Credential Vault</h2>

OAuth tokens are sensitive. LogicaProxy includes a built-in encrypted vault for credential storage.

| Feature | Detail |
|---------|--------|
| Encryption | AES-256-GCM |
| Key derivation | SHA-256 from passphrase |
| Expiration | Automatic token expiry and purge |
| Concurrency | Thread-safe with RWMutex |
| Verification | Key fingerprint display |
| Operations | Store, Retrieve, Delete, List, PurgeExpired |

No more plaintext JSON files with your OAuth tokens sitting on disk.

---

<h2 id="quick-start">Quick Start</h2>

### From Source

```bash
git clone https://github.com/Rovemark/LogicaProxy.git
cd LogicaProxy
go build -o logicaproxy ./cmd/server
./logicaproxy -config config.example.yaml
```

### Docker

```bash
docker build -t logicaproxy .
docker run -p 8317:8317 \
  -v ./config.yaml:/LogicaProxy/config.yaml \
  -v ./auths:/root/.logicaproxy \
  logicaproxy
```

### Docker Compose

```bash
docker-compose up -d
```

### First-Time Setup

1. Copy and edit the config:
```bash
cp config.example.yaml config.yaml
```

2. Add your API key:
```yaml
host: "127.0.0.1"
port: 8317
api-keys:
  - "your-secret-key"
auth-dir: "~/.logicaproxy"
```

3. Login to a provider:
```bash
# Claude Code
./logicaproxy -claude-login

# Gemini CLI
./logicaproxy -gemini-login

# OpenAI Codex
./logicaproxy -openai-login

# Qwen Code
./logicaproxy -qwen-login
```

4. Start the proxy:
```bash
./logicaproxy -config config.yaml
```

5. Open the dashboard:
```
http://localhost:8317/dashboard/
```

---

<h2 id="api-reference">API Reference</h2>

### AI Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/chat/completions` | `POST` | OpenAI-compatible chat completions |
| `/v1/messages` | `POST` | Claude Messages API |
| `/v1/messages/count_tokens` | `POST` | Claude token counting |
| `/v1/completions` | `POST` | Legacy completions |
| `/v1/responses` | `POST` `WS` | OpenAI Responses API |
| `/v1/responses/compact` | `POST` | Compact responses |
| `/v1/models` | `GET` | List available models |
| `/v1beta/models/:model` | `POST` | Gemini generate content |

### System Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | `GET` | Health check — returns status, version, uptime |
| `/ready` | `GET` | Readiness probe for Kubernetes |
| `/metrics` | `GET` | Prometheus-compatible metrics |
| `/dashboard/` | `GET` | Web dashboard UI |
| `/management.html` | `GET` | Management control panel |

### Amp CLI Endpoints

| Endpoint | Description |
|----------|-------------|
| `/api/provider/{provider}/v1/messages` | Provider-specific messages |
| `/api/provider/{provider}/v1/chat/completions` | Provider-specific completions |
| `/api/provider/{provider}/v1beta/models/...` | Provider-specific generate |

---

<h2 id="sdk">Go SDK</h2>

Embed LogicaProxy in your own Go application:

```go
package main

import "github.com/Rovemark/LogicaProxy/v6/sdk/cliproxy"

func main() {
    service := cliproxy.NewBuilder().
        WithConfig("config.yaml").
        Build()

    if err := service.Start(); err != nil {
        panic(err)
    }
}
```

### SDK Documentation

| Doc | Description |
|-----|-------------|
| [Usage Guide](docs/sdk-usage.md) | Getting started with the SDK |
| [Advanced](docs/sdk-advanced.md) | Custom executors and translators |
| [Access Control](docs/sdk-access.md) | Authentication and authorization |
| [Watcher](docs/sdk-watcher.md) | Hot-reload configuration |
| [Examples](examples/) | Working code examples |

---

<h2 id="deployment">Deployment</h2>

### Monitoring with Prometheus

Scrape the `/metrics` endpoint:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'logicaproxy'
    static_configs:
      - targets: ['localhost:8317']
```

Available metrics:

```
logicaproxy_uptime_seconds 3600.00
logicaproxy_requests_total 1542
logicaproxy_errors_total 3
logicaproxy_endpoint_requests_total{method="POST",path="/v1/messages"} 1200
logicaproxy_endpoint_errors_total{method="POST",path="/v1/messages"} 0
logicaproxy_endpoint_avg_latency_ms{method="POST",path="/v1/messages"} 2450
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: logicaproxy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: logicaproxy
  template:
    metadata:
      labels:
        app: logicaproxy
    spec:
      containers:
      - name: logicaproxy
        image: rovemark/logicaproxy:latest
        ports:
        - containerPort: 8317
        livenessProbe:
          httpGet:
            path: /health
            port: 8317
          initialDelaySeconds: 5
        readinessProbe:
          httpGet:
            path: /ready
            port: 8317
          initialDelaySeconds: 3
        volumeMounts:
        - name: config
          mountPath: /LogicaProxy/config.yaml
          subPath: config.yaml
      volumes:
      - name: config
        configMap:
          name: logicaproxy-config
```

### Process Manager (PM2)

```javascript
// ecosystem.config.js
module.exports = {
  apps: [{
    name: 'logicaproxy',
    script: '/path/to/logicaproxy',
    args: '-config /path/to/config.yaml',
    autorestart: true,
    restart_delay: 5000,
  }]
};
```

### Systemd

```ini
[Unit]
Description=LogicaProxy
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/logicaproxy -config /etc/logicaproxy/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## Architecture

```
Client Request
     │
     ▼
┌─────────────────┐
│   API Gateway    │  Gin HTTP server, CORS, auth middleware
├─────────────────┤
│   AI Router      │  Task classification → model selection
├─────────────────┤
│   Metrics        │  Request counting, latency tracking, Prometheus export
├─────────────────┤
│   Translator     │  OpenAI ↔ Claude ↔ Gemini format conversion
├─────────────────┤
│   Executor       │  Provider-specific OAuth + request forwarding
├─────────────────┤
│   Vault          │  AES-256-GCM encrypted credential storage
├─────────────────┤
│   Watcher        │  Hot-reload config and auth files (fsnotify)
└─────────────────┘
```

---

## Configuration Reference

```yaml
# Server
host: "127.0.0.1"            # Bind address
port: 8317                    # Listen port
debug: false                  # Debug logging

# Authentication
api-keys:                     # API keys for client auth
  - "your-key-here"
auth-dir: "~/.logicaproxy"    # OAuth token storage directory

# Request handling
request-retry: 2              # Retry failed requests
max-retry-interval: 30        # Max retry delay (seconds)

# AI Router
ai-router:
  enabled: true
  default-model: ""           # Fallback model
  routes:
    code: "claude-sonnet-4-6"
    reasoning: "claude-opus-4-6"
    fast: "claude-haiku-4-5-20251001"

# Logging
request-log: true             # Log requests
error-logs-max-files: 10      # Max error log files

# Advanced
commercial-mode: false        # Disable middleware for high concurrency
websocket-auth: false         # Require auth for WebSocket connections
```

See [`config.example.yaml`](config.example.yaml) for the complete reference.

---

## Contributing

Contributions are welcome. For major changes, please open an issue first.

```bash
# Fork and clone
git clone https://github.com/your-username/LogicaProxy.git
cd LogicaProxy

# Create feature branch
git checkout -b feature/your-feature

# Run tests
go test ./...

# Commit and push
git commit -m "feat: your feature"
git push origin feature/your-feature

# Open a Pull Request
```

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

<p align="center">
  <strong>LogicaProxy</strong> by <a href="https://github.com/Rovemark">Rovemark</a><br/>
  <sub>Open-source infrastructure for AI developers.</sub>
</p>
