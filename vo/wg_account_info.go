package vo

type WGAccountInfo struct {
	Status string `json:"status"`
	Meta   struct {
		Count  int         `json:"count"`
		Hidden interface{} `json:"hidden"`
	} `json:"meta"`
	Data map[int]struct {
		Statistics struct {
			Pvp struct {
				Wins            int `json:"wins"`
				Battles         int `json:"battles"`
				DamageDealt     int `json:"damage_dealt"`
				Xp              int `json:"xp"`
				Frags           int `json:"frags"`
				SurvivedBattles int `json:"survived_battles"`
			} `json:"pvp"`
		} `json:"statistics"`
	} `json:"data"`
}
