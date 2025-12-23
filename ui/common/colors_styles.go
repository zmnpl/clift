package common

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var HEADER_FOOTER_HEIGHT = 4
var V_MARGINS = 2

func GetContentHeight(screenHeight int) int {
	return screenHeight - HEADER_FOOTER_HEIGHT - V_MARGINS
}

var WINDOW_HEIGHT = 10
var WINDOW_WIDTH = 10

var FilterCursorStyle = lipgloss.NewStyle().Background(Theme.Yellow)
var FilterPromptStyle = lipgloss.NewStyle().Foreground(Theme.Yellow)

func ListItemStyle() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	//d.Styles.SelectedTitle.Inherit(lipgloss.NewStyle().Foreground(theme.red).BorderForeground(lipgloss.Color(theme.red)))
	//d.Styles.SelectedDesc.Inherit(lipgloss.NewStyle().Foreground(theme.yellow).BorderForeground(lipgloss.Color(theme.red)))

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(Theme.Red).BorderForeground(lipgloss.Color(Theme.Red))
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(Theme.Yellow).BorderForeground(lipgloss.Color(Theme.Red))

	return d
}

// kitty themes here
// https://github.com/kovidgoyal/kitty-themes/blob/master/themes

type Colorscheme struct {
	Foreground           lipgloss.Color
	Background           lipgloss.Color
	Selection_foreground lipgloss.Color
	Selection_background lipgloss.Color
	Cursor               lipgloss.Color
	Cursor_text_color    lipgloss.Color
	Url_color            lipgloss.Color

	Active_border_color   lipgloss.Color
	Inactive_border_color lipgloss.Color
	Bell_border_color     lipgloss.Color
	Visual_bell_color     lipgloss.Color

	Black   lipgloss.Color
	Red     lipgloss.Color
	Green   lipgloss.Color
	Yellow  lipgloss.Color
	Blue    lipgloss.Color
	Magenta lipgloss.Color
	Cyan    lipgloss.Color
	White   lipgloss.Color

	Bright_black   lipgloss.Color
	Bright_red     lipgloss.Color
	Bright_green   lipgloss.Color
	Bright_yellow  lipgloss.Color
	Bright_blue    lipgloss.Color
	Bright_magenta lipgloss.Color
	Bright_cyan    lipgloss.Color
	Bright_white   lipgloss.Color
}

var Hachikoo = Colorscheme{
	Black:   lipgloss.Color("#181818"),
	Red:     lipgloss.Color("#960042"),
	Green:   lipgloss.Color("#FF0000"),
	Yellow:  lipgloss.Color("#FF5D05"),
	Blue:    lipgloss.Color("#FF2044"),
	Magenta: lipgloss.Color("#FFEDCF"),
	Cyan:    lipgloss.Color("#6F0027"),
	White:   lipgloss.Color("#FFDAF1"),

	Bright_black:   lipgloss.Color("#333333"),
	Bright_red:     lipgloss.Color("#870300"),
	Bright_green:   lipgloss.Color("#690000"),
	Bright_yellow:  lipgloss.Color("#6F2700"),
	Bright_blue:    lipgloss.Color("#333333"),
	Bright_magenta: lipgloss.Color("#FFFFB5"),
	Bright_cyan:    lipgloss.Color("#F50056"),
	Bright_white:   lipgloss.Color("#FFE6DA"),
}

var Blackmetal = Colorscheme{
	Black:   lipgloss.Color("#000000"),
	Red:     lipgloss.Color("#5f8787"),
	Green:   lipgloss.Color("#dd9999"),
	Yellow:  lipgloss.Color("#a06666"),
	Blue:    lipgloss.Color("#888888"),
	Magenta: lipgloss.Color("#999999"),
	Cyan:    lipgloss.Color("#aaaaaa"),
	White:   lipgloss.Color("#c1c1c1"),

	Bright_black:   lipgloss.Color("#333333"),
	Bright_red:     lipgloss.Color("#5f8787"),
	Bright_green:   lipgloss.Color("#dd9999"),
	Bright_yellow:  lipgloss.Color("#a06666"),
	Bright_blue:    lipgloss.Color("#888888"),
	Bright_magenta: lipgloss.Color("#999999"),
	Bright_cyan:    lipgloss.Color("#aaaaaa"),
	Bright_white:   lipgloss.Color("#c1c1c1"),
}

var Terafox = Colorscheme{
	Foreground: lipgloss.Color("#e6eaea"),

	Black:   lipgloss.Color("#2f3239"),
	Red:     lipgloss.Color("#e85c51"),
	Green:   lipgloss.Color("#7aa4a1"),
	Yellow:  lipgloss.Color("#fda47f"),
	Blue:    lipgloss.Color("#5a93aa"),
	Magenta: lipgloss.Color("#ad5c7c"),
	Cyan:    lipgloss.Color("#a1cdd8"),
	White:   lipgloss.Color("#ebebeb"),

	Bright_black:   lipgloss.Color("#4e5157"),
	Bright_red:     lipgloss.Color("#eb746b"),
	Bright_green:   lipgloss.Color("#8eb2af"),
	Bright_yellow:  lipgloss.Color("#fdb292"),
	Bright_blue:    lipgloss.Color("#73a3b7"),
	Bright_magenta: lipgloss.Color("#b97490"),
	Bright_cyan:    lipgloss.Color("#afd4de"),
	Bright_white:   lipgloss.Color("#eeeeee"),
}

var Theme = Terafox

var (
	StatusGood   = lipgloss.NewStyle().Foreground(Theme.Cyan).Width(7).Align(lipgloss.Left)
	StatusBad    = lipgloss.NewStyle().Foreground(Theme.Magenta).Width(7).Align(lipgloss.Left)
	StatusCenter = lipgloss.NewStyle().Foreground(Theme.Cyan).Width(70).Align(lipgloss.Left)

	PrimaryExStlye = lipgloss.NewStyle().Foreground(Theme.Cyan)
	FocusedStyle   = lipgloss.NewStyle().Foreground(Theme.Red)
	BlurredStyle   = lipgloss.NewStyle().Foreground(Theme.Bright_yellow)
	HeaderStyle    = lipgloss.NewStyle().Foreground(Theme.Red).Width(80)
	NoStyle        = lipgloss.NewStyle()
	Margin         = lipgloss.NewStyle().Margin(1, 1, 1, 1)
)
