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

1. Copy `api-usage-plugin.go` to `internal/tui/components/chat/api_usage.go`
2. Follow the integration steps in `integration-guide.md`

## Integration

See [integration-guide.md](integration-guide.md) for detailed integration instructions.

## Requirements

- Go 1.24+
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)

## License

MIT
