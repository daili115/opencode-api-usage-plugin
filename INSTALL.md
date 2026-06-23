# OpenCode API Usage Plugin - AI Auto-Install Prompt

## Usage

Copy the prompt below into OpenCode's input box, and the AI will automatically complete the plugin installation and integration.

---

## Install Prompt

```
Please help me install the OpenCode API usage statistics plugin.

## Plugin Source

Repository: https://github.com/daili115/opencode-api-usage-plugin

## Installation Steps

Please complete the installation automatically following these steps:

### 1. Download Plugin File

Download from https://raw.githubusercontent.com/daili115/opencode-api-usage-plugin/main/api_usage.go and save to `internal/tui/components/chat/api_usage.go`

### 2. Modify editor.go

Find the `internal/tui/components/chat/editor.go` file and make the following modifications:

#### 2.1 Add apiUsage field to editorCmp struct

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

#### 2.2 Modify NewEditorCmp function

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

#### 2.3 Modify View method

At the end of the existing View method, add API statistics rendering:

```go
func (m *editorCmp) View() string {
    // ... original code remains unchanged ...
    
    // Add the following code at the end of the method
    apiUsageView := m.apiUsage.View()
    if apiUsageView != "" {
        return lipgloss.JoinVertical(lipgloss.Top,
            editorView,  // This is the original editor view
            apiUsageView,
        )
    }
    return editorView
}
```

#### 2.4 Modify SetSize method

```go
func (m *editorCmp) SetSize(width, height int) tea.Cmd {
    m.width = width
    m.height = height
    m.textarea.SetWidth(width - 3)
    m.textarea.SetHeight(height - 3)  // Reserve 3 lines for statistics
    m.textarea.SetWidth(width)
    m.apiUsage.SetSize(width, 3)  // Add this line
    return nil
}
```

#### 2.5 Add APIUsageMsg handling in Update method

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
}
```

### 3. Add Statistics Update at LLM Response

Find the code that handles LLM responses (usually in `internal/llm/agent/` directory), and add after receiving the response:

```go
// When receiving LLM response, extract usage information and send
if resp.Usage != nil {
    // Send API usage statistics message
    // Use your project's message passing mechanism here
}
```

### 4. Verify Compilation

Run `go build` to ensure the code compiles successfully.

## Plugin Features

After installation, you will see below the input box:

- **Token Usage Progress Bar** - Shows current token usage and limits
- **Hit Rate Progress Bar** - Shows cache hit rate
- **Detailed Statistics** - Prompt/Completion Token count and cost

Shortcut: `Ctrl+U` to toggle display/hide

## Notes

1. Ensure the package name in `api_usage.go` is `chat`
2. If compilation errors occur, check if the `time` package is imported correctly
3. Statistics data updates need to extract usage information from LLM Provider's response

Please execute the above steps to complete the installation.
```

---

## Quick Install Commands

If you are using an AI that supports tool calls (like Claude Code), you can use the following commands directly:

```bash
# 1. Download plugin file
curl -o internal/tui/components/chat/api_usage.go \
  https://raw.githubusercontent.com/daili115/opencode-api-usage-plugin/main/api_usage.go

# 2. Let AI help you modify editor.go
# Copy the prompt above into OpenCode input box

# 3. Compile and verify
go build ./...
```

---

## Manual Installation Checklist

- [ ] `api_usage.go` copied to `internal/tui/components/chat/`
- [ ] `editorCmp` struct added `apiUsage` field
- [ ] `NewEditorCmp` function initialized `apiUsage`
- [ ] `View()` method renders API statistics
- [ ] `SetSize()` reserves space for statistics
- [ ] `Update()` method handles `APIUsageMsg`
- [ ] LLM response handling sends `APIUsageMsg`
- [ ] `go build` compiles successfully

---

## Troubleshooting

### Compilation Error: undefined: apiUsageCmp
**Cause**: `api_usage.go` file not imported correctly
**Solution**: Ensure the file is in `internal/tui/components/chat/` directory with package name `chat`

### Compilation Error: undefined: APIUsageMsg
**Cause**: Message type not defined
**Solution**: Check if `api_usage.go` contains `APIUsageMsg` struct definition

### Statistics Display Blank
**Cause**: No `APIUsageMsg` received
**Solution**: Add message sending logic in LLM response handling

### Style Display Abnormal
**Cause**: Color value or style definition issue
**Solution**: Check `lipgloss` style definitions, ensure color value format is correct
