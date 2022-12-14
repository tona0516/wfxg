package repo

import (
	"net/url"
	"strconv"
	"strings"
	"wfxg/vo"
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

type Wargaming struct {
}

func (w *Wargaming) GetAccountInfo(accountIDs []int) (vo.WGAccountInfo, error) {
	accountIDsString := make([]string, 0)
	for i := range accountIDs {
		accountIDsString = append(accountIDsString, strconv.Itoa(accountIDs[i]))
	}
	u := buildUrl(
		"/wows/account/info/",
		map[string]string{
			"application_id": appid,
			"account_id":     strings.Join(accountIDsString, ","),
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

	client := ApiClient[vo.WGAccountInfo]{}
	return client.GetRequest(u.String())
}

func (w *Wargaming) GetAccountList(accountNames []string) (vo.WGAccountList, error) {
	u := buildUrl(
		"/wows/account/list/",
		map[string]string{
			"application_id": "3bd34ff346625bf01cc8ba6a9204dd16",
			"search":         strings.Join(accountNames, ","),
			"fields":         strings.Join([]string{"account_id", "nickname"}, ","),
			"type":           "exact",
		},
	)

	client := ApiClient[vo.WGAccountList]{}
	return client.GetRequest(u.String())
}

func (w *Wargaming) GetClansAccountInfo(accountIDs []int) (vo.WGClansAccountInfo, error) {
	accountIDsString := make([]string, 0)
	for i := range accountIDs {
		accountIDsString = append(accountIDsString, strconv.Itoa(accountIDs[i]))
	}

	u := buildUrl(
		"/wows/clans/accountinfo/",
		map[string]string{
			"application_id": appid,
			"account_id":     strings.Join(accountIDsString, ","),
			"fields":         "clan_id",
		},
	)

	client := ApiClient[vo.WGClansAccountInfo]{}
	return client.GetRequest(u.String())
}

func (w *Wargaming) GetClansInfo(clanIDs []int) (vo.WGClansInfo, error) {
	clanIDsString := make([]string, 0)
	for i := range clanIDs {
		clanIDsString = append(clanIDsString, strconv.Itoa(clanIDs[i]))
	}

	u := buildUrl(
		"/wows/clans/info/",
		map[string]string{
			"application_id": appid,
			"clan_id":        strings.Join(clanIDsString, ","),
			"fields":         "tag",
		},
	)

	client := ApiClient[vo.WGClansInfo]{}
	return client.GetRequest(u.String())
}

func (w *Wargaming) GetEncyclopediaShips(pageNo int) (vo.WGEncyclopediaShips, error) {
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

	client := ApiClient[vo.WGEncyclopediaShips]{}
	return client.GetRequest(u.String())
}

func (w *Wargaming) GetShipsStats(accountID int) (vo.WGShipsStats, error) {
	u := buildUrl(
		"/wows/ships/stats/",
		map[string]string{
			"application_id": appid,
			"account_id":     strconv.Itoa(accountID),
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

	client := ApiClient[vo.WGShipsStats]{}
	return client.GetRequest(u.String())
}
