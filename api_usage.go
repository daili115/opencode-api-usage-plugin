package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// APIUsageStats stores API usage statistics information
type APIUsageStats struct {
	// Token usage
	PromptTokens     int64
	CompletionTokens int64
	TotalTokens      int64

	// Hit rate statistics
	CacheHits   int64
	CacheMisses int64
	HitRate     float64

	// Limit information
	LimitTotal     int64
	LimitRemaining int64
	LimitResetTime time.Time

	// Cost
	TotalCost float64
}

// APIUsageMsg is used to update API usage statistics messages
type APIUsageMsg struct {
	Stats APIUsageStats
}

// apiUsageCmp is the API usage statistics component
type apiUsageCmp struct {
	width   int
	height  int
	stats   APIUsageStats
	visible bool
}

// APIUsageKeyMaps defines shortcut keys
type APIUsageKeyMaps struct {
	Toggle key.Binding
}

var apiUsageKeys = APIUsageKeyMaps{
	Toggle: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "toggle API usage stats"),
	),
}

// NewAPIUsageCmp creates a new API usage statistics component
func NewAPIUsageCmp() tea.Model {
	return &apiUsageCmp{
		visible: true,
		stats:   APIUsageStats{},
	}
}

func (m *apiUsageCmp) Init() tea.Cmd {
	return nil
}

func (m *apiUsageCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, apiUsageKeys.Toggle) {
			m.visible = !m.visible
			return m, nil
		}
	case APIUsageMsg:
		m.stats = msg.Stats
		// Calculate hit rate
		if m.stats.CacheHits+m.stats.CacheMisses > 0 {
			m.stats.HitRate = float64(m.stats.CacheHits) / float64(m.stats.CacheHits+m.stats.CacheMisses) * 100
		}
		return m, nil
	}
	return m, nil
}

// View renders API usage statistics information
func (m *apiUsageCmp) View() string {
	if !m.visible {
		return ""
	}

	// Use clean modern style
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1)

	// Token usage progress bar
	tokenBar := m.renderProgressBar(
		m.stats.PromptTokens+m.stats.CompletionTokens,
		m.stats.LimitTotal,
		40,
		"#00D9FF", // Cyan
		"#333333",
	)

	// Hit rate progress bar
	hitRateBar := m.renderProgressBar(
		int64(m.stats.HitRate),
		100,
		20,
		"#00FF88", // Green
		"#333333",
	)

	var parts []string

	// Token usage display
	if m.stats.LimitTotal > 0 {
		remaining := m.stats.LimitTotal - m.stats.TotalTokens
		if remaining < 0 {
			remaining = 0
		}
		percentage := float64(remaining) / float64(m.stats.LimitTotal) * 100

		resetTime := ""
		if !m.stats.LimitResetTime.IsZero() {
			resetTime = m.stats.LimitResetTime.Format("Jan 2")
		}

		parts = append(parts, style.Render(
			fmt.Sprintf("Token: %s  Remaining %.0f%% (Reset: %s)",
				tokenBar,
				percentage,
				resetTime,
			),
		))
	} else {
		parts = append(parts, style.Render(
			fmt.Sprintf("Token: %s  %d/%d",
				tokenBar,
				m.stats.TotalTokens,
				m.stats.LimitTotal,
			),
		))
	}

	// Hit rate display
	parts = append(parts, style.Render(
		fmt.Sprintf("Hit Rate: %s  %.1f%% (%d/%d)",
			hitRateBar,
			m.stats.HitRate,
			m.stats.CacheHits,
			m.stats.CacheHits+m.stats.CacheMisses,
		),
	))

	// Detailed statistics
	detailStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Padding(0, 1)

	parts = append(parts, detailStyle.Render(
		fmt.Sprintf("Prompt: %d  |  Completion: %d  |  Cost: $%.4f",
			m.stats.PromptTokens,
			m.stats.CompletionTokens,
			m.stats.TotalCost,
		),
	))

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderProgressBar renders a progress bar
func (m *apiUsageCmp) renderProgressBar(current, total int64, width int, fillColor, emptyColor string) string {
	if total <= 0 {
		return strings.Repeat("░", width)
	}

	percentage := float64(current) / float64(total)
	if percentage > 1 {
		percentage = 1
	}

	filled := int(float64(width) * percentage)
	empty := width - filled

	var bar strings.Builder

	// Filled part
	if filled > 0 {
		fillStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fillColor))
		bar.WriteString(fillStyle.Render(strings.Repeat("█", filled)))
	}

	// Empty part
	if empty > 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(emptyColor))
		bar.WriteString(emptyStyle.Render(strings.Repeat("░", empty)))
	}

	return bar.String()
}

func (m *apiUsageCmp) SetSize(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	return nil
}

func (m *apiUsageCmp) GetSize() (int, int) {
	return m.width, m.height
}

func (m *apiUsageCmp) BindingKeys() []key.Binding {
	return []key.Binding{apiUsageKeys.Toggle}
}

// SetStats sets statistics information (called externally)
func (m *apiUsageCmp) SetStats(stats APIUsageStats) tea.Cmd {
	return func() tea.Msg {
		return APIUsageMsg{Stats: stats}
	}
}
