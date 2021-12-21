package main

import (
  "bufio"
  "flag"
  "log"
  "os"
  "regexp"
  "strconv"
)

// structs!
type Point struct {
  X int
  Y int
}

type Line struct {
  Start Point
  End Point
}

func min(a int, b int) int {
  if a < b {
    return a
  }
  return b
}

func max(a int, b int) int {
  if a > b {
    return a
  }
  return b
}

func parseDataFile(filepath string) ([]Line, Point) {
  ventLineFinder := regexp.MustCompile(`^(\d+),(\d+) -> (\d+),(\d+)$`)
  var lines []Line
  bounds := Point{}

  log.Printf("reading file %v", filepath)
  file, err := os.Open(filepath)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    // secret sauce
    lineParts := ventLineFinder.FindStringSubmatch(scanner.Text())
    if lineParts == nil {
      log.Fatal("unable to parse vent line " + scanner.Text())
    }

    startX, err := strconv.Atoi(lineParts[1])
    if err != nil {
      log.Fatal(err)
    }
    startY, err := strconv.Atoi(lineParts[2])
    if err != nil {
      log.Fatal(err)
    }
    start := Point{X:startX, Y:startY}
    // func max(a Point, b Point) Point {} ?
    if start.X > bounds.X {
      bounds.X = start.X
    }
    if start.Y > bounds.Y {
      bounds.Y = start.Y
    }

    endX, err := strconv.Atoi(lineParts[3])
    if err != nil {
      log.Fatal(err)
    }
    endY, err := strconv.Atoi(lineParts[4])
    if err != nil {
      log.Fatal(err)
    }
    end := Point{X:endX, Y:endY}
    if end.X > bounds.X {
      bounds.X = end.X
    }
    if end.Y > bounds.Y {
      bounds.Y = end.Y
    }

    line := Line{Start:start, End:end}
    // log.Printf("line: %v", line)
    lines = append(lines, line)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  bounds.X += 1
  bounds.Y += 1
  return lines, bounds
}

func buildMap(bounds Point, lines []Line) [][]int {
  var voxels [][]int
  for rowIndex := 0; rowIndex < bounds.Y; rowIndex++ {
    row := make([]int, bounds.X)
    voxels = append(voxels, row)
  }

  for _, line := range(lines) {
    if line.Start.X == line.End.X {
      // horizontal
      lower := min(line.Start.Y, line.End.Y)
      higher := max(line.Start.Y, line.End.Y)
      for i := lower; i <= higher; i++ {
        voxels[i][line.Start.X] += 1
      }
    } else if line.Start.Y == line.End.Y {
      // vertical
      lower := min(line.Start.X, line.End.X)
      higher := max(line.Start.X, line.End.X)
      for i := lower; i <= higher; i++ {
        voxels[line.Start.Y][i] += 1
      }
    } else {
      // 45 diag
      // TODO: could this be simpler? ex: tried to min & max too soon?
      xStart := min(line.Start.X, line.End.X)
      xEnd := max(line.Start.X, line.End.X)
      yStart := min(line.Start.Y, line.End.Y)
      yEnd := max(line.Start.Y, line.End.Y)

      diffX := xEnd - xStart
      diffY := yEnd - yStart
      if (diffX != diffY) {
        log.Fatal("non-45 line in data %v", line)
      }

      xDirection := 1
      if (xStart != line.Start.X) {
        xDirection = -1
        xStart = line.Start.X
        xEnd = line.End.X
      }
      yDirection := 1
      if (yStart != line.Start.Y) {
        yDirection = -1
        yStart = line.Start.Y
        yEnd = line.End.Y
      }

      // log.Printf("line: %v", line)
      // log.Printf("start & end X: %v, %v", xStart, xEnd)
      // log.Printf("start & end Y: %v, %v", yStart, yEnd)
      // log.Printf("direction X & Y: %v, %v", xDirection, yDirection)

      for i := 0; i <= diffX; i++ {
        voxels[yStart + i * yDirection][xStart + i * xDirection] += 1
      }
    }
  }

  return voxels
}

func printMap(voxels [][]int) {
  for _, row := range(voxels) {
    var txt = ""
    for _, voxel := range(row) {
      if voxel == 0 {
        txt = txt + "."
      } else {
        txt = txt + strconv.Itoa(voxel)
      }
    }
    log.Print(txt)
  }
}

func countPointsOverlapping(voxels [][]int, threshold int) int {
  count := 0
  for _, row := range(voxels) {
    for _, voxel := range(row) {
      if voxel >= threshold {
        count += 1
      }
    }
  }
  return count
}

func main() {
  // experiment with using STDIN on a later day
  var filepath string
  flag.StringVar(&filepath, "data", "./vents.txt", "path to file containing vent data")
  flag.Parse()

  // returned in order of usefulness
  lines, bounds := parseDataFile(filepath)
  // log.Printf("lines: %v", lines)
  log.Printf("bounds: %v", bounds)

  // order of importance? or is consistency with parseDataFile better?
  voxelMap := buildMap(bounds, lines)
  printMap(voxelMap)

  dangerPointCount := countPointsOverlapping(voxelMap, 2)
  log.Printf("number of dangerous points: %v", dangerPointCount)
}
