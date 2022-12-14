package service

import (
	"fmt"
	"sync"
	"wfxg/domain"
	"wfxg/repo"
	"wfxg/vo"
)

type StatsService struct {
	InstallPath string
}

func (s *StatsService) GetsStats() {
	wargaming := repo.Wargaming{}
	numbers := repo.Numbers{}
	local := repo.Local{}

	tempArenaInfo, err := local.GetTempArenaInfo(s.InstallPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	accountIDs, accountList, err := fetchAccountList(&wargaming, tempArenaInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	accountInfoResult := make(chan vo.Result[vo.WGAccountInfo])
	shipStatsResult := make(chan vo.Result[map[int]vo.WGShipsStats])
	clanTagResult := make(chan vo.Result[map[int]string])
	shipInfoResult := make(chan vo.Result[map[int]vo.ShipInfo])
	expectedStatsResult := make(chan vo.Result[vo.NSExpectedStats])

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
		accountInfo.Result,
		accountList,
		clanTag.Result,
		shipStats.Result,
		shipInfo.Result,
		expectedStats.Result,
	)
}

func fetchAccountList(wargaming *repo.Wargaming, tempArenaInfo vo.TempArenaInfo) ([]int, vo.WGAccountList, error) {
	accountNames := tempArenaInfo.AccountNames()

	accountList, err := wargaming.GetAccountList(accountNames)
	if err != nil {
		return nil, accountList, err
	}

	accountIDs := accountList.AccountIDs()

	return accountIDs, accountList, nil
}

func fetchAccountInfo(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[vo.WGAccountInfo]) {
	accountInfo, err := wargaming.GetAccountInfo(accountIDs)
	if err != nil {
		result <- vo.Result[vo.WGAccountInfo]{Result: accountInfo, Error: err}
		return
	}

	result <- vo.Result[vo.WGAccountInfo]{Result: accountInfo, Error: nil}
}

func fetchShipStats(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[map[int]vo.WGShipsStats]) {
	shipStatsMap := make(map[int]vo.WGShipsStats)
	limit := make(chan struct{}, 5)
	wg := sync.WaitGroup{}
	for i := range accountIDs {
		limit <- struct{}{}
		wg.Add(1)
		go func(accountID int) {
			defer func() {
				wg.Done()
				<-limit
			}()

			shipStats, err := wargaming.GetShipsStats(accountID)
			if err != nil {
				result <- vo.Result[map[int]vo.WGShipsStats]{Result: shipStatsMap, Error: err}
				return
			}

			shipStatsMap[accountID] = shipStats
		}(accountIDs[i])
	}
	wg.Wait()

	result <- vo.Result[map[int]vo.WGShipsStats]{Result: shipStatsMap, Error: nil}
}

func fetchClanTag(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[map[int]string]) {
	clanTagMap := make(map[int]string)

	clansAccountInfo, err := wargaming.GetClansAccountInfo(accountIDs)
	if err != nil {
		result <- vo.Result[map[int]string]{Result: clanTagMap, Error: err}
		return
	}

	clanIDs := clansAccountInfo.ClanIDs()

	clansInfo, err := wargaming.GetClansInfo(clanIDs)
	if err != nil {
		result <- vo.Result[map[int]string]{Result: clanTagMap, Error: err}
		return
	}

	for i := range accountIDs {
		accountID := accountIDs[i]
		clanID := clansAccountInfo.Data[accountID].ClanID
		clanTag := clansInfo.Data[clanID].Tag
		clanTagMap[accountID] = clanTag
	}

	result <- vo.Result[map[int]string]{Result: clanTagMap, Error: nil}
}

func fetchShipInfo(wargaming *repo.Wargaming, result chan vo.Result[map[int]vo.ShipInfo]) {
	shipInfoMap := make(map[int]vo.ShipInfo, 0)
	res, err := wargaming.GetEncyclopediaShips(1)
	if err != nil {
		result <- vo.Result[map[int]vo.ShipInfo]{Result: shipInfoMap, Error: err}
		return
	}
	pageTotal := res.Meta.PageTotal

	var mu sync.Mutex
	limit := make(chan struct{}, 5)
	wg := sync.WaitGroup{}
	for i := 1; i <= pageTotal; i++ {
		limit <- struct{}{}
		wg.Add(1)
		go func(pageNo int) {
			defer func() {
				wg.Done()
				<-limit
			}()

			encyclopediaShips, err := wargaming.GetEncyclopediaShips(pageNo)
			if err != nil {
				result <- vo.Result[map[int]vo.ShipInfo]{Result: shipInfoMap, Error: err}
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
		}(i)
	}
	wg.Wait()

	result <- vo.Result[map[int]vo.ShipInfo]{Result: shipInfoMap, Error: nil}
}

func fetchExpectedStats(numbers *repo.Numbers, result chan vo.Result[vo.NSExpectedStats]) {
	expectedStats, err := numbers.Get()
	if err != nil {
		result <- vo.Result[vo.NSExpectedStats]{Result: *expectedStats, Error: err}
		return
	}

	result <- vo.Result[vo.NSExpectedStats]{Result: *expectedStats, Error: err}
}

func calculateAvgTier(accountID int, shipInfo map[int]vo.ShipInfo, shipStats map[int]vo.WGShipsStats) float64 {
	sum := 0
	battles := 0
	playerShipStats := shipStats[accountID].Data[accountID]
	for i := range playerShipStats {
		shipID := playerShipStats[i].ShipID
		tier := shipInfo[shipID].Tier
		sum += playerShipStats[i].Pvp.Battles * tier
		battles += playerShipStats[i].Pvp.Battles
	}

	if battles == 0 {
		return 0
	} else {
		return float64(sum) / float64(battles)
	}
}

func compose(
	tempArenaInfo vo.TempArenaInfo,
	accountInfo vo.WGAccountInfo,
	accountList vo.WGAccountList,
	clanTag map[int]string,
	shipStats map[int]vo.WGShipsStats,
	shipInfo map[int]vo.ShipInfo,
	expectedStats vo.NSExpectedStats,
) {
	friends := make([]vo.Player, 0)
	enemies := make([]vo.Player, 0)
	rating := domain.Rating{}

	for i := range tempArenaInfo.Vehicles {
		vehicle := tempArenaInfo.Vehicles[i]
		playerShipInfo := shipInfo[vehicle.ShipID]

		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clanTag[accountID]

		var summaryStats domain.SummaryStats
		playerAccountInfo := accountInfo.Data[accountID]
		for k := range shipStats[accountID].Data[accountID] {
			playerShipStats := shipStats[accountID].Data[accountID][k]
			if playerShipStats.ShipID == vehicle.ShipID {
				summaryStats = domain.SummaryStats{
					Player: domain.Stats{
						Battles:         playerAccountInfo.Statistics.Pvp.Battles,
						SurvivedBattles: playerAccountInfo.Statistics.Pvp.SurvivedBattles,
						DamageDealt:     playerAccountInfo.Statistics.Pvp.DamageDealt,
						Xp:              playerAccountInfo.Statistics.Pvp.Xp,
						Frags:           playerAccountInfo.Statistics.Pvp.Frags,
						Wins:            playerAccountInfo.Statistics.Pvp.Wins,
					},
					Ship: domain.Stats{
						Battles:         playerShipStats.Pvp.Battles,
						SurvivedBattles: playerShipStats.Pvp.SurvivedBattles,
						DamageDealt:     playerShipStats.Pvp.DamageDealt,
						Xp:              playerShipStats.Pvp.Xp,
						Frags:           playerShipStats.Pvp.Frags,
						Wins:            playerShipStats.Pvp.Wins,
					},
				}
				break
			}
		}

		expectedShipStats := expectedStats.Data[vehicle.ShipID]

		player := vo.Player{
			ShipInfo: vo.PlayerShipInfo{
				Name:   playerShipInfo.Name,
				Nation: playerShipInfo.Nation,
				Tier:   playerShipInfo.Tier,
				Type:   playerShipInfo.Type,
			},
			ShipStats: vo.PlayerShipStats{
				Battles:   summaryStats.Player.Battles,
				AvgDamage: int(summaryStats.PlayerAvgDamage()),
				AvgExp:    int(summaryStats.PlayerAvgExp()),
				WinRate:   float32(summaryStats.PlayerWinRate()),
				KdRate:    float32(summaryStats.PlayerKdRate()),
				CombatPower: rating.CombatPower(
					summaryStats.ShipAvgDamage(),
					summaryStats.ShipKdRate(),
					summaryStats.ShipAvgExp(),
					playerShipInfo.Tier,
					playerShipInfo.Type,
				),
				PersonalRating: rating.PersonalRating(
					summaryStats.ShipAvgDamage(),
					summaryStats.ShipAvgFrags(),
					summaryStats.ShipWinRate(),
					expectedShipStats.AverageDamageDealt,
					expectedShipStats.AverageFrags,
					expectedShipStats.WinRate,
				),
			},
			PlayerInfo: vo.PlayerPlayerInfo{
				Name: nickname,
				Clan: clan,
			},
			PlayerStats: vo.PlayerPlayerStats{
				Battles:   summaryStats.Player.Battles,
				AvgDamage: int(summaryStats.PlayerAvgDamage()),
				AvgExp:    int(summaryStats.PlayerAvgExp()),
				WinRate:   float32(summaryStats.PlayerWinRate()),
				KdRate:    float32(summaryStats.PlayerKdRate()),
				AvgTier:   float32(calculateAvgTier(accountID, shipInfo, shipStats)),
			},
		}

		if vehicle.Relation == 0 || vehicle.Relation == 1 {
			friends = append(friends, player)
		} else {
			enemies = append(enemies, player)
		}

		fmt.Println(player)
	}
}
