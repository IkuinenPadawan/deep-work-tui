package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"os"
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
	focused     int
	err         error
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
			{"Other Work", parseTime("10:00"), parseTime("12:00")},
			{"Meeting", parseTime("12:00"), parseTime("14:00")},
			{"Deep Work", parseTime("14:00"), parseTime("16:00")},
		},
		selected:    make(map[int]struct{}),
		showInput:   false,
		inputFields: []textinput.Model{taskNameInput, startTimeInput, endTimeInput},
		focused:     0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.showInput {
				m.showInput = false
			} else {
				return m, tea.Quit
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
			m.showInput = true
			m.focused = 0
			for i := range m.inputFields {
				if i == m.focused {
					m.inputFields[i].Focus()
				} else {
					m.inputFields[i].Blur()
				}
			}
			return m, nil

		case "enter":
			if m.showInput {
				name := m.inputFields[0].Value()
				start := m.inputFields[1].Value()
				end := m.inputFields[2].Value()

				if name != "" && isValidTime(start) && isValidTime(end) {
					m.timeblocks = append(m.timeblocks, Timeblock{task: name, starttime: parseTime(start), endtime: parseTime(end)})
					m.showInput = false
					for i := range m.inputFields {
						m.inputFields[i].SetValue("")
					}
				} else {
					m.err = fmt.Errorf("invalid input")
				}
				return m, nil
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

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	for i, timeblock := range m.timeblocks {
		var blockView *strings.Builder = &strings.Builder{}
		blockView.WriteString(fmt.Sprintf("%s-%s",
			timeblock.starttime.Format("15:04"), timeblock.endtime.Format("15:04")))
		blockView.WriteString(" ")
		blockView.WriteString(timeblock.task)
		blockText := blockView.String()

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, blockText))
	}

	if m.showInput {
		fmt.Fprintln(&b, "\nAdd a time block:")
		for i, input := range m.inputFields {
			indicator := " " // Default no indicator
			if i == m.focused {
				indicator = ">" // Highlight focused input
			}
			fmt.Fprintf(&b, "  %s %s\n", indicator, input.View())
		}
		b.WriteString("\n [ tab: Cycle Focus | enter: Save | q: close ]")
	} else {
		b.WriteString("\n [ a: Add new timeblock | j: Down | k: Up | enter/space: Toggle select | q: Quit ]")
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
