package main

import (
	"fmt"
	"wfxg/repository"
)

func main() {
	wargaming := repository.Wargaming{}
	numbers := repository.Numbers{}
	local := repository.Local{}

	res, err := wargaming.GetAccountInfo([]string{"2010342809", "2030131054"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

	res2, err := wargaming.GetAccountList([]string{"tonango", "MTDroine"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res2)

	res3, err := wargaming.GetEncyclopediaShips(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res3)

	res4, err := wargaming.GetShipsStats("2010342809")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res4)

	res5, err := numbers.Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res5)

	f, err := local.GetTempArenaInfo("./")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(f)
}
