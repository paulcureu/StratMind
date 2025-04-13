package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var roundTimelines = make(map[string]*PlayerTimeline)

type MapZone struct {
	Name  string  `json:"name"`
	XMin  float64 `json:"x_min"`
	XMax  float64 `json:"x_max"`
	YMin  float64 `json:"y_min"`
	YMax  float64 `json:"y_max"`
}

func LoadZones(path string) []MapZone {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("❌ Nu pot deschide fișierul cu zone: %v", err)
	}
	defer file.Close()

	var zones []MapZone
	json.NewDecoder(file).Decode(&zones)
	return zones
}

func GetZoneName(x, y float64, zones []MapZone) string {
	for _, z := range zones {
		if x >= z.XMin && x <= z.XMax && y >= z.YMin && y <= z.YMax {
			return z.Name
		}
	}
	return "Unknown"
}

func TeamToString(t common.Team) string {
	switch t {
	case common.TeamTerrorists:
		return "Terrorists"
	case common.TeamCounterTerrorists:
		return "CounterTerrorists"
	case common.TeamSpectators:
		return "Spectators"
	default:
		return "Unknown"
	}
}

func TrackRounds(demoPath string, outputPath string) {
	zones := LoadZones("D:\\VSCode\\StratMind\\zones\\mirage_zones.json")

	file, err := os.Open(demoPath)
	if err != nil {
		log.Fatalf("❌ Nu am putut deschide fișierul DEM: %v", err)
	}
	defer file.Close()

	parser := dem.NewParser(file)
	defer parser.Close()

	gameState := parser.GameState()

	var timelines []PlayerTimeline
	currentRound := 0

	parser.RegisterEventHandler(func(e events.RoundStart) {
		currentRound++
		fmt.Printf("🔁 Începe runda #%d\n", currentRound)

		roundTimelines = make(map[string]*PlayerTimeline)

		for _, p := range gameState.Participants().Playing() {
			pl := &PlayerTimeline{
				RoundNumber: currentRound,
				SteamID:     fmt.Sprintf("%d", p.SteamID64),
				Nickname:    p.Name,
				Team:        TeamToString(p.Team),
				Side: map[common.Team]string{
					common.TeamTerrorists:        "T",
					common.TeamCounterTerrorists: "CT",
				}[p.Team],
				StartTick: gameState.IngameTick(),
				StartTime: parser.CurrentTime().Seconds(),
				StartPosition: Position{
					X: float64(p.Position().X),
					Y: float64(p.Position().Y),
					Z: float64(p.Position().Z),
				},
			}
			roundTimelines[pl.SteamID] = pl
		}
	})

	// 🟡 Tracking pe fiecare tick
	parser.RegisterNetMessageHandler(func(_ interface{}) {
		for _, p := range gameState.Participants().Playing() {
			if p == nil || p.IsBot || !p.IsAlive() {
				continue
			}

			steamID := fmt.Sprintf("%d", p.SteamID64)
			pl, ok := roundTimelines[steamID]
			if !ok {
				continue
			}

			tick := gameState.IngameTick()
			if tick%10 != 0 {
    			return
			}

			t := parser.CurrentTime().Seconds()
			x := float64(p.Position().X)
			y := float64(p.Position().Y)
			z := float64(p.Position().Z)

			zone := GetZoneName(x, y, zones)


			pl.Path = append(pl.Path, PositionTick{
				Tick:   tick,
				Time:   t,
				Position: Position{
					X: x,
					Y: y,
					Z: z,
				},
				Action:     zone,
				WeaponHeld: "", // optional completat
				IsScoped:   p.IsScoped(),
				IsDucking:  p.IsDucking(),
				IsMoving:   p.Velocity().X != 0 || p.Velocity().Y != 0,
				IsAirborne: p.IsAirborne(),
				HP:         p.Health(),
				InCombat:   false,
			})
		}
	})

	parser.RegisterEventHandler(func(e events.RoundEnd) {
		for _, pl := range roundTimelines {
			timelines = append(timelines, *pl)
		}
	})

	if err := parser.ParseToEnd(); err != nil {
		log.Fatalf("❌ Eroare la parsarea demo-ului: %v", err)
	}

	saveTimelinesAsJSON(timelines, outputPath)
	fmt.Println("✅ Timeline-urile au fost extrase și salvate.")
}

func saveTimelinesAsJSON(timelines []PlayerTimeline, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("❌ Eroare la salvarea fișierului JSON: %v", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(timelines)
	if err != nil {
		log.Fatalf("❌ Eroare la scrierea JSON-ului: %v", err)
	}
}
