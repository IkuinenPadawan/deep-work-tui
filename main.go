package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"sort"
	"strings"
	"time"
)

type Timeblock struct {
	task      string
	starttime time.Time
	endtime   time.Time
}

type model struct {
	timeblocks  []Timeblock
	cursor      int
	selected    map[int]struct{}
	inputFields []textinput.Model
	showInput   bool
	editing     bool
	editIndex   int
	focused     int
	err         error
	lastKey     string
}

var (
	blockStyle        = lipgloss.NewStyle().PaddingLeft(2)
	timeStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	taskStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("241"))
	freeTimeStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	selectedTimeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
)

func getSeamlessBlockStyle(isFirst, isLast bool) lipgloss.Style {
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

func parseTime(t string) time.Time {
	parsed, _ := time.Parse("15:04", t)
	return parsed
}

func isValidTime(t string) bool {
	_, err := time.Parse("15:04", t)
	return err == nil
}

func initialModel() model {
	taskNameInput := textinput.New()
	taskNameInput.Placeholder = "Task Name"
	taskNameInput.CharLimit = 50

	startTimeInput := textinput.New()
	startTimeInput.Placeholder = "Start Time (HH:mm)"
	startTimeInput.CharLimit = 5

	endTimeInput := textinput.New()
	endTimeInput.Placeholder = "End Time (HH:mm)"
	endTimeInput.CharLimit = 5
	return model{
		timeblocks: []Timeblock{
			{"Deep Work", parseTime("07:00"), parseTime("10:00")},
			{"Email", parseTime("10:00"), parseTime("10:30")},
			{"Other Work", parseTime("10:30"), parseTime("12:00")},
			{"Meeting", parseTime("12:00"), parseTime("14:00")},
			{"Deep Work", parseTime("14:00"), parseTime("16:00")},
		},
		selected:    make(map[int]struct{}),
		showInput:   false,
		editing:     false,
		inputFields: []textinput.Model{taskNameInput, startTimeInput, endTimeInput},
		focused:     0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) enterAddMode() {
	m.showInput = true
	m.focused = 0
	for i := range m.inputFields {
		if i == m.focused {
			m.inputFields[i].Focus()
		} else {
			m.inputFields[i].Blur()
		}
	}
}

func (m *model) saveAdd() {
	name := m.inputFields[0].Value()
	start := m.inputFields[1].Value()
	end := m.inputFields[2].Value()

	if name != "" && isValidTime(start) && isValidTime(end) {
		m.timeblocks = append(m.timeblocks, Timeblock{task: name, starttime: parseTime(start), endtime: parseTime(end)})
		m.showInput = false
		m.clearInputFields()
	}
}

func (m *model) enterEditMode() {
	if len(m.timeblocks) == 0 {
		return
	}

	m.editing = true
	m.editIndex = m.cursor

	timeblock := m.timeblocks[m.editIndex]
	m.inputFields[0].SetValue(timeblock.task)
	m.inputFields[1].SetValue(timeblock.starttime.Format("15:04"))
	m.inputFields[2].SetValue(timeblock.endtime.Format("15:04"))
}

func (m *model) saveEdit() {
	if !m.editing || m.editIndex < 0 || m.editIndex >= len(m.timeblocks) {
		return
	}

	m.timeblocks[m.editIndex].task = m.inputFields[0].Value()
	m.timeblocks[m.editIndex].starttime, _ = time.Parse("15:04", m.inputFields[1].Value())
	m.timeblocks[m.editIndex].endtime, _ = time.Parse("15:04", m.inputFields[2].Value())

	m.editing = false
	m.editIndex = -1

	m.clearInputFields()
}

func (m *model) cancelEdit() {
	m.editing = false
	m.editIndex = -1
	m.clearInputFields()
}

func (m *model) clearInputFields() {
	for i := range m.inputFields {
		m.inputFields[i].SetValue("")
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "d" && m.lastKey == "d" && !m.editing && !m.showInput {
			indexToRemove := m.cursor
			m.timeblocks = append(m.timeblocks[:indexToRemove], m.timeblocks[indexToRemove+1:]...)
			m.lastKey = ""
			return m, nil
		}

		// Store the current key as the last key
		m.lastKey = msg.String()
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.showInput {
				m.showInput = false
			} else {
				return m, tea.Quit
			}

		case "esc":
			if m.editing {
				m.cancelEdit()
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.timeblocks)-1 {
				m.cursor++
			}

		case "a":
			m.enterAddMode()
			return m, nil

		case "e":
			m.enterEditMode()
			return m, nil

		case "enter":
			if m.showInput && m.editing == false {
				m.saveAdd()
				return m, nil
			} else if m.editing == true {
				m.saveEdit()
				return m, nil
			} else {
				m.err = fmt.Errorf("invalid input")
			}

		case "tab":
			m.focused = (m.focused + 1) % len(m.inputFields)

			for i := range m.inputFields {
				if i == m.focused {
					m.inputFields[i].Focus()
				} else {
					m.inputFields[i].Blur()
				}
			}
			return m, nil
		}
	}

	if m.showInput {
		updatedInput, cmd := m.inputFields[m.focused].Update(msg)
		m.inputFields[m.focused] = updatedInput
		return m, cmd
	}

	if m.editing {
		for i := range m.inputFields {
			var cmd tea.Cmd
			m.inputFields[i], cmd = m.inputFields[i].Update(msg)
			_ = cmd
		}
	}

	sort.Slice(m.timeblocks, func(i, j int) bool {
		return m.timeblocks[i].starttime.Before(m.timeblocks[j].starttime)
	})

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	for i, timeblock := range m.timeblocks {
		duration := timeblock.endtime.Sub(timeblock.starttime).Minutes()
		lines := int(duration / 30)

		var styleToUse lipgloss.Style
		if m.cursor == i {
			styleToUse = selectedTimeStyle
		} else {
			styleToUse = timeStyle
		}

		var blockView *strings.Builder = &strings.Builder{}
		blockView.WriteString(styleToUse.Render(fmt.Sprintf("%s-%s",
			timeblock.starttime.Format("15:04"), timeblock.endtime.Format("15:04"))))
		blockView.WriteString("\n")
		blockView.WriteString(styleToUse.Render(timeblock.task))
		blockText := blockView.String()

		for j := 0; j < lines-1; j++ {
			blockText += "\n"
		}

		isFirst := i == 0
		isLast := i == len(m.timeblocks)-1

		style := getSeamlessBlockStyle(isFirst, isLast)
		b.WriteString(fmt.Sprintf("%s \n", style.Render(blockText)))
	}

	if m.editing {
		fmt.Fprintln(&b, "\nEdit a time block:")
		for i, input := range m.inputFields {
			indicator := " "
			if i == m.focused {
				indicator = ">"
			}
			fmt.Fprintf(&b, "  %s %s\n", indicator, input.View())
		}
		b.WriteString("\n [ tab: Cycle Focus | enter: Save | esc: Cancel ]")
	}

	if m.showInput {
		fmt.Fprintln(&b, "\nAdd a time block:")
		for i, input := range m.inputFields {
			indicator := " "
			if i == m.focused {
				indicator = ">"
			}
			fmt.Fprintf(&b, "  %s %s\n", indicator, input.View())
		}
		b.WriteString("\n [ tab: Cycle Focus | enter: Save | q: close ]")
	}

	if !m.showInput && !m.editing {
		b.WriteString("\n [ a: Add new time block | e: Edit selected time block | dd: Delete selected time block | j: Down | k: Up | enter/space: Toggle select | q: Quit ]")
	}
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
