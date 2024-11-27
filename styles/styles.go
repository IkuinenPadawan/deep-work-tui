package styles

import "github.com/charmbracelet/lipgloss"

var (
	BlockStyle        = lipgloss.NewStyle().PaddingLeft(2)
	TimeStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	TaskStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("241"))
	FreeTime          = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	SelectedTimeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
)

func GetSeamlessBlockStyle(isFirst, isLast bool) lipgloss.Style {
	var border lipgloss.Border
	if isFirst {
		border = lipgloss.Border{
			Top:         "─",
			Bottom:      " ",
			Left:        "│",
			Right:       "│",
			TopLeft:     " ",
			TopRight:    " ",
			BottomLeft:  "│",
			BottomRight: "│",
		}
	} else if isLast {
		border = lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "│",
			TopRight:    "│",
			BottomLeft:  " ",
			BottomRight: " ",
		}
	} else {
		border = lipgloss.Border{
			Top:         "─",
			Bottom:      " ",
			Left:        "│",
			Right:       "│",
			TopLeft:     "│",
			TopRight:    "│",
			BottomLeft:  "│",
			BottomRight: "│",
		}
	}
	baseStyle := lipgloss.NewStyle().
		BorderStyle(border).
		Padding(0, 1)

	return baseStyle.BorderForeground(lipgloss.Color("#AAAAAA"))
}
