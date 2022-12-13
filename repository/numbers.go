package repository

type ExpectedStats struct {
	Time int `json:"time"`
	Data map[string]interface {
		// 3655251664が"[]"のためマッピングできない
		// AverageDamageDealt float64 `json:"average_damage_dealt"`
		// AverageFrags       float64 `json:"average_frags"`
		// WinRate            float64 `json:"win_rate"`
	} `json:"data"`
}

type Numbers struct {
}

func (n *Numbers) Get() (ExpectedStats, error) {
	client := ApiClient[ExpectedStats]{}
	return client.GetRequest("https://api.wows-numbers.com/personal/rating/expected/json/")
}
