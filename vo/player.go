package vo

type Player struct {
	ShipInfo struct {
		Name     string
		Nation   string
		Tier     int
		Type     string
		StatsURL string
	}
	ShipStats struct {
		Battles        int
		AvgDamage      int
		AvgExp         int
		WinRate        float32
		CombatPower    int
		PersonalRating int
	}
	PlayerInfo struct {
		Name     string
		Clan     string
		IsHidden string
		StatsURL string
	}
	PlayerStats struct {
		Battles   int
		AvgDamage int
		AvgExp    int
		WinRate   float32
		KdRate    float32
		avgTier   float32
	}
}
