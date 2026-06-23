# OpenCode API Usage Statistics Plugin Integration Guide

## Overview

This plugin adds an API usage statistics panel to OpenCode, displayed below the input box, including:
- **Token Usage**: Visual display with progress bar
- **Hit Rate**: Real-time cache hit rate statistics
- **Cost Estimation**: Session-based API cost estimation

## File Structure

```
opencode/
├── internal/
│   └── tui/
│       └── components/
│           └── chat/
│               ├── editor.go          # Original editor component
│               └── api_usage.go       # New API statistics component
```

## Integration Steps

### 1. Create API Statistics Component

Copy the `api_usage.go` file to the `internal/tui/components/chat/` directory.

### 2. Modify editor.go

Add the API statistics component to the `editorCmp` struct:

```go
type editorCmp struct {
    width       int
    height      int
    app         *app.App
    session     session.Session
    textarea    textarea.Model
    attachments []message.Attachment
    deleteMode  bool
    apiUsage    *apiUsageCmp  // Add this line
}
```

### 3. Modify NewEditorCmp Function

```go
func NewEditorCmp(app *app.App) tea.Model {
    ta := CreateTextArea(nil)
    return &editorCmp{
        app:      app,
        textarea: ta,
        apiUsage: NewAPIUsageCmp().(*apiUsageCmp),  // Add this line
    }
}
```

### 4. Modify View Method

In the `View()` method, add API statistics information to the original return:

```go
func (m *editorCmp) View() string {
    t := theme.CurrentTheme()
    style := lipgloss.NewStyle().
        Padding(0, 0, 0, 1).
        Bold(true).
        Foreground(t.Primary())

    var editorView string
    if len(m.attachments) == 0 {
        editorView = lipgloss.JoinHorizontal(lipgloss.Top, style.Render(">"), m.textarea.View())
    } else {
        m.textarea.SetHeight(m.height - 1)
        editorView = lipgloss.JoinVertical(lipgloss.Top,
            m.attachmentsContent(),
            lipgloss.JoinHorizontal(lipgloss.Top, style.Render(">"),
                m.textarea.View()),
        )
    }

    // Add API usage statistics
    apiUsageView := m.apiUsage.View()
    if apiUsageView != "" {
        return lipgloss.JoinVertical(lipgloss.Top,
            editorView,
            apiUsageView,
        )
    }

    return editorView
}
```

### 5. Modify SetSize Method

```go
func (m *editorCmp) SetSize(width, height int) tea.Cmd {
    m.width = width
    m.height = height
    // Reserve space for API statistics (3 lines)
    m.textarea.SetWidth(width - 3)
    m.textarea.SetHeight(height - 3)
    m.textarea.SetWidth(width)
    m.apiUsage.SetSize(width, 3)  // Add this line
    return nil
}
```

### 6. Add Update Message Handling

Add handling for API usage statistics messages in the `Update` method:

```go
func (m *editorCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    switch msg := msg.(type) {
    // ... original code ...
    
    case APIUsageMsg:
        // Forward API usage statistics message
        _, cmd := m.apiUsage.Update(msg)
        return m, cmd
    }
    
    // ... original code ...
    
    m.textarea, cmd = m.textarea.Update(msg)
    return m, cmd
}
```

### 7. Update Statistics at LLM Call Site

In `internal/llm/agent/` or related files, send APIUsageMsg when receiving LLM response:

```go
// In agent's response handling
func (a *Agent) handleResponse(resp *llm.Response) {
    // ... original processing logic ...
    
    // Send API usage statistics
    if resp.Usage != nil {
        pubsub.Publish(APIUsageMsg{
            Stats: APIUsageStats{
                PromptTokens:     resp.Usage.PromptTokens,
                CompletionTokens: resp.Usage.CompletionTokens,
                TotalTokens:      resp.Usage.TotalTokens,
                // Other statistics...
            },
        })
    }
}
```

## Shortcut Keys

- `Ctrl+U`: Toggle API usage statistics panel display/hide

## Custom Styles

You can customize the appearance by modifying color constants in `api_usage.go`:

```go
const (
    TokenBarColor  = "#00D9FF"  // Token progress bar color (cyan)
    HitRateColor   = "#00FF88"  // Hit rate progress bar color (green)
    EmptyBarColor  = "#333333"  // Empty part color
    TextColor      = "#888888"  // Text color
    DetailColor    = "#666666"  // Detail information color
)
```

## Data Flow

```
LLM Provider Response
    ↓
Extract Usage information (Tokens, Cost)
    ↓
Send to PubSub / Message Bus
    ↓
Editor Component receives APIUsageMsg
    ↓
Update apiUsageCmp state
    ↓
Re-render View
```

## Notes

1. **Performance**: Statistics updates should not block the main UI thread
2. **Accuracy**: Actual token count should be based on LLM Provider return
3. **Cache Statistics**: Hit rate requires LLM Provider to support cache header information
4. **Limit Information**: Monthly limits require user configuration or API retrieval
