package vo

type WGEncyclopediaShips struct {
	Status string `json:"status"`
	Meta   struct {
		Count     int `json:"count"`
		PageTotal int `json:"page_total"`
		Total     int `json:"total"`
		Limit     int `json:"limit"`
		Page      int `json:"page"`
	} `json:"meta"`
	Data map[int]struct {
		Tier   int    `json:"tier"`
		Type   string `json:"type"`
		Name   string `json:"name"`
		Nation string `json:"nation"`
	} `json:"data"`
}
