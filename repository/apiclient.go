package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiClient[T any] struct {
}

func (a *ApiClient[T]) GetRequest(url string) (T, error) {
	res, err := http.Get(url)
	if res != nil {
		defer res.Body.Close()
	}

	var empty T
	if err != nil {
		return empty, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return empty, err
	}

	fmt.Println(string(body))

	var response T
	err = json.Unmarshal(body, &response)
	if err != nil {
		return empty, err
	}

	return response, nil
}
