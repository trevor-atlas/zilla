package mvc

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/trevor-atlas/zilla/jira"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func main() {
	t := textinput.New()
	t.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	service := jira.NewService()

	initialModel := Model{
		textInput:  t,
		spinner:    s,
		typing:     true,
		jiraClient: &service,
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
	jiraClient *jira.ClientService

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

func (m *Model) fetchIssues() tea.Cmd {
	return func() tea.Msg {
		issues, err := m.jiraClient.GetIssues()
		if err != nil {
			return GotIssues{Err: err}
		}

		return GotIssues{Issues: issues}
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

		if err := msg.Err; err != nil {
			m.err = err
			return m, nil
		}

		m.issues = msg.Issues
		return m, nil
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
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

	return m, nil
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
	items := []list.Item{}
	for _, issue := range m.issues.Issues {
		append(items, item{title: issue.Key, desc: issue.Fields.Summary})
	}
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Issues"

	return docStyle.Render(m.list.View())
}
