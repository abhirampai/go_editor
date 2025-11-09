package editor

/*
#include "termbox2.h"
*/
import "C"

type Theme struct {
	Name         string
	Background   C.uintattr_t
	Foreground   C.uintattr_t
	LineNumber   C.uintattr_t
	StatusBarBg  C.uintattr_t
	StatusBarFg  C.uintattr_t
	StatusModeBg C.uintattr_t
	StatusModeFg C.uintattr_t
	StatusInfoBg C.uintattr_t
	StatusInfoFg C.uintattr_t
	SelectionBg  C.uintattr_t
	PopupBg      C.uintattr_t
	PopupFg      C.uintattr_t
	PopupTitleBg C.uintattr_t
	PopupTitleFg C.uintattr_t
}

const (
	ColorDefault = C.TB_DEFAULT
	ColorBlack   = C.TB_BLACK
	ColorRed     = C.TB_RED
	ColorGreen   = C.TB_GREEN
	ColorYellow  = C.TB_YELLOW
	ColorBlue    = C.TB_BLUE
	ColorMagenta = C.TB_MAGENTA
	ColorCyan    = C.TB_CYAN
	ColorWhite   = C.TB_WHITE
)

var OneDarkTheme = Theme{
	Name:         "One Dark",
	Background:   ColorDefault,
	Foreground:   ColorWhite,
	LineNumber:   ColorBlue,
	StatusBarBg:  ColorBlack,
	StatusBarFg:  ColorWhite,
	StatusModeBg: ColorGreen,
	StatusModeFg: ColorBlack,
	StatusInfoBg: ColorBlack,
	StatusInfoFg: ColorWhite,
	SelectionBg:  ColorBlue,
	PopupBg:      ColorBlack,
	PopupFg:      ColorWhite,
	PopupTitleBg: ColorBlack,
	PopupTitleFg: ColorYellow,
}

var CurrentTheme = OneDarkTheme

var SolarizedDark = Theme{
	Name:         "Solarized Dark",
	Background:   ColorBlack,
	Foreground:   ColorWhite,
	LineNumber:   ColorCyan,
	StatusBarBg:  ColorBlue,
	StatusBarFg:  ColorWhite,
	StatusModeBg: ColorYellow,
	StatusModeFg: ColorBlack,
	StatusInfoBg: ColorBlack,
	StatusInfoFg: ColorWhite,
	SelectionBg:  ColorGreen,
	PopupBg:      ColorBlack,
	PopupFg:      ColorWhite,
	PopupTitleBg: ColorBlue,
	PopupTitleFg: ColorYellow,
}

var DraculaTheme = Theme{
	Name:         "Dracula",
	Background:   ColorBlack,
	Foreground:   ColorWhite,
	LineNumber:   ColorMagenta,
	StatusBarBg:  ColorMagenta,
	StatusBarFg:  ColorWhite,
	StatusModeBg: ColorMagenta,
	StatusModeFg: ColorBlack,
	StatusInfoBg: ColorBlack,
	StatusInfoFg: ColorWhite,
	SelectionBg:  ColorBlue,
	PopupBg:      ColorBlack,
	PopupFg:      ColorWhite,
	PopupTitleBg: ColorMagenta,
	PopupTitleFg: ColorYellow,
}

var GruvboxDark = Theme{
	Name:         "Gruvbox Dark",
	Background:   ColorBlack,
	Foreground:   ColorYellow,
	LineNumber:   ColorRed,
	StatusBarBg:  ColorBlack,
	StatusBarFg:  ColorYellow,
	StatusModeBg: ColorGreen,
	StatusModeFg: ColorBlack,
	StatusInfoBg: ColorBlack,
	StatusInfoFg: ColorYellow,
	SelectionBg:  ColorBlue,
	PopupBg:      ColorBlack,
	PopupFg:      ColorYellow,
	PopupTitleBg: ColorBlack,
	PopupTitleFg: ColorRed,
}

var MonokaiPro = Theme{
	Name:         "Monokai Pro",
	Background:   ColorBlack,
	Foreground:   ColorWhite,
	LineNumber:   ColorMagenta,
	StatusBarBg:  ColorBlack,
	StatusBarFg:  ColorGreen,
	StatusModeBg: ColorGreen,
	StatusModeFg: ColorBlack,
	StatusInfoBg: ColorBlack,
	StatusInfoFg: ColorWhite,
	SelectionBg:  ColorMagenta,
	PopupBg:      ColorBlack,
	PopupFg:      ColorWhite,
	PopupTitleBg: ColorBlack,
	PopupTitleFg: ColorGreen,
}

var NordDark = Theme{
	Name:         "Nord Dark",
	Background:   ColorBlack,
	Foreground:   ColorCyan,
	LineNumber:   ColorBlue,
	StatusBarBg:  ColorBlue,
	StatusBarFg:  ColorWhite,
	StatusModeBg: ColorCyan,
	StatusModeFg: ColorBlack,
	StatusInfoBg: ColorBlack,
	StatusInfoFg: ColorCyan,
	SelectionBg:  ColorBlue,
	PopupBg:      ColorBlack,
	PopupFg:      ColorCyan,
	PopupTitleBg: ColorBlue,
	PopupTitleFg: ColorWhite,
}

var TokyoNight = Theme{
	Name:         "Tokyo Night",
	Background:   ColorBlack,
	Foreground:   ColorWhite,
	LineNumber:   ColorMagenta,
	StatusBarBg:  ColorBlue,
	StatusBarFg:  ColorWhite,
	StatusModeBg: ColorMagenta,
	StatusModeFg: ColorBlack,
	StatusInfoBg: ColorBlack,
	StatusInfoFg: ColorWhite,
	SelectionBg:  ColorBlue,
	PopupBg:      ColorBlack,
	PopupFg:      ColorWhite,
	PopupTitleBg: ColorBlue,
	PopupTitleFg: ColorMagenta,
}

var MaterialDark = Theme{
	Name:         "Material Dark",
	Background:   ColorBlack,
	Foreground:   ColorWhite,
	LineNumber:   ColorBlue,
	StatusBarBg:  ColorBlue,
	StatusBarFg:  ColorWhite,
	StatusModeBg: ColorCyan,
	StatusModeFg: ColorBlack,
	StatusInfoBg: ColorBlack,
	StatusInfoFg: ColorWhite,
	SelectionBg:  ColorBlue,
	PopupBg:      ColorBlack,
	PopupFg:      ColorWhite,
	PopupTitleBg: ColorBlue,
	PopupTitleFg: ColorCyan,
}

var Themes = map[string]Theme{
	"one-dark":       OneDarkTheme,
	"solarized-dark": SolarizedDark,
	"dracula":        DraculaTheme,
	"gruvbox":        GruvboxDark,
	"monokai-pro":    MonokaiPro,
	"nord-dark":      NordDark,
	"tokyo-night":    TokyoNight,
	"material-dark":  MaterialDark,
}

func GetThemeNames() []string {
	names := make([]string, 0, len(Themes))
	for k := range Themes {
		if CurrentTheme.Name == Themes[k].Name {
			continue
		}
		names = append(names, k)
	}
	return names
}

func SetTheme(name string) bool {
	if t, ok := Themes[name]; ok {
		CurrentTheme = t
		return true
	}
	return false
}

func ApplySettingsTheme() {
	if editSettings == nil {
		return
	}
	if editSettings.Theme == "" {
		return
	}
	SetTheme(editSettings.Theme)
}

func SetThemeAndSave(name string) bool {
	if !SetTheme(name) {
		return false
	}
	if editSettings == nil {
		editSettings = &Settings{
			TabSize: 4,
			Theme:   name,
		}
	} else {
		editSettings.Theme = name
	}

	_ = SaveSettings(editSettings)
	return true
}

func GetCurrentThemeKey() string {
	for k, t := range Themes {
		if t.Name == CurrentTheme.Name {
			return k
		}
	}
	return ""
}
