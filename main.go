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
	timeblocks []Timeblock
	cursor     int
	selected   map[int]struct{}
	taskInput  textinput.Model
	startTime  textinput.Model
	endTime    textinput.Model
	focused    int
}

func parseTime(t string) time.Time {
	parsed, _ := time.Parse("15:04", t)
	return parsed
}

func initialModel() model {
	return model{
		timeblocks: []Timeblock{
			{"Deep Work", parseTime("07:00"), parseTime("10:00")},
			{"Other Work", parseTime("10:00"), parseTime("12:00")},
			{"Meeting", parseTime("12:00"), parseTime("14:00")},
			{"Deep Work", parseTime("14:00"), parseTime("16:00")},
		},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.timeblocks)-1 {
				m.cursor++
			}

		case "a":
			m.timeblocks = append(m.timeblocks, Timeblock{"New work", parseTime("16:00"), parseTime("18:00")})
			return m, nil

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
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

	b.WriteString("\n [ a: Add new timeblock | j: Down | k: Up | enter/space: Toggle select | q: Quit ]")

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
