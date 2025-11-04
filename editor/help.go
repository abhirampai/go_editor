package editor

type KeyBinding struct {
	Key         string
	Mode        string
	Description string
}

func GetKeyBindings() []KeyBinding {
	return []KeyBinding{
		{Key: "ESC", Mode: "Any", Description: "Exit insert mode / Close popover"},
		{Key: "i", Mode: "Visual", Description: "Enter insert mode"},
		{Key: "w", Mode: "Visual", Description: "Save file"},
		{Key: "q", Mode: "Visual", Description: "Quit editor"},
		{Key: "h", Mode: "Visual", Description: "Show this help"},
		{Key: "c", Mode: "Visual", Description: "Copy current line"},
		{Key: "p", Mode: "Visual", Description: "Paste copied line below"},
		{Key: "d", Mode: "Visual", Description: "Delete current line"},
		{Key: "s", Mode: "Visual", Description: "Save current buffer for undo"},
		{Key: "l", Mode: "Visual", Description: "Load saved buffer (undo)"},
		{Key: "Left", Mode: "Any", Description: "Move cursor left"},
		{Key: "Right", Mode: "Any", Description: "Move cursor right"},
		{Key: "Up", Mode: "Any", Description: "Move cursor up"},
		{Key: "Down", Mode: "Any", Description: "Move cursor down"},
		{Key: "Home", Mode: "Any", Description: "Move to start of line"},
		{Key: "End", Mode: "Any", Description: "Move to end of line"},
		{Key: "PgUp", Mode: "Any", Description: "Move up by quarter page"},
		{Key: "PgDn", Mode: "Any", Description: "Move down by quarter page"},
		{Key: "Enter", Mode: "Insert", Description: "Insert new line"},
		{Key: "Backspace", Mode: "Insert", Description: "Delete character"},
		{Key: "Tab", Mode: "Insert", Description: "Insert tab"},
	}
}

func FormatKeyBindingsHelp() []string {
	bindings := GetKeyBindings()
	var lines []string

	// Header with proper centering
	lines = append(lines, "    --- Welcome to GoEditor Help ---")
	lines = append(lines, "Editor Key Bindings:")
	lines = append(lines, "")

	// Calculate column widths
	keyWidth := 12  // Increased to accommodate longer key names
	modeWidth := 10 // Increased for better spacing

	// Format each binding with proper column alignment
	for _, binding := range bindings {
		line := "    "

		keyStr := binding.Key
		for i := len(keyStr); i < keyWidth-2; i++ {
			line += " "
		}
		line += keyStr
		line += "  " // Consistent spacing after key

		line += binding.Mode
		for i := len(binding.Mode); i < modeWidth; i++ {
			line += " "
		}

		line += "  " + binding.Description
		lines = append(lines, line)
	}

	return lines
}
