package vo

type AccountList struct {
	Status string `json:"status"`
	Meta   struct {
		Count  int         `json:"count"`
		Hidden interface{} `json:"hidden"`
	} `json:"meta"`
	Data []struct {
		NickName  string `json:"nickname"`
		AccountID int    `json:"account_id"`
	} `json:"data"`
}
