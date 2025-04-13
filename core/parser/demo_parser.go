package parser

import (
	"log"
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
)

func ParseDemo(path string) {
    f, err := os.Open(path)
    if err != nil {
        log.Fatalf("Nu am putut deschide demo-ul: %v", err)
    }
    defer f.Close()

    parser := dem.NewParser(f)
    defer parser.Close()

    err = parser.ParseToEnd()
    if err != nil {
        log.Fatalf("Eroare la parsat: %v", err)
    }
}
