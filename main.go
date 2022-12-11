package main

import (
	"fmt"
	"wfxg/repository"
)

func main() {
	wargaming := repository.Wargaming{}
	accountInfo, err := wargaming.GetAccountInfo([]string{"2010342809", "2030131054"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(accountInfo)
}
