package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// timeline-urile pe rundă
var roundTimelines = make(map[string]*PlayerTimeline)
var currentRound = 0

// vom ține un flag: freezeTimeOver = false până se termină freeze time
var freezeTimeOver = false

// Mapă pentru a preveni duplicarea la același tick (dacă primim multiple net messages):
var lastTickSaved = make(map[string]int)

// TrackRounds parsează DEM-ul și produce un JSON cu timeline-uri
func TrackRounds(demoPath string, outputPath string) error {
	// Încarcă zonele din fișierul JSON (modifică calea dacă fișierul tău e în altă parte)
	zones := LoadZones("D:/VSCode/StratMind/zones/mirage_zones.json")

	file, err := os.Open(demoPath)
	if err != nil {
		return fmt.Errorf("nu pot deschide fișierul DEM: %w", err)
	}
	defer file.Close()

	parser := dem.NewParser(file)
	defer parser.Close()

	gameState := parser.GameState()

	var allTimelines []PlayerTimeline

	// ====================================
	// RoundStart -> reset & freezeTimeOver = false
	// ====================================
	parser.RegisterEventHandler(func(e events.RoundStart) {
		currentRound++
		fmt.Printf("Începe runda #%d\n", currentRound)

		roundTimelines = make(map[string]*PlayerTimeline)
		freezeTimeOver = false

		for _, p := range gameState.Participants().Playing() {
			if p == nil || p.IsBot {
				continue
			}
			steamID := fmt.Sprintf("%d", p.SteamID64)
			pl := &PlayerTimeline{
				RoundNumber: currentRound,
				SteamID:     steamID,
				Nickname:    p.Name,
				Team:        TeamToString(p.Team),
				Side:        sideString(p.Team),
				StartTick:   gameState.IngameTick(),
				StartTime:   parser.CurrentTime().Seconds(),
				StartPosition: Position{
					X: float64(p.Position().X),
					Y: float64(p.Position().Y),
					Z: float64(p.Position().Z),
				},
				UtilityThrown: UtilityStats{
					Smokes:   []string{},
					Molotovs: []string{},
				},
			}
			roundTimelines[steamID] = pl
		}
	})

	// ====================================
	// RoundFreezetimeEnd -> freezeTimeOver = true
	// ====================================
	parser.RegisterEventHandler(func(e events.RoundFreezetimeEnd) {
		freezeTimeOver = true
		fmt.Println("Freeze time s-a încheiat, acum colectăm date reale.")
	})

	// ====================================
	// Colectare date: net messages (mișcare, poziție, etc.)
	// ====================================
	parser.RegisterNetMessageHandler(func(_ interface{}) {
		// dacă suntem încă în freeze time, nu colectăm
		if !freezeTimeOver {
			return
		}

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

			// exemplu: colectăm la fiecare 10 tick-uri
			if tick%5 != 0 {
				continue
			}

			// Verificăm dacă am salvat deja pentru jucătorul ăsta la tick-ul curent
			if lastTickSaved[steamID] == tick {
				// deja salvat -> nu mai adăugăm iarăși
				continue
			}
			// Marchează tick-ul ca salvat
			lastTickSaved[steamID] = tick

			// Construim intrarea pt path
			pt := PositionTick{
				Tick: tick,
				Time: parser.CurrentTime().Seconds(),
				Position: Position{
					float64(p.Position().X),
					float64(p.Position().Y),
					float64(p.Position().Z),
				},
				Zone: GetZoneName(
					float64(p.Position().X),
					float64(p.Position().Y),
					float64(p.Position().Z),
					zones,
				),
				
				
				Action:        getPlayerAction(p),
				WeaponHeld:    getWeaponName(p),
				IsScoped:      p.IsScoped(),
				IsDucking:     p.IsDucking(),
				IsMoving:      (p.Velocity().X != 0 || p.Velocity().Y != 0),
				IsAirborne:    p.IsAirborne(),
				HP:            p.Health(),
				ViewAngle:     float64(p.ViewDirectionX()),
				InCombat:      false,
				NearTeammates: countNearbyTeammates(p, gameState),
			}
			pl.Path = append(pl.Path, pt)
		}
	})

	// ====================================
	// Utility: flash, smoke, he, molotov
	// ====================================
	parser.RegisterEventHandler(func(e events.FlashExplode) {
		if e.Thrower == nil {
			return
		}
		if tl, ok := roundTimelines[fmt.Sprintf("%d", e.Thrower.SteamID64)]; ok {
			tl.UtilityThrown.Flashes++
		}
	})

	parser.RegisterEventHandler(func(e events.SmokeStart) {
		if e.Thrower == nil {
			return
		}
		if tl, ok := roundTimelines[fmt.Sprintf("%d", e.Thrower.SteamID64)]; ok {
			pos := fmt.Sprintf("%.0f,%.0f", e.Position.X, e.Position.Y)
			tl.UtilityThrown.Smokes = append(tl.UtilityThrown.Smokes, pos)
		}
	})

	parser.RegisterEventHandler(func(e events.FireGrenadeStart) {
		if e.Thrower == nil {
			return
		}
		if tl, ok := roundTimelines[fmt.Sprintf("%d", e.Thrower.SteamID64)]; ok {
			pos := fmt.Sprintf("%.0f,%.0f", e.Position.X, e.Position.Y)
			tl.UtilityThrown.Molotovs = append(tl.UtilityThrown.Molotovs, pos)
		}
	})

	parser.RegisterEventHandler(func(e events.HeExplode) {
		if e.Thrower == nil {
			return
		}
		if tl, ok := roundTimelines[fmt.Sprintf("%d", e.Thrower.SteamID64)]; ok {
			tl.UtilityThrown.HEGrenades++
		}
	})

	// ====================================
	// Kills & Deaths
	// ====================================
	parser.RegisterEventHandler(func(e events.Kill) {
		if e.Killer == nil || e.Victim == nil {
			return
		}
		killerID := fmt.Sprintf("%d", e.Killer.SteamID64)
		victimID := fmt.Sprintf("%d", e.Victim.SteamID64)

		if tl, ok := roundTimelines[killerID]; ok {
			ke := KillEvent{
				Tick:     gameState.IngameTick(),
				Time:     parser.CurrentTime().Seconds(),
				Victim:   e.Victim.Name,
				Weapon:   e.Weapon.String(),
				Headshot: e.IsHeadshot,
			}
			tl.KillEvents = append(tl.KillEvents, ke)
		}

		if tl, ok := roundTimelines[victimID]; ok {
			tl.DeathEvent = &DeathEvent{
				Tick:   gameState.IngameTick(),
				Time:   parser.CurrentTime().Seconds(),
				Killer: e.Killer.Name,
			}
		}
	})

	// ====================================
	// RoundEnd -> finalize timelines
	// ====================================
	parser.RegisterEventHandler(func(e events.RoundEnd) {
		for _, pl := range roundTimelines {
			pl.Survived = (pl.DeathEvent == nil)
			if pl.Survived {
				pl.TimeAlive = parser.CurrentTime().Seconds() - pl.StartTime
			} else {
				pl.TimeAlive = pl.DeathEvent.Time - pl.StartTime
			}
			allTimelines = append(allTimelines, *pl)
		}
	})

	// parse to end
	if err := parser.ParseToEnd(); err != nil {
		return fmt.Errorf("eroare la parsare dem: %w", err)
	}

	// salvăm totul în JSON
	if err := saveTimelinesJSON(allTimelines, outputPath); err != nil {
		return err
	}

	fmt.Println("Timeline salvat în:", outputPath)
	return nil
}

// ----------------------------------------------------------------------------------
// Alte funcții helper
// ----------------------------------------------------------------------------------

// E.g. scrie fișierul JSON final
func saveTimelinesJSON(timelines []PlayerTimeline, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("nu pot crea fișier JSON %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(timelines); err != nil {
		return fmt.Errorf("eroare scriere JSON: %w", err)
	}
	return nil
}

// Convertește Team la string
func TeamToString(t common.Team) string {
	switch t {
	case common.TeamTerrorists:
		return "Terrorists"
	case common.TeamCounterTerrorists:
		return "CounterTerrorists"
	default:
		return "Unknown"
	}
}

// Convertește Team la "T"/"CT"
func sideString(t common.Team) string {
	switch t {
	case common.TeamTerrorists:
		return "T"
	case common.TeamCounterTerrorists:
		return "CT"
	default:
		return ""
	}
}

// Returnează numele armei
func getWeaponName(p *common.Player) string {
	if p.ActiveWeapon() != nil {
		return p.ActiveWeapon().Type.String()
	}
	return ""
}

// Determină starea jucătorului (Idle, Walk, Crouch etc.)
func getPlayerAction(p *common.Player) string {
	vx := p.Velocity().X
	vy := p.Velocity().Y
	moving := (vx != 0 || vy != 0)

	switch {
	case p.IsAirborne():
		return "Jump"
	case p.IsDucking() && moving:
		return "Crouch Walk"
	case p.IsDucking():
		return "Crouch"
	case p.IsScoped() && moving:
		return "Scoped Walk"
	case p.IsScoped():
		return "Scoped"
	case moving:
		return "Walk"
	default:
		return "Idle"
	}
}

// Numără coechipierii apropiați (sub 500 units)
func countNearbyTeammates(player *common.Player, gs dem.GameState) int {
	if player == nil {
		return 0
	}
	count := 0
	for _, teammate := range gs.Participants().Playing() {
		if teammate == nil || teammate == player {
			continue
		}
		if teammate.Team == player.Team && player.Team != common.TeamUnassigned {
			dist := player.Position().Distance(teammate.Position())
			if dist < 500 {
				count++
			}
		}
	}
	return count
}

// ------------------------
// Zone logic
// ------------------------
type MapZone struct {
	Name  string  `json:"name"`
	XMin  float64 `json:"x_min"`
	XMax  float64 `json:"x_max"`
	YMin  float64 `json:"y_min"`
	YMax  float64 `json:"y_max"`
	ZMin  float64 `json:"z_min"`
	ZMax  float64 `json:"z_max"`
}

// Încarcă zonele dintr-un fișier JSON
func LoadZones(path string) []MapZone {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("nu pot deschide fișierul cu zone: %v", err)
	}
	defer file.Close()

	var zones []MapZone
	if err := json.NewDecoder(file).Decode(&zones); err != nil {
		log.Fatalf("eroare decodare zone: %v", err)
	}
	log.Printf("Au fost încărcate %d zone din %s\n", len(zones), path)
	return zones
}

// Determină numele zonei pe baza X, Y
func GetZoneName(x, y, z float64, zones []MapZone) string {
	for _, zone := range zones {
		if x >= zone.XMin && x <= zone.XMax &&
		   y >= zone.YMin && y <= zone.YMax &&
		   z >= zone.ZMin && z <= zone.ZMax {
			return zone.Name
		}
	}
	return "Unknown"
}


