package main

import (
	"deep-work-tui/cmd"
	"deep-work-tui/models"
	"deep-work-tui/styles"
	"deep-work-tui/utils"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"sort"
	"strings"
	"time"
)

type model struct {
	timeblocks     []models.Timeblock
	cursor         int
	selected       map[int]struct{}
	inputFields    []textinput.Model
	shutdownInput  []textinput.Model
	adding         bool
	editing        bool
	shutdown       bool
	editIndex      int
	focused        int
	err            error
	lastKey        string
	highlightIndex int
}

func initialModel(argstimeblocks []models.Timeblock) model {
	taskNameInput := textinput.New()
	taskNameInput.Placeholder = "Task Name"
	taskNameInput.CharLimit = 50

	startTimeInput := textinput.New()
	startTimeInput.Placeholder = "Start Time (HH:mm)"
	startTimeInput.CharLimit = 5

	endTimeInput := textinput.New()
	endTimeInput.Placeholder = "End Time (HH:mm)"
	endTimeInput.CharLimit = 5

	shutdownInput := textinput.New()
	shutdownInput.Placeholder = "SHUTDOWN COMPLETE"
	shutdownInput.CharLimit = 17
	var timeblocks []models.Timeblock
	if argstimeblocks != nil {
		timeblocks = argstimeblocks
	} else {
		timeblocks = []models.Timeblock{
			{"Deep Work", utils.ParseTime("07:00"), utils.ParseTime("10:00")},
			{"Email", utils.ParseTime("10:00"), utils.ParseTime("10:30")},
			{"Other Work", utils.ParseTime("10:30"), utils.ParseTime("12:00")},
			{"Meeting", utils.ParseTime("12:00"), utils.ParseTime("14:00")},
			{"Deep Work", utils.ParseTime("14:00"), utils.ParseTime("16:00")},
		}
	}

	return model{
		timeblocks:    timeblocks,
		selected:      make(map[int]struct{}),
		adding:        false,
		editing:       false,
		inputFields:   []textinput.Model{taskNameInput, startTimeInput, endTimeInput},
		focused:       0,
		shutdown:      false,
		shutdownInput: []textinput.Model{shutdownInput},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tick())
}
func tick() tea.Cmd {
	return tea.Tick(time.Minute, func(t time.Time) tea.Msg {
		return t
	})
}

func (m *model) enterAddMode() {
	m.adding = true
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

	if name != "" && utils.IsValidTime(start) && utils.IsValidTime(end) {
		m.timeblocks = append(m.timeblocks, models.Timeblock{Task: name, Starttime: utils.ParseTime(start), Endtime: utils.ParseTime(end)})
		m.adding = false
		m.clearInputFields()
	}
}

func (m *model) cancelAdd() {
	m.adding = false
	m.clearInputFields()
}

func (m *model) enterEditMode() {
	if len(m.timeblocks) == 0 {
		return
	}

	m.editing = true
	m.editIndex = m.cursor

	timeblock := m.timeblocks[m.editIndex]
	m.inputFields[0].SetValue(timeblock.Task)
	m.inputFields[1].SetValue(timeblock.Starttime.Format("15:04"))
	m.inputFields[2].SetValue(timeblock.Endtime.Format("15:04"))
}

func (m *model) saveEdit() {
	if !m.editing || m.editIndex < 0 || m.editIndex >= len(m.timeblocks) {
		return
	}

	m.timeblocks[m.editIndex].Task = m.inputFields[0].Value()
	m.timeblocks[m.editIndex].Starttime, _ = time.Parse("15:04", m.inputFields[1].Value())
	m.timeblocks[m.editIndex].Endtime, _ = time.Parse("15:04", m.inputFields[2].Value())

	m.editing = false
	m.editIndex = -1

	m.clearInputFields()
}

func (m *model) cancelEdit() {
	m.editing = false
	m.editIndex = -1
	m.clearInputFields()
}

func (m *model) cancelShutdown() {
	m.shutdown = false
	m.shutdownInput[0].SetValue("")
}

func (m *model) enterShutdownMode() {
	m.shutdown = true
	m.shutdownInput[0].SetValue("")
	m.shutdownInput[0].Focus()
}

func (m *model) clearInputFields() {
	for i := range m.inputFields {
		m.inputFields[i].SetValue("")
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case time.Time:
		now := msg

		for i, tb := range m.timeblocks {
			nowHour, nowMin, _ := now.Clock()
			startHour, startMin, _ := tb.Starttime.Clock()
			endHour, endMin, _ := tb.Endtime.Clock()

			if (nowHour > startHour || (nowHour == startHour && nowMin >= startMin)) &&
				(nowHour < endHour || (nowHour == endHour && nowMin < endMin)) {
				m.highlightIndex = i
				break
			}
		}
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.adding {
				m.cancelAdd()
			} else {
				return m, tea.Quit
			}

		case "esc":
			if m.editing {
				m.cancelEdit()
			} else if m.shutdown {
				m.cancelShutdown()
			} else if m.adding {
				m.cancelAdd()
			}

		case "up", "k":
			if m.cursor > 0 && !m.editing && !m.adding && !m.shutdown {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.timeblocks)-1 && !m.editing && !m.adding && !m.shutdown {
				m.cursor++
			}

		case "a":
			if !m.editing && !m.shutdown && !m.adding {
				m.enterAddMode()
				return m, nil
			}

		case "e":
			if !m.adding && !m.shutdown && !m.editing {
				m.enterEditMode()
				return m, nil
			}

		case "s":
			if !m.adding && !m.editing && !m.shutdown {
				m.enterShutdownMode()
				return m, nil
			}

		case "d":
			if m.lastKey == "d" && !m.editing && !m.adding {
				// Remove the timeblock at the current cursor position
				if m.cursor >= 0 && m.cursor < len(m.timeblocks) {
					indexToRemove := m.cursor
					m.timeblocks = append(m.timeblocks[:indexToRemove], m.timeblocks[indexToRemove+1:]...)

					// Adjust cursor if necessary
					if m.cursor >= len(m.timeblocks) {
						m.cursor = len(m.timeblocks) - 1
					}
				}
				m.lastKey = "" // Reset the last key
				return m, nil
			}
			m.lastKey = "d" // Store this key as the last pressed

		case "enter":
			if m.adding && m.editing == false {
				m.saveAdd()
				return m, nil
			} else if m.editing == true {
				m.saveEdit()
				return m, nil
			} else if m.shutdown == true {
				return m, tea.Quit
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

		case "ctrl-a":
			if m.cursor >= 0 && m.cursor < len(m.timeblocks) && !m.editing && !m.adding {
				tb := &m.timeblocks[m.cursor]
				tb.Starttime = tb.Starttime.Add(15 * time.Minute)
				tb.Endtime = tb.Endtime.Add(15 * time.Minute)
			}
			return m, nil
		}

	}

	if m.adding {
		updatedInput, cmd := m.inputFields[m.focused].Update(msg)
		m.inputFields[m.focused] = updatedInput
		return m, cmd
	}

	if m.shutdown {
		updatedInput, cmd := m.shutdownInput[0].Update(msg)
		m.shutdownInput[0] = updatedInput
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
		return m.timeblocks[i].Starttime.Before(m.timeblocks[j].Starttime)
	})

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	for i, timeblock := range m.timeblocks {
		duration := timeblock.Endtime.Sub(timeblock.Starttime).Minutes()
		lines := int(duration / 30)

		var styleToUse lipgloss.Style
		if m.cursor == i {
			styleToUse = styles.SelectedTimeStyle
		} else if m.highlightIndex == i {
			styleToUse = styles.HighlightStyle
		} else {
			styleToUse = styles.TimeStyle
		}
		var blockView *strings.Builder = &strings.Builder{}
		blockView.WriteString(styleToUse.Render(fmt.Sprintf("%s-%s",
			timeblock.Starttime.Format("15:04"), timeblock.Endtime.Format("15:04"))))
		blockView.WriteString("\n")
		blockView.WriteString(styleToUse.Render(timeblock.Task))
		blockText := blockView.String()

		for j := 0; j < lines-1; j++ {
			blockText += "\n"
		}

		isFirst := i == 0
		isLast := i == len(m.timeblocks)-1

		style := styles.GetSeamlessBlockStyle(isFirst, isLast)
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

	if m.adding {
		fmt.Fprintln(&b, "\nAdd a time block:")
		for i, input := range m.inputFields {
			indicator := " "
			if i == m.focused {
				indicator = ">"
			}
			fmt.Fprintf(&b, "  %s %s\n", indicator, input.View())
		}
		b.WriteString("\n [ tab: Cycle Focus | enter: Save | esc: Cancel ]")
	}

	if !m.adding && !m.editing && !m.shutdown {
		b.WriteString("\n [ a: Add new time block | e: Edit time block | dd: Delete time block | j: Down | k: Up | s: shutdown ]")
	}

	if m.shutdown {
		fmt.Fprintln(&b, "\nEnter SHUTDOWN COMPLETE to end the day:")
		fmt.Fprintln(&b, m.shutdownInput[0].View())
		b.WriteString("\n enter: SHUTDOWN | esc: Cancel ]")
	}

	return b.String()
}

func main() {
	timeblocks, err := cmd.ParseArgs()
	if err != nil {
		fmt.Printf("Error parsing blocks: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(timeblocks))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
