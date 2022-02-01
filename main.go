package main

//import (
//	"fmt"
//	"github.com/trevor-atlas/zilla/jira"
//)
//
//func main() {
//	//cached := util.GetCachedIssues()
//	jiraClient := jira.NewService()
//	issues := jiraClient.GetMappedCustomFields()
//	fmt.Printf("\n%#v \n", issues)
//}

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/trevor-atlas/zilla/jira"
	"os"
	"strings"
)

type item struct {
	title, desc string
}

var docStyle = lipgloss.NewStyle()
var style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#44EEFF")).
	Border(lipgloss.NormalBorder(), false, false, false, true).
	PaddingTop(2)

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func main() {
	t := textinput.New()
	t.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	service := jira.NewService()

	items := []list.Item{}
	initialModel := Model{
		textInput:  t,
		spinner:    s,
		typing:     true,
		jiraClient: service,
		list:       list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	err := tea.NewProgram(initialModel, tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Model struct {
	textInput  textinput.Model
	spinner    spinner.Model
	jiraClient jira.ClientService

	typing  bool
	loading bool
	err     error
	issues  jira.JiraIssues
	list    list.Model
}

type GotIssues struct {
	Err    error
	Issues jira.JiraIssues
}

func (m Model) fetchIssues() tea.Cmd {
	return func() tea.Msg {
		issues, err := m.jiraClient.GetIssues(context.Background())
		if err != nil {
			return GotIssues{Err: err}
		}

		return GotIssues{Issues: *issues}
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.typing {
				query := strings.TrimSpace(m.textInput.Value())
				if query != "" {
					m.typing = false
					m.loading = true
					return m, tea.Batch(
						spinner.Tick,
						m.fetchIssues(),
					)
				}
			}

		case "esc":
			if !m.typing && !m.loading {
				m.typing = true
				m.err = nil
				return m, nil
			}
		}

	case GotIssues:
		m.loading = false
		m.typing = false

		if err := msg.Err; err != nil {
			m.err = err
			return m, nil
		}

		m.issues = msg.Issues
		m.list.Title = "Issues"

		for i, issue := range m.issues.Issues {
			m.list.InsertItem(i, item{title: issue.Key, desc: issue.Fields.Summary})
		}
		m.list.InsertItem(0, item{title: "COM-2156", desc: "IE-11 babel config support"})
		m.list.InsertItem(1, item{title: "COM-3121", desc: "IE-11 polyfills"})
		m.list.InsertItem(2, item{title: "COM-4199", desc: "convert to typescript"})
		m.list.InsertItem(3, item{title: "COM-2156", desc: "IE-11 babel config support"})
		m.list.InsertItem(4, item{title: "COM-3121", desc: "IE-11 polyfills"})
		m.list.InsertItem(5, item{title: "COM-4129", desc: "convert to typescript"})
		m.list.InsertItem(6, item{title: "COM-2156", desc: "IE-11 babel config support"})
		m.list.InsertItem(7, item{title: "COM-3121", desc: "IE-11 polyfills"})
		m.list.InsertItem(8, item{title: "COM-4139", desc: "convert to typescript"})
		m.list.InsertItem(9, item{title: "COM-2156", desc: "IE-11 babel config support"})
		m.list.InsertItem(10, item{title: "COM-3121", desc: "IE-11 polyfills"})
		m.list.InsertItem(11, item{title: "COM-4199", desc: "convert to typescript"})

		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize((msg.Width/3)-1, msg.Height)
		style.Width((msg.Width / 3) * 2).Height(msg.Height)
		return m, nil

	}

	if m.typing {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	if m.loading {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.typing {
		return fmt.Sprintf("Enter location:\n%s", m.textInput.View())
	}

	if m.loading {
		return fmt.Sprintf("%s fetching issues... please wait.", m.spinner.View())
	}

	if err := m.err; err != nil {
		return fmt.Sprintf("Could not fetch issues: %v", err)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), style.Render("Ticket details"))
}
