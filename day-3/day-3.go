package main

import (
  "bufio"
  "log"
  "os"
)

func parseReport(filePath string) ([]int, int) {
  var entries []int
  var numberOfBits int

  file, err := os.Open(filePath)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    if numberOfBits == 0 {
      numberOfBits = len(scanner.Text())
    }
    entry := 0
    for _, bitRaw := range scanner.Text() {
      entry = entry << 1
      // ASCII-esque representation of bits, so 48 == 0, 49 == 1
      if bitRaw == 49 {
        entry += 1
      }
    }
    entries = append(entries, entry)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  return entries, numberOfBits
}

func highBitCounts(numberOfBits int, entries []int) []int {
  highs := make([]int, numberOfBits)
  // initialized to 0s
  for _, entry := range(entries) {
    for i := 0; i < numberOfBits; i++ {
      if entry & (1 << i) != 0 {
        highs[i] += 1
      }
    }
  }
  return highs
}

func calculateOxygenGeneratorRating(numberOfBits int, oxygenCandidates []int) int {
  // log.Printf("oxygenCandidates: %v", oxygenCandidates)
  // log.Printf("oxygenCandidates: %b", oxygenCandidates)

  for bitOffset := numberOfBits - 1; bitOffset >= 0; bitOffset-- {
    var nextCandidates []int
    popularBits := highBitCounts(numberOfBits, oxygenCandidates)
    // log.Printf("popularBits: %v", popularBits)
    // the only difference from calculateScrubberRating, is there a clean way to combine?
    keepHighBit := popularBits[bitOffset] * 2 >= len(oxygenCandidates)

    for _, candidate := range(oxygenCandidates) {
      if (1 << bitOffset & candidate != 0) == keepHighBit {
        nextCandidates = append(nextCandidates, candidate)
      }
    }
    oxygenCandidates = nextCandidates[:]
    // log.Printf("oxygenCandidates: %b", oxygenCandidates)
    if len(oxygenCandidates) == 1 {
      // log.Printf("oxygenCandidates down to 1")
      return oxygenCandidates[0]
    }
  }

  // FIXME: return error? panic?
  return 0
}

func calculateScrubberRating(numberOfBits int, candidates []int) int {
  // log.Printf("candidates: %v", candidates)
  // log.Printf("candidates: %b", candidates)

  for bitOffset := numberOfBits - 1; bitOffset >= 0; bitOffset-- {
    var nextCandidates []int
    popularBits := highBitCounts(numberOfBits, candidates)
    // log.Printf("popularBits: %v", popularBits)
    // the only difference from calculateOxygenGeneratorRating, is there a clean way to combine?
    keepHighBit := popularBits[bitOffset] * 2 < len(candidates)

    for _, candidate := range(candidates) {
      if (1 << bitOffset & candidate != 0) == keepHighBit {
        nextCandidates = append(nextCandidates, candidate)
      }
    }
    candidates = nextCandidates[:]
    // log.Printf("candidates: %b", candidates)
    if len(candidates) == 1 {
      // log.Printf("candidates down to 1")
      return candidates[0]
    }
  }

  // FIXME: return error? panic?
  return 0
}

func main() {
  entries, numberOfBits := parseReport("./diagnostic-report.txt")
  highs := highBitCounts(numberOfBits, entries)

  // log.Printf("highs: %v", highs);

  gammaRate := 0
  epsilonRate := 0

  lines := len(entries)
  // log.Printf("lines: %v", lines);
  for i, highCount := range highs {
    if (highCount * 2 == lines) {
      log.Fatal("tie case not defined")
    } else if (highCount * 2 > lines) {
      gammaRate += 1 << i
    } else {
      epsilonRate += 1 << i
    }
  }

  log.Println("Power:")
  log.Printf(" gammaRate: %v %b", gammaRate, gammaRate)
  log.Printf(" epsilonRate: %v %b", epsilonRate, epsilonRate)
  log.Printf(" power consumption: %v", gammaRate * epsilonRate)

  oxygenGeneratorRating := calculateOxygenGeneratorRating(numberOfBits, entries[:])
  scrubberRating := calculateScrubberRating(numberOfBits, entries[:])
  log.Println("Air:")
  log.Printf(" Oxygen Generator Rating: %v %b", oxygenGeneratorRating, oxygenGeneratorRating)
  log.Printf(" CO2 Scrubber Rating: %v %b", scrubberRating, scrubberRating)
  log.Printf(" Life Support Rating: %v = %v * %v", oxygenGeneratorRating * scrubberRating, oxygenGeneratorRating, scrubberRating)
}
