# Go-Editor

A lightweight terminal-based text editor written in Go using termbox2 for terminal manipulation.

## Features

- Terminal-based user interface
- Basic text display with colored syntax highlighting
- Configurable tab size
- Support for displaying special characters and tabs
- Terminal resize handling
- Status bar: shows editor mode (INSERT / VISUAL), file status (filename + line count + saved/modified), copy/undo indicators, cursor position (Ln/Col), and current tab size.
- Mode display: the editor now shows whether it's in INSERT or VISUAL mode in the status bar.
- Copy/Undo indicators: the status bar will show when there's content in the copy buffer or undo buffer.
- Open file by argument: you can pass a filename as the first CLI argument to open that file on startup; otherwise an untitled buffer is created.

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

## Build

```bash
go build -o go_editor main.go
```

## Run

```bash
./go_editor
```

## Configuration

The editor uses a configuration file located at `~/.gocodeeditor/settings.json`. If the file doesn't exist, it will be created automatically with default settings.

### Default settings:

```json
{
  "tab_size": 4
}
```

### Configuration Options:

- `tab_size`: Number of spaces to display for a tab character (default: 4)

## Keyboard Shortcuts

- `ESC`: Exit the editor

## Contributing

Feel free to open issues or submit pull requests for improvements and bug fixes.

## License

This project is open source. Please check the repository for license details.