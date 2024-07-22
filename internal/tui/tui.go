package tui

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Have to make types for bubbletea to understand when it comes to the update function
type sessionState uint
type errMsg error

// list types
type item string

const (
	inventoryHeight              = 10
	messageView     sessionState = iota
	listView
)

var (
	// Inventory
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)

	// unfocused model
	unfocusedModelStyle = lipgloss.NewStyle().
				Width(30). // Be sure to change the values for the textarea and viewport in the message model if you change these
				Height(10).
				Align(lipgloss.Left, lipgloss.Bottom). // Sets alignment of content within the model
				BorderStyle(lipgloss.HiddenBorder())

		// Focused model
	focusedModelStyle = lipgloss.NewStyle().
				Width(30).
				Height(10).
				Align(lipgloss.Left, lipgloss.Bottom).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type model struct {
	state sessionState

	// List model attributes
	inventory list.Model
	choice    string

	// Message model attributes
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func newModel() model {
	// Messages
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 15                                // Needs editing
	ta.SetWidth(30)                                  // Same as {model}Style width
	ta.SetHeight(1)                                  // Cause I want just one line for users to enter messsages (I believe this adds 1 BELOW the viewport, making the message model have more height... see viewMessage)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle() // Remove cursor line styling
	ta.ShowLineNumbers = false
	vp := viewport.New(30, 10)
	ta.KeyMap.InsertNewline.SetEnabled(false)

	// Inventory
	items := []list.Item{
		item("Ramen"),
		item("Tomato Soup"),
		item("Hamburgers"),
		item("Cheeseburgers"),
		item("Currywurst"),
	}
	const defaultWidth = 20
	l := list.New(items, itemDelegate{}, defaultWidth, inventoryHeight)
	l.Title = "Inventory"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return model{
		state:       messageView,
		inventory:   l,
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

// Init commands
func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		// If the given command is a key
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == messageView {
				m.state = listView
			} else {
				m.state = messageView
			}
		case "enter":
			if m.state == listView {
				i, ok := m.inventory.SelectedItem().(item)
				if ok {
					m.choice = string(i)
				}
			} else {
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.textarea.Reset()
				m.viewport.GotoBottom()
			}
		}

		// Update whichever model is focused
		switch m.state {
		case listView:
			m.inventory, cmd = m.inventory.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.textarea, cmd = m.textarea.Update(msg)
			cmds = append(cmds, cmd)
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

// For Viewing two things at once when entering the View() function
func (m model) viewMessage() string {
	return fmt.Sprintf(
		"%s\n%s",
		m.viewport.View(),
		m.textarea.View(), // This is what causes the height increase I believe...
	)
}

func (m model) View() string {
	var s string
	if m.state == listView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(m.inventory.View()), unfocusedModelStyle.Render(fmt.Sprintf("%4s", m.viewMessage())))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, unfocusedModelStyle.Render(m.inventory.View()), focusedModelStyle.Render(fmt.Sprintf("%4s", m.viewMessage()))) // viewMessage is needed so we can work with this line here nicely
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • q: exit\n"))
	return s
}

func Start() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	// p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
