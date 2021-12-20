package main

import (
  "bufio"
  "log"
  "regexp"
  "strconv"
  "os"
)

func main() {
  depth := 0
  relativeHorizTravel := 0

  file, err := os.Open("./commands.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  commandFinder := regexp.MustCompile(`^(up|down|forward) (\d+)$`)

  for scanner.Scan() {
    commandParts := commandFinder.FindStringSubmatch(scanner.Text())

    if commandParts == nil {
      log.Fatal("unable to parse command" + scanner.Text())
    }

    commandAmount, err := strconv.Atoi(commandParts[2])
    if err != nil {
      log.Fatal(err)
    }

    switch commandParts[1] {
    case "forward":
      relativeHorizTravel += commandAmount
    case "up":
      depth -= commandAmount
    case "down":
      depth += commandAmount
    default:
      log.Fatal("unknown command " +commandParts[1])
    }

    log.Printf("current position: %v over, %v down", relativeHorizTravel, depth)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  log.Printf("%v = %v * %v", relativeHorizTravel * depth, relativeHorizTravel, depth)
}
