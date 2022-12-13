package repository

import (
	"net/url"
	"strconv"
	"strings"
)

const appid = "3bd34ff346625bf01cc8ba6a9204dd16"

func buildUrl(path string, query map[string]string) *url.URL {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "api.worldofwarships.asia"
	u.Path = path
	q := u.Query()
	for key, value := range query {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u
}

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

type EncyclopediaShips struct {
	Status string `json:"status"`
	Meta   struct {
		Count     int `json:"count"`
		PageTotal int `json:"page_total"`
		Total     int `json:"total"`
		Limit     int `json:"limit"`
		Page      int `json:"page"`
	} `json:"meta"`
	Data map[string]struct {
		Tier   int    `json:"tier"`
		Type   string `json:"type"`
		Name   string `json:"name"`
		Nation string `json:"nation"`
	} `json:"data"`
}

type ShipsStats struct {
	Status string `json:"status"`
	Meta   struct {
		Count  int         `json:"count"`
		Hidden interface{} `json:"hidden"`
	} `json:"meta"`
	Data map[int][]struct {
		Pvp struct {
			Wins            int `json:"wins"`
			Battles         int `json:"battles"`
			DamageDealt     int `json:"damage_dealt"`
			Xp              int `json:"xp"`
			Frags           int `json:"frags"`
			SurvivedBattles int `json:"survived_battles"`
		} `json:"pvp"`
		ShipID int `json:"ship_id"`
	} `json:"data"`
}

type Wargaming struct {
}

func (w *Wargaming) GetAccountInfo(accountID []string) (AccountInfo, error) {
	u := buildUrl(
		"/wows/account/info/",
		map[string]string{
			"application_id": appid,
			"account_id":     strings.Join(accountID, ","),
			"fields": strings.Join([]string{
				"statistics.pvp.xp",
				"statistics.pvp.survived_battles",
				"statistics.pvp.battles",
				"statistics.pvp.frags",
				"statistics.pvp.wins",
				"statistics.pvp.damage_dealt",
			}, ","),
		},
	)

	client := ApiClient[AccountInfo]{}
	return client.GetRequest(u.String())
}

func (w *Wargaming) GetAccountList(accountNames []string) (AccountList, error) {
	u := buildUrl(
		"/wows/account/list/",
		map[string]string{
			"application_id": "3bd34ff346625bf01cc8ba6a9204dd16",
			"search":         strings.Join(accountNames, ","),
			"fields":         strings.Join([]string{"account_id", "nickname"}, ","),
			"type":           "exact",
		},
	)

	client := ApiClient[AccountList]{}
	return client.GetRequest(u.String())
}

func (w *Wargaming) GetEncyclopediaShips(pageNo int) (EncyclopediaShips, error) {
	u := buildUrl(
		"/wows/encyclopedia/ships/",
		map[string]string{
			"application_id": appid,
			"fields": strings.Join([]string{
				"name",
				"tier",
				"type",
				"nation",
			}, ","),
			"language": "en",
			"page_no":  strconv.Itoa(pageNo),
		},
	)

	client := ApiClient[EncyclopediaShips]{}
	return client.GetRequest(u.String())
}

func (w *Wargaming) GetShipsStats(accountID string) (ShipsStats, error) {
	u := buildUrl(
		"/wows/ships/stats/",
		map[string]string{
			"application_id": appid,
			"account_id":     accountID,
			"fields": strings.Join([]string{
				"ship_id",
				"pvp.wins",
				"pvp.battles",
				"pvp.damage_dealt",
				"pvp.xp",
				"pvp.frags",
				"pvp.survived_battles",
			}, ","),
		},
	)

	client := ApiClient[ShipsStats]{}
	return client.GetRequest(u.String())
}
