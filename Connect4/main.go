package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var board [][]byte
var userInput string
var playerWin bool
var compWin bool

func setupBoard() {

	board = make([][]byte, 6)
	for index, _ := range board {
		board[index] = make([]byte, 7)
	}

	for i, row := range board {
		for j, _ := range row {
			board[i][j] = '_'
		}
	}
}

func printBoard() {
	fmt.Println("  0 1 2 3 4 5 6")
	for i, row := range board {
		fmt.Print("|")
		for j := range row {
			fmt.Printf(" %v", string(board[i][j]))
		}
		fmt.Print(" |\n")
	}
	fmt.Println()
}

func validInput(input string) (bool, int) {
	num, err := strconv.Atoi(input)
	if err != nil {
		return false, -1
	} else if num > 6 || num < 0 {
		return false, -1
	} else if board[0][num] != '_' {
		return false, -1
	}
	return true, num
}

func placePiece(col_index int, isPlayer bool) (int, int) {
	row_index := 0

	for row_index < 5 && board[row_index+1][col_index] == '_' {
		row_index++
	}

	if board[row_index][col_index] != '_' {
		return -1, -1
	}

	switch isPlayer {
	case true:
		board[row_index][col_index] = 'X'
		return row_index, col_index
	default:
		board[row_index][col_index] = 'O'
		return row_index, col_index
	}
}

/*
	The intent here is to have the comp choose a space to play in a semi-intelligent way,
	to better represent a human opponent

	For now, I will implement the ai to use random-select for the sake of simplicity in figuring out
	the rest of the game's logic
*/

func compChooseSpace() (int, int) {
	return placePiece(randomCol(), false)
}

func randomCol() int {
	return rand.Intn(7)
}

/*
	We want to check to the left and right of the given position to see
	if a 4-in-a-row is made
*/

func checkLeftRight(row_index int, col_index int, token byte) bool {
	leftFinds := 0
	col_num := col_index - 1
	for col_num >= 0 && board[row_index][col_num] == token {
		leftFinds++
		col_num--
	}

	rightFinds := 0
	col_num = col_index + 1
	for col_num <= 6 && board[row_index][col_num] == token {
		rightFinds++
		col_num++
	}

	return (leftFinds + rightFinds) >= 3
}

/*
	Check down from a space to see if a 4-in-a-row connection
	has been made
*/

func checkDown(row_index int, col_index int, token byte) bool {

	downFinds := 0
	row_num := row_index + 1
	for row_num <= 5 && board[row_num][col_index] == token {
		downFinds++
		row_num++
	}

	return downFinds >= 3
}

func checkUpLeftDownRight(row_index int, col_index int, token byte) bool {
	row_num := row_index - 1
	col_num := col_index - 1
	upLeftFinds := 0
	downRightFinds := 0
	for row_num >= 0 && col_num >= 0 && board[row_num][col_num] == token {
		upLeftFinds++
		row_num--
		col_num--
	} //Check Up Left

	row_num = row_index + 1
	col_num = col_index + 1
	for row_num <= 5 && col_num <= 6 && board[row_num][col_num] == token {
		downRightFinds++
		row_num++
		col_num++
	} //Check Down Right

	return (upLeftFinds + downRightFinds) >= 3
}

func checkUpRightDownLeft(row_index int, col_index int, token byte) bool {
	upRightFinds := 0
	downLeftFinds := 0
	row_num := row_index - 1
	col_num := col_index + 1
	for row_num >= 0 && col_num <= 6 && board[row_num][col_num] == token {
		upRightFinds++
		row_num--
		col_num++
	} //Check Up Right

	row_num = row_index + 1
	col_num = col_index - 1
	for row_num <= 5 && col_num >= 0 && board[row_num][col_num] == token {
		downLeftFinds++
		row_num++
		col_num--
	} //Check Down Left

	return (upRightFinds + downLeftFinds) >= 3
}

func checkFor4(row_index int, col_index int, isPlayer bool) bool {
	var token byte
	if isPlayer {
		token = 'X'
	} else {
		token = 'O'
	}
	return checkLeftRight(row_index, col_index, token) || checkDown(row_index, col_index, token) || checkUpLeftDownRight(row_index, col_index, token) || checkUpRightDownLeft(row_index, col_index, token)
}

func main() {
	rand.Seed(time.Now().UnixMilli())
	setupBoard()
	printBoard()

	for {
		fmt.Print("Enter a column coordinate (0-6): ")
		fmt.Scanln(&userInput)
		isValid, result := validInput(userInput)
		if isValid {
			res_row, res_col := placePiece(result, true)
			printBoard()
			if checkFor4(res_row, res_col, true) {
				fmt.Println("4 in a Row, Player Wins")
				break
			}
			time.Sleep(time.Second * 2)

			comp_row, comp_col := compChooseSpace()
			for {
				if comp_row == -1 {
					comp_row, comp_col = compChooseSpace()
				} else {
					break
				}
			}
			printBoard()
			if checkFor4(comp_row, comp_col, false) {
				fmt.Println("4 in a Row, Comp Wins")
				break
			}
			time.Sleep(time.Second * 2)
		} else {
			fmt.Println("Invalid")
		}
	}
}
