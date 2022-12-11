package repository

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AccountInfo struct {
	Status string `json:"status"`
	Meta   struct {
		Count  int         `json:"count"`
		Hidden interface{} `json:"hidden"`
	} `json:"meta"`
	Data map[string]struct {
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

type Wargaming struct {
}

func (w *Wargaming) GetAccountInfo(accountID []string) (AccountInfo, error) {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "api.worldofwarships.asia"
	u.Path = "wows/account/info/"
	q := u.Query()
	q.Set("application_id", "3bd34ff346625bf01cc8ba6a9204dd16")
	q.Set("account_id", strings.Join(accountID, ","))
	q.Set("fields", strings.Join([]string{
		"statistics.pvp.xp",
		"statistics.pvp.survived_battles",
		"statistics.pvp.battles",
		"statistics.pvp.frags",
		"statistics.pvp.wins",
		"statistics.pvp.damage_dealt",
	}, ","))
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return AccountInfo{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return AccountInfo{}, err
	}

	var accountInfo AccountInfo
	err = json.Unmarshal(body, &accountInfo)
	if err != nil {
		return AccountInfo{}, err
	}

	return accountInfo, nil
}
