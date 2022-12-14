package vo

type WGAccountList struct {
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

func (w *WGAccountList) AccountIDs() []int {
	accountIDs := make([]int, 0)
	for i := range w.Data {
		accountID := w.Data[i].AccountID
		if accountID != 0 {
			accountIDs = append(accountIDs, w.Data[i].AccountID)
		}
	}
	return accountIDs
}
