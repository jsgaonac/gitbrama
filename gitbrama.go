package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"strings"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func getBranches() ([]string, error) {
	cmd := exec.Command("git", "branch")

	stdout, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	branchesStr := string(stdout)

	branches := []string{}

	for _, branch := range strings.Split(strings.TrimSpace(branchesStr), "\n") {
		if !strings.Contains(branch, "*") {
			branches = append(branches, strings.TrimSpace(branch))
		}
	}

	return branches, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	s := "Select branches to delete?\n\n"

	for i, choice := range m.choices {

		cursor := ""

		if m.cursor == i {
			cursor = ">"
		}

		checked := " "

		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress 'q' to quit"

	return s
}

func main() {

	branches, err := getBranches()

	if err != nil {
		fmt.Println("Error getting branches:", err)
		os.Exit(1)
	}

	for _, branch := range branches {
		fmt.Println(branch)
	}

	initialModel := model{
		choices:  []string{},
		selected: make(map[int]struct{}),
	}

	p := tea.NewProgram(initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
