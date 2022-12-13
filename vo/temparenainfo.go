package vo

type TempArenaInfo struct {
	Vehicles []struct {
		ShipID   int    `json:"shipId"`
		Relation int    `json:"relation"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
	} `json:"vehicles"`
}
