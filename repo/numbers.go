package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"wfxg/vo"
)

type Numbers struct {
}

func (n *Numbers) Get() (vo.ExpectedStats, error) {
	res, err := http.Get("https://api.wows-numbers.com/personal/rating/expected/json/")
	if res != nil {
		defer res.Body.Close()
	}

	var empty vo.ExpectedStats
	if err != nil {
		fmt.Println(err)
		return empty, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return empty, err
	}

	data1 := make(map[string]interface{})

	err = json.Unmarshal(body, &data1)
	if err != nil {
		fmt.Println(err)
		return empty, err
	}

	time := data1["time"].(float64)
	data2 := data1["data"].(map[string]interface{})
	data := make(map[int]vo.ExpectedStatsData)
	for key, value := range data2 {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		valueMap, ok := value.(vo.ExpectedStatsData)
		if !ok {
			continue
		}

		data[keyInt] = valueMap
	}

	response := vo.ExpectedStats{
		Time: int(time),
		Data: data,
	}

	return response, nil
}
