package main

import (
  "bufio"
  "log"
  "os"
)

func main() {
  lines := 0
  var highs []int

  file, err := os.Open("./diagnostic-report.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    lines += 1
    // similar to looking for nil but also handles slices of zero length
    // `var highs []int` vs `highs := make([]int, 0)`
    if len(highs) == 0 {
      highs = make([]int, len(scanner.Text()))
      // each entry initialized to 0
    }
    for i, bitRaw := range scanner.Text() {
      // log.Printf("%v: %v", i, bitRaw)
      // ASCII representation of bits, so 48 == 0, 49 == 1
      if bitRaw == 49 {
        highs[i] += 1
      }
    }
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  log.Printf("lines: %v", lines);
  log.Printf("highs: %v", highs);

  gammaRate := 0
  epsilonRate := 0

  for _, highCount := range highs {
    gammaRate = gammaRate << 1
    epsilonRate = epsilonRate << 1

    if (highCount * 2 == lines) {
      log.Fatal("tie case not defined")
    } else if (highCount * 2 > lines) {
      gammaRate += 1
    } else {
      epsilonRate += 1
    }
  }

  log.Printf("gammaRate: %v", gammaRate)
  log.Printf("epsilonRate: %v", epsilonRate)

  log.Printf("power consumption: %v", gammaRate * epsilonRate)
}
