package tui

import (
	"explore/internal/commander"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Custom types
type errMsg error
type item string
type inventoryUpdateMsg struct{} // Specific type for our inventory channel signal
type commanderResponseMsg struct{}

var (
	GameCommander *commander.Commander

    viewPaneWidth = 60

	// Inventory
	titleStyle      = lipgloss.NewStyle().MarginLeft(2)
	itemStyle       = lipgloss.NewStyle().PaddingLeft(4)
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	inventoryStyle  = lipgloss.NewStyle().
			Width(20). // Be sure to change the values for the textarea and viewport in the message model if you change these
			Height(23).
			Align(lipgloss.Left, lipgloss.Top). // Sets alignment of content within the model
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#ffffff"))

		// Focused model
	viewportStyle = lipgloss.NewStyle().
			Width(viewPaneWidth).
			Height(19).
			Align(lipgloss.Left, lipgloss.Top).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#ffffff"))

	textareaStyle = lipgloss.NewStyle().
			Width(60).
			Height(1).
			Align(lipgloss.Left, lipgloss.Top).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type model struct {
	// List model attributes
	inventory        list.Model
	inventoryChanges <-chan struct{} // Channel for inventory changes - Same type as commanders - can only recieve...

	// Message model attributes
	viewport        viewport.Model
	messages        []string
	textarea        textarea.Model
	senderStyle     lipgloss.Style
	err             error
	messageResponse <-chan struct{}
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
	fmt.Fprint(w, fn(str))
}

func newModel(invChanges <-chan struct{}, commanderResponse <-chan struct{}) model {
	// Messages
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 15                                // Needs editing
	ta.SetWidth(60)                                  // Same as {model}Style width
	ta.SetHeight(1)                                  // Cause I want just one line for users to enter messsages (I believe this adds 1 BELOW the viewport, making the message model have more height... see viewMessage)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle() // Remove cursor line styling
	ta.ShowLineNumbers = false
	vp := viewport.New(viewPaneWidth, 20)
	ta.KeyMap.InsertNewline.SetEnabled(false)

	// Inventory
	items := []list.Item{}
	for _, baseItem := range GameCommander.GetCurrPlayerInv() {
		items = append(items, item(baseItem))
	}

	const defaultWidth = 20
	const defaultHeight = 21 // Cause the textarea adds 1 height as well
	l := list.New(items, itemDelegate{}, defaultWidth, defaultHeight)
	l.Title = "Inventory"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle

	return model{
		inventory:        l,
		inventoryChanges: invChanges,
		textarea:         ta,
		messages:         []string{},
		viewport:         vp,
		senderStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:              nil,
		messageResponse:  commanderResponse,
	}
}

// Init commands
func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink)
}

// This method returns just returns our custom signal type for our update function to pick it up
func InventoryUpdateCmd() tea.Cmd {
	return func() tea.Msg {
		return inventoryUpdateMsg{}
	}
}

func CommanderResponseCmd() tea.Cmd {
	return func() tea.Msg {
		return commanderResponseMsg{}
	}
}

func appendMsg(m *model, msg string) {
	var result []string
	var currentChunk strings.Builder

	words := strings.Fields(msg)
	for _, word := range words {
		// Check if adding this word would exceed the maxSize
		if currentChunk.Len()+len(word)+1 > viewPaneWidth {
			// If the current chunk is not empty, add it to the result
			if currentChunk.Len() > 0 {
				result = append(result, currentChunk.String())
				currentChunk.Reset()
			}
		}

		// Add the word to the current chunk
		if currentChunk.Len() > 0 {
			currentChunk.WriteString(" ")
		}
		currentChunk.WriteString(word)
	}

	// Add the last chunk if it's not empty
	if currentChunk.Len() > 0 {
		result = append(result, currentChunk.String())
	}

    for _, message := range result {
        m.messages = append(m.messages, message)
    }

    m.viewport.SetContent(strings.Join(m.messages, "\n"))
    m.viewport.GotoBottom()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	select {
	case <-m.messageResponse:
		cmds = append(cmds, CommanderResponseCmd())
	default:
		// No inventory change
	}

	switch msg := msg.(type) {

    case commanderResponseMsg:
        appendMsg(&m, m.senderStyle.Render("God: ")+GameCommander.Response)

	case inventoryUpdateMsg: // Inventory change notification was made
		items := []list.Item{}
		for _, newItem := range GameCommander.GetCurrPlayerInv() {
			items = append(items, item(newItem))
		}
		// Update inventory list
		m.inventory = list.New(items, itemDelegate{}, 20, 21)
		m.inventory.Title = "Inventory"
		m.inventory.SetShowStatusBar(false)
		m.inventory.SetFilteringEnabled(false)
		m.inventory.SetShowHelp(false)
		m.inventory.Styles.Title = titleStyle
		m.inventory.Styles.PaginationStyle = paginationStyle

	case tea.KeyMsg:
		// If the given command is a key
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			appendMsg(&m, m.senderStyle.Render(GameCommander.GetCurrPlayerName()+": ")+m.textarea.Value())
			GameCommander.PlayerCommand(m.textarea.Value())
            m.textarea.Reset()
		}
	}

	select {
	case <-m.inventoryChanges: // Waits for a signal to be in the buffer - if there is one it calls the command which gives us a inventoryUpdateMsg
		cmds = append(cmds, InventoryUpdateCmd())
	default:
		// No inventory change
	}

	// Update text area and viewport no matter what
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string
	s += lipgloss.JoinHorizontal(
		lipgloss.Top,
		inventoryStyle.Render(m.inventory.View()),
		lipgloss.JoinVertical(
			lipgloss.Top,
			viewportStyle.Render(m.viewport.View()),
			textareaStyle.Render(m.textarea.View())))

	s += helpStyle.Render(fmt.Sprintf("\nq: exit\n"))
	return s
}

func Start(cmder *commander.Commander) error {
	GameCommander = cmder
	p := tea.NewProgram(newModel(GameCommander.InventoryChangeChannel, GameCommander.ResponseChannel), tea.WithAltScreen()) // Send the inv channel to the model to be monitored
	// p := tea.NewProgram(newModel(GameCommander.InventoryChangeChannel, GameCommander.ResponseChannel))

	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
