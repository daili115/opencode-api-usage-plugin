package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// This is an example showing how to integrate the API usage component
// with a Bubble Tea application similar to OpenCode's editor

type exampleModel struct {
	editor   textarea.Model
	apiUsage *apiUsageCmp
	width    int
	height   int
}

func initialModel() tea.Model {
	ta := textarea.New()
	ta.Placeholder = "Ask anything... \"Fix a TODO in the codebase\""
	ta.ShowLineNumbers = false
	ta.Focus()

	return &exampleModel{
		editor:   ta,
		apiUsage: NewAPIUsageCmp().(*apiUsageCmp),
	}
}

func (m *exampleModel) Init() tea.Cmd {
	return tea.Batch(
		m.editor.Init(),
		m.apiUsage.Init(),
		simulateAPIUsage(),
	)
}

func (m *exampleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.editor.SetWidth(msg.Width)
		m.editor.SetHeight(msg.Height - 4)
		m.apiUsage.SetSize(msg.Width, 3)

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Update editor
	editor, cmd := m.editor.Update(msg)
	m.editor = editor
	cmds = append(cmds, cmd)

	// Update API usage statistics
	apiUsage, cmd := m.apiUsage.Update(msg)
	m.apiUsage = apiUsage.(*apiUsageCmp)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *exampleModel) View() string {
	// Editor view
	editorView := m.editor.View()

	// API usage statistics view
	statsView := m.apiUsage.View()

	// Combined view: editor on top, statistics below
	return lipgloss.JoinVertical(
		lipgloss.Top,
		editorView,
		statsView,
	)
}

func simulateAPIUsage() tea.Cmd {
	return func() tea.Msg {
		// Simulate API usage statistics
		return APIUsageMsg{
			Stats: APIUsageStats{
				PromptTokens:     1500,
				CompletionTokens: 800,
				TotalTokens:      2300,
				CacheHits:        45,
				CacheMisses:      5,
				HitRate:          90.0,
				LimitTotal:       10000,
				LimitRemaining:   7700,
				LimitResetTime:   time.Now().Add(30 * 24 * time.Hour),
				TotalCost:        0.0234,
			},
		}
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
