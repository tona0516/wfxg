package vo

type PlayerShipInfo struct {
	Name     string
	Nation   string
	Tier     int
	Type     string
	StatsURL string
}

type PlayerShipStats struct {
	Battles        int
	AvgDamage      int
	AvgExp         int
	WinRate        float32
	KdRate         float32
	CombatPower    int
	PersonalRating int
}

type PlayerPlayerInfo struct {
	Name     string
	Clan     string
	IsHidden string
	StatsURL string
}

type PlayerPlayerStats struct {
	Battles   int
	AvgDamage int
	AvgExp    int
	WinRate   float32
	KdRate    float32
	AvgTier   float32
}

type Player struct {
	ShipInfo    PlayerShipInfo
	ShipStats   PlayerShipStats
	PlayerInfo  PlayerPlayerInfo
	PlayerStats PlayerPlayerStats
}
