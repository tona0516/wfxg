package main

import (
	"fmt"
	"strings"
	"sync"
	"wfxg/domain"
	"wfxg/repo"
	"wfxg/vo"
)

func main() {
	wargaming := repo.Wargaming{}
	numbers := repo.Numbers{}
	local := repo.Local{}

	tempArenaInfo, err := local.GetTempArenaInfo("./")
	if err != nil {
		fmt.Println(err)
		return
	}

	accountIDs, accountList, err := fetchAccount(&wargaming, tempArenaInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	accountInfoResult := make(chan vo.Result[*vo.AccountInfo], 1)
	shipStatsResult := make(chan vo.Result[map[int]vo.ShipsStats], 1)
	clanTagResult := make(chan vo.Result[map[int]string], 1)
	shipInfoResult := make(chan vo.Result[map[int]vo.ShipInfo], 1)
	expectedStatsResult := make(chan vo.Result[*vo.ExpectedStats], 1)

	go fetchAccountInfo(&wargaming, accountIDs, accountInfoResult)
	go fetchShipStats(&wargaming, accountIDs, shipStatsResult)
	go fetchClanTag(&wargaming, accountIDs, clanTagResult)
	go fetchShipInfo(&wargaming, shipInfoResult)
	go fetchExpectedStats(&numbers, expectedStatsResult)

	accountInfo := <-accountInfoResult
	if accountInfo.Error != nil {
		fmt.Println(accountInfo.Error)
		return
	}
	shipStats := <-shipStatsResult
	if shipStats.Error != nil {
		fmt.Println(shipStats.Error)
		return
	}
	clanTag := <-clanTagResult
	if clanTag.Error != nil {
		fmt.Println(clanTag.Error)
		return
	}
	shipInfo := <-shipInfoResult
	if shipInfo.Error != nil {
		fmt.Println(shipInfo.Error)
		return
	}
	expectedStats := <-expectedStatsResult
	if expectedStats.Error != nil {
		fmt.Println(expectedStats.Error)
		return
	}

	compose(
		tempArenaInfo,
		*accountInfo.Result,
		*accountList,
		clanTag.Result,
		shipStats.Result,
		shipInfo.Result,
		*expectedStats.Result,
	)
}

func fetchAccount(wargaming *repo.Wargaming, tempArenaInfo vo.TempArenaInfo) ([]int, *vo.AccountList, error) {
	accountNames := make([]string, 0)
	for i := range tempArenaInfo.Vehicles {
		vehicle := tempArenaInfo.Vehicles[i]
		if strings.HasPrefix(vehicle.Name, ":") && strings.HasSuffix(vehicle.Name, ":") {
			continue
		}

		accountNames = append(accountNames, vehicle.Name)
	}

	accountList, err := wargaming.GetAccountList(accountNames)
	if err != nil {
		return nil, nil, err
	}

	accountIDs := make([]int, 0)
	for i := range accountList.Data {
		accountIDs = append(accountIDs, accountList.Data[i].AccountID)
	}

	return accountIDs, &accountList, nil
}

func fetchAccountInfo(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[*vo.AccountInfo]) {
	accountInfo, err := wargaming.GetAccountInfo(accountIDs)
	if err != nil {
		result <- vo.Result[*vo.AccountInfo]{nil, err}
		return
	}

	result <- vo.Result[*vo.AccountInfo]{&accountInfo, err}
}

func fetchShipStats(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[map[int]vo.ShipsStats]) {
	shipStatsMap := make(map[int]vo.ShipsStats)
	limit := make(chan struct{}, 5)
	wg := sync.WaitGroup{}
	for i := range accountIDs {
		wg.Add(1)
		accountID := accountIDs[i]
		go func() {
			defer func() {
				<-limit
				wg.Done()
			}()

			limit <- struct{}{}
			shipStats, err := wargaming.GetShipsStats(accountID)
			if err != nil {
				result <- vo.Result[map[int]vo.ShipsStats]{nil, err}
				return
			}

			shipStatsMap[accountID] = shipStats
		}()
	}
	wg.Wait()

	result <- vo.Result[map[int]vo.ShipsStats]{shipStatsMap, nil}
}

func fetchClanTag(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[map[int]string]) {
	clansAccountInfo, err := wargaming.GetClansAccountInfo(accountIDs)
	if err != nil {
		result <- vo.Result[map[int]string]{nil, err}
		return
	}

	clanIDs := make([]int, 0)
	for i := range clansAccountInfo.Data {
		clanID := clansAccountInfo.Data[i].ClanID
		if clanID != 0 {
			clanIDs = append(clanIDs, clansAccountInfo.Data[i].ClanID)
		}
	}

	clansInfo, err := wargaming.GetClansInfo(clanIDs)
	if err != nil {
		result <- vo.Result[map[int]string]{nil, err}
		return
	}

	clanTagMap := make(map[int]string)
	for i := range accountIDs {
		accountID := accountIDs[i]
		clanID := clansAccountInfo.Data[accountID].ClanID
		clanTag := clansInfo.Data[clanID].Tag
		clanTagMap[accountID] = clanTag
	}

	result <- vo.Result[map[int]string]{clanTagMap, err}
}

func fetchShipInfo(wargaming *repo.Wargaming, result chan vo.Result[map[int]vo.ShipInfo]) {
	shipInfoMap := make(map[int]vo.ShipInfo, 0)
	res, err := wargaming.GetEncyclopediaShips(1)
	if err != nil {
		result <- vo.Result[map[int]vo.ShipInfo]{nil, err}
		return
	}
	pageTotal := res.Meta.PageTotal

	limit := make(chan struct{}, 5)
	wg := sync.WaitGroup{}
	var mu sync.Mutex
	for i := 1; i <= pageTotal; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer func() {
				<-limit
				wg.Done()
			}()

			limit <- struct{}{}
			encyclopediaShips, err := wargaming.GetEncyclopediaShips(i)
			if err != nil {
				result <- vo.Result[map[int]vo.ShipInfo]{nil, err}
				return
			}

			for shipID, shipInfo := range encyclopediaShips.Data {
				mu.Lock()
				shipInfoMap[shipID] = vo.ShipInfo{
					Name:   shipInfo.Name,
					Tier:   shipInfo.Tier,
					Type:   shipInfo.Type,
					Nation: shipInfo.Nation,
				}
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	result <- vo.Result[map[int]vo.ShipInfo]{shipInfoMap, err}
}

func fetchExpectedStats(numbers *repo.Numbers, result chan vo.Result[*vo.ExpectedStats]) {
	expectedStats, err := numbers.Get()
	if err != nil {
		result <- vo.Result[*vo.ExpectedStats]{nil, err}
		return
	}

	result <- vo.Result[*vo.ExpectedStats]{expectedStats, err}
}

func compose(
	tempArenaInfo vo.TempArenaInfo,
	accountInfo vo.AccountInfo,
	accountList vo.AccountList,
	clanTag map[int]string,
	shipStats map[int]vo.ShipsStats,
	shipInfo map[int]vo.ShipInfo,
	expectedStats vo.ExpectedStats,
) {
	// friends := make([]vo.Player, 0)
	// enemies := make([]vo.Player, 0)
	rating := domain.Rating{}

	for i := range tempArenaInfo.Vehicles {
		vehicle := tempArenaInfo.Vehicles[i]
		playerShipInfo := shipInfo[vehicle.ShipID]

		var accountID int
		for j := range accountList.Data {
			item := accountList.Data[j]
			if item.NickName == vehicle.Name {
				accountID = item.AccountID
				break
			}
		}

		playerAccountInfo := accountInfo.Data[accountID]
		var playerAvgDamage int
		var playerKdRate float32
		var playerAvgExp int
		var playerWinRate float32
		if playerAccountInfo.Statistics.Pvp.Battles != 0 {
			playerAvgDamage = playerAccountInfo.Statistics.Pvp.DamageDealt / playerAccountInfo.Statistics.Pvp.Battles
			playerKdRate = float32(playerAccountInfo.Statistics.Pvp.Frags) / float32(playerAccountInfo.Statistics.Pvp.Battles-playerAccountInfo.Statistics.Pvp.SurvivedBattles)
			playerAvgExp = playerAccountInfo.Statistics.Pvp.Xp / playerAccountInfo.Statistics.Pvp.Battles
			playerWinRate = float32(playerAccountInfo.Statistics.Pvp.Wins) / float32(playerAccountInfo.Statistics.Pvp.Battles) * 100
		}

		var playerShipAvgDamage int
		var playerShipKdRate float32
		var playerShipAvgExp int
		var playerShipWinRate float32
		var playerShipAvgFrags float32
		for k := range shipStats[accountID].Data[accountID] {
			playerShipStats := shipStats[accountID].Data[accountID][k]
			if playerShipStats.ShipID == vehicle.ShipID {
				if playerShipStats.Pvp.Battles != 0 {
					playerShipAvgDamage = playerShipStats.Pvp.DamageDealt / playerShipStats.Pvp.Battles
					playerShipKdRate = float32(playerShipStats.Pvp.Frags) / float32(playerShipStats.Pvp.Battles-playerShipStats.Pvp.SurvivedBattles)
					playerShipAvgExp = playerShipStats.Pvp.Xp / playerShipStats.Pvp.Battles
					playerShipWinRate = float32(playerShipStats.Pvp.Wins) / float32(playerShipStats.Pvp.Battles) * 100
					playerShipAvgFrags = float32(playerShipStats.Pvp.Frags) / float32(playerShipStats.Pvp.Battles)
				}
				break
			}
		}

		combatPower := rating.CombatPower(
			float64(playerShipAvgDamage),
			float64(playerKdRate),
			float64(playerAvgExp),
			playerShipInfo.Tier,
			playerShipInfo.Type,
		)

		expectedShipStats := expectedStats.Data[vehicle.ShipID]
		personalRating := rating.PersonalRating(
			float64(playerShipAvgDamage),
			float64(playerShipAvgFrags),
			float64(playerShipWinRate),
			expectedShipStats.AverageDamageDealt,
			expectedShipStats.AverageFrags,
			expectedShipStats.WinRate,
		)

		fmt.Println(vehicle.Name, accountID, playerShipInfo.Name)
		fmt.Println(combatPower, personalRating)
		fmt.Println(playerAvgDamage, playerKdRate, playerAvgExp, playerWinRate)
		fmt.Println(playerShipAvgDamage, playerShipKdRate, playerShipAvgExp, playerShipWinRate)
		fmt.Println()
	}
}
