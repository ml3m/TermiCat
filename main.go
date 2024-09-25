package main

import (
	"os"
    "fmt"
    "time"
	"strings"
    "golang.org/x/term"
	"TermiCat/asciiart" // my ascii 
    "github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
)

func splitLines(content string) []string {return strings.Split(content, "\n")}

type model struct {
    currentFrame   int       // Track the current animation frame
    frames         []string  // Store the ASCII frames
    focusedButton  int       // Track the index of the focused button
    buttonLabels   []string  // Store button labels
    actionMessage  string    // Message to display when a button is pressed
}

func CenterEngine(content string) {
	// Get terminal size
	width, height, _ := term.GetSize(int(os.Stdout.Fd())) // _ is err ignored for now 

	// Split the content into lines
	lines := splitLines(content)
	contentHeight := len(lines)

	// Find the maximum width of the content
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Calculate center position
	x := (width - maxWidth) / 2
	y := (height - contentHeight) / 2

	// Print each line at the center position
	for i, line := range lines {
		fmt.Printf("\033[%d;%dH%s", y+i+1, x+1, line) // Add 1 because terminal positions are 1-based
	}
	fmt.Println("\033[0m") // Reset terminal color
}

// Define styles using Lipgloss
var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).MarginBottom(1)
	coinsStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Align(lipgloss.Center)
	levelStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Align(lipgloss.Center)
	catStyle          = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Align(lipgloss.Center).Height(8).Width(18)
	actionsStyle      = lipgloss.NewStyle().Align(lipgloss.Center).MarginTop(1)
	buttonStyle       = lipgloss.NewStyle().Padding(0, 2).Border(lipgloss.NormalBorder()).Align(lipgloss.Center)
	focusedButtonStyle = lipgloss.NewStyle().Padding(0, 2).Border(lipgloss.DoubleBorder()).Align(lipgloss.Center)
	centerStyle       = lipgloss.NewStyle().Align(lipgloss.Center).MarginBottom(2)
	actionMessageStyle = lipgloss.NewStyle().Align(lipgloss.Center).MarginTop(1)
)


func (m model) Init() tea.Cmd {
	// Start the animation loop with a delay
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit (q or ctrl+c)
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// Handle button focus navigation
		switch msg.String() {
		case "left":
			m.focusedButton = (m.focusedButton - 1 + len(m.buttonLabels)) % len(m.buttonLabels) // Up
		case "right":
			m.focusedButton = (m.focusedButton + 1) % len(m.buttonLabels) // Down
		case "enter":
			// Set the action message based on the focused button
			m.actionMessage = fmt.Sprintf("You pressed: %s", m.buttonLabels[m.focusedButton])
            // handling for each button goes here, probably a switch case :p
            /*
            // idk yet this part 
            switch m.buttonLabels[m.focusedButton]*/
		}
	case time.Time:
		// Handle animation frame updates
		m.currentFrame = (m.currentFrame + 1) % len(m.frames)
		// Continue the animation loop
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return t
		})
	}

	return m, nil
}

func (m model) View() string {
	// layout using Lipgloss
	header := titleStyle.Render("Termicat")
	coins := coinsStyle.Render("Coins: 100")
	level := levelStyle.Render("Lv: 5")

	// Get the current frame of the cat animation
	cat := catStyle.Render(m.frames[m.currentFrame])

	// Define the buttons and highlight the focused one
	buttons := make([]string, len(m.buttonLabels))
	for i, label := range m.buttonLabels {
		if i == m.focusedButton {
			buttons[i] = focusedButtonStyle.Render(label) // Focused button
		} else {
			buttons[i] = buttonStyle.Render(label) // Regular button
		}
	}

	actionText := actionsStyle.Render("More actions if needed")
	actionMessage := actionMessageStyle.Render(m.actionMessage)

	// Combine everything into the final layout
	layout := lipgloss.JoinVertical(lipgloss.Center,
		header,
		lipgloss.JoinHorizontal(lipgloss.Center, coins, level),
		cat,
		lipgloss.JoinHorizontal(lipgloss.Top, buttons...), // Join buttons
		actionText,
		actionMessage, // Render the action message below
	)

	return layout
}

// Move cursor to the home position and clear the screen
func clearTerminal() {fmt.Print("\033[H\033[2J") }

/********************************************************************************/
//                                    main 
/********************************************************************************/

func main() {
    clearTerminal()

    // simple cat blinking animation
	frames := asciiart.GetFrames()
	m := model{
		frames:       frames,
		buttonLabels: []string{"Feed", "Games", "Clean", "Sleep"}, // Define button labels
	}

	// Create new program object and run
	p := tea.NewProgram(m)
	p.Run() 
}
