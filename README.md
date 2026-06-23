# OpenCode API Usage Plugin

A plugin for [OpenCode](https://github.com/opencode-ai/opencode) that displays API usage statistics below the input box, including token usage, cache hit rate, and cost information.

## Features

- **Token Usage Progress Bar**: Visual display of current token consumption and limits
- **Cache Hit Rate**: Real-time statistics showing cache efficiency
- **Cost Tracking**: Session-based API cost estimation
- **Toggle Display**: Show/hide statistics with `Ctrl+U`

## Preview

```
> Ask anything... "Fix a TODO in the codebase"

Token: ████████████████████░░░░░░░░░░░░░░░░░░░░  Remaining 77% (Reset: Jul 23)
Hit Rate: ████████████████████░░░░  90.0% (45/50)
Prompt: 1500  |  Completion: 800  |  Cost: $0.0234
```

## Installation

### Method 1: One-Line Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/daili115/opencode-api-usage-plugin/main/install.sh | bash
```

This will automatically:
1. Download the plugin file
2. Backup your `editor.go`
3. Apply all necessary modifications
4. Verify compilation

### Method 2: AI Auto-Install

Copy the prompt from [INSTALL.md](INSTALL.md) into OpenCode's input box, and the AI will automatically install the plugin for you.

### Method 3: Manual Install

1. Copy `api_usage.go` to `internal/tui/components/chat/api_usage.go`
2. Follow the integration steps in [integration-guide.md](integration-guide.md)

### Method 4: Quick Install

```bash
# Download plugin file
curl -o internal/tui/components/chat/api_usage.go \
  https://raw.githubusercontent.com/daili115/opencode-api-usage-plugin/main/api_usage.go

# Then follow the integration guide to modify editor.go
```

## Quick Start

Run the example:

```bash
git clone https://github.com/daili115/opencode-api-usage-plugin.git
cd opencode-api-usage-plugin
go mod tidy
go run example.go
```

## Files

| File | Description |
|------|-------------|
| `api_usage.go` | Core plugin code |
| `example.go` | Standalone demo |
| `integration-guide.md` | Detailed integration guide |
| `INSTALL.md` | AI auto-install prompt |
| `install.sh` | One-line install script |

## Requirements

- Go 1.24+
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)

## License

MIT
