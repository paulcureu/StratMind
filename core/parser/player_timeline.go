package parser

type PlayerTimeline struct {
	RoundNumber int    `json:"round_number"`
	MatchID     string `json:"match_id,omitempty"`
	SteamID     string `json:"steam_id"`
	Nickname    string `json:"nickname"`
	Team        string `json:"team"`
	Side        string `json:"side"` // "T" sau "CT"

	RoleInRound    string  `json:"role_in_round"`
	RoleConfidence float64 `json:"role_confidence"`

	StartPosition Position       `json:"start_position"`
	StartTick     int            `json:"start_tick"`
	StartTime     float64        `json:"start_time"`
	Path          []PositionTick `json:"path"`

	UtilityThrown UtilityStats `json:"utility_thrown"`
	FlashAssists  int          `json:"flash_assists"`

	KillEvents       []KillEvent `json:"kill_events"`
	DeathEvent       *DeathEvent `json:"death_event,omitempty"`
	TradedTeammate   bool        `json:"traded_teammate"`
	TradedByTeammate bool        `json:"traded_by_teammate"`

	Survived        bool    `json:"survived"`
	TimeAlive       float64 `json:"time_alive"`
	WasFirstContact bool    `json:"was_first_contact"`
	Clutched        bool    `json:"clutched"`

	ImpactRating float64 `json:"impact_rating"`
	Notes        string  `json:"notes"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type PositionTick struct {
	Tick          int      `json:"tick"`
	Time          float64  `json:"time"`
	Position      Position `json:"position"`
	Zone          string   `json:"zone"` // ← Zone-ul jucătorului
	Action        string   `json:"action"`
	WeaponHeld    string   `json:"weapon_held"`
	IsScoped      bool     `json:"is_scoped"`
	IsDucking     bool     `json:"is_ducking"`
	IsMoving      bool     `json:"is_moving"`
	IsAirborne    bool     `json:"is_airborne"`
	HP            int      `json:"hp"`
	ViewAngle     float64  `json:"view_angle,omitempty"`
	InCombat      bool     `json:"in_combat"`
	NearTeammates int      `json:"near_teammates,omitempty"`
}

type UtilityStats struct {
	Smokes     []string `json:"smokes"`
	Flashes    int      `json:"flashes"`
	Molotovs   []string `json:"molotovs"`
	HEGrenades int      `json:"he_grenades"`
}

type KillEvent struct {
	Tick       int     `json:"tick"`
	Time       float64 `json:"time"`
	Victim     string  `json:"victim"`
	Weapon     string  `json:"weapon"`
	Headshot   bool    `json:"headshot"`
	AssistedBy string  `json:"assisted_by,omitempty"`
}

type DeathEvent struct {
	Tick     int     `json:"tick"`
	Time     float64 `json:"time"`
	Killer   string  `json:"killer"`
	TradedBy string  `json:"traded_by,omitempty"`
}
