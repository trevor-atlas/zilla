package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/trevor-atlas/zilla/jira"
	"github.com/trevor-atlas/zilla/util"
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
	app := util.New()
	service := jira.NewService(app)
	initialModel := createModel(app, service)

	err := tea.NewProgram(initialModel, tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func createModel(app *util.Zilla, service jira.ClientService) Model {
	t := textinput.New()
	t.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	var items []list.Item
	model := Model{
		app:        *app,
		textInput:  t,
		spinner:    s,
		viewport:   viewport.New(0, 0),
		typing:     true,
		jiraClient: service,
		list:       list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	return model
}

type Model struct {
	app        util.Zilla
	textInput  textinput.Model
	spinner    spinner.Model
	jiraClient jira.ClientService

	viewport viewport.Model
	ready    bool
	typing   bool
	loading  bool
	err      error
	issues   jira.JiraIssues
	list     list.Model
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
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
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

		var issues jira.JiraIssues
		for i, _ := range m.issues.Issues {
			mockissue := jira.JiraIssue{
				ID:   fmt.Sprintf("%s", i),
				Self: fmt.Sprintf("%s", i),
				Key:  fmt.Sprintf("ABC-%s", i),
				Fields: jira.IssueFields{
					Summary:     "an issue summary",
					Created:     nil,
					Updated:     nil,
					Description: "Bacon ipsum dolor amet pig turkey bresaola, jowl fatback venison t-bone andouille. Boudin pork belly chicken meatball, short ribs shankle pork t-bone cow biltong doner. Brisket short ribs ribeye frankfurter pork loin buffalo shank picanha tenderloin turducken boudin pig. Picanha cupim ham ham hock burgdoggen pancetta chicken spare ribs salami landjaeger sausage brisket bacon kevin tenderloin.",
					Reporter: jira.IssueUser{
						Active:       false,
						TimeZone:     "",
						DisplayName:  "",
						Name:         "",
						EmailAddress: "",
						AvatarUrls:   nil,
						AccountId:    "",
						Key:          "",
						Self:         "",
					},
					Assignee: jira.IssueUser{
						Active:       false,
						TimeZone:     "",
						DisplayName:  "",
						Name:         "",
						EmailAddress: "",
						AvatarUrls:   nil,
						AccountId:    "",
						Key:          "",
						Self:         "",
					},
					Comment:   jira.IssueComments{},
					Priority:  jira.IssuePriority{},
					IssueType: jira.IssueType{},
					Status:    jira.IssueStatus{},
					Project:   jira.IssueProject{},
				},
			}
			issues.Issues = append(issues.Issues, mockissue)

			m.list.InsertItem(i, item{title: mockissue.Key, desc: mockissue.Fields.Summary})

			m.issues = issues //msg.Issues
			m.list.Title = "Issues"
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width/3-1, msg.Height)
		contentWidth := (msg.Width / 3) * 2
		style.Width(contentWidth).Height(msg.Height)

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(contentWidth, msg.Height)
			//m.viewport.YPosition = headerHeight
			//m.viewport.HighPerformanceRendering = true
			m.viewport.SetContent(m.issues.Issues[m.list.Cursor()].Fields.Description)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			//m.viewport.YPosition = headerHeight + 1

		} else {
			m.viewport.Width = contentWidth
			m.viewport.Height = msg.Height
		}
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

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return fmt.Sprintf("\ninitializing %s", m.spinner.View())
	}
	if m.typing {
		return fmt.Sprintf("Enter location:\n%s", m.textInput.View())
	}

	if m.loading {
		return fmt.Sprintf("%s fetching issues... please wait.", m.spinner.View())
	}

	if err := m.err; err != nil {
		return fmt.Sprintf("Could not fetch issues: %v", err)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.viewport.View())
}
