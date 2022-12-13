package vo

type ClansAccountInfo struct {
	Status string `json:"status"`
	Data   map[int]struct {
		ClanID int `json:"clan_id"`
	} `json:"data"`
}
