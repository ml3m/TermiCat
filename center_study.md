
// idk if this will ever work ...

func splitLines(content string) []string {return strings.Split(content, "\n")}


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
