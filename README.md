# Go-Editor

A lightweight terminal-based text editor written in Go using termbox2 for terminal manipulation. Inspired by Vim's modal editing, it provides a simple yet efficient text editing experience.

## Features

### Core Features
- Terminal-based user interface with modal editing (Visual and Insert modes)
- Built-in file browser for navigating and opening files
- Syntax highlighting with colored text display
- Full terminal window utilization with automatic resize handling
- Configuration system with JSON-based settings

### Editor Interface
- Comprehensive status bar showing:
  - Current mode (VISUAL/INSERT)
  - File information (name, line count, modified status)
  - Copy and undo buffer indicators
  - Cursor position (line and column)
  - Tab size setting
- Built-in help system (press 'h' to view)
- Clear visual indicators for tabs and special characters

### Text Manipulation
- File editing with Visual and Insert modes
- Full cursor movement support (arrows, Home/End, PgUp/PgDn)
- Tab character support with configurable width
- Basic text operations (insert, delete, newline)
- File operations (open, save)

### File Management
- Interactive file browser for navigating directories and opening files
- Open files from command line
- Create new files
- Save files with write protection
- Modified file indicator

## Prerequisites

- Go installed on your system
- GCC compiler for building termbox2
- Make utility

## Installing termbox2
Install termbox2 in the editor folder:

```bash
cd editor
git clone https://github.com/termbox/termbox2.git
cd termbox2
make
sudo make install
```

If you already have a local checkout that includes `editor/termbox2`, you can build that copy instead by running `make` inside `editor/termbox2`.

## Build

```bash
go build -o go_editor main.go
```

## Run

```bash
./go_editor
```

## Installation and Setup

### Prerequisites
- Go installed on your system
- GCC compiler for building termbox2
- Make utility

### Building from Source

1. **Clone the repository:**
   ```bash
   git clone https://github.com/abhirampai/go_editor.git
   cd go_editor
   ```

2. **Install termbox2:**
   ```bash
   cd editor
   git clone https://github.com/termbox/termbox2.git
   cd termbox2
   make
   sudo make install
   cd ../..
   ```

3. **Build the editor:**
   ```bash
   go build -o go_editor main.go
   ```

4. **Run the editor:**
   ```bash
   # Create new file
   ./go_editor
   
   # Open existing file
   ./go_editor path/to/file
   ```

## Configuration

The editor uses a configuration file located at `~/.gocodeeditor/settings.json`. If the file doesn't exist, it will be created automatically with default settings.

### Configuration File Location
- Unix/Linux/macOS: `~/.gocodeeditor/settings.json`

### Default Settings

```json
{
  "tab_size": 4
}
```

### Available Settings

| Setting    | Description                                    | Default |
|------------|------------------------------------------------|---------|
| `tab_size` | Number of spaces to display for a tab character | 4       |

### Modifying Settings
1. Open the settings file: `~/.gocodeeditor/settings.json`
2. Modify the values as needed
3. Save the file
4. Restart the editor for changes to take effect

## Keyboard Shortcuts

### Mode Control
- `ESC`: Exit insert mode / Close popover
- `i`: Enter insert mode (from visual mode)

### File Operations
- `w`: Save current file
- `q`: Quit editor

### Navigation (Any Mode)
- `←` or `Left Arrow`: Move cursor left
- `→` or `Right Arrow`: Move cursor right
- `↑` or `Up Arrow`: Move cursor up
- `↓` or `Down Arrow`: Move cursor down
- `Home`: Move to start of line
- `End`: Move to end of line
- `PgUp`: Move up by quarter page
- `PgDn`: Move down by quarter page

### Text Manipulation (Insert Mode)
- `Enter`: Insert new line
- `Backspace`: Delete character
- `Tab`: Insert tab

### Text Manipulation (Visual Mode)
- `c`: Copy current line
- `p`: Paste copied line
- `d`: Delete current line
- `s`: Save current buffer for undo
- `l`: Load saved buffer (undo)

### File Browser
- `o`: Open file browser modal
  - `↑/↓`: Navigate through files and directories
  - `Enter`: Open selected file or enter directory
  - `ESC`: Close file browser

### Help and Information
- `h`: Show help popover with key bindings

## Status Bar Information

The status bar at the bottom of the editor provides important information:

### Left Side
- Editor Mode: Shows "-- VISUAL --" or "-- INSERT --"
- File Status: Shows filename, line count, and modified/saved status
- Buffer Indicators: Shows [Copy] and [Undo] when content is available

### Right Side
- Cursor Position: Shows current line and column numbers
- Tab Size: Shows current tab size setting

## Contributing

Feel free to open issues or submit pull requests for improvements and bug fixes.

## License

This project is open source. Please check the repository for license details.