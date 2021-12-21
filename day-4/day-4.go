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

func calcBoardWinningScore(draws []int, board [5][5]int) (int, int) {
  // map the board's squares to the first draw to mark it off
  // log.Printf("board: %v", board)
  var boardDraws [5][5]int
  for rowIndex, row := range(board) {
    var rowDraws [5]int
    for colIndex, square := range(row) {
      rowDraws[colIndex] = indexIn(draws, square)
    }
    boardDraws[rowIndex] = rowDraws
  }
  // log.Printf("boardDraws: %v", boardDraws)

  // find the row or column with the lowest draws required Min(Max(squareDraws))
  // rows
  bestWinningDraw := math.MaxInt
  for _, rowDraws := range(boardDraws) {
    drawsToWinRow := -1
    for _, squareDraw := range(rowDraws) {
      if drawsToWinRow < squareDraw && squareDraw != -1 {
        drawsToWinRow = squareDraw
      }
    }
    // log.Printf("drawsToWinRow: %v", drawsToWinRow)
    if drawsToWinRow < bestWinningDraw && drawsToWinRow != -1 {
      bestWinningDraw = drawsToWinRow
    }
  }
  // log.Printf("bestWinningDraw: %v", bestWinningDraw)
  // columns
  for colIndex := 0; colIndex < 5; colIndex++ {
    drawsToWinCol := -1
    for rowIndex := 0; rowIndex < 5; rowIndex++ {
      squareDraw := boardDraws[rowIndex][colIndex]
      if (drawsToWinCol < squareDraw && squareDraw != -1) {
        drawsToWinCol = squareDraw
      }
    }
    // log.Printf("drawsToWinCol: %v", drawsToWinCol)
    if drawsToWinCol < bestWinningDraw && drawsToWinCol != -1 {
      bestWinningDraw = drawsToWinCol
    }
  }
  // log.Printf("bestWinningDraw: %v", bestWinningDraw)

  // now calculate the score
  // sum all unmarked squares
  // that is, sum the square values if their draw index is over the winning draw
  sum := 0
  for rowIndex, rowDraws := range(boardDraws) {
    for colIndex, squareDraw := range(rowDraws) {
      if (squareDraw > bestWinningDraw || squareDraw == -1) {
        sum += board[rowIndex][colIndex]
      }
    }
  }
  // log.Printf("sum of unmarked squares: %v", sum)
  // log.Printf("last number drawn to win: %v", draws[bestWinningDraw])
  score := sum * draws[bestWinningDraw]
  // log.Printf("winning score for board: %v", score)
  return bestWinningDraw + 1, score
}

func main() {
  draws, boards := parseGame("./game.txt")

  log.Printf("draws: %v", draws);
  log.Printf("boards: %v", boards);

  bestWinningDraw := math.MaxInt
  winningScore := -1
  bestBoardIndex := -1
  for boardIndex, boards := range(boards) {
    drawsToWin, score := calcBoardWinningScore(draws, boards)
    if (drawsToWin < bestWinningDraw && drawsToWin != -1) {
      bestWinningDraw = drawsToWin
      winningScore = score
      bestBoardIndex = boardIndex
    }
  }

  if bestWinningDraw == math.MaxInt {
    log.Println("no winning boards")
    os.Exit(0)
  }

  // humans start at 1
  log.Printf(
    "board %v wins in %v draws with score %v",
    bestBoardIndex + 1,
    bestWinningDraw,
    winningScore,
  )
}
