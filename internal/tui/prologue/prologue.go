package prologue

import (
	"fmt"
	"log"
    "strings"

	"github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Have to make types for bubbletea to understand when it comes to the update function
type sessionState uint
type errMsg error

const (
	messageView   sessionState = iota
	spinnerView
)

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}

    // unfocused model
	modelStyle = lipgloss.NewStyle().
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

	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type model struct {
	state   sessionState

    // Spinner model attributes
	spinner spinner.Model
	index   int

    // Message model attributes
    viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error

}

func newModel() model {
    ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 15 // Needs editing

	ta.SetWidth(30) // Same as {model}Style width
	ta.SetHeight(1) // Cause I want just one line for users to enter messsages (I believe this adds 1 BELOW the viewport, making the message model have more height... see viewMessage)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 10)

	ta.KeyMap.InsertNewline.SetEnabled(false)

    return model {
            state: messageView,
            spinner: spinner.New(),
            textarea:    ta,
            messages:    []string{},
            viewport:    vp,
            senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
            err:         nil,
        }
}

func (m model) Init() tea.Cmd {
	// start the timer and spinner on program start
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == messageView {
				m.state = spinnerView
			} else {
				m.state = messageView
			}
		case "n":
			if m.state == spinnerView {
				m.Next()
				m.resetSpinner()
				cmds = append(cmds, m.spinner.Tick)
			}
        case "enter":
            m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

		switch m.state {
		// update whichever model is focused
		case spinnerView:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.textarea, cmd = m.textarea.Update(msg)
			cmds = append(cmds, cmd)
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
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
	model := m.currentFocusedModel()
	if m.state == messageView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(fmt.Sprintf("%4s", m.viewMessage())), modelStyle.Render(m.spinner.View())) // viewMessage is needed so we can work with this line here nicely
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", m.viewMessage())), focusedModelStyle.Render(m.spinner.View()))
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	return s
}

func (m model) currentFocusedModel() string {
	if m.state == messageView {
		return "message"
	}
	return "spinner"
}

// Spinner BS
func (m *model) Next() {
	if m.index == len(spinners)-1 {
		m.index = 0
	} else {
		m.index++
	}
}

func (m *model) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[m.index]
}
// ---

func Start() {
	// p := tea.NewProgram(newModel(), tea.WithAltScreen())
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
