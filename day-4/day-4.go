package main

import (
  "bufio"
  "regexp"
  "log"
  "math"
  "os"
  "strconv"
  "strings"
)

func parseGame(filePath string) ([]int, [][5][5]int) {
  var draws []int
  var boards [][5][5]int

  file, err := os.Open(filePath)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  // draws
  scanner.Scan()
  for _, raw := range(strings.Split(scanner.Text(), ",")) {
    draw, err := strconv.Atoi(raw)
    if err != nil {
      log.Fatal(err)
    }
    draws = append(draws, draw)
  }

  // boards
  squareFinder := regexp.MustCompile(`^\s*(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)$`)
  for scanner.Scan() {
    var board [5][5]int
    for rowIndex:=0; rowIndex<5; rowIndex++ {
      row := [5]int{0,0,0,0,0}
      scanner.Scan()
      squareParts := squareFinder.FindStringSubmatch(scanner.Text())
      for columnIndex := 1; columnIndex <= 5; columnIndex++ {
        square, err := strconv.Atoi(squareParts[columnIndex])
        if err != nil {
          log.Fatal(err)
        }
        row[columnIndex - 1] = square
      }
      board[rowIndex] = row
    }
    boards = append(boards, board)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  return draws, boards
}

func indexIn(list []int, val int) int {
  for index, listValue := range(list) {
    if (val == listValue) {
      return index
    }
  }
  return -1
}

func findLeastDrawsToWinBoard(draws []int, board [5][5]int) (int, [5]int) {
  // first check each row
  bestRowDraws := math.MaxInt
  var bestRow [5]int
  for _, row := range(board) {
    drawsToWinRow := -1
    // var squareDrawsRequired []int
    for _, square := range(row) {
      // log.Printf("square %v", square)
      // log.Printf("index? %v", indexIn(draws, square))
      squareDrawsRequired := indexIn(draws, square)
      if (drawsToWinRow < squareDrawsRequired && squareDrawsRequired > 0) {
        drawsToWinRow = squareDrawsRequired
      }
    }
    // log.Printf("drawsToWinRow: %v", drawsToWinRow)
    if (drawsToWinRow < bestRowDraws && drawsToWinRow > 0) {
      bestRowDraws = drawsToWinRow
      bestRow = row
    }
  }
  // log.Printf("bestRow: %v %v", bestRowDraws, bestRow)

  // then check each column
  bestColDraws := math.MaxInt
  var bestCol [5]int
  for colIndex := 0; colIndex < 5; colIndex++ {
    var colSquares [5]int
    drawsToWinCol := -1
    for rowIndex := 0; rowIndex < 5; rowIndex++ {
      square := board[rowIndex][colIndex]
      colSquares[rowIndex] = square
      squareDrawsRequired := indexIn(draws, square)
      // log.Printf("square, drawsRequired: %v, %v", square, squareDrawsRequired)
      if (drawsToWinCol < squareDrawsRequired && squareDrawsRequired >= 0) {
        drawsToWinCol = squareDrawsRequired
      }
    }
    // log.Printf("drawsToWinCol: %v", drawsToWinCol)
    if (drawsToWinCol < bestColDraws && drawsToWinCol >= 0) {
      bestColDraws = drawsToWinCol
      bestCol = colSquares
    }
  }
  // log.Printf("bestCol: %v %v", bestColDraws, bestCol)

  if (bestRowDraws == math.MaxInt) {
    if (bestColDraws == math.MaxInt) {
      // no winning rows or columns
      return -1, [5]int{-1, -1, -1, -1, -1}
    }
    return bestRowDraws, bestRow
  }
  return bestColDraws, bestCol
}

func main() {
  draws, boards := parseGame("./game.txt")

  log.Printf("draws: %v", draws);
  log.Printf("boards: %v", boards);

  bestDrawCount := math.MaxInt
  bestBoardIndex := -1
  var bestSquares [5]int
  for boardIndex, board := range(boards) {
    drawCount, squares := findLeastDrawsToWinBoard(draws, board)
    if (drawCount < bestDrawCount && drawCount >= 0) {
      bestDrawCount = drawCount
      bestBoardIndex = boardIndex
      bestSquares = squares
    }
  }

  if bestDrawCount == math.MaxInt {
    log.Println("no winning boards")
    os.Exit(0)
  }

  // most of the time we're lying, using the index instead of the actual count
  // so here we need to turn the draw index into the value a human expects
  // similar for board index, though we only use it locally
  log.Printf("board %v wins in %v draws", bestBoardIndex + 1, bestDrawCount + 1)
  squaresSum := 0
  for _, square := range(bestSquares) {
    squaresSum += square
  }
  log.Printf("winning score: %v = SUM(%v)", squaresSum, bestSquares)
}
