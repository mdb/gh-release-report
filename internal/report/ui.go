package report

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// AssetData holds download statistics for a single release asset.
type AssetData struct {
	Name          string
	DownloadCount int
}

// ReleaseData holds all information about a release needed for display.
type ReleaseData struct {
	RepoFullName string
	TagName      string
	PublishedAt  *time.Time
	URL          string
	Assets       []AssetData
	TotalCount   int
}

// Lipgloss styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")).
			Background(lipgloss.Color("0")).
			MarginBottom(1).
			Bold(true)

	urlStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true).
			MarginBottom(1).
			Underline(true)

	barStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14"))

	totalStyle = lipgloss.NewStyle().
			MarginTop(1).
			Foreground(lipgloss.Color("13"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("7")).
			Padding(0, 1)
)

// model implements the bubbletea Model interface for displaying release data.
type model struct {
	content string
}

// newModel creates a new model with pre-rendered release data.
func newModel(data *ReleaseData) model {
	return model{
		content: renderContent(data),
	}
}

// Init returns tea.Quit immediately since this is a one-shot display.
func (m model) Init() tea.Cmd {
	return tea.Quit
}

// Update handles messages; returns tea.Quit for any message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

// View renders the model content wrapped in a box.
func (m model) View() string {
	return boxStyle.Render(m.content) + "\n"
}

// renderContent builds the full display string from release data.
func renderContent(data *ReleaseData) string {
	var lines []string

	// Title line
	title := fmt.Sprintf("%s %s", data.RepoFullName, data.TagName)
	lines = append(lines, titleStyle.Render(title))

	// Published date
	if data.PublishedAt != nil {
		lines = append(lines, fmt.Sprintf("Published %s", data.PublishedAt))
	} else {
		lines = append(lines, "Published <nil>")
	}

	// URL
	lines = append(lines, urlStyle.Render(data.URL))

	// Bar chart or "no assets" message
	if len(data.Assets) == 0 {
		lines = append(lines, "No release assets")
	} else {
		lines = append(lines, renderBarChart(data.Assets))
	}

	// Total downloads
	p := message.NewPrinter(language.English)
	formattedTotal := p.Sprintf("%d", data.TotalCount)
	lines = append(lines, totalStyle.Render(formattedTotal)+" downloads")

	return strings.Join(lines, "\n")
}

// renderBarChart creates a horizontal bar chart from asset data.
func renderBarChart(assets []AssetData) string {
	if len(assets) == 0 {
		return ""
	}

	// Find the maximum download count and longest label
	maxCount := 0
	maxLabelLen := 0
	for _, asset := range assets {
		if asset.DownloadCount > maxCount {
			maxCount = asset.DownloadCount
		}
		if len(asset.Name) > maxLabelLen {
			maxLabelLen = len(asset.Name)
		}
	}

	// Build bars
	const maxBarWidth = 75
	var lines []string
	for _, asset := range assets {
		// Calculate bar width proportionally
		barWidth := 0
		if maxCount > 0 {
			barWidth = (asset.DownloadCount * maxBarWidth) / maxCount
		}
		if barWidth == 0 && asset.DownloadCount > 0 {
			barWidth = 1 // Minimum 1 character for non-zero counts
		}

		// Build the bar
		bar := strings.Repeat("█", barWidth)

		// Format: label (padded) | bar | count
		label := fmt.Sprintf("%-*s", maxLabelLen, asset.Name)
		line := fmt.Sprintf("%s %s %d",
			labelStyle.Render(label),
			barStyle.Render(bar),
			asset.DownloadCount,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
