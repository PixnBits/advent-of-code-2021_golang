package main

import (
  "bufio"
  "log"
  "strconv"
  "os"
)

func main() {
  deeperCount := 0
  shallowerCount := 0

  file, err := os.Open("./depths.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  var previous int = 0
  for scanner.Scan() {
    lineDepth, err := strconv.Atoi(scanner.Text())
    if err != nil {
      log.Fatal(err)
    }
    // log.Print(lineDepth)

    if previous == 0 {
      previous = lineDepth
      log.Printf("%v (N/A - no previous measurement)", lineDepth)
      continue
    }

    if lineDepth < previous {
      shallowerCount += 1
      log.Printf("%v (decreased)", lineDepth)
    } else if lineDepth > previous {
      deeperCount += 1
      log.Printf("%v (increased)", lineDepth)
    } else {
      log.Printf("%v (unchanged)", lineDepth)
    }
    previous = lineDepth
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  log.Printf("summary: %v deepers, %v shallowers", deeperCount, shallowerCount)
}
