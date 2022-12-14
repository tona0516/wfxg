package main

import "wfxg/service"

func main() {
	statsService := service.StatsService{
		InstallPath: "./",
	}
	statsService.GetsStats()
}
