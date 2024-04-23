package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type SquareStyles struct {
	BorderColor      lipgloss.Color
	HoverBorderColor lipgloss.Color
	Styling          lipgloss.Style
}

type Square struct {
	// Check Symbol: "X" or "O" or "" if not checked
	Checker string
	Hovered bool
	Styles  SquareStyles
}

func (s *Square) New(hovered bool, checker string) *Square {
	defaultBorderColor := lipgloss.Color("238")
	hoverBorderColor := lipgloss.Color("255")

	borderColor := defaultBorderColor

	if hovered {
		borderColor = hoverBorderColor
	}

	return &Square{
		Checker: checker,
		Hovered: hovered,
		Styles: SquareStyles{
			BorderColor:      defaultBorderColor,
			HoverBorderColor: hoverBorderColor,
			Styling:          lipgloss.NewStyle().BorderForeground(borderColor).BorderStyle(lipgloss.ThickBorder()).Bold(true).Width(3).Height(1).Align(lipgloss.Center),
		},
	}
}

func GetCurrentPlayer(turn int, particpants []string) string {
	return particpants[turn%len(particpants)]
}

func CheckRows(m *model) bool {
	for _, row := range m.board {
		isWinningRow := true

		for _, col := range row {
			if row[0] == " " {
				isWinningRow = false
				break
			}

			if col != row[0] {
				isWinningRow = false
				break
			}
		}

		if isWinningRow {
			return true
		}
	}

	return false
}

func CheckCols(m *model) bool {
	for col := 0; col < len(m.board[0]); col++ {
		isWinningCol := true

		for row := 0; row < len(m.board); row++ {
			if m.board[0][col] == " " {
				isWinningCol = false
				break
			}

			if m.board[row][col] != m.board[0][col] {
				isWinningCol = false
				break
			}
		}

		if isWinningCol {
			return true
		}

	}

	return false
}

func CheckDiagonals(m *model) bool {
	dLeftWin := true
	dRightWin := true

	for ri, row := range m.board {
		if !dLeftWin && !dRightWin {
			return false
		}

		maxCol := len(row) - 1

		if m.board[0][0] == " " {
			dLeftWin = false
		}
		if m.board[0][maxCol] == " " {
			dRightWin = false
		}

		cLeft := row[ri]
		cRight := row[maxCol-ri]

		if cLeft != m.board[0][0] {
			dLeftWin = false
		}

		if cRight != m.board[0][maxCol] {
			dRightWin = false

		}
	}

	return dLeftWin || dRightWin
}

func CheckWin(m *model) bool {
	return CheckRows(m) || CheckCols(m) || CheckDiagonals(m)
}

// model, win, error (eg: occupied)
func RunGameTurn(m *model) (*model, bool, error) {
	selectedSquare := m.board[m.selY][m.selX]

	if selectedSquare != " " {
		return nil, false, fmt.Errorf("occupied")
	}

	m.board[m.selY][m.selX] = GetCurrentPlayer(m.turn, m.participants)

	if CheckWin(m) {
		return m, true, nil
	}

	m.turn += 1

	return m, false, nil
}
