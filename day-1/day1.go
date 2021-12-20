package main

import (
  "bufio"
  "log"
  "strconv"
  "os"
)

// TODO: pass by reference to improve performance?
func sum(window []int) int {
  total := 0
  for _, num := range window {
    total += num
  }
  return total
}

func main() {
  deeperCount := 0
  shallowerCount := 0

  file, err := os.Open("./depths.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  previous := 0
  windowBacking := [3]int{0, 0, 0}
  window := windowBacking[:]

  for scanner.Scan() {
    lineDepth, err := strconv.Atoi(scanner.Text())
    if err != nil {
      log.Fatal(err)
    }

    // shift off the front, push on to the end
    window = window[1:]
    window = append(window, lineDepth)

    if window[0] == 0 {
      log.Printf("%v (still filling the window)", lineDepth)
      continue
    }

    current := sum(window)

    if previous == 0 {
      previous = current
      log.Printf("%v (N/A - no previous measurement)", lineDepth)
      continue
    }

    if current < previous {
      shallowerCount += 1
      log.Printf("%v (decreased) %v from %v", lineDepth, current, previous)
    } else if current > previous {
      deeperCount += 1
      log.Printf("%v (increased) %v from %v", lineDepth, current, previous)
    } else {
      log.Printf("%v (no change) %v from %v", lineDepth, current, previous)
    }
    previous = current
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  log.Printf("summary: %v deepers, %v shallowers", deeperCount, shallowerCount)
}
