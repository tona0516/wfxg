package vo

type ExpectedStats struct {
	Time int                       `json:"time"`
	Data map[int]ExpectedStatsData `json:"data"`
}

type ExpectedStatsData struct {
	AverageDamageDealt float64 `json:"average_damage_dealt"`
	AverageFrags       float64 `json:"average_frags"`
	WinRate            float64 `json:"win_rate"`
}
