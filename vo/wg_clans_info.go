package vo

type WGClansInfo struct {
	Status string `json:"status"`
	Data   map[int]struct {
		Tag string `json:"tag"`
	} `json:"data"`
}
