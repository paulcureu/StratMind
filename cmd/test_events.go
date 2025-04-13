package main

import (
	"fmt"
	"log"
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func main() {
    f, err := os.Open("D:\\VSCode\\StratMind\\demos\\falcons-vs-faze-m3-mirage.dem") //adresa dem.
    if err != nil {
        log.Fatalf("❌ Nu pot deschide fișierul DEM: %v", err)
    }
    defer f.Close()

    parser := dem.NewParser(f)
    defer parser.Close()

    // Ex: handler pentru Kill
    parser.RegisterEventHandler(func(e events.Kill) {
        fmt.Printf("🔫 %s → %s cu %s\n", e.Killer, e.Victim, e.Weapon)
    })

    // Parse demo
    if err := parser.ParseToEnd(); err != nil {
        log.Fatalf("❌ Eroare la parsare: %v", err)
    }

    fmt.Println("✅ Parsing finalizat cu succes!")
}
