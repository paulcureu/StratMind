package main

import (
	"fmt"
	"log"

	// Import your local parser code using module name + folder path
	"stratmind/core/parser"
)

func main() {
    // CHANGE THIS to your real .dem file
    demoPath := "D:\\VSCode\\StratMind\\demos\\falcons-vs-faze-m3-mirage.dem"
    outputPath := "core/export/timelines.json"

    fmt.Println("Starting parser on:", demoPath)
    err := parser.TrackRounds(demoPath, outputPath)
    if err != nil {
        log.Fatalf("Error parsing demo: %v", err)
    }

    fmt.Println("Done! JSON saved to:", outputPath)
}
