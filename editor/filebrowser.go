package editor

import (
	"os"
	"path/filepath"
	"sort"
)

type FileEntry struct {
	Name     string
	IsDir    bool
	Path     string
	Selected bool
}

type FileBrowser struct {
	CurrentPath string
	Entries     []FileEntry
	Cursor      int
	Scroll      int
}

func NewFileBrowser() *FileBrowser {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}
	fb := &FileBrowser{
		CurrentPath: pwd,
		Cursor:      0,
		Scroll:      0,
	}
	fb.RefreshEntries()
	return fb
}

func (fb *FileBrowser) RefreshEntries() error {
	entries, err := os.ReadDir(fb.CurrentPath)
	if err != nil {
		return err
	}

	fb.Entries = make([]FileEntry, 0)

	if fb.CurrentPath != "/" {
		fb.Entries = append(fb.Entries, FileEntry{
			Name:  "..",
			IsDir: true,
			Path:  filepath.Join(fb.CurrentPath, ".."),
		})
	}

	for _, entry := range entries {
		fullPath := filepath.Join(fb.CurrentPath, entry.Name())
		fb.Entries = append(fb.Entries, FileEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Path:  fullPath,
		})
	}

	sort.Slice(fb.Entries, func(i, j int) bool {
		if fb.Entries[i].IsDir == fb.Entries[j].IsDir {
			return fb.Entries[i].Name < fb.Entries[j].Name
		}
		return fb.Entries[i].IsDir
	})

	return nil
}

func (fb *FileBrowser) MoveUp() {
	if fb.Cursor > 0 {
		fb.Cursor--
		if fb.Cursor < fb.Scroll {
			fb.Scroll = fb.Cursor
		}
	}
}

func (fb *FileBrowser) MoveDown() {
	if fb.Cursor < len(fb.Entries)-1 {
		fb.Cursor++
	}
}

func (fb *FileBrowser) Enter() (string, bool, error) {
	if fb.Cursor >= len(fb.Entries) {
		return "", false, nil
	}

	selected := fb.Entries[fb.Cursor]
	if selected.IsDir {
		fb.CurrentPath = selected.Path
		err := fb.RefreshEntries()
		fb.Cursor = 0
		fb.Scroll = 0
		return "", true, err
	}

	return selected.Path, false, nil
}
