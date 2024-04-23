package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	participants   []string
	turn           int
	board          [3][3]string
	selX           int
	selY           int
	viewportWidth  int
	viewportHeight int
}

func initialModel() model {
	return model{
		participants: []string{"X", "O"},
		turn:         0,
		board: [3][3]string{
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "}},
		selX: 1,
		selY: 1,
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Tic Tac Toe")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		{
			m.viewportWidth = msg.Width
			m.viewportHeight = msg.Height
		}

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case tea.KeyCtrlC.String(), "q":
			{
				return m, tea.Quit
			}
		case tea.KeyUp.String():
			{
				if m.selY-1 >= 0 {
					m.selY -= 1
				}

				return m, nil
			}
		case tea.KeyDown.String():
			{
				if m.selY+1 < len(m.board) {
					m.selY += 1
				}

				return m, nil
			}
		case tea.KeyLeft.String():
			{
				if m.selX-1 >= 0 {
					m.selX -= 1
				}

				return m, nil
			}
		case tea.KeyRight.String():
			{
				if m.selX+1 < len(m.board) {
					m.selX += 1
				}

				return m, nil
			}
		// Enter and spacebar
		case tea.KeyEnter.String(), tea.KeySpace.String():
			{
				updatedModel, gameWon, err := RunGameTurn(&m)

				if err != nil {
					return m, nil
				}

				if gameWon {
					return m, tea.Quit
				}

				return updatedModel, tea.SetWindowTitle(fmt.Sprintf("T: %s, R: %d - Tic Tac Toe", GetCurrentPlayer(m.turn, m.participants), m.turn))
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) View() string {
	rows := make([]string, 3)

	for ri, row := range m.board {
		cells := make([]string, 3)

		for ci, square := range row {
			if m.selX == ci && m.selY == ri {
				cells[ci] = new(Square).New(true, square).Styles.Styling.Render(square)
			} else {
				cells[ci] = new(Square).New(false, square).Styles.Styling.Render(square)
			}
		}

		rows[ri] = lipgloss.JoinHorizontal(lipgloss.Top, cells...)
	}

	board := lipgloss.JoinVertical(lipgloss.Left, rows...)
	turnTitle := lipgloss.NewStyle().Bold(true).MarginBottom(1).Render(fmt.Sprintf("Current player: %s\nTurn no: %d", GetCurrentPlayer(m.turn, m.participants), m.turn+1))

	return lipgloss.JoinHorizontal(lipgloss.Center, strings.Repeat( /* CENTERS BOARD */ " ", max(0, (m.viewportWidth-lipgloss.Width(board))/2)), lipgloss.JoinVertical(lipgloss.Left, turnTitle, board))
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
