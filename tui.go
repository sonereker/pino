package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type view int

const (
	listView view = iota
	detailView
)

type model struct {
	list       list.Model
	viewport   viewport.Model
	activeView view
	selected   *Pino
	width      int
	height     int
	ready      bool
}

var keys = struct {
	Enter  key.Binding
	Back   key.Binding
	Delete key.Binding
	Quit   key.Binding
}{
	Enter:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "view")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Delete: key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete")),
	Quit:   key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
}

func newModel(pinos []Pino) model {
	items := make([]list.Item, len(pinos))
	for i, p := range pinos {
		items[i] = p
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(colorAmber).
		BorderLeftForeground(colorAmber)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(colorMuted).
		BorderLeftForeground(colorAmber)

	l := list.New(items, delegate, 0, 0)
	l.Title = "🍕 pino"
	l.Styles.Title = titleStyle
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(colorBlue)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(colorAmber)
	l.SetFilteringEnabled(true)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Enter, keys.Delete}
	}

	return model{list: l}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		if !m.ready {
			m.viewport = viewport.New(msg.Width-h, msg.Height-v-4)
			m.viewport.Style = viewportStyle
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - h
			m.viewport.Height = msg.Height - v - 4
		}
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch m.activeView {
		case listView:
			switch {
			case key.Matches(msg, keys.Enter):
				if item, ok := m.list.SelectedItem().(Pino); ok {
					m.selected = &item
					m.activeView = detailView
					m.viewport.SetContent(renderDetail(item, m.viewport.Width-2))
					m.viewport.GotoTop()
					return m, nil
				}
			case key.Matches(msg, keys.Delete):
				if item, ok := m.list.SelectedItem().(Pino); ok {
					_ = deletePino(item.Filename)
					idx := m.list.Index()
					m.list.RemoveItem(idx)
					return m, m.list.NewStatusMessage(
						statusStyle.Render("deleted: " + item.Summary),
					)
				}
			case key.Matches(msg, keys.Quit):
				return m, tea.Quit
			}

		case detailView:
			switch {
			case key.Matches(msg, keys.Back), key.Matches(msg, keys.Quit):
				m.activeView = listView
				m.selected = nil
				return m, nil
			}
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	}

	switch m.activeView {
	case listView:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	case detailView:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.activeView {
	case detailView:
		header := detailTitleStyle.Render(m.selected.Summary)
		meta := metaStyle.Render(fmt.Sprintf(
			"%s  •  %s",
			m.selected.Filename,
			m.selected.CreatedAt.Format("Jan 02, 2006 15:04"),
		))
		footer := helpStyle.Render("esc: back  j/k: scroll  q: quit")
		return appStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left, header, meta, m.viewport.View(), footer),
		)
	default:
		return appStyle.Render(m.list.View())
	}
}

func renderDetail(p Pino, width int) string {
	wrap := lipgloss.NewStyle().Width(width)

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("PROMPT"))
	b.WriteString("\n\n")
	b.WriteString(wrap.Render(contentStyle.Render(p.Prompt)))

	if p.Plan != "" {
		b.WriteString("\n\n")
		b.WriteString(sectionTitleStyle.Render("PLAN"))
		b.WriteString("\n\n")
		b.WriteString(wrap.Render(contentStyle.Render(p.Plan)))
	}

	return b.String()
}
