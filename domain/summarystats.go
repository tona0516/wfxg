package domain

type Stats struct {
	Battles         int
	SurvivedBattles int
	DamageDealt     int
	Frags           int
	Xp              int
	Wins            int
}
type SummaryStats struct {
	Ship   Stats
	Player Stats
}

func (s *SummaryStats) ShipAvgDamage() float64 {
	if s.Ship.Battles != 0 {
		return float64(s.Ship.DamageDealt) / float64(s.Ship.Battles)
	}
	return 0
}

func (s *SummaryStats) ShipAvgExp() float64 {
	if s.Ship.Battles != 0 {
		return float64(s.Ship.Xp) / float64(s.Ship.Battles)
	}
	return 0
}

func (s *SummaryStats) ShipKdRate() float64 {
	if s.Ship.Battles-s.Ship.SurvivedBattles > 0 {
		return float64(s.Ship.Frags) / float64(s.Ship.Battles-s.Ship.SurvivedBattles)
	}
	return 0
}

func (s *SummaryStats) ShipAvgFrags() float64 {
	if s.Ship.Battles > 0 {
		return float64(s.Ship.Frags) / float64(s.Ship.Battles)
	}
	return 0
}

func (s *SummaryStats) ShipWinRate() float64 {
	if s.Ship.Battles != 0 {
		return float64(s.Ship.Wins) / float64(s.Ship.Battles) * 100
	}
	return 0
}

func (s *SummaryStats) PlayerAvgDamage() float64 {
	if s.Player.Battles > 0 {
		return float64(s.Player.DamageDealt) / float64(s.Player.Battles)
	}
	return 0
}

func (s *SummaryStats) PlayerAvgExp() float64 {
	if s.Player.Battles > 0 {
		return float64(s.Player.Xp) / float64(s.Player.Battles)
	}
	return 0
}

func (s *SummaryStats) PlayerKdRate() float64 {
	if s.Player.Battles-s.Player.SurvivedBattles > 0 {
		return float64(s.Player.Frags) / float64(s.Player.Battles-s.Player.SurvivedBattles)
	}
	return 0
}

func (s *SummaryStats) PlayerWinRate() float64 {
	if s.Player.Battles != 0 {
		return float64(s.Player.Wins) / float64(s.Player.Battles) * 100
	}
	return 0
}
